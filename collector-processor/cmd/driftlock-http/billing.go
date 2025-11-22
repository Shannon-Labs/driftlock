package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/billingportal/session"
	checkoutsession "github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/webhook"
)

func init() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
}

func billingCheckoutHandler(store *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, r, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
			return
		}

		tc, ok := tenantFromContext(r.Context())
		if !ok {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}

		// Parse request body for plan selection
		var req struct {
			Plan string `json:"plan"`
		}
		// Allow empty body for backward compatibility (default to basic)
		if r.Body != nil {
			json.NewDecoder(r.Body).Decode(&req)
		}

		// Determine Price ID based on plan
		var priceID string
		switch req.Plan {
		case "pro":
			priceID = os.Getenv("STRIPE_PRICE_ID_PRO")
		case "basic":
			priceID = os.Getenv("STRIPE_PRICE_ID_BASIC")
		default:
			// Default to basic if not specified, or handle error
			// For now default to basic as the entry paid tier
			priceID = os.Getenv("STRIPE_PRICE_ID_BASIC")
			if priceID == "" {
				// Fallback to Pro if Basic not set (migration path)
				priceID = os.Getenv("STRIPE_PRICE_ID_PRO")
			}
		}

		if priceID == "" {
			writeError(w, r, http.StatusInternalServerError, fmt.Errorf("pricing configuration missing for plan: %s", req.Plan))
			return
		}

		domain := os.Getenv("DOMAIN_URL")
		if domain == "" {
			domain = "https://driftlock.net"
		}

		params := &stripe.CheckoutSessionParams{
			CustomerEmail: stripe.String(tc.Tenant.Name + "@driftlock.net"), // Ideally use actual email if available
			LineItems: []*stripe.CheckoutSessionLineItemParams{
				{
					Price:    stripe.String(priceID),
					Quantity: stripe.Int64(1),
				},
			},
			Mode:              stripe.String(string(stripe.CheckoutSessionModeSubscription)),
			SuccessURL:        stripe.String(domain + "/dashboard?success=true&session_id={CHECKOUT_SESSION_ID}"),
			CancelURL:         stripe.String(domain + "/dashboard?canceled=true"),
			ClientReferenceID: stripe.String(tc.Tenant.ID.String()),
			Metadata: map[string]string{
				"tenant_id": tc.Tenant.ID.String(),
			},
		}

		s, err := checkoutsession.New(params)
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, fmt.Errorf("stripe session creation failed: %w", err))
			return
		}

		writeJSON(w, r, http.StatusOK, map[string]string{
			"url": s.URL,
		})
	}
}

func billingPortalHandler(store *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ensure only authenticated users (via Firebase Auth for dashboard) access this
		// Or allow via API key if we want CLI users to manage billing?
		// For now, assume dashboard user session via Firebase Auth middleware

		tc, ok := tenantFromContext(r.Context())
		if !ok {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}

		// We need the Stripe Customer ID from the tenant record
		// If they don't have one, we can't send them to portal (they need to checkout first)
		// We'll need to fetch the full tenant record if it's not fully populated in context

		// TODO: Fetch full tenant details including stripe_customer_id
		var customerID string
		err := store.pool.QueryRow(r.Context(), "SELECT stripe_customer_id FROM tenants WHERE id = $1", tc.Tenant.ID).Scan(&customerID)
		if err != nil || customerID == "" {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("no billing account found; please subscribe first"))
			return
		}

		domain := os.Getenv("DOMAIN_URL")
		if domain == "" {
			domain = "https://driftlock.net"
		}

		// Create Portal Session
		// Ideally set STRIPE_CUSTOMER_PORTAL_ID in env if using a specific configuration
		params := &stripe.BillingPortalSessionParams{
			Customer:  stripe.String(customerID),
			ReturnURL: stripe.String(domain + "/dashboard"),
		}

		ps, err := session.New(params)
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, fmt.Errorf("stripe portal session failed: %w", err))
			return
		}

		// Redirect or return URL
		http.Redirect(w, r, ps.URL, http.StatusSeeOther)
	}
}

func billingWebhookHandler(store *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const MaxBodyBytes = int64(65536)
		r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
		event, err := webhook.ConstructEvent(payload, r.Header.Get("Stripe-Signature"), endpointSecret)
		if err != nil {
			log.Printf("Error verifying webhook signature: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		switch event.Type {
		case "checkout.session.completed":
			var session stripe.CheckoutSession
			err := json.Unmarshal(event.Data.Raw, &session)
			if err != nil {
				log.Printf("Error parsing webhook JSON: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			handleCheckoutSessionCompleted(store, session)
		case "customer.subscription.updated", "customer.subscription.deleted":
			var subscription stripe.Subscription
			err := json.Unmarshal(event.Data.Raw, &subscription)
			if err != nil {
				log.Printf("Error parsing webhook JSON: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			handleSubscriptionUpdated(store, subscription)
		}

		w.WriteHeader(http.StatusOK)
	}
}

func handleCheckoutSessionCompleted(store *store, session stripe.CheckoutSession) {
	tenantIDStr := session.ClientReferenceID
	if tenantIDStr == "" {
		tenantIDStr = session.Metadata["tenant_id"]
	}
	if tenantIDStr == "" {
		log.Printf("No tenant_id in session %s", session.ID)
		return
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		log.Printf("Invalid tenant_id %s: %v", tenantIDStr, err)
		return
	}

	customerID := session.Customer.ID
	subscriptionID := session.Subscription.ID

	ctx := context.Background()
	if err := store.updateTenantStripeInfo(ctx, tenantID, customerID, subscriptionID, "active"); err != nil {
		log.Printf("Failed to update tenant %s with stripe info: %v", tenantID, err)
	} else {
		log.Printf("Updated tenant %s with subscription %s", tenantID, subscriptionID)
	}
}

func handleSubscriptionUpdated(store *store, subscription stripe.Subscription) {
	// We need to find the tenant by subscription ID or customer ID
	// Since we don't have a direct lookup by subscription ID in our cache (yet),
	// we might need to add a DB query or rely on metadata if we added it to the subscription.
	// For now, let's assume we can find it via customer ID if we stored it.

	status := string(subscription.Status)
	customerID := subscription.Customer.ID

	ctx := context.Background()
	if err := store.updateTenantStatusByStripeID(ctx, customerID, status); err != nil {
		log.Printf("Failed to update tenant status for customer %s: %v", customerID, err)
	}
}

// Store methods for billing

func (s *store) updateTenantStripeInfo(ctx context.Context, tenantID uuid.UUID, customerID, subscriptionID, status string) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE tenants 
		SET stripe_customer_id = $2, stripe_subscription_id = $3, stripe_status = $4, plan = 'pro', updated_at = NOW()
		WHERE id = $1`,
		tenantID, customerID, subscriptionID, status)
	if err != nil {
		return err
	}
	// Invalidate/Refresh cache
	return s.loadCache(ctx)
}

func (s *store) updateTenantStatusByStripeID(ctx context.Context, customerID, status string) error {
	// Also revert plan if not active?
	plan := "pro"
	if status != "active" && status != "trialing" {
		plan = "pilot" // Revert to free tier
	}

	_, err := s.pool.Exec(ctx, `
		UPDATE tenants 
		SET stripe_status = $2, plan = $3, updated_at = NOW()
		WHERE stripe_customer_id = $1`,
		customerID, status, plan)
	if err != nil {
		return err
	}
	return s.loadCache(ctx)
}

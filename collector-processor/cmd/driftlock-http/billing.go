package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/billingportal/session"
	checkoutsession "github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/subscription"
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
		// But return error for malformed JSON
		if r.Body != nil {
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err != io.EOF {
				writeJSON(w, r, http.StatusBadRequest, map[string]string{
					"error":   "invalid JSON body",
					"details": err.Error(),
				})
				return
			}
		}

		// Determine Price ID based on plan
		var priceID string
		switch req.Plan {
		case "orbit", "enterprise", "horizon": // Enterprise tier ($299/mo, 25M events)
			priceID = os.Getenv("STRIPE_PRICE_ID_ENTERPRISE")
		case "tensor", "sentinel", "lock", "transistor", "pro": // Pro tier ($100/mo)
			priceID = os.Getenv("STRIPE_PRICE_ID_PRO")
		case "radar", "signal", "basic": // Standard tier ($15/mo)
			priceID = os.Getenv("STRIPE_PRICE_ID_BASIC")
		default:
			// Default to radar (basic) as the entry paid tier
			priceID = os.Getenv("STRIPE_PRICE_ID_BASIC")
			if priceID == "" {
				// Fallback to Tensor (Pro) if Radar not set (migration path)
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
				"plan":      req.Plan,
			},
			SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
				TrialPeriodDays: stripe.Int64(14),
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

func billingWebhookHandler(store *store, emailer *emailService, webhookStore *WebhookEventStore) http.HandlerFunc {
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

		// Store event for durability before processing
		// This ensures we can retry if processing fails
		if webhookStore != nil {
			ctx := r.Context()
			eventID, isNew, err := webhookStore.StoreEvent(ctx, event.ID, string(event.Type), event.Data.Raw)
			if err != nil {
				log.Printf("Failed to store webhook event %s: %v", event.ID, err)
				// Fall through to synchronous processing as fallback
			} else if !isNew {
				// Duplicate event - already processed or in queue
				log.Printf("Webhook event %s already stored, skipping", event.ID)
				w.WriteHeader(http.StatusOK)
				return
			} else {
				// Event stored successfully - return 200 immediately
				// The retry worker will process it
				log.Printf("Stored webhook event %s (type: %s) for processing", eventID, event.Type)
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		// Fallback: synchronous processing if webhookStore is nil or storage failed
		// This maintains backward compatibility and handles edge cases
		switch event.Type {
		case "checkout.session.completed":
			var sess stripe.CheckoutSession
			err := json.Unmarshal(event.Data.Raw, &sess)
			if err != nil {
				log.Printf("Error parsing webhook JSON: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			handleCheckoutSessionCompleted(store, sess)

		case "customer.subscription.updated", "customer.subscription.deleted":
			var sub stripe.Subscription
			err := json.Unmarshal(event.Data.Raw, &sub)
			if err != nil {
				log.Printf("Error parsing webhook JSON: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			handleSubscriptionUpdated(store, sub)

		case "customer.subscription.trial_will_end":
			var sub stripe.Subscription
			err := json.Unmarshal(event.Data.Raw, &sub)
			if err != nil {
				log.Printf("Error parsing trial_will_end webhook JSON: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			handleTrialWillEnd(store, emailer, sub)

		case "invoice.payment_failed":
			var inv stripe.Invoice
			err := json.Unmarshal(event.Data.Raw, &inv)
			if err != nil {
				log.Printf("Error parsing payment_failed webhook JSON: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			handlePaymentFailed(store, emailer, inv)

		case "invoice.payment_succeeded":
			var inv stripe.Invoice
			err := json.Unmarshal(event.Data.Raw, &inv)
			if err != nil {
				log.Printf("Error parsing payment_succeeded webhook JSON: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			handlePaymentSucceeded(store, inv)
		}

		w.WriteHeader(http.StatusOK)
	}
}

func handleCheckoutSessionCompleted(store *store, sess stripe.CheckoutSession) {
	tenantIDStr := sess.ClientReferenceID
	if tenantIDStr == "" {
		tenantIDStr = sess.Metadata["tenant_id"]
	}
	if tenantIDStr == "" {
		log.Printf("No tenant_id in session %s", sess.ID)
		return
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		log.Printf("Invalid tenant_id %s: %v", tenantIDStr, err)
		return
	}

	customerID := sess.Customer.ID
	subscriptionID := sess.Subscription.ID

	// Retrieve plan from metadata, default to "tensor" (previously pro/transistor) if missing
	plan := sess.Metadata["plan"]
	if plan == "" {
		plan = "tensor"
	}

	// Fetch subscription to get trial_end (Stripe is source of truth)
	var trialEndsAt *time.Time
	status := "active"
	if subscriptionID != "" {
		sub, err := subscription.Get(subscriptionID, nil)
		if err == nil && sub.TrialEnd > 0 {
			t := time.Unix(sub.TrialEnd, 0)
			trialEndsAt = &t
			status = "trialing"
		}
	}

	ctx := context.Background()
	if err := store.updateTenantStripeInfo(ctx, tenantID, customerID, subscriptionID, status, plan, trialEndsAt); err != nil {
		log.Printf("Failed to update tenant %s with stripe info: %v", tenantID, err)
	} else {
		log.Printf("Updated tenant %s with subscription %s (plan: %s, trial_ends: %v)", tenantID, subscriptionID, plan, trialEndsAt)
	}
}

func handleSubscriptionUpdated(store *store, sub stripe.Subscription) {
	// We need to find the tenant by subscription ID or customer ID
	// Since we don't have a direct lookup by subscription ID in our cache (yet),
	// we might need to add a DB query or rely on metadata if we added it to the subscription.
	// For now, let's assume we can find it via customer ID if we stored it.

	status := string(sub.Status)
	customerID := sub.Customer.ID

	// Determine plan from subscription items
	plan := "signal" // Default fallthrough
	if len(sub.Items.Data) > 0 {
		priceID := sub.Items.Data[0].Price.ID
		if priceID == os.Getenv("STRIPE_PRICE_ID_ENTERPRISE") {
			plan = "orbit" // Enterprise tier ($299/mo, 25M events)
		} else if priceID == os.Getenv("STRIPE_PRICE_ID_PRO") {
			plan = "tensor"
		} else if priceID == os.Getenv("STRIPE_PRICE_ID_BASIC") {
			plan = "signal"
		}
	}

	// Extract trial end from Stripe (source of truth)
	var trialEnd *time.Time
	if sub.TrialEnd > 0 {
		t := time.Unix(sub.TrialEnd, 0)
		trialEnd = &t
	}

	ctx := context.Background()
	if err := store.updateTenantStatusByStripeID(ctx, customerID, status, plan, trialEnd); err != nil {
		log.Printf("Failed to update tenant status for customer %s: %v", customerID, err)
	}
}

func handleTrialWillEnd(store *store, emailer *emailService, sub stripe.Subscription) {
	customerID := sub.Customer.ID

	ctx := context.Background()
	tenant, err := store.getTenantByStripeCustomer(ctx, customerID)
	if err != nil {
		log.Printf("Failed to find tenant for customer %s: %v", customerID, err)
		return
	}

	// Send trial ending reminder email (Stripe sends this 3 days before trial ends)
	if emailer != nil && tenant.Email != "" {
		emailer.sendTrialEndingEmail(tenant.Email, tenant.Name, 3)
	}

	log.Printf("Sent trial ending reminder to tenant %s (%s)", tenant.ID, tenant.Email)
}

func handlePaymentFailed(store *store, emailer *emailService, inv stripe.Invoice) {
	customerID := inv.Customer.ID

	ctx := context.Background()
	tenant, err := store.getTenantByStripeCustomer(ctx, customerID)
	if err != nil {
		log.Printf("Failed to find tenant for customer %s: %v", customerID, err)
		return
	}

	// Set 7-day grace period
	gracePeriodEnd := time.Now().Add(7 * 24 * time.Hour)
	if err := store.setGracePeriod(ctx, tenant.ID, gracePeriodEnd); err != nil {
		log.Printf("Failed to set grace period for tenant %s: %v", tenant.ID, err)
		return
	}

	// Send payment failed email
	if emailer != nil && tenant.Email != "" {
		emailer.sendPaymentFailedEmail(tenant.Email, tenant.Name, gracePeriodEnd)
	}

	log.Printf("Set grace period for tenant %s until %s", tenant.ID, gracePeriodEnd.Format(time.RFC3339))
}

func handlePaymentSucceeded(store *store, inv stripe.Invoice) {
	customerID := inv.Customer.ID

	ctx := context.Background()
	tenant, err := store.getTenantByStripeCustomer(ctx, customerID)
	if err != nil {
		log.Printf("Failed to find tenant for customer %s: %v", customerID, err)
		return
	}

	// Clear grace period and reset failure count
	if err := store.clearGracePeriod(ctx, tenant.ID); err != nil {
		log.Printf("Failed to clear grace period for tenant %s: %v", tenant.ID, err)
		return
	}

	log.Printf("Payment succeeded for tenant %s, cleared grace period", tenant.ID)
}

// Store methods for billing

func (s *store) updateTenantStripeInfo(ctx context.Context, tenantID uuid.UUID, customerID, subscriptionID, status, plan string, trialEndsAt *time.Time) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE tenants
		SET stripe_customer_id = $2, stripe_subscription_id = $3, stripe_status = $4, plan = $5, trial_ends_at = $6, plan_started_at = NOW(), updated_at = NOW()
		WHERE id = $1`,
		tenantID, customerID, subscriptionID, status, plan, trialEndsAt)
	if err != nil {
		return err
	}
	// Invalidate/Refresh cache
	return s.loadCache(ctx)
}

func (s *store) updateTenantStatusByStripeID(ctx context.Context, customerID, status, plan string, trialEnd *time.Time) error {
	// Also revert plan if not active?
	targetPlan := plan
	if status != "active" && status != "trialing" {
		targetPlan = "pilot" // Revert to free tier
	}

	_, err := s.pool.Exec(ctx, `
		UPDATE tenants
		SET stripe_status = $2, plan = $3, trial_ends_at = $4, updated_at = NOW()
		WHERE stripe_customer_id = $1`,
		customerID, status, targetPlan, trialEnd)
	if err != nil {
		return err
	}
	return s.loadCache(ctx)
}

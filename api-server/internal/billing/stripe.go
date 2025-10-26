package billing

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/client"
	"github.com/your-org/driftlock/api-server/internal/models"
)

// StripeService handles all Stripe-related operations
type StripeService struct {
	svc *client.API
}

// NewStripeService creates a new Stripe service instance
func NewStripeService() (*StripeService, error) {
	// Set Stripe API key from environment
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	if stripeKey == "" {
		return nil, errors.New("STRIPE_SECRET_KEY environment variable is required")
	}

	stripe.Key = stripeKey

	svc := &StripeService{
		svc: &client.API{},
	}
	svc.svc.Init(stripeKey, nil)

	return svc, nil
}

// CreateCustomer creates a new customer in Stripe
func (s *StripeService) CreateCustomer(ctx context.Context, tenant *models.Tenant) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Name:  &tenant.Name,
		Email: &tenant.Email, // Note: Email field should be added to Tenant model
		Description: stripe.String(fmt.Sprintf("Driftlock tenant: %s", tenant.Name)),
		Metadata: map[string]string{
			"tenant_id": tenant.ID,
			"domain":    tenant.Domain,
		},
	}

	customer, err := s.svc.Customers.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer in Stripe: %w", err)
	}

	return customer, nil
}

// CreateSubscription creates a new subscription for a customer
func (s *StripeService) CreateSubscription(ctx context.Context, customerID, planID string) (*stripe.Subscription, error) {
	params := &stripe.SubscriptionParams{
		Customer: stripe.String(customerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(planID),
			},
		},
		// Set trial period for trial plan
		TrialPeriodDays: stripe.Int64(14), // Default 14-day trial
	}

	subscription, err := s.svc.Subscriptions.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	return subscription, nil
}

// CreateCheckoutSession creates a checkout session for payment
func (s *StripeService) CreateCheckoutSession(ctx context.Context, tenant *models.Tenant, planID string) (*stripe.CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(planID),
				Quantity: stripe.Int64(1),
			},
		},
		CustomerCreation: stripe.String(string(stripe.CheckoutSessionCustomerCreationAlways)),
		CustomerEmail:    &tenant.Email, // Use tenant's contact email
		SuccessURL:       stripe.String(fmt.Sprintf("https://%s/dashboard?session_id={CHECKOUT_SESSION_ID}", tenant.Domain)),
		CancelURL:        stripe.String(fmt.Sprintf("https://%s/billing", tenant.Domain)),
		Metadata: map[string]string{
			"tenant_id": tenant.ID,
			"plan":      planID,
		},
	}

	session, err := s.svc.CheckoutSessions.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create checkout session: %w", err)
	}

	return session, nil
}

// UpdateSubscription updates an existing subscription
func (s *StripeService) UpdateSubscription(ctx context.Context, subscriptionID, newPlanID string) (*stripe.Subscription, error) {
	params := &stripe.SubscriptionParams{
		Items: []*stripe.SubscriptionItemsParams{
			{
				ID:   stripe.String(subscriptionID),
				Plan: stripe.String(newPlanID),
			},
		},
	}

	subscription, err := s.svc.Subscriptions.Update(subscriptionID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}

	return subscription, nil
}

// CancelSubscription cancels a subscription
func (s *StripeService) CancelSubscription(ctx context.Context, subscriptionID string) error {
	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true), // Cancel at end of billing period
	}

	_, err := s.svc.Subscriptions.Update(subscriptionID, params)
	if err != nil {
		return fmt.Errorf("failed to cancel subscription: %w", err)
	}

	return nil
}

// GetSubscription retrieves subscription details
func (s *StripeService) GetSubscription(ctx context.Context, subscriptionID string) (*stripe.Subscription, error) {
	params := &stripe.SubscriptionParams{}
	subscription, err := s.svc.Subscriptions.Get(subscriptionID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return subscription, nil
}

// HandleWebhook processes Stripe webhooks
func (s *StripeService) HandleWebhook(ctx context.Context, payload []byte, signature string) error {
	// Get webhook signing secret from environment
	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if webhookSecret == "" {
		return errors.New("STRIPE_WEBHOOK_SECRET environment variable is required")
	}

	// Construct event from payload
	event, err := stripe.ConstructEvent(payload, signature, webhookSecret)
	if err != nil {
		return fmt.Errorf("failed to verify webhook signature: %w", err)
	}

	// Handle the event based on its type
	switch event.Type {
	case "customer.subscription.created":
		s.handleSubscriptionCreated(ctx, event.Data.Object)
	case "customer.subscription.updated":
		s.handleSubscriptionUpdated(ctx, event.Data.Object)
	case "customer.subscription.deleted":
		s.handleSubscriptionDeleted(ctx, event.Data.Object)
	case "checkout.session.completed":
		s.handleCheckoutSessionCompleted(ctx, event.Data.Object)
	default:
		slog.Info("Unhandled Stripe event type", "type", event.Type)
	}

	return nil
}

// handleSubscriptionCreated handles subscription creation event
func (s *StripeService) handleSubscriptionCreated(ctx context.Context, object interface{}) {
	// Extract subscription details and update tenant status in our database
	subscription, ok := object.(*stripe.Subscription)
	if !ok {
		slog.Error("Failed to parse subscription object from webhook")
		return
	}

	// Update tenant status based on subscription status
	tenantID := subscription.Metadata["tenant_id"]
	if tenantID != "" {
		// In a real implementation, you would update the tenant in your database
		slog.Info("Subscription created", "tenant_id", tenantID, "status", subscription.Status)
	}
}

// handleSubscriptionUpdated handles subscription update event
func (s *StripeService) handleSubscriptionUpdated(ctx context.Context, object interface{}) {
	subscription, ok := object.(*stripe.Subscription)
	if !ok {
		slog.Error("Failed to parse subscription object from webhook")
		return
	}

	tenantID := subscription.Metadata["tenant_id"]
	if tenantID != "" {
		// Update tenant based on new subscription details
		slog.Info("Subscription updated", "tenant_id", tenantID, "status", subscription.Status)
	}
}

// handleSubscriptionDeleted handles subscription cancellation event
func (s *StripeService) handleSubscriptionDeleted(ctx context.Context, object interface{}) {
	subscription, ok := object.(*stripe.Subscription)
	if !ok {
		slog.Error("Failed to parse subscription object from webhook")
		return
	}

	tenantID := subscription.Metadata["tenant_id"]
	if tenantID != "" {
		// Update tenant status (likely to suspended)
		slog.Info("Subscription canceled", "tenant_id", tenantID)
	}
}

// handleCheckoutSessionCompleted handles successful checkout completion
func (s *StripeService) handleCheckoutSessionCompleted(ctx context.Context, object interface{}) {
	session, ok := object.(*stripe.CheckoutSession)
	if !ok {
		slog.Error("Failed to parse checkout session object from webhook")
		return
	}

	tenantID := session.Metadata["tenant_id"]
	if tenantID != "" {
		// Mark checkout as completed and update tenant
		slog.Info("Checkout completed", "tenant_id", tenantID, "customer", session.Customer.ID)
	}
}

// CreateUsageRecord creates a usage record for metered billing
func (s *StripeService) CreateUsageRecord(ctx context.Context, subscriptionItemID string, quantity int64, timestamp time.Time) error {
	params := &stripe.UsageRecordParams{
		SubscriptionItem: stripe.String(subscriptionItemID),
		Quantity:         stripe.Int64(quantity),
		Timestamp:        stripe.Int64(timestamp.Unix()),
		Action:           stripe.String(string(stripe.UsageRecordActionIncrement)),
	}

	_, err := s.svc.UsageRecords.New(params)
	if err != nil {
		return fmt.Errorf("failed to create usage record: %w", err)
	}

	return nil
}
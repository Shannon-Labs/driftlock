package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Shannon-Labs/driftlock/collector-processor/cmd/driftlock-http/plans"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/subscription"
)

// WebhookRetryWorker processes failed webhook events
type WebhookRetryWorker struct {
	store        *store
	webhookStore *WebhookEventStore
	emailer      *emailService
	interval     time.Duration
	batchSize    int
	stopCh       chan struct{}
	wg           sync.WaitGroup
}

// NewWebhookRetryWorker creates a new retry worker
func NewWebhookRetryWorker(store *store, webhookStore *WebhookEventStore, emailer *emailService, interval time.Duration, batchSize int) *WebhookRetryWorker {
	return &WebhookRetryWorker{
		store:        store,
		webhookStore: webhookStore,
		emailer:      emailer,
		interval:     interval,
		batchSize:    batchSize,
		stopCh:       make(chan struct{}),
	}
}

// Start begins the retry worker loop
func (w *WebhookRetryWorker) Start(ctx context.Context) {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		log.Printf("Webhook retry worker started (interval: %s, batch: %d)", w.interval, w.batchSize)

		// Process immediately on start
		w.processRetryBatch(ctx)

		for {
			select {
			case <-ticker.C:
				w.processRetryBatch(ctx)
			case <-w.stopCh:
				log.Println("Webhook retry worker stopping...")
				return
			case <-ctx.Done():
				log.Println("Webhook retry worker context canceled")
				return
			}
		}
	}()
}

// Stop gracefully stops the retry worker
func (w *WebhookRetryWorker) Stop() {
	close(w.stopCh)
	w.wg.Wait()
}

// processRetryBatch fetches and processes a batch of events due for retry
func (w *WebhookRetryWorker) processRetryBatch(ctx context.Context) {
	events, err := w.webhookStore.GetEventsForRetry(ctx, w.batchSize)
	if err != nil {
		log.Printf("Failed to fetch events for retry: %v", err)
		return
	}

	if len(events) == 0 {
		return
	}

	log.Printf("Processing %d webhook events for retry", len(events))

	for _, event := range events {
		// Mark as processing first
		if err := w.webhookStore.MarkProcessing(ctx, event.ID); err != nil {
			log.Printf("Failed to mark event %s as processing: %v", event.ID, err)
			continue
		}

		// Process the event
		if err := w.processEvent(ctx, event); err != nil {
			log.Printf("Failed to process event %s (type: %s, retry: %d): %v",
				event.ID, event.EventType, event.RetryCount, err)
			if markErr := w.webhookStore.MarkFailed(ctx, event.ID, err.Error()); markErr != nil {
				log.Printf("Failed to mark event %s as failed: %v", event.ID, markErr)
			}
		} else {
			if err := w.webhookStore.MarkCompleted(ctx, event.ID); err != nil {
				log.Printf("Failed to mark event %s as completed: %v", event.ID, err)
			} else {
				log.Printf("Successfully processed webhook event %s (type: %s)", event.ID, event.EventType)
			}
		}
	}

	// Check for new dead letter events and alert
	w.checkDeadLetterAlerts(ctx)
}

// processEvent handles a single webhook event
func (w *WebhookRetryWorker) processEvent(ctx context.Context, event *WebhookEvent) error {
	switch event.EventType {
	case "checkout.session.completed":
		var sess stripe.CheckoutSession
		if err := json.Unmarshal(event.EventData, &sess); err != nil {
			return fmt.Errorf("failed to unmarshal checkout session: %w", err)
		}
		return w.handleCheckoutSessionCompletedWithError(ctx, sess)

	case "customer.subscription.updated", "customer.subscription.deleted":
		var sub stripe.Subscription
		if err := json.Unmarshal(event.EventData, &sub); err != nil {
			return fmt.Errorf("failed to unmarshal subscription: %w", err)
		}
		return w.handleSubscriptionUpdatedWithError(ctx, sub)

	case "customer.subscription.trial_will_end":
		var sub stripe.Subscription
		if err := json.Unmarshal(event.EventData, &sub); err != nil {
			return fmt.Errorf("failed to unmarshal subscription: %w", err)
		}
		return w.handleTrialWillEndWithError(ctx, sub)

	case "invoice.payment_failed":
		var inv stripe.Invoice
		if err := json.Unmarshal(event.EventData, &inv); err != nil {
			return fmt.Errorf("failed to unmarshal invoice: %w", err)
		}
		return w.handlePaymentFailedWithError(ctx, inv)

	case "invoice.payment_succeeded":
		var inv stripe.Invoice
		if err := json.Unmarshal(event.EventData, &inv); err != nil {
			return fmt.Errorf("failed to unmarshal invoice: %w", err)
		}
		return w.handlePaymentSucceededWithError(ctx, inv)

	default:
		// Unknown event types are considered successful (we don't need to retry)
		log.Printf("Unknown webhook event type: %s", event.EventType)
		return nil
	}
}

// Error-returning versions of webhook handlers

func (w *WebhookRetryWorker) handleCheckoutSessionCompletedWithError(ctx context.Context, sess stripe.CheckoutSession) error {
	tenantIDStr := sess.ClientReferenceID
	if tenantIDStr == "" {
		tenantIDStr = sess.Metadata["tenant_id"]
	}
	if tenantIDStr == "" {
		return fmt.Errorf("no tenant_id in session %s", sess.ID)
	}

	tenantID, err := parseUUID(tenantIDStr)
	if err != nil {
		return fmt.Errorf("invalid tenant_id %s: %w", tenantIDStr, err)
	}

	customerID := sess.Customer.ID
	subscriptionID := sess.Subscription.ID

	// Normalize plan to canonical name.
	// If the plan is missing or maps to Pulse (free tier), default to Tensor
	// since this is a paid checkout session - free users don't go through checkout.
	plan, _ := plans.NormalizePlan(sess.Metadata["plan"])
	if plan == plans.Pulse {
		plan = plans.Tensor
	}

	// Fetch subscription for trial info
	var trialEndsAt *time.Time
	status := "active"
	if subscriptionID != "" {
		sub, err := getStripeSubscription(subscriptionID)
		if err == nil && sub.TrialEnd > 0 {
			t := time.Unix(sub.TrialEnd, 0)
			trialEndsAt = &t
			status = "trialing"
		}
	}

	if err := w.store.updateTenantStripeInfo(ctx, tenantID, customerID, subscriptionID, status, plan, trialEndsAt); err != nil {
		return fmt.Errorf("failed to update tenant %s: %w", tenantID, err)
	}

	return nil
}

func (w *WebhookRetryWorker) handleSubscriptionUpdatedWithError(ctx context.Context, sub stripe.Subscription) error {
	customerID := sub.Customer.ID
	status := string(sub.Status)

	// Use canonical plan names
	plan := plans.Radar // Default to basic paid tier
	if len(sub.Items.Data) > 0 {
		priceID := sub.Items.Data[0].Price.ID
		if priceID == getEnvPriceIDPro() {
			plan = plans.Tensor
		} else if priceID == getEnvPriceIDBasic() {
			plan = plans.Radar
		}
	}

	var trialEnd *time.Time
	if sub.TrialEnd > 0 {
		t := time.Unix(sub.TrialEnd, 0)
		trialEnd = &t
	}

	if err := w.store.updateTenantStatusByStripeID(ctx, customerID, status, plan, trialEnd); err != nil {
		return fmt.Errorf("failed to update tenant for customer %s: %w", customerID, err)
	}

	return nil
}

func (w *WebhookRetryWorker) handleTrialWillEndWithError(ctx context.Context, sub stripe.Subscription) error {
	customerID := sub.Customer.ID

	tenant, err := w.store.getTenantByStripeCustomer(ctx, customerID)
	if err != nil {
		return fmt.Errorf("failed to find tenant for customer %s: %w", customerID, err)
	}

	if w.emailer != nil && tenant.Email != "" {
		w.emailer.sendTrialEndingEmail(tenant.Email, tenant.Name, 3)
	}

	return nil
}

func (w *WebhookRetryWorker) handlePaymentFailedWithError(ctx context.Context, inv stripe.Invoice) error {
	customerID := inv.Customer.ID

	tenant, err := w.store.getTenantByStripeCustomer(ctx, customerID)
	if err != nil {
		return fmt.Errorf("failed to find tenant for customer %s: %w", customerID, err)
	}

	gracePeriodEnd := time.Now().Add(7 * 24 * time.Hour)
	if err := w.store.setGracePeriod(ctx, tenant.ID, gracePeriodEnd); err != nil {
		return fmt.Errorf("failed to set grace period for tenant %s: %w", tenant.ID, err)
	}

	if w.emailer != nil && tenant.Email != "" {
		w.emailer.sendPaymentFailedEmail(tenant.Email, tenant.Name, gracePeriodEnd)
	}

	return nil
}

func (w *WebhookRetryWorker) handlePaymentSucceededWithError(ctx context.Context, inv stripe.Invoice) error {
	customerID := inv.Customer.ID

	tenant, err := w.store.getTenantByStripeCustomer(ctx, customerID)
	if err != nil {
		return fmt.Errorf("failed to find tenant for customer %s: %w", customerID, err)
	}

	if err := w.store.clearGracePeriod(ctx, tenant.ID); err != nil {
		return fmt.Errorf("failed to clear grace period for tenant %s: %w", tenant.ID, err)
	}

	return nil
}

// checkDeadLetterAlerts sends alerts for new dead letter events
func (w *WebhookRetryWorker) checkDeadLetterAlerts(ctx context.Context) {
	// Check for dead letter events in the last hour
	events, err := w.webhookStore.GetDeadLetterEvents(ctx, time.Now().Add(-1*time.Hour), 10)
	if err != nil {
		log.Printf("Failed to check dead letter events: %v", err)
		return
	}

	for _, event := range events {
		log.Printf("ALERT: Webhook event %s (type: %s) moved to dead letter after %d retries. Last error: %s",
			event.ID, event.EventType, event.RetryCount, stringOrEmpty(event.LastError))

		// Send alert email if configured
		alertEmail := getEnvAlertEmail()
		if w.emailer != nil && alertEmail != "" {
			w.emailer.sendAdminAlert(
				alertEmail,
				fmt.Sprintf("Dead Letter Webhook: %s", event.EventType),
				fmt.Sprintf("Webhook event failed after maximum retries.\n\nEvent ID: %s\nStripe Event: %s\nType: %s\nRetries: %d\nLast Error: %s\nCreated: %s",
					event.ID, event.StripeEventID, event.EventType, event.RetryCount,
					stringOrEmpty(event.LastError), event.CreatedAt.Format(time.RFC3339)),
			)
		}
	}
}

// Helper functions

func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

func getStripeSubscription(subscriptionID string) (*stripe.Subscription, error) {
	return subscription.Get(subscriptionID, nil)
}

func getEnvPriceIDPro() string {
	return os.Getenv("STRIPE_PRICE_ID_PRO")
}

func getEnvPriceIDBasic() string {
	return os.Getenv("STRIPE_PRICE_ID_BASIC")
}

func getEnvAlertEmail() string {
	return os.Getenv("ALERT_EMAIL")
}

func stringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

package main

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
)

// BillingCronWorker handles scheduled billing tasks like grace period downgrades
// and trial expiry notifications.
type BillingCronWorker struct {
	store   *store
	emailer *emailService
}

// NewBillingCronWorker creates a new billing cron worker.
func NewBillingCronWorker(store *store, emailer *emailService) *BillingCronWorker {
	return &BillingCronWorker{
		store:   store,
		emailer: emailer,
	}
}

// Start begins the daily billing cron jobs.
// Runs at startup and then every 24 hours.
func (w *BillingCronWorker) Start(ctx context.Context) {
	// Run immediately on startup
	w.runDailyTasks(ctx)

	// Then run every 24 hours
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Billing cron worker stopping")
			return
		case <-ticker.C:
			w.runDailyTasks(ctx)
		}
	}
}

func (w *BillingCronWorker) runDailyTasks(ctx context.Context) {
	log.Println("Running daily billing tasks...")

	// 1. Process expired grace periods
	w.processExpiredGracePeriods(ctx)

	// 2. Send trial ending reminders
	w.sendTrialEndingReminders(ctx)

	log.Println("Daily billing tasks completed")
}

// processExpiredGracePeriods finds tenants whose grace period has expired
// and downgrades them to the free tier.
func (w *BillingCronWorker) processExpiredGracePeriods(ctx context.Context) {
	tenants, err := w.store.getExpiredGracePeriodTenants(ctx)
	if err != nil {
		log.Printf("ERROR: Failed to query expired grace periods: %v", err)
		return
	}

	if len(tenants) == 0 {
		return
	}

	log.Printf("Processing %d expired grace period tenant(s)", len(tenants))

	for _, t := range tenants {
		if err := w.store.downgradeToFreeTier(ctx, t.ID); err != nil {
			log.Printf("ERROR: Failed to downgrade tenant %s: %v", t.ID, err)
			continue
		}

		log.Printf("Downgraded tenant %s (%s) from %s to free tier due to expired grace period",
			t.ID, t.Name, t.Plan)

		// Send notification email
		if w.emailer != nil && t.Email != "" {
			w.emailer.sendGraceExpiredEmail(t.Email, t.Name)
		}
	}
}

// sendTrialEndingReminders sends reminder emails to tenants whose trial ends in 3 days.
func (w *BillingCronWorker) sendTrialEndingReminders(ctx context.Context) {
	tenants, err := w.store.getTrialsEndingSoon(ctx, 3)
	if err != nil {
		log.Printf("ERROR: Failed to query ending trials: %v", err)
		return
	}

	if len(tenants) == 0 {
		return
	}

	log.Printf("Sending trial ending reminders to %d tenant(s)", len(tenants))

	for _, t := range tenants {
		if t.Email == "" {
			continue
		}

		if w.emailer != nil {
			w.emailer.sendTrialEndingEmail(t.Email, t.Name, 3)
		}

		// Mark reminder sent to prevent duplicates
		if err := w.store.markTrialReminderSent(ctx, t.ID); err != nil {
			log.Printf("ERROR: Failed to mark trial reminder sent for %s: %v", t.ID, err)
		}
	}
}

// TenantBillingInfo contains minimal tenant info for billing operations.
type TenantBillingInfo struct {
	ID    uuid.UUID
	Name  string
	Email string
	Plan  string
}

// getExpiredGracePeriodTenants returns tenants with expired grace periods.
func (s *store) getExpiredGracePeriodTenants(ctx context.Context) ([]TenantBillingInfo, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, name, COALESCE(email, '') as email, plan
		FROM tenants
		WHERE grace_period_ends_at < NOW()
		  AND grace_period_ends_at IS NOT NULL
		  AND plan != 'free'
		  AND plan != 'pulse'
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenants []TenantBillingInfo
	for rows.Next() {
		var t TenantBillingInfo
		if err := rows.Scan(&t.ID, &t.Name, &t.Email, &t.Plan); err != nil {
			return nil, err
		}
		tenants = append(tenants, t)
	}
	return tenants, rows.Err()
}

// downgradeToFreeTier downgrades a tenant to the free tier and clears billing fields.
func (s *store) downgradeToFreeTier(ctx context.Context, tenantID uuid.UUID) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE tenants
		SET plan = 'pulse',
		    grace_period_ends_at = NULL,
		    stripe_status = 'canceled',
		    updated_at = NOW()
		WHERE id = $1
	`, tenantID)
	return err
}

// getTrialsEndingSoon returns tenants whose trial ends within the specified days
// and haven't been sent a reminder yet.
func (s *store) getTrialsEndingSoon(ctx context.Context, daysOut int) ([]TenantBillingInfo, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, name, COALESCE(email, '') as email, plan
		FROM tenants
		WHERE trial_ends_at BETWEEN NOW() AND NOW() + $1 * INTERVAL '1 day'
		  AND COALESCE(trial_reminder_sent, false) = false
		  AND status = 'active'
	`, daysOut)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenants []TenantBillingInfo
	for rows.Next() {
		var t TenantBillingInfo
		if err := rows.Scan(&t.ID, &t.Name, &t.Email, &t.Plan); err != nil {
			return nil, err
		}
		tenants = append(tenants, t)
	}
	return tenants, rows.Err()
}

// markTrialReminderSent marks that a trial ending reminder was sent.
func (s *store) markTrialReminderSent(ctx context.Context, tenantID uuid.UUID) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE tenants
		SET trial_reminder_sent = true,
		    updated_at = NOW()
		WHERE id = $1
	`, tenantID)
	return err
}

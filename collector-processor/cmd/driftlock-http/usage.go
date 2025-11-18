package main

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
)

// Usage tracking for monitoring tenant activity and enforcing plan limits

type usageTracker struct {
	store *store
}

func newUsageTracker(store *store) *usageTracker {
	return &usageTracker{store: store}
}

// TrackUsage records usage metrics for a tenant
func (u *usageTracker) TrackUsage(ctx context.Context, tenantID, streamID uuid.UUID, eventCount, anomalyCount int) {
	go func() {
		trackCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := u.recordUsage(trackCtx, tenantID, streamID, eventCount, anomalyCount); err != nil {
			log.Printf("[usage] failed to track usage for tenant %s: %v", tenantID, err)
		}
	}()
}

func (u *usageTracker) recordUsage(ctx context.Context, tenantID, streamID uuid.UUID, eventCount, anomalyCount int) error {
	today := time.Now().UTC().Truncate(24 * time.Hour)

	// Upsert usage metrics
	_, err := u.store.pool.Exec(ctx, `
		INSERT INTO usage_metrics (tenant_id, stream_id, metric_date, event_count, anomaly_count, api_request_count)
		VALUES ($1, $2, $3, $4, $5, 1)
		ON CONFLICT (tenant_id, stream_id, metric_date)
		DO UPDATE SET
			event_count = usage_metrics.event_count + EXCLUDED.event_count,
			anomaly_count = usage_metrics.anomaly_count + EXCLUDED.anomaly_count,
			api_request_count = usage_metrics.api_request_count + 1,
			updated_at = NOW()`,
		tenantID, streamID, today, eventCount, anomalyCount)

	return err
}

// GetUsageSummary returns usage summary for a tenant
func (u *usageTracker) GetUsageSummary(ctx context.Context, tenantID uuid.UUID, days int) (*UsageSummary, error) {
	var summary UsageSummary

	// Get total usage for the period
	err := u.store.pool.QueryRow(ctx, `
		SELECT
			COALESCE(SUM(event_count), 0),
			COALESCE(SUM(anomaly_count), 0),
			COALESCE(SUM(api_request_count), 0)
		FROM usage_metrics
		WHERE tenant_id = $1 AND metric_date >= NOW() - INTERVAL '1 day' * $2`,
		tenantID, days).Scan(&summary.EventCount, &summary.AnomalyCount, &summary.APIRequestCount)

	if err != nil {
		return nil, err
	}

	summary.TenantID = tenantID.String()
	summary.Period = days

	return &summary, nil
}

// GetCurrentMonthUsage returns usage for the current billing month
func (u *usageTracker) GetCurrentMonthUsage(ctx context.Context, tenantID uuid.UUID) (*UsageSummary, error) {
	var summary UsageSummary

	// Get usage since the first of the month
	err := u.store.pool.QueryRow(ctx, `
		SELECT
			COALESCE(SUM(event_count), 0),
			COALESCE(SUM(anomaly_count), 0),
			COALESCE(SUM(api_request_count), 0)
		FROM usage_metrics
		WHERE tenant_id = $1 AND metric_date >= DATE_TRUNC('month', NOW())`,
		tenantID).Scan(&summary.EventCount, &summary.AnomalyCount, &summary.APIRequestCount)

	if err != nil {
		return nil, err
	}

	summary.TenantID = tenantID.String()
	summary.Period = -1 // indicates current month

	return &summary, nil
}

// CheckPlanLimits checks if tenant is approaching or exceeding plan limits
func (u *usageTracker) CheckPlanLimits(ctx context.Context, tenantID uuid.UUID) (*PlanLimitStatus, error) {
	// Get tenant's plan
	var plan string
	err := u.store.pool.QueryRow(ctx, `SELECT plan FROM tenants WHERE id = $1`, tenantID).Scan(&plan)
	if err != nil {
		return nil, err
	}

	// Get current month usage
	usage, err := u.GetCurrentMonthUsage(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// Define plan limits
	limits := getPlanLimits(plan)

	status := &PlanLimitStatus{
		TenantID:      tenantID.String(),
		Plan:          plan,
		EventCount:    usage.EventCount,
		EventLimit:    limits.EventLimit,
		UsagePercent:  float64(usage.EventCount) / float64(limits.EventLimit) * 100,
		AtWarning:     false,
		AtLimit:       false,
		AtOverage:     false,
	}

	// Check thresholds
	if status.UsagePercent >= 80 {
		status.AtWarning = true
	}
	if status.UsagePercent >= 100 {
		status.AtLimit = true
	}
	if status.UsagePercent >= 120 {
		status.AtOverage = true
	}

	return status, nil
}

// UsageSummary represents aggregated usage metrics
type UsageSummary struct {
	TenantID        string `json:"tenant_id"`
	EventCount      int64  `json:"event_count"`
	AnomalyCount    int64  `json:"anomaly_count"`
	APIRequestCount int64  `json:"api_request_count"`
	Period          int    `json:"period_days"`
}

// PlanLimitStatus represents the current plan limit status
type PlanLimitStatus struct {
	TenantID     string  `json:"tenant_id"`
	Plan         string  `json:"plan"`
	EventCount   int64   `json:"event_count"`
	EventLimit   int64   `json:"event_limit"`
	UsagePercent float64 `json:"usage_percent"`
	AtWarning    bool    `json:"at_warning"`
	AtLimit      bool    `json:"at_limit"`
	AtOverage    bool    `json:"at_overage"`
}

// PlanLimits defines limits for each plan
type PlanLimits struct {
	EventLimit int64
	StreamLimit int
	RetentionDays int
}

func getPlanLimits(plan string) PlanLimits {
	switch plan {
	case "trial":
		return PlanLimits{EventLimit: 10000, StreamLimit: 1, RetentionDays: 14}
	case "starter":
		return PlanLimits{EventLimit: 500000, StreamLimit: 5, RetentionDays: 30}
	case "growth":
		return PlanLimits{EventLimit: 5000000, StreamLimit: 20, RetentionDays: 90}
	case "enterprise":
		return PlanLimits{EventLimit: 100000000, StreamLimit: 100, RetentionDays: 365}
	default:
		return PlanLimits{EventLimit: 10000, StreamLimit: 1, RetentionDays: 14}
	}
}

// DailyUsageAggregator aggregates usage metrics daily (to be run as a cron job)
type DailyUsageAggregator struct {
	store        *store
	emailService *emailService
}

func newDailyUsageAggregator(store *store, emailService *emailService) *DailyUsageAggregator {
	return &DailyUsageAggregator{
		store:        store,
		emailService: emailService,
	}
}

// Run performs daily usage aggregation and sends warnings
func (a *DailyUsageAggregator) Run(ctx context.Context) error {
	log.Printf("[usage] Starting daily usage aggregation")

	// Get all tenants
	rows, err := a.store.pool.Query(ctx, `
		SELECT id, name, email, plan, created_at
		FROM tenants
		WHERE status = 'active'`)
	if err != nil {
		return err
	}
	defer rows.Close()

	tracker := newUsageTracker(a.store)
	now := time.Now()

	for rows.Next() {
		var (
			tenantID  uuid.UUID
			name      string
			email     *string
			plan      string
			createdAt time.Time
		)
		if err := rows.Scan(&tenantID, &name, &email, &plan, &createdAt); err != nil {
			log.Printf("[usage] error scanning tenant: %v", err)
			continue
		}

		// Check plan limits
		status, err := tracker.CheckPlanLimits(ctx, tenantID)
		if err != nil {
			log.Printf("[usage] error checking limits for %s: %v", tenantID, err)
			continue
		}

		// Send warnings if needed
		if email != nil && *email != "" {
			// Usage warning at 80%
			if status.AtWarning && !status.AtLimit {
				log.Printf("[usage] Tenant %s at %.1f%% usage", tenantID, status.UsagePercent)
			}

			// Check trial expiration
			if plan == "trial" {
				daysActive := int(now.Sub(createdAt).Hours() / 24)
				daysRemaining := 14 - daysActive
				if daysRemaining == 7 || daysRemaining == 3 || daysRemaining == 1 {
					if a.emailService.enabled {
						if err := a.emailService.SendTrialExpiring(ctx, *email, name, daysRemaining); err != nil {
							log.Printf("[usage] error sending trial warning to %s: %v", *email, err)
						}
					}
				}
			}
		}
	}

	log.Printf("[usage] Daily usage aggregation completed")
	return rows.Err()
}

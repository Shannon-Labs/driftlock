package main

import (
	"context"
	"log"
	"time"

	"github.com/Shannon-Labs/driftlock/collector-processor/cmd/driftlock-http/plans"
	"github.com/google/uuid"
)

type usageTracker struct {
	store   *store
	emailer *emailService
	// Cache usage to avoid checking DB on every request
	// For MVP, we'll check DB every N requests or just rely on async
	// Let's just do async for now
}

func newUsageTracker(store *store, emailer *emailService) *usageTracker {
	return &usageTracker{
		store:   store,
		emailer: emailer,
	}
}

func (ut *usageTracker) track(ctx context.Context, tenantID, streamID uuid.UUID, plan string, eventCount, anomalyCount int) {
	// Use a detached context for background work so it doesn't get cancelled with the request
	bgCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Increment usage
	if err := ut.store.incrementUsage(bgCtx, tenantID, streamID, eventCount, 1, anomalyCount); err != nil {
		log.Printf("ERROR: Failed to increment usage for tenant %s: %v", tenantID, err)
		return
	}

	// 2. Check limits (soft enforcement)
	// Normalize plan name and get limit from canonical source
	normalizedPlan, _ := plans.NormalizePlan(plan)
	limit := plans.GetLimit(normalizedPlan)

	totalUsage, err := ut.store.getMonthlyUsage(bgCtx, tenantID)
	if err != nil {
		log.Printf("ERROR: Failed to get usage for tenant %s: %v", tenantID, err)
		return
	}

	usagePercent := float64(totalUsage) / float64(limit)

	if usagePercent >= 1.0 {
		// Over limit
		log.Printf("WARNING: Tenant %s (Plan: %s) is OVER LIMIT (%d/%d events)", tenantID, normalizedPlan, totalUsage, limit)
		// TODO(P2): Rate limit or send email (once per day?) - requires state tracking
	} else if usagePercent >= 0.8 {
		// Near limit
		log.Printf("INFO: Tenant %s (Plan: %s) is near limit (%.1f%% - %d/%d events)", tenantID, normalizedPlan, usagePercent*100, totalUsage, limit)
	}
}

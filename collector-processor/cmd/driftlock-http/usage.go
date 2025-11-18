package main

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
)

// Plan limits (monthly events)
var planLimits = map[string]int64{
	"trial":      10_000,
	"starter":    500_000,
	"growth":     5_000_000,
	"enterprise": 1_000_000_000, // Custom
}

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
	// Don't check on every single request to reduce load?
	// For now, check every request but async.

	limit, ok := planLimits[plan]
	if !ok {
		limit = planLimits["trial"] // Default
	}

	totalUsage, err := ut.store.getMonthlyUsage(bgCtx, tenantID)
	if err != nil {
		log.Printf("ERROR: Failed to get usage for tenant %s: %v", tenantID, err)
		return
	}

	usagePercent := float64(totalUsage) / float64(limit)

	if usagePercent >= 1.0 {
		// Over limit
		log.Printf("WARNING: Tenant %s (Plan: %s) is OVER LIMIT (%d/%d events)", tenantID, plan, totalUsage, limit)
		// TODO: Rate limit or send email (once per day?)
		// Implementing "send email once" requires state. skipping for MVP.
	} else if usagePercent >= 0.8 {
		// Near limit
		log.Printf("INFO: Tenant %s (Plan: %s) is near limit (%.1f%% - %d/%d events)", tenantID, plan, usagePercent*100, totalUsage, limit)
	}
}

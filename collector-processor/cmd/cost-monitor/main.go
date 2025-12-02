package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	log.Println("Starting Driftlock Cost Monitor...")

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	ctx := context.Background()

	// 1. Check Global Hourly Spend
	checkGlobalHourlySpend(ctx, pool)

	// 2. Check for Margin Health
	checkMarginHealth(ctx, pool)
}

func checkGlobalHourlySpend(ctx context.Context, pool *pgxpool.Pool) {
	threshold := 50.0 // $50/hour warning

	var hourlySpend float64
	err := pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(cost_usd), 0)
		FROM ai_usage
		WHERE created_at > NOW() - INTERVAL '1 hour'
	`).Scan(&hourlySpend)

	if err != nil {
		log.Printf("Error checking hourly spend: %v", err)
		return
	}

	if hourlySpend > threshold {
		log.Printf("ALERT: High global hourly spend: $%.2f (Threshold: $%.2f)", hourlySpend, threshold)
		// In a real app, send Slack/Email alert here
	} else {
		log.Printf("Global hourly spend OK: $%.2f", hourlySpend)
	}
}

func checkMarginHealth(ctx context.Context, pool *pgxpool.Pool) {
	var marginPct float64
	err := pool.QueryRow(ctx, `
		SELECT
			CASE WHEN SUM(total_charge_usd) = 0 THEN 0
			ELSE (SUM(total_charge_usd - cost_usd) / SUM(total_charge_usd)) * 100
			END
		FROM ai_usage
		WHERE created_at > NOW() - INTERVAL '24 hours'
	`).Scan(&marginPct)

	if err != nil {
		log.Printf("Error checking margin health: %v", err)
		return
	}

	if marginPct < 10.0 {
		log.Printf("ALERT: Low margin detected: %.2f%% (Target: >15%%)", marginPct)
	} else {
		log.Printf("Margin health OK: %.2f%%", marginPct)
	}
}

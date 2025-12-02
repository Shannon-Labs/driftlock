package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shannon-Labs/driftlock/collector-processor/internal/ai"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	log.Println("Starting Driftlock Worker...")

	// Connect to database
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	// Initialize services
	batchService := ai.NewBatchService(pool)

	// Start cron scheduler
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go runBatchProcessor(ctx, batchService)

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down Driftlock Worker...")
}

func runBatchProcessor(ctx context.Context, batchService *ai.BatchService) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	log.Println("Batch processor started (interval: 5m)")

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			log.Println("Running batch processing job...")
			if err := processBatches(ctx, batchService); err != nil {
				log.Printf("Batch processing failed: %v", err)
			}
		}
	}
}

func processBatches(ctx context.Context, batchService *ai.BatchService) error {
	// 1. Fetch pending requests
	requests, err := batchService.GetPendingRequests(ctx, 100) // Process 100 at a time
	if err != nil {
		return err
	}

	if len(requests) == 0 {
		return nil
	}

	log.Printf("Found %d pending requests", len(requests))

	// 2. Group by model (simple implementation for now)
	// In a real implementation, we would use smart_router.OptimizeBatch logic here
	ids := make([]string, len(requests))
	for i, req := range requests {
		ids[i] = req.ID
	}

	// 3. Send to AI provider (Mock for now)
	// In a real implementation, this would call Anthropic Batch API
	batchID := "batch_" + time.Now().Format("20060102150405")
	log.Printf("Simulating batch submission: %s", batchID)

	// 4. Update status
	if err := batchService.MarkAsProcessing(ctx, ids, batchID); err != nil {
		return err
	}

	log.Printf("Successfully marked %d requests as processing", len(requests))
	return nil
}

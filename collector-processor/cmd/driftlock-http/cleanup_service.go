package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shannon-Labs/driftlock/collector-processor/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RunCleanupService starts the retention cleanup service
func RunCleanupService(pool *pgxpool.Pool) {
	logger := log.New(os.Stdout, "[CLEANUP] ", log.LstdFlags)

	// Create retention manager with default policy
	policy := storage.DefaultRetentionPolicy()
	policy.CleanupInterval = 24 * time.Hour // Run daily

	manager := storage.NewRetentionManager(pool, policy, logger)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Start the cleanup scheduler in a goroutine
	go manager.ScheduleCleanup(ctx)

	// Run initial cleanup
	logger.Println("Running initial cleanup on startup...")
	if err := manager.RunCleanup(ctx); err != nil {
		logger.Printf("Initial cleanup failed: %v", err)
	} else {
		logger.Println("Initial cleanup completed successfully")
	}

	// Show storage estimate
	if estimate, err := manager.GetStorageEstimate(ctx); err == nil {
		logger.Println("Current storage estimate by tier:")
		for tier, size := range estimate {
			logger.Printf("  %s: %s", tier, formatBytes(size))
		}
	}

	// Wait for shutdown signal
	<-sigCh
	logger.Println("Shutdown signal received, stopping cleanup service...")
	cancel()

	// Give some time for cleanup to finish
	time.Sleep(5 * time.Second)
	logger.Println("Cleanup service stopped")
}

// RunOnce runs cleanup once and exits
func RunOnce(pool *pgxpool.Pool, dryRun bool) {
	logger := log.New(os.Stdout, "[CLEANUP] ", log.LstdFlags)

	policy := storage.DefaultRetentionPolicy()
	if dryRun {
		logger.Println("DRY RUN MODE - No data will be deleted")
		policy.SoftDelete = true
	}

	manager := storage.NewRetentionManager(pool, policy, logger)

	ctx := context.Background()

	if err := manager.RunCleanup(ctx); err != nil {
		logger.Fatalf("Cleanup failed: %v", err)
	}

	logger.Println("Cleanup completed successfully")
}

// formatBytes formats bytes in human readable format
func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

// cleanupCommand handles the cleanup CLI command
func cleanupCommand(store *store) {
	var (
		dryRun = flag.Bool("dry-run", false, "Show what would be deleted without actually deleting")
		once   = flag.Bool("once", false, "Run cleanup once and exit")
	)
	flag.Parse()

	if *once {
		// Run once and exit
		RunOnce(store.pool, *dryRun)
		return
	}

	// Start the service
	RunCleanupService(store.pool)
}

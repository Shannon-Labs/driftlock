package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RetentionPolicy defines data retention settings
type RetentionPolicy struct {
	FreeTierDays     int           // Days for free tier (7)
	BasicTierDays    int           // Days for basic tier (30)
	ProTierDays      int           // Days for pro tier (90)
	ArchiveAfterDays int           // Days after which to archive to cold storage
	CleanupInterval  time.Duration // How often to run cleanup
	SoftDelete       bool          // Whether to soft delete or hard delete
}

// DefaultRetentionPolicy returns the default retention policy
func DefaultRetentionPolicy() *RetentionPolicy {
	return &RetentionPolicy{
		FreeTierDays:     7,   // Trial: 7 days
		BasicTierDays:    30,  // Radar: 30 days
		ProTierDays:      90,  // Tensor: 90 days
		ArchiveAfterDays: 365, // Orbit: Archive after 1 year
		CleanupInterval:  24 * time.Hour,
		SoftDelete:       true,
	}
}

// RetentionManager handles data retention and cleanup
type RetentionManager struct {
	pool   *pgxpool.Pool
	policy *RetentionPolicy
	logger *log.Logger
}

// NewRetentionManager creates a new retention manager
func NewRetentionManager(pool *pgxpool.Pool, policy *RetentionPolicy, logger *log.Logger) *RetentionManager {
	if policy == nil {
		policy = DefaultRetentionPolicy()
	}
	return &RetentionManager{
		pool:   pool,
		policy: policy,
		logger: logger,
	}
}

// RunCleanup executes the retention cleanup process
func (rm *RetentionManager) RunCleanup(ctx context.Context) error {
	rm.logger.Println("Starting retention cleanup process")

	// 1. Handle free tier data (7 days)
	if err := rm.cleanupTierData(ctx, "trial", rm.policy.FreeTierDays, true); err != nil {
		return fmt.Errorf("failed to cleanup free tier data: %w", err)
	}

	// 2. Handle basic tier data (30 days)
	if err := rm.cleanupTierData(ctx, "radar", rm.policy.BasicTierDays, false); err != nil {
		return fmt.Errorf("failed to cleanup basic tier data: %w", err)
	}

	// 3. Handle pro tier data (90 days)
	if err := rm.cleanupTierData(ctx, "tensor", rm.policy.ProTierDays, false); err != nil {
		return fmt.Errorf("failed to cleanup pro tier data: %w", err)
	}

	// 4. Archive old data from all tiers
	if err := rm.archiveOldData(ctx, rm.policy.ArchiveAfterDays); err != nil {
		return fmt.Errorf("failed to archive old data: %w", err)
	}

	// 5. Clean up orphaned records
	if err := rm.cleanupOrphanedRecords(ctx); err != nil {
		return fmt.Errorf("failed to cleanup orphaned records: %w", err)
	}

	// 6. Update statistics
	if err := rm.updateRetentionStats(ctx); err != nil {
		return fmt.Errorf("failed to update retention stats: %w", err)
	}

	rm.logger.Println("Retention cleanup process completed successfully")
	return nil
}

// cleanupTierData removes data for a specific tier
func (rm *RetentionManager) cleanupTierData(ctx context.Context, plan string, retentionDays int, hardDelete bool) error {
	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)

	rm.logger.Printf("Cleaning up %s tier data older than %s (hard delete: %v)",
		plan, cutoffDate.Format(time.RFC3339), hardDelete)

	// Get all tenants for this plan
	tenants, err := rm.getTenantsByPlan(ctx, plan)
	if err != nil {
		return err
	}

	totalDeleted := int64(0)

	for _, tenantID := range tenants {
		deleted, err := rm.deleteTenantData(ctx, tenantID, cutoffDate, hardDelete)
		if err != nil {
			rm.logger.Printf("Error cleaning up tenant %s: %v", tenantID, err)
			continue
		}
		totalDeleted += deleted
	}

	rm.logger.Printf("Deleted %d records for %s tier", totalDeleted, plan)
	return nil
}

// deleteTenantData removes data for a specific tenant
func (rm *RetentionManager) deleteTenantData(ctx context.Context, tenantID uuid.UUID, cutoffDate time.Time, hardDelete bool) (int64, error) {
	tx, err := rm.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var totalDeleted int64

	// Delete from streams (this will cascade to events and anomalies due to foreign keys)
	query := fmt.Sprintf(`
		DELETE FROM streams
		WHERE tenant_id = $1 AND created_at < $2
	`)

	if rm.policy.SoftDelete && !hardDelete {
		// Soft delete by marking as deleted
		query = fmt.Sprintf(`
			UPDATE streams
			SET deleted_at = NOW()
			WHERE tenant_id = $1 AND created_at < $2 AND deleted_at IS NULL
		`)
	}

	result, err := tx.Exec(ctx, query, tenantID, cutoffDate)
	if err != nil {
		return 0, err
	}

	totalDeleted = result.RowsAffected()

	// Also clean up usage metrics
	_, err = tx.Exec(ctx, `
		DELETE FROM usage_metrics
		WHERE tenant_id = $1 AND date < $2
	`, tenantID, cutoffDate.Format("2006-01-02"))
	if err != nil {
		return 0, err
	}

	// Clean up AI usage for free tier (always hard delete for cost)
	if hardDelete {
		_, err = tx.Exec(ctx, `
			DELETE FROM ai_usage
			WHERE tenant_id = $1 AND created_at < $2
		`, tenantID, cutoffDate)
		if err != nil {
			return 0, err
		}
	}

	// Clean up expired AI usage limits
	_, err = tx.Exec(ctx, `
		DELETE FROM ai_usage_limits
		WHERE tenant_id = $1 AND window_start < $2
	`, tenantID, time.Now().AddDate(0, 0, -30))
	if err != nil {
		return 0, err
	}

	return totalDeleted, tx.Commit(ctx)
}

// archiveOldData moves old data to cold storage
func (rm *RetentionManager) archiveOldData(ctx context.Context, archiveAfterDays int) error {
	cutoffDate := time.Now().AddDate(0, 0, -archiveAfterDays)

	rm.logger.Printf("Archiving data older than %s", cutoffDate.Format(time.RFC3339))

	// In a real implementation, this would move data to a cheaper storage solution
	// For now, we'll just mark it as archived
	_, err := rm.pool.Exec(ctx, `
		UPDATE streams
		SET archived = true, archived_at = NOW()
		WHERE created_at < $1 AND archived = false
	`, cutoffDate)

	if err != nil {
		return err
	}

	rm.logger.Println("Data archival completed")
	return nil
}

// cleanupOrphanedRecords removes records that reference deleted entities
func (rm *RetentionManager) cleanupOrphanedRecords(ctx context.Context) error {
	// Clean up usage_metrics with no matching stream
	_, err := rm.pool.Exec(ctx, `
		DELETE FROM usage_metrics
		WHERE stream_id NOT IN (SELECT id FROM streams WHERE deleted_at IS NULL)
	`)
	if err != nil {
		return err
	}

	// Clean up anomalies with no matching stream
	_, err = rm.pool.Exec(ctx, `
		DELETE FROM anomalies
		WHERE stream_id NOT IN (SELECT id FROM streams WHERE deleted_at IS NULL)
	`)
	if err != nil {
		return err
	}

	return nil
}

// updateRetentionStats updates statistics for monitoring
func (rm *RetentionManager) updateRetentionStats(ctx context.Context) error {
	stats := struct {
		TotalStreams     int64 `json:"total_streams"`
		DeletedStreams   int64 `json:"deleted_streams"`
		ArchivedStreams  int64 `json:"archived_streams"`
		TotalEvents      int64 `json:"total_events"`
		TotalAIUsage     int64 `json:"total_ai_usage"`
		StorageSizeBytes int64 `json:"storage_size_bytes"`
	}{}

	// Get storage statistics
	err := rm.pool.QueryRow(ctx, `
		SELECT
			COUNT(*)::bigint,
			COUNT(CASE WHEN deleted_at IS NOT NULL THEN 1 END)::bigint,
			COUNT(CASE WHEN archived = true THEN 1 END)::bigint
		FROM streams
	`).Scan(&stats.TotalStreams, &stats.DeletedStreams, &stats.ArchivedStreams)
	if err != nil {
		return err
	}

	// Log statistics for monitoring
	rm.logger.Printf("Retention stats: Total=%d, Deleted=%d, Archived=%d",
		stats.TotalStreams, stats.DeletedStreams, stats.ArchivedStreams)

	return nil
}

// getTenantsByPlan returns all tenant IDs for a given plan
func (rm *RetentionManager) getTenantsByPlan(ctx context.Context, plan string) ([]uuid.UUID, error) {
	rows, err := rm.pool.Query(ctx, `
		SELECT id FROM tenants WHERE plan = $1
	`, plan)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenants []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		tenants = append(tenants, id)
	}

	return tenants, rows.Err()
}

// ScheduleCleanup sets up a recurring cleanup job
func (rm *RetentionManager) ScheduleCleanup(ctx context.Context) {
	ticker := time.NewTicker(rm.policy.CleanupInterval)
	defer ticker.Stop()

	rm.logger.Printf("Scheduled cleanup job to run every %v", rm.policy.CleanupInterval)

	for {
		select {
		case <-ctx.Done():
			rm.logger.Println("Cleanup scheduler stopped")
			return
		case <-ticker.C:
			if err := rm.RunCleanup(ctx); err != nil {
				rm.logger.Printf("Cleanup job failed: %v", err)
			}
		}
	}
}

// GetStorageEstimate returns estimated storage usage by tier
func (rm *RetentionManager) GetStorageEstimate(ctx context.Context) (map[string]int64, error) {
	query := `
		SELECT
			t.plan,
			COUNT(s.id) as stream_count,
			COALESCE(SUM(s.event_count), 0) as event_count
		FROM tenants t
		LEFT JOIN streams s ON t.id = s.tenant_id AND s.deleted_at IS NULL
		GROUP BY t.plan
	`

	rows, err := rm.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	estimate := make(map[string]int64)

	// Average sizes (rough estimates)
	const avgEventSize = 1024       // 1KB per event
	const avgStreamSize = 10 * 1024 // 10KB per stream metadata

	for rows.Next() {
		var plan string
		var streamCount, eventCount int64

		if err := rows.Scan(&plan, &streamCount, &eventCount); err != nil {
			return nil, err
		}

		totalSize := streamCount*avgStreamSize + eventCount*avgEventSize
		estimate[plan] = totalSize
	}

	return estimate, nil
}

// CleanupTenant completely removes all data for a tenant (for GDPR compliance)
func (rm *RetentionManager) CleanupTenant(ctx context.Context, tenantID uuid.UUID) error {
	rm.logger.Printf("Starting complete cleanup for tenant %s", tenantID)

	tx, err := rm.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Delete in order of dependencies to respect foreign keys
	tables := []string{
		"ai_usage",
		"ai_usage_limits",
		"usage_metrics",
		"anomalies",
		"events",
		"streams",
		"api_keys",
	}

	for _, table := range tables {
		_, err := tx.Exec(ctx, fmt.Sprintf(`
			DELETE FROM %s WHERE tenant_id = $1
		`, table), tenantID)
		if err != nil {
			return fmt.Errorf("failed to delete from %s: %w", table, err)
		}
	}

	// Finally delete the tenant
	_, err = tx.Exec(ctx, `
		DELETE FROM tenants WHERE id = $1
	`, tenantID)
	if err != nil {
		return err
	}

	rm.logger.Printf("Successfully cleaned up all data for tenant %s", tenantID)
	return tx.Commit(ctx)
}

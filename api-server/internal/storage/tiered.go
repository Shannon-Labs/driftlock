package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/Shannon-Labs/driftlock/api-server/internal/models"
)

// TieredStorage provides hot/warm/cold storage tiering
type TieredStorage struct {
	hotStorage   StorageInterface      // PostgreSQL for hot data (last 7 days)
	warmStorage  StorageInterface      // ClickHouse for warm data (8-90 days) 
	coldStorage  StorageInterface      // S3/MinIO for cold data (90+ days)
	tierConfig   TierConfig
}

// TierConfig holds configuration for storage tiering
type TierConfig struct {
	HotRetentionDays  int           // Days to keep in hot storage
	WarmRetentionDays int           // Days to keep in warm storage (after hot)
	ArchiveInterval   time.Duration // How often to archive data
	Compression       CompressionConfig // Compression configuration for this tier
}

// CompressionConfig holds compression settings
type CompressionConfig struct {
	Enabled      bool
	Algorithm    string // openzl, zstd, lz4, gzip
	Level        int
	MinSizeBytes int    // Only compress if data exceeds this size
}

// NewTieredStorage creates a new tiered storage instance
func NewTieredStorage(hot, warm, cold StorageInterface, config TierConfig) *TieredStorage {
	return &TieredStorage{
		hotStorage:  hot,
		warmStorage: warm,
		coldStorage: cold,
		tierConfig:  config,
	}
}

// GetAnomaly retrieves an anomaly, checking all tiers in order
func (ts *TieredStorage) GetAnomaly(ctx context.Context, id uuid.UUID) (*models.Anomaly, error) {
	// Try hot storage first
	anomaly, err := ts.hotStorage.GetAnomaly(ctx, id)
	if err == nil && anomaly != nil {
		return anomaly, nil
	}

	// Try warm storage
	anomaly, err = ts.warmStorage.GetAnomaly(ctx, id)
	if err == nil && anomaly != nil {
		return anomaly, nil
	}

	// Try cold storage
	anomaly, err = ts.coldStorage.GetAnomaly(ctx, id)
	if err != nil {
		return nil, err
	}

	return anomaly, nil
}

// ListAnomalies retrieves anomalies across tiers
func (ts *TieredStorage) ListAnomalies(ctx context.Context, filter *models.AnomalyFilter) (*models.AnomalyListResponse, error) {
	// Get hot data
	hotResponse, err := ts.hotStorage.ListAnomalies(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("hot storage error: %w", err)
	}

	var allAnomalies []models.Anomaly
	allAnomalies = append(allAnomalies, hotResponse.Anomalies...)

	// Get warm data if needed
	if filter.StartTime == nil || filter.StartTime.Before(time.Now().AddDate(0, 0, -ts.tierConfig.HotRetentionDays)) {
		warmFilter := *filter // Copy filter
		// Restrict to warm storage date range
		if warmFilter.EndTime == nil {
			endTime := time.Now().AddDate(0, 0, -ts.tierConfig.HotRetentionDays)
			warmFilter.EndTime = &endTime
		}
		if warmFilter.StartTime == nil || warmFilter.StartTime.Before(time.Now().AddDate(0, 0, -ts.tierConfig.HotRetentionDays-ts.tierConfig.WarmRetentionDays)) {
			startTime := time.Now().AddDate(0, 0, -ts.tierConfig.HotRetentionDays-ts.tierConfig.WarmRetentionDays)
			warmFilter.StartTime = &startTime
		}
		
		warmResponse, err := ts.warmStorage.ListAnomalies(ctx, &warmFilter)
		if err != nil {
			return nil, fmt.Errorf("warm storage error: %w", err)
		}
		allAnomalies = append(allAnomalies, warmResponse.Anomalies...)
	}

	// Get cold data if needed
	if filter.StartTime == nil || filter.StartTime.Before(time.Now().AddDate(0, 0, -ts.tierConfig.HotRetentionDays-ts.tierConfig.WarmRetentionDays)) {
		coldFilter := *filter // Copy filter
		// Restrict to cold storage date range
		if coldFilter.StartTime == nil || coldFilter.StartTime.Before(time.Now().AddDate(0, 0, -ts.tierConfig.HotRetentionDays-ts.tierConfig.WarmRetentionDays)) {
			startTime := time.Now().AddDate(0, 0, -ts.tierConfig.HotRetentionDays-ts.tierConfig.WarmRetentionDays)
			coldFilter.StartTime = &startTime
		}
		
		coldResponse, err := ts.coldStorage.ListAnomalies(ctx, &coldFilter)
		if err != nil {
			return nil, fmt.Errorf("cold storage error: %w", err)
		}
		allAnomalies = append(allAnomalies, coldResponse.Anomalies...)
	}

	// Apply client-side filtering since data comes from multiple sources
	filtered := applyFilterModels(allAnomalies, filter)

	// Apply pagination
	start := filter.Offset
	end := start + filter.Limit
	if start > len(filtered) {
		start = len(filtered)
	}
	if end > len(filtered) {
		end = len(filtered)
	}

	paginated := filtered[start:end]

	return &models.AnomalyListResponse{
		Anomalies: paginated,
		Total:     len(filtered),
		Limit:     filter.Limit,
		Offset:    filter.Offset,
		HasMore:   end < len(filtered),
	}, nil
}

// CreateAnomaly stores anomaly in hot storage
func (ts *TieredStorage) CreateAnomaly(ctx context.Context, create *models.AnomalyCreate) (*models.Anomaly, error) {
	return ts.hotStorage.CreateAnomaly(ctx, create)
}

// ArchiveOldData moves old data from hot to warm storage
func (ts *TieredStorage) ArchiveOldData(ctx context.Context) error {
	// Define cutoff time for "old" data
	cutoffTime := time.Now().AddDate(0, 0, -ts.tierConfig.HotRetentionDays)

	// Find anomalies older than cutoff in hot storage
	filter := &models.AnomalyFilter{
		EndTime: &cutoffTime,
		Limit:   1000, // Process in batches
	}
	
	for {
		response, err := ts.hotStorage.ListAnomalies(ctx, filter)
		if err != nil {
			return fmt.Errorf("error listing old anomalies: %w", err)
		}
		
		if len(response.Anomalies) == 0 {
			break // No more old data to archive
		}
		
		// Move each anomaly to warm storage
		for _, anomaly := range response.Anomalies {
			// Create in warm storage
			anomalyCreate := &models.AnomalyCreate{
				StreamType:          anomaly.StreamType,
				Timestamp:           anomaly.Timestamp,
				NCDScore:            anomaly.NCDScore,
				PValue:              anomaly.PValue,
				GlassBoxExplanation: anomaly.GlassBoxExplanation,
				CompressionBaseline: anomaly.CompressionBaseline,
				CompressionWindow:   anomaly.CompressionWindow,
				CompressionCombined: anomaly.CompressionCombined,
				ConfidenceLevel:     anomaly.ConfidenceLevel,
				BaselineData:        anomaly.BaselineData,
				WindowData:          anomaly.WindowData,
				Metadata:            anomaly.Metadata,
				Tags:                anomaly.Tags,
			}
			
			_, err := ts.warmStorage.CreateAnomaly(ctx, anomalyCreate)
			if err != nil {
				return fmt.Errorf("error archiving anomaly %s: %w", anomaly.ID, err)
			}
			
			// Delete from hot storage after successful archive
			// Note: This would require a Delete method in the Storage interface that's not defined yet
			// For now, we'll just mark as archived or implement this later
		}
		
		// Continue with next batch (pagination would be needed)
		break // For now, just process first batch
	}
	
	return nil
}

// Helper function to apply filters to slice
func applyFilterModels(anomalies []models.Anomaly, filter *models.AnomalyFilter) []models.Anomaly {
	var filtered []models.Anomaly
	for _, anomaly := range anomalies {
		// Check all filter conditions
		if filter.StreamType != nil && *filter.StreamType != anomaly.StreamType {
			continue
		}
		if filter.Status != nil && *filter.Status != anomaly.Status {
			continue
		}
		if filter.MinNCDScore != nil && anomaly.NCDScore < *filter.MinNCDScore {
			continue
		}
		if filter.MaxPValue != nil && anomaly.PValue > *filter.MaxPValue {
			continue
		}
		if filter.StartTime != nil && anomaly.Timestamp.Before(*filter.StartTime) {
			continue
		}
		if filter.EndTime != nil && anomaly.Timestamp.After(*filter.EndTime) {
			continue
		}
		if filter.OnlySignificant && !isStatisticallySignificantModel(anomaly) {
			continue
		}
		
		filtered = append(filtered, anomaly)
	}
	return filtered
}

// Helper function to check if anomaly is statistically significant
func isStatisticallySignificantModel(anomaly models.Anomaly) bool {
	// Placeholder - in real implementation, use proper statistical significance check
	return anomaly.PValue < 0.05
}

// UpdateAnomalyStatus updates the status of an anomaly
// In a tiered storage system, this would search all tiers for the anomaly and update it in the appropriate tier
func (ts *TieredStorage) UpdateAnomalyStatus(ctx context.Context, id uuid.UUID, update *models.AnomalyUpdate, username string) error {
	// First, try to find the anomaly in any tier
	anomaly, err := ts.GetAnomaly(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find anomaly %s: %w", id, err)
	}
	if anomaly == nil {
		return fmt.Errorf("anomaly %s not found", id)
	}

	// Update should happen in the tier where the anomaly currently exists
	// For this simplified implementation, we'll update in the hot storage since we know
	// that's where recent items are
	return ts.hotStorage.UpdateAnomalyStatus(ctx, id, update, username)
}

// StartArchiveWorker starts a background worker for automatic archival
func (ts *TieredStorage) StartArchiveWorker(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(ts.tierConfig.ArchiveInterval)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				if err := ts.ArchiveOldData(ctx); err != nil {
					// Log error but continue
					fmt.Printf("Archive worker error: %v\n", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

// Close closes all underlying storage connections
func (ts *TieredStorage) Close() error {
	var firstErr error
	
	if err := ts.hotStorage.Close(); err != nil && firstErr == nil {
		firstErr = fmt.Errorf("hot storage close error: %w", err)
	}
	
	if err := ts.warmStorage.Close(); err != nil && firstErr == nil {
		firstErr = fmt.Errorf("warm storage close error: %w", err)
	}
	
	if err := ts.coldStorage.Close(); err != nil && firstErr == nil {
		firstErr = fmt.Errorf("cold storage close error: %w", err)
	}
	
	return firstErr
}

// Ping checks connectivity to all storage tiers
func (ts *TieredStorage) Ping(ctx context.Context) error {
	if err := ts.hotStorage.Ping(ctx); err != nil {
		return fmt.Errorf("hot storage ping failed: %w", err)
	}
	
	if err := ts.warmStorage.Ping(ctx); err != nil {
		return fmt.Errorf("warm storage ping failed: %w", err)
	}
	
	if err := ts.coldStorage.Ping(ctx); err != nil {
		return fmt.Errorf("cold storage ping failed: %w", err)
	}
	
	return nil
}

// Config methods - Configuration is not tiered, always use hot storage
func (ts *TieredStorage) GetActiveConfig(ctx context.Context) (*models.DetectionConfig, error) {
	return ts.hotStorage.GetActiveConfig(ctx)
}

func (ts *TieredStorage) UpdateConfig(ctx context.Context, update *models.DetectionConfigUpdate) (*models.DetectionConfig, error) {
	return ts.hotStorage.UpdateConfig(ctx, update)
}

// Metrics methods
func (ts *TieredStorage) RecordPerformanceMetric(ctx context.Context, metric *models.PerformanceMetric) error {
	return ts.hotStorage.RecordPerformanceMetric(ctx, metric)
}

// DB returns the hot storage database connection (for compatibility)
func (ts *TieredStorage) DB() *sql.DB {
	return ts.hotStorage.DB()
}
package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/Hmbown/driftlock/api-server/internal/compression"
	"github.com/Hmbown/driftlock/api-server/internal/models"
)

// Storage provides database operations for Driftlock
type Storage struct {
	db           *sql.DB
	compression  *compression.CompressionManager
}

// New creates a new Storage instance
func New(db *sql.DB) *Storage {
	return &Storage{
		db: db,
		compression: compression.NewCompressionManager(nil), // Use default compression strategy
	}
}

// NewPostgres creates a new PostgreSQL connection
func NewPostgres(connString string) (*Storage, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

    // Configure connection pool (production-friendly defaults)
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)
    db.SetConnMaxIdleTime(10 * time.Minute)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Storage{db: db}, nil
}

// Close closes the database connection
func (s *Storage) Close() error {
	return s.db.Close()
}

// Ping verifies database connectivity
func (s *Storage) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

// DB returns the underlying database connection
func (s *Storage) DB() *sql.DB {
	return s.db
}

// === Anomaly Operations ===

// CreateAnomaly inserts a new anomaly into the database
func (s *Storage) CreateAnomaly(ctx context.Context, create *models.AnomalyCreate) (*models.Anomaly, error) {
	baselineJSON, _ := json.Marshal(create.BaselineData)
	windowJSON, _ := json.Marshal(create.WindowData)
	metadataJSON, _ := json.Marshal(create.Metadata)

	compressionRatioChange := ((create.CompressionWindow - create.CompressionBaseline) / create.CompressionBaseline) * 100
	isSignificant := create.PValue < 0.05

	query := `
		INSERT INTO anomalies (
			timestamp, stream_type, ncd_score, p_value,
			glass_box_explanation, compression_baseline, compression_window,
			compression_combined, compression_ratio_change, confidence_level,
			is_statistically_significant, baseline_data, window_data, metadata, tags
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, created_at, updated_at, status
	`

	var anomaly models.Anomaly
	err := s.db.QueryRowContext(
		ctx, query,
		create.Timestamp, create.StreamType, create.NCDScore, create.PValue,
		create.GlassBoxExplanation, create.CompressionBaseline, create.CompressionWindow,
		create.CompressionCombined, compressionRatioChange, create.ConfidenceLevel,
		isSignificant, baselineJSON, windowJSON, metadataJSON, pq.Array(create.Tags),
	).Scan(&anomaly.ID, &anomaly.CreatedAt, &anomaly.UpdatedAt, &anomaly.Status)

	if err != nil {
		return nil, fmt.Errorf("failed to create anomaly: %w", err)
	}

	// Populate fields from create request
	anomaly.Timestamp = create.Timestamp
	anomaly.StreamType = create.StreamType
	anomaly.NCDScore = create.NCDScore
	anomaly.PValue = create.PValue
	anomaly.GlassBoxExplanation = create.GlassBoxExplanation
	anomaly.CompressionBaseline = create.CompressionBaseline
	anomaly.CompressionWindow = create.CompressionWindow
	anomaly.CompressionCombined = create.CompressionCombined
	anomaly.CompressionRatioChange = compressionRatioChange
	anomaly.ConfidenceLevel = create.ConfidenceLevel
	anomaly.IsStatisticallySignificant = isSignificant
	anomaly.BaselineData = create.BaselineData
	anomaly.WindowData = create.WindowData
	anomaly.Metadata = create.Metadata
	anomaly.Tags = create.Tags

	return &anomaly, nil
}

// GetAnomaly retrieves a single anomaly by ID
func (s *Storage) GetAnomaly(ctx context.Context, id uuid.UUID) (*models.Anomaly, error) {
	query := `
		SELECT
			id, timestamp, stream_type, ncd_score, p_value, status,
			glass_box_explanation, detailed_explanation,
			compression_baseline, compression_window, compression_combined, compression_ratio_change,
			baseline_entropy, window_entropy, entropy_change,
			confidence_level, is_statistically_significant,
			baseline_data, window_data, metadata, tags,
			acknowledged_by, acknowledged_at, dismissed_by, dismissed_at, notes,
			created_at, updated_at
		FROM anomalies
		WHERE id = $1
	`

	var anomaly models.Anomaly
	var baselineJSON, windowJSON, metadataJSON []byte

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&anomaly.ID, &anomaly.Timestamp, &anomaly.StreamType, &anomaly.NCDScore, &anomaly.PValue, &anomaly.Status,
		&anomaly.GlassBoxExplanation, &anomaly.DetailedExplanation,
		&anomaly.CompressionBaseline, &anomaly.CompressionWindow, &anomaly.CompressionCombined, &anomaly.CompressionRatioChange,
		&anomaly.BaselineEntropy, &anomaly.WindowEntropy, &anomaly.EntropyChange,
		&anomaly.ConfidenceLevel, &anomaly.IsStatisticallySignificant,
		&baselineJSON, &windowJSON, &metadataJSON, pq.Array(&anomaly.Tags),
		&anomaly.AcknowledgedBy, &anomaly.AcknowledgedAt, &anomaly.DismissedBy, &anomaly.DismissedAt, &anomaly.Notes,
		&anomaly.CreatedAt, &anomaly.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("anomaly not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get anomaly: %w", err)
	}

	// Unmarshal JSON fields
	if len(baselineJSON) > 0 {
		json.Unmarshal(baselineJSON, &anomaly.BaselineData)
	}
	if len(windowJSON) > 0 {
		json.Unmarshal(windowJSON, &anomaly.WindowData)
	}
	if len(metadataJSON) > 0 {
		json.Unmarshal(metadataJSON, &anomaly.Metadata)
	}

	return &anomaly, nil
}

// ListAnomalies retrieves anomalies with filtering and pagination
func (s *Storage) ListAnomalies(ctx context.Context, filter *models.AnomalyFilter) (*models.AnomalyListResponse, error) {
	// Build WHERE clause
	conditions := []string{}
	args := []interface{}{}
	argIndex := 1

	if filter.StreamType != nil {
		conditions = append(conditions, fmt.Sprintf("stream_type = $%d", argIndex))
		args = append(args, *filter.StreamType)
		argIndex++
	}

	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *filter.Status)
		argIndex++
	}

	if filter.MinNCDScore != nil {
		conditions = append(conditions, fmt.Sprintf("ncd_score >= $%d", argIndex))
		args = append(args, *filter.MinNCDScore)
		argIndex++
	}

	if filter.MaxPValue != nil {
		conditions = append(conditions, fmt.Sprintf("p_value <= $%d", argIndex))
		args = append(args, *filter.MaxPValue)
		argIndex++
	}

	if filter.StartTime != nil {
		conditions = append(conditions, fmt.Sprintf("timestamp >= $%d", argIndex))
		args = append(args, *filter.StartTime)
		argIndex++
	}

	if filter.EndTime != nil {
		conditions = append(conditions, fmt.Sprintf("timestamp <= $%d", argIndex))
		args = append(args, *filter.EndTime)
		argIndex++
	}

	if filter.OnlySignificant {
		conditions = append(conditions, "is_statistically_significant = true")
	}

	if len(filter.Tags) > 0 {
		conditions = append(conditions, fmt.Sprintf("tags && $%d", argIndex))
		args = append(args, pq.Array(filter.Tags))
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM anomalies %s", whereClause)
	var total int
	err := s.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count anomalies: %w", err)
	}

	// Set defaults for pagination
	limit := filter.Limit
	if limit == 0 {
		limit = 50
	}
	offset := filter.Offset

	// Query anomalies
	query := fmt.Sprintf(`
		SELECT
			id, timestamp, stream_type, ncd_score, p_value, status,
			glass_box_explanation, compression_baseline, compression_window, compression_combined,
			compression_ratio_change, confidence_level, is_statistically_significant,
			tags, created_at, updated_at
		FROM anomalies
		%s
		ORDER BY timestamp DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list anomalies: %w", err)
	}
	defer rows.Close()

	anomalies := []models.Anomaly{}
	for rows.Next() {
		var a models.Anomaly
		err := rows.Scan(
			&a.ID, &a.Timestamp, &a.StreamType, &a.NCDScore, &a.PValue, &a.Status,
			&a.GlassBoxExplanation, &a.CompressionBaseline, &a.CompressionWindow, &a.CompressionCombined,
			&a.CompressionRatioChange, &a.ConfidenceLevel, &a.IsStatisticallySignificant,
			pq.Array(&a.Tags), &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan anomaly: %w", err)
		}
		anomalies = append(anomalies, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating anomalies: %w", err)
	}

	return &models.AnomalyListResponse{
		Anomalies: anomalies,
		Total:     total,
		Limit:     limit,
		Offset:    offset,
		HasMore:   offset+limit < total,
	}, nil
}

// UpdateAnomalyStatus updates the status of an anomaly
func (s *Storage) UpdateAnomalyStatus(ctx context.Context, id uuid.UUID, update *models.AnomalyUpdate, username string) error {
	now := time.Now()

	query := `
		UPDATE anomalies
		SET status = $1,
		    notes = COALESCE($2, notes),
		    acknowledged_by = CASE WHEN $1 = 'acknowledged' THEN $3 ELSE acknowledged_by END,
		    acknowledged_at = CASE WHEN $1 = 'acknowledged' THEN $4 ELSE acknowledged_at END,
		    dismissed_by = CASE WHEN $1 = 'dismissed' THEN $3 ELSE dismissed_by END,
		    dismissed_at = CASE WHEN $1 = 'dismissed' THEN $4 ELSE dismissed_at END
		WHERE id = $5
	`

	result, err := s.db.ExecContext(ctx, query, update.Status, update.Notes, username, now, id)
	if err != nil {
		return fmt.Errorf("failed to update anomaly: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("anomaly not found: %s", id)
	}

	return nil
}

// === Configuration Operations ===

// GetActiveConfig retrieves the active detection configuration
func (s *Storage) GetActiveConfig(ctx context.Context) (*models.DetectionConfig, error) {
	query := `
		SELECT id, ncd_threshold, p_value_threshold, baseline_size, window_size, hop_size,
		       stream_overrides, created_at, updated_at, created_by, notes, is_active
		FROM detection_config
		WHERE is_active = true
		LIMIT 1
	`

	var config models.DetectionConfig
	var overridesJSON []byte

	err := s.db.QueryRowContext(ctx, query).Scan(
		&config.ID, &config.NCDThreshold, &config.PValueThreshold,
		&config.BaselineSize, &config.WindowSize, &config.HopSize,
		&overridesJSON, &config.CreatedAt, &config.UpdatedAt,
		&config.CreatedBy, &config.Notes, &config.IsActive,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no active configuration found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	if len(overridesJSON) > 0 {
		json.Unmarshal(overridesJSON, &config.StreamOverrides)
	}

	return &config, nil
}

// UpdateConfig updates the active detection configuration
func (s *Storage) UpdateConfig(ctx context.Context, update *models.DetectionConfigUpdate) (*models.DetectionConfig, error) {
	// Build update query dynamically
	updates := []string{}
	args := []interface{}{}
	argIndex := 1

	if update.NCDThreshold != nil {
		updates = append(updates, fmt.Sprintf("ncd_threshold = $%d", argIndex))
		args = append(args, *update.NCDThreshold)
		argIndex++
	}

	if update.PValueThreshold != nil {
		updates = append(updates, fmt.Sprintf("p_value_threshold = $%d", argIndex))
		args = append(args, *update.PValueThreshold)
		argIndex++
	}

	if update.BaselineSize != nil {
		updates = append(updates, fmt.Sprintf("baseline_size = $%d", argIndex))
		args = append(args, *update.BaselineSize)
		argIndex++
	}

	if update.WindowSize != nil {
		updates = append(updates, fmt.Sprintf("window_size = $%d", argIndex))
		args = append(args, *update.WindowSize)
		argIndex++
	}

	if update.HopSize != nil {
		updates = append(updates, fmt.Sprintf("hop_size = $%d", argIndex))
		args = append(args, *update.HopSize)
		argIndex++
	}

	if update.StreamOverrides != nil {
		overridesJSON, _ := json.Marshal(update.StreamOverrides)
		updates = append(updates, fmt.Sprintf("stream_overrides = $%d", argIndex))
		args = append(args, overridesJSON)
		argIndex++
	}

	if update.Notes != nil {
		updates = append(updates, fmt.Sprintf("notes = $%d", argIndex))
		args = append(args, *update.Notes)
		argIndex++
	}

	if len(updates) == 0 {
		return s.GetActiveConfig(ctx)
	}

	query := fmt.Sprintf(`
		UPDATE detection_config
		SET %s
		WHERE is_active = true
		RETURNING id
	`, strings.Join(updates, ", "))

	var id int
	err := s.db.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to update config: %w", err)
	}

	return s.GetActiveConfig(ctx)
}

// RecordPerformanceMetric logs a performance metric
func (s *Storage) RecordPerformanceMetric(ctx context.Context, metric *models.PerformanceMetric) error {
	metadataJSON, _ := json.Marshal(metric.Metadata)

	query := `
		INSERT INTO performance_metrics (metric_type, endpoint, duration_ms, success, error_message, metadata)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := s.db.ExecContext(
		ctx, query,
		metric.MetricType, metric.Endpoint, metric.DurationMs,
		metric.Success, metric.ErrorMessage, metadataJSON,
	)

	if err != nil {
		return fmt.Errorf("failed to record metric: %w", err)
	}

	return nil
}

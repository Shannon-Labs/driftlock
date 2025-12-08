package main

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// AutoTuneConfig holds configuration for the auto-tuning algorithm
type AutoTuneConfig struct {
	MinFeedbackSamples   int           // Minimum feedback samples before tuning
	FalsePositiveTarget  float64       // Target false positive rate (e.g., 0.05 = 5%)
	LearningRate         float64       // How aggressively to adjust (0.1 = 10% of delta)
	MaxAdjustmentPercent float64       // Max single adjustment (e.g., 0.20 = 20%)
	CooldownPeriod       time.Duration // Time between auto-tune runs
}

var defaultAutoTuneConfig = AutoTuneConfig{
	MinFeedbackSamples:   20,
	FalsePositiveTarget:  0.05,
	LearningRate:         0.15,
	MaxAdjustmentPercent: 0.25,
	CooldownPeriod:       1 * time.Hour,
}

// shouldThrottleAutoTune enforces a cooldown between auto-tune evaluations.
func shouldThrottleAutoTune(lastTune *time.Time, cfg AutoTuneConfig, now time.Time) bool {
	if lastTune == nil {
		return false
	}
	return now.Sub(*lastTune) < cfg.CooldownPeriod
}

// FeedbackStats summarizes user feedback for a stream
type FeedbackStats struct {
	TotalFeedback     int
	FalsePositives    int
	Confirmed         int
	Dismissed         int
	FalsePositiveRate float64
	AvgFPNCD          float64 // Average NCD of false positives
	AvgConfirmedNCD   float64 // Average NCD of confirmed anomalies
}

// computeAutoTuneAdjustment calculates threshold adjustments based on feedback
func computeAutoTuneAdjustment(stats FeedbackStats, currentNCD, currentPValue float64, cfg AutoTuneConfig) (newNCD, newPValue float64, reason string, shouldAdjust bool) {
	// Default: no change
	newNCD = currentNCD
	newPValue = currentPValue
	shouldAdjust = false

	if stats.TotalFeedback < cfg.MinFeedbackSamples {
		return newNCD, newPValue, "insufficient_feedback", false
	}

	// Calculate adjustments based on false positive rate
	fpDelta := stats.FalsePositiveRate - cfg.FalsePositiveTarget

	if fpDelta > 0.1 { // FP rate too high (>10% above target)
		// Increase thresholds to reduce sensitivity
		// NCD: higher threshold = fewer detections
		ncdAdjust := math.Min(currentNCD*cfg.MaxAdjustmentPercent, fpDelta*cfg.LearningRate)
		newNCD = currentNCD + ncdAdjust

		// P-value: lower threshold = more strict significance
		pvalueAdjust := math.Min(currentPValue*cfg.MaxAdjustmentPercent, fpDelta*cfg.LearningRate*0.5)
		newPValue = math.Max(0.001, currentPValue-pvalueAdjust)

		reason = "high_false_positive_rate"
		shouldAdjust = true

	} else if stats.FalsePositiveRate < cfg.FalsePositiveTarget*0.5 && stats.Confirmed > 0 { // FP rate very low (< half target), might be missing anomalies
		// Slightly decrease thresholds to increase sensitivity
		// Use a minimum adjustment of 5% of current threshold to ensure measurable change
		ncdAdjust := currentNCD * cfg.LearningRate * 0.5 // ~7.5% reduction with default 0.15 learning rate
		newNCD = math.Max(0.1, currentNCD-ncdAdjust)

		reason = "low_detection_rate"
		shouldAdjust = true
	}

	// Use feedback-specific NCD values to refine threshold
	if stats.FalsePositives > 5 && stats.AvgFPNCD > 0 {
		// If false positives have lower NCD than current threshold, raise it
		if stats.AvgFPNCD > currentNCD*0.8 {
			// FPs are close to or above threshold - threshold might be too low
			suggestedNCD := stats.AvgFPNCD * 1.1 // Set above FP average
			if suggestedNCD > newNCD {
				newNCD = math.Min(newNCD*1.2, suggestedNCD)
				shouldAdjust = true
				reason = "fp_ncd_boundary"
			}
		}
	}

	// Clamp to reasonable bounds
	newNCD = math.Max(0.1, math.Min(0.8, newNCD))
	newPValue = math.Max(0.001, math.Min(0.2, newPValue))

	// Only adjust if change is significant
	ncdChanged := math.Abs(newNCD-currentNCD) > 0.01
	pvalueChanged := math.Abs(newPValue-currentPValue) > 0.005
	if !ncdChanged && !pvalueChanged {
		shouldAdjust = false
	}

	return newNCD, newPValue, reason, shouldAdjust
}

// FeedbackRecord represents a single feedback entry for database operations
type FeedbackRecord struct {
	ID                    uuid.UUID
	AnomalyID             uuid.UUID
	StreamID              uuid.UUID
	TenantID              uuid.UUID
	FeedbackType          string
	NCDAtDetection        float64
	PValueAtDetection     float64
	ConfidenceAtDetection float64
	FeedbackReason        *string
	CreatedAt             time.Time
	CreatedBy             string
}

// TuneHistoryRecord represents a threshold adjustment for audit trail
type TuneHistoryRecord struct {
	ID         uuid.UUID
	StreamID   uuid.UUID
	TuneType   string // 'ncd', 'pvalue', 'baseline', 'window'
	OldValue   *float64
	NewValue   float64
	Reason     string
	Confidence float64
	CreatedAt  time.Time
}

// recordFeedback stores user feedback on an anomaly
func (s *store) recordFeedback(ctx context.Context, feedback FeedbackRecord) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO anomaly_feedback
		(anomaly_id, stream_id, tenant_id, feedback_type,
		 ncd_at_detection, pvalue_at_detection, confidence_at_detection,
		 feedback_reason, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		feedback.AnomalyID, feedback.StreamID, feedback.TenantID, feedback.FeedbackType,
		feedback.NCDAtDetection, feedback.PValueAtDetection, feedback.ConfidenceAtDetection,
		feedback.FeedbackReason, feedback.CreatedBy)
	return err
}

// getFeedbackStats retrieves feedback statistics for a stream
func (s *store) getFeedbackStats(ctx context.Context, streamID uuid.UUID, since time.Time) (*FeedbackStats, error) {
	stats := &FeedbackStats{}

	err := s.pool.QueryRow(ctx, `
		SELECT
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE feedback_type = 'false_positive') as false_positives,
			COUNT(*) FILTER (WHERE feedback_type = 'confirmed') as confirmed,
			COUNT(*) FILTER (WHERE feedback_type = 'dismissed') as dismissed,
			COALESCE(AVG(ncd_at_detection) FILTER (WHERE feedback_type = 'false_positive'), 0) as avg_fp_ncd,
			COALESCE(AVG(ncd_at_detection) FILTER (WHERE feedback_type = 'confirmed'), 0) as avg_confirmed_ncd
		FROM anomaly_feedback
		WHERE stream_id = $1 AND created_at >= $2`,
		streamID, since).Scan(
		&stats.TotalFeedback, &stats.FalsePositives, &stats.Confirmed, &stats.Dismissed,
		&stats.AvgFPNCD, &stats.AvgConfirmedNCD)

	if err != nil {
		return nil, err
	}

	if stats.TotalFeedback > 0 {
		stats.FalsePositiveRate = float64(stats.FalsePositives) / float64(stats.TotalFeedback)
	}

	return stats, nil
}

// recordTuneHistory records a threshold adjustment for audit
func (s *store) recordTuneHistory(ctx context.Context, history TuneHistoryRecord) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO threshold_tune_history (stream_id, tune_type, old_value, new_value, reason, confidence)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		history.StreamID, history.TuneType, history.OldValue, history.NewValue, history.Reason, history.Confidence)
	return err
}

// StreamTuningSettings holds tuning-related settings for a stream
type StreamTuningSettings struct {
	DetectionProfile      string
	AutoTuneEnabled       bool
	AdaptiveWindowEnabled bool
	NCDThreshold          float64
	PValueThreshold       float64
	TunedNCDThreshold     *float64
	TunedPValueThreshold  *float64
	BaselineSize          int
	WindowSize            int
	LastTuneCheck         *time.Time
}

// getStreamTuningSettings retrieves tuning settings for a stream
func (s *store) getStreamTuningSettings(ctx context.Context, streamID uuid.UUID) (*StreamTuningSettings, error) {
	settings := &StreamTuningSettings{}

	// Get stream settings with defaults applied
	err := s.pool.QueryRow(ctx, `
		SELECT
			COALESCE(detection_profile, 'balanced'),
			COALESCE(auto_tune_enabled, false),
			COALESCE(adaptive_window_enabled, false),
			COALESCE(tuned_ncd_threshold, NULL),
			COALESCE(tuned_pvalue_threshold, NULL)
		FROM streams
		WHERE id = $1`,
		streamID).Scan(
		&settings.DetectionProfile, &settings.AutoTuneEnabled, &settings.AdaptiveWindowEnabled,
		&settings.TunedNCDThreshold, &settings.TunedPValueThreshold)

	if err != nil {
		return nil, err
	}

	// Get profile defaults
	profile := DetectionProfile(settings.DetectionProfile)
	defaults := GetProfileDefaults(profile)

	// Apply tuned values or profile defaults
	if settings.TunedNCDThreshold != nil {
		settings.NCDThreshold = *settings.TunedNCDThreshold
	} else {
		settings.NCDThreshold = defaults.NCDThreshold
	}
	if settings.TunedPValueThreshold != nil {
		settings.PValueThreshold = *settings.TunedPValueThreshold
	} else {
		settings.PValueThreshold = defaults.PValueThreshold
	}
	settings.BaselineSize = defaults.BaselineSize
	settings.WindowSize = defaults.WindowSize

	return settings, nil
}

// getLastTuneTime returns the most recent tune history timestamp, if any.
func (s *store) getLastTuneTime(ctx context.Context, streamID uuid.UUID) (*time.Time, error) {
	var ts time.Time
	err := s.pool.QueryRow(ctx, `
		SELECT created_at
		FROM threshold_tune_history
		WHERE stream_id = $1
		ORDER BY created_at DESC
		LIMIT 1`,
		streamID).Scan(&ts)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &ts, nil
}

// applyAutoTune runs the auto-tune algorithm for a stream and applies adjustments
func (s *store) applyAutoTune(ctx context.Context, streamID uuid.UUID) error {
	// Get current settings
	settings, err := s.getStreamTuningSettings(ctx, streamID)
	if err != nil {
		return err
	}

	if !settings.AutoTuneEnabled {
		return nil
	}

	// Respect cooldown to avoid thrashing from frequent feedback submissions
	if lastTune, err := s.getLastTuneTime(ctx, streamID); err == nil {
		if shouldThrottleAutoTune(lastTune, defaultAutoTuneConfig, time.Now()) {
			return nil
		}
	} else {
		return err
	}

	// Get feedback from last 30 days
	stats, err := s.getFeedbackStats(ctx, streamID, time.Now().AddDate(0, 0, -30))
	if err != nil {
		return err
	}

	currentNCD := settings.NCDThreshold
	currentPValue := settings.PValueThreshold

	newNCD, newPValue, reason, shouldAdjust := computeAutoTuneAdjustment(
		*stats, currentNCD, currentPValue, defaultAutoTuneConfig)

	if !shouldAdjust {
		return nil
	}

	// Apply changes and record history
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	ncdChanged := math.Abs(newNCD-currentNCD) > 0.01
	pvalueChanged := math.Abs(newPValue-currentPValue) > 0.005

	if ncdChanged {
		_, err = tx.Exec(ctx, `
			UPDATE streams SET tuned_ncd_threshold = $1, detection_profile = 'custom' WHERE id = $2`,
			newNCD, streamID)
		if err != nil {
			return err
		}

		confidence := 1.0 - stats.FalsePositiveRate
		err = s.recordTuneHistoryTx(ctx, tx, TuneHistoryRecord{
			StreamID:   streamID,
			TuneType:   "ncd",
			OldValue:   &currentNCD,
			NewValue:   newNCD,
			Reason:     reason,
			Confidence: confidence,
		})
		if err != nil {
			return err
		}
	}

	if pvalueChanged {
		_, err = tx.Exec(ctx, `
			UPDATE streams SET tuned_pvalue_threshold = $1, detection_profile = 'custom' WHERE id = $2`,
			newPValue, streamID)
		if err != nil {
			return err
		}

		confidence := 1.0 - stats.FalsePositiveRate
		err = s.recordTuneHistoryTx(ctx, tx, TuneHistoryRecord{
			StreamID:   streamID,
			TuneType:   "pvalue",
			OldValue:   &currentPValue,
			NewValue:   newPValue,
			Reason:     reason,
			Confidence: confidence,
		})
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// recordTuneHistoryTx records tune history within a transaction
func (s *store) recordTuneHistoryTx(ctx context.Context, tx pgx.Tx, history TuneHistoryRecord) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO threshold_tune_history (stream_id, tune_type, old_value, new_value, reason, confidence)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		history.StreamID, history.TuneType, history.OldValue, history.NewValue, history.Reason, history.Confidence)
	return err
}

// applyRustRecommendation applies Rust's data-driven threshold recommendation
// This runs after each detection, not triggered by user feedback
func (s *store) applyRustRecommendation(ctx context.Context, streamID uuid.UUID, recommendedNCD float64, stabilityScore float64) error {
	if recommendedNCD <= 0 {
		return nil // Invalid recommendation
	}

	// Get current settings
	settings, err := s.getStreamTuningSettings(ctx, streamID)
	if err != nil {
		return err
	}

	if !settings.AutoTuneEnabled {
		return nil
	}

	// Respect cooldown - same as feedback-based tuning
	if lastTune, err := s.getLastTuneTime(ctx, streamID); err == nil {
		if shouldThrottleAutoTune(lastTune, defaultAutoTuneConfig, time.Now()) {
			return nil
		}
	} else {
		return err
	}

	// Get current threshold (tuned or profile default)
	currentNCD := settings.NCDThreshold

	// Only adjust if recommendation differs by >2%
	diff := math.Abs(recommendedNCD - currentNCD)
	if diff < 0.02 {
		return nil
	}

	// Clamp to reasonable bounds
	newNCD := math.Max(0.1, math.Min(0.8, recommendedNCD))

	// Apply and record
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		UPDATE streams SET tuned_ncd_threshold = $1, detection_profile = 'custom' WHERE id = $2`,
		newNCD, streamID)
	if err != nil {
		return err
	}

	err = s.recordTuneHistoryTx(ctx, tx, TuneHistoryRecord{
		StreamID:   streamID,
		TuneType:   "ncd",
		OldValue:   &currentNCD,
		NewValue:   newNCD,
		Reason:     "rust_recommendation",
		Confidence: stabilityScore, // Use stability score as confidence
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// getTuneHistory retrieves recent tune history for a stream
func (s *store) getTuneHistory(ctx context.Context, streamID uuid.UUID, limit int) ([]TuneHistoryRecord, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, stream_id, tune_type, old_value, new_value, reason, confidence, created_at
		FROM threshold_tune_history
		WHERE stream_id = $1
		ORDER BY created_at DESC
		LIMIT $2`,
		streamID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []TuneHistoryRecord
	for rows.Next() {
		var h TuneHistoryRecord
		if err := rows.Scan(&h.ID, &h.StreamID, &h.TuneType, &h.OldValue, &h.NewValue, &h.Reason, &h.Confidence, &h.CreatedAt); err != nil {
			return nil, err
		}
		history = append(history, h)
	}

	return history, rows.Err()
}

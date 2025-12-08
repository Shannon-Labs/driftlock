-- +goose Up
-- Adaptive Sliding Scales: Detection Profiles, Auto-Tuning, and Adaptive Windows

-- Add detection profile and adaptive settings to streams table
ALTER TABLE streams
ADD COLUMN IF NOT EXISTS detection_profile TEXT NOT NULL DEFAULT 'balanced'
    CHECK (detection_profile IN ('sensitive', 'balanced', 'strict', 'custom')),
ADD COLUMN IF NOT EXISTS auto_tune_enabled BOOLEAN NOT NULL DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS adaptive_window_enabled BOOLEAN NOT NULL DEFAULT FALSE,
-- Tuned threshold overrides (NULL means use profile default)
ADD COLUMN IF NOT EXISTS tuned_ncd_threshold NUMERIC(8,4),
ADD COLUMN IF NOT EXISTS tuned_pvalue_threshold NUMERIC(8,4),
-- Adaptive window bounds
ADD COLUMN IF NOT EXISTS adaptive_baseline_min INTEGER DEFAULT 100,
ADD COLUMN IF NOT EXISTS adaptive_baseline_max INTEGER DEFAULT 2000,
ADD COLUMN IF NOT EXISTS adaptive_window_min INTEGER DEFAULT 10,
ADD COLUMN IF NOT EXISTS adaptive_window_max INTEGER DEFAULT 200;

-- Store user feedback on anomalies (for auto-tuning)
CREATE TABLE IF NOT EXISTS anomaly_feedback (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    anomaly_id UUID NOT NULL REFERENCES anomalies(id) ON DELETE CASCADE,
    stream_id UUID NOT NULL REFERENCES streams(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    feedback_type TEXT NOT NULL CHECK (feedback_type IN ('false_positive', 'confirmed', 'dismissed')),
    -- Capture metrics at feedback time for learning
    ncd_at_detection NUMERIC(8,4) NOT NULL,
    pvalue_at_detection NUMERIC(8,4) NOT NULL,
    confidence_at_detection NUMERIC(8,4) NOT NULL,
    -- User context
    feedback_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by TEXT NOT NULL DEFAULT 'api'
);

CREATE INDEX IF NOT EXISTS anomaly_feedback_stream_idx
ON anomaly_feedback (stream_id, created_at DESC);

CREATE INDEX IF NOT EXISTS anomaly_feedback_tenant_idx
ON anomaly_feedback (tenant_id, created_at DESC);

-- Stream calibration statistics (for adaptive windows)
CREATE TABLE IF NOT EXISTS stream_statistics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stream_id UUID NOT NULL REFERENCES streams(id) ON DELETE CASCADE,
    -- Event frequency metrics
    avg_events_per_hour NUMERIC(12,2),
    avg_event_size_bytes INTEGER,
    event_size_variance NUMERIC(12,2),
    -- Pattern complexity metrics
    avg_baseline_entropy NUMERIC(8,4),
    avg_ncd_self NUMERIC(8,4),  -- NCD of normal data against itself
    pattern_diversity_score NUMERIC(8,4),  -- Higher = more varied patterns
    -- Computed recommendations
    recommended_baseline_size INTEGER,
    recommended_window_size INTEGER,
    -- Metadata
    sample_count INTEGER NOT NULL DEFAULT 0,
    last_computed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (stream_id)
);

-- Auto-tune history (audit trail)
CREATE TABLE IF NOT EXISTS threshold_tune_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stream_id UUID NOT NULL REFERENCES streams(id) ON DELETE CASCADE,
    tune_type TEXT NOT NULL CHECK (tune_type IN ('ncd', 'pvalue', 'baseline', 'window')),
    old_value NUMERIC(12,4),
    new_value NUMERIC(12,4) NOT NULL,
    reason TEXT NOT NULL,  -- 'false_positive_rate', 'detection_rate', 'event_frequency', etc.
    confidence NUMERIC(6,4),  -- How confident we are in this adjustment
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS threshold_tune_history_stream_idx
ON threshold_tune_history (stream_id, created_at DESC);

-- Documentation comments
COMMENT ON COLUMN streams.detection_profile IS 'Sensitivity preset: sensitive (low thresholds), balanced (default), strict (high thresholds), custom (use tuned values)';
COMMENT ON COLUMN streams.auto_tune_enabled IS 'Whether to automatically adjust thresholds based on user feedback';
COMMENT ON COLUMN streams.adaptive_window_enabled IS 'Whether to automatically adjust window sizes based on data characteristics';
COMMENT ON COLUMN streams.tuned_ncd_threshold IS 'Auto-tuned or user-set NCD threshold override (NULL = use profile default)';
COMMENT ON COLUMN streams.tuned_pvalue_threshold IS 'Auto-tuned or user-set p-value threshold override (NULL = use profile default)';
COMMENT ON TABLE anomaly_feedback IS 'User feedback on detected anomalies for threshold auto-tuning';
COMMENT ON TABLE stream_statistics IS 'Computed statistics for adaptive window sizing';
COMMENT ON TABLE threshold_tune_history IS 'Audit trail of automatic threshold adjustments';

-- +goose Down
DROP TABLE IF EXISTS threshold_tune_history;
DROP TABLE IF EXISTS stream_statistics;
DROP TABLE IF EXISTS anomaly_feedback;

ALTER TABLE streams
DROP COLUMN IF EXISTS detection_profile,
DROP COLUMN IF EXISTS auto_tune_enabled,
DROP COLUMN IF EXISTS adaptive_window_enabled,
DROP COLUMN IF EXISTS tuned_ncd_threshold,
DROP COLUMN IF EXISTS tuned_pvalue_threshold,
DROP COLUMN IF EXISTS adaptive_baseline_min,
DROP COLUMN IF EXISTS adaptive_baseline_max,
DROP COLUMN IF EXISTS adaptive_window_min,
DROP COLUMN IF EXISTS adaptive_window_max;

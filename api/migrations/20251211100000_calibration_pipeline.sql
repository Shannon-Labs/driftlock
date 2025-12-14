-- +goose Up
-- Full Calibration Pipeline: DB-driven thresholds with benchmark integration
-- Enables auto-calibration during warmup and feedback-based continuous improvement

-- Profile calibration settings (replaces hardcoded ProfileThresholds)
-- Default values validated on PaySim (AUPRC=1.0) and BAF (AUPRC=0.187) benchmarks
CREATE TABLE IF NOT EXISTS profile_calibrations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_name TEXT NOT NULL UNIQUE,  -- 'sensitive', 'balanced', 'strict', 'financial_fraud', etc.

    -- Detection thresholds
    ncd_threshold NUMERIC(8,4) NOT NULL,
    p_value_threshold NUMERIC(8,4) NOT NULL,
    composite_threshold NUMERIC(8,4) NOT NULL,

    -- Composite score weights (must sum to ~1.0)
    ncd_weight NUMERIC(6,4) NOT NULL DEFAULT 0.5,
    p_value_weight NUMERIC(6,4) NOT NULL DEFAULT 0.25,
    compression_weight NUMERIC(6,4) NOT NULL DEFAULT 0.25,

    -- Window configuration
    baseline_size INTEGER NOT NULL,
    window_size INTEGER NOT NULL,
    permutation_count INTEGER NOT NULL DEFAULT 100,

    -- Adaptive settings
    adaptive_target_fpr NUMERIC(6,4) DEFAULT 0.01,
    require_statistical_significance BOOLEAN NOT NULL DEFAULT TRUE,

    -- Source tracking
    source TEXT NOT NULL DEFAULT 'default',  -- 'default', 'benchmark', 'feedback', 'manual'
    benchmark_auprc NUMERIC(6,4),    -- from benchmark evaluation
    benchmark_f1 NUMERIC(6,4),
    benchmark_dataset TEXT,          -- e.g., 'paysim', 'baf'

    -- Metadata
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by TEXT,

    -- Validation
    CONSTRAINT weights_sum_check CHECK (
        ABS(ncd_weight + p_value_weight + compression_weight - 1.0) < 0.02
    ),
    CONSTRAINT valid_thresholds CHECK (
        ncd_threshold >= 0 AND ncd_threshold <= 1 AND
        p_value_threshold >= 0 AND p_value_threshold <= 1 AND
        composite_threshold >= 0 AND composite_threshold <= 1
    )
);

-- Seed default profile calibrations
INSERT INTO profile_calibrations (
    profile_name, ncd_threshold, p_value_threshold, composite_threshold,
    ncd_weight, p_value_weight, compression_weight,
    baseline_size, window_size, source
) VALUES
('sensitive', 0.20, 0.10, 0.55, 0.5, 0.25, 0.25, 200, 30, 'default'),
('balanced', 0.30, 0.05, 0.60, 0.5, 0.25, 0.25, 400, 50, 'default'),
('strict', 0.40, 0.01, 0.70, 0.5, 0.25, 0.25, 500, 100, 'default'),
-- Financial fraud profile from PaySim benchmark
('financial_fraud', 0.55, 0.08, 0.83, 0.5, 0.25, 0.25, 200, 8, 'benchmark')
ON CONFLICT (profile_name) DO NOTHING;

-- Stream-specific calibration from warmup (auto-calibration results)
CREATE TABLE IF NOT EXISTS stream_calibrations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stream_id UUID NOT NULL REFERENCES streams(id) ON DELETE CASCADE,

    -- Calibration settings
    calibration_method TEXT NOT NULL DEFAULT 'fpr_target',  -- 'fpr_target', 'f1_max', 'manual'
    target_fpr NUMERIC(6,4),           -- for fpr_target method

    -- Calibrated thresholds
    calibrated_threshold NUMERIC(8,4) NOT NULL,
    calibrated_ncd_threshold NUMERIC(8,4),
    calibrated_pvalue_threshold NUMERIC(8,4),

    -- Custom weights (NULL = use profile defaults)
    ncd_weight NUMERIC(6,4),
    p_value_weight NUMERIC(6,4),
    compression_weight NUMERIC(6,4),

    -- Calibration statistics
    warmup_sample_count INTEGER NOT NULL DEFAULT 0,
    warmup_score_mean NUMERIC(8,4),
    warmup_score_stddev NUMERIC(8,4),
    warmup_score_p95 NUMERIC(8,4),
    warmup_score_p99 NUMERIC(8,4),

    -- Observed performance (updated from feedback)
    observed_fpr NUMERIC(6,4),
    observed_f1 NUMERIC(6,4),
    observed_precision NUMERIC(6,4),
    observed_recall NUMERIC(6,4),

    -- Metadata
    calibrated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_validated_at TIMESTAMPTZ,
    validation_sample_count INTEGER DEFAULT 0,

    UNIQUE(stream_id)
);

CREATE INDEX IF NOT EXISTS stream_calibrations_stream_idx
ON stream_calibrations (stream_id);

-- Feedback aggregation for learning (complements anomaly_feedback)
CREATE TABLE IF NOT EXISTS feedback_statistics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Scope (one of these should be set)
    stream_id UUID REFERENCES streams(id) ON DELETE CASCADE,
    profile_name TEXT,  -- NULL if stream-specific, profile name otherwise
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,

    -- Time period
    period_start TIMESTAMPTZ NOT NULL,
    period_end TIMESTAMPTZ NOT NULL,

    -- Counts
    total_detections INTEGER NOT NULL DEFAULT 0,
    confirmed_count INTEGER NOT NULL DEFAULT 0,
    false_positive_count INTEGER NOT NULL DEFAULT 0,
    dismissed_count INTEGER NOT NULL DEFAULT 0,

    -- Performance metrics
    observed_precision NUMERIC(6,4),  -- confirmed / (confirmed + false_positive)
    observed_recall NUMERIC(6,4),     -- requires ground truth, often NULL

    -- Score distribution during period
    avg_composite_score NUMERIC(8,4),
    score_stddev NUMERIC(8,4),
    score_min NUMERIC(8,4),
    score_max NUMERIC(8,4),

    -- Recommendations based on feedback
    recommended_threshold NUMERIC(8,4),
    recommendation_confidence NUMERIC(6,4),

    -- Metadata
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    -- Ensure we have at least one scope
    CONSTRAINT feedback_stats_scope CHECK (
        stream_id IS NOT NULL OR profile_name IS NOT NULL
    )
);

CREATE INDEX IF NOT EXISTS feedback_statistics_stream_idx
ON feedback_statistics (stream_id, period_end DESC)
WHERE stream_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS feedback_statistics_profile_idx
ON feedback_statistics (profile_name, period_end DESC)
WHERE profile_name IS NOT NULL;

-- Add composite threshold and weights to streams table
ALTER TABLE streams
ADD COLUMN IF NOT EXISTS composite_threshold NUMERIC(8,4),
ADD COLUMN IF NOT EXISTS ncd_weight NUMERIC(6,4),
ADD COLUMN IF NOT EXISTS p_value_weight NUMERIC(6,4),
ADD COLUMN IF NOT EXISTS compression_weight NUMERIC(6,4),
ADD COLUMN IF NOT EXISTS calibration_method TEXT DEFAULT 'fpr_target'
    CHECK (calibration_method IN ('fpr_target', 'f1_max', 'manual', NULL));

-- Add composite score to anomalies for analysis
ALTER TABLE anomalies
ADD COLUMN IF NOT EXISTS composite_score NUMERIC(8,4);

-- Documentation comments
COMMENT ON TABLE profile_calibrations IS 'Default detection thresholds and weights per profile, sourced from benchmarks or manual tuning';
COMMENT ON TABLE stream_calibrations IS 'Stream-specific calibrated thresholds from auto-calibration during warmup';
COMMENT ON TABLE feedback_statistics IS 'Aggregated feedback statistics for continuous threshold learning';
COMMENT ON COLUMN profile_calibrations.benchmark_auprc IS 'Area Under PR Curve from benchmark evaluation (1.0 = perfect)';
COMMENT ON COLUMN profile_calibrations.source IS 'Where calibration came from: default, benchmark, feedback, manual';
COMMENT ON COLUMN streams.composite_threshold IS 'Stream-specific composite threshold override (NULL = use profile default)';
COMMENT ON COLUMN streams.calibration_method IS 'How threshold was calibrated: fpr_target (unsupervised), f1_max (supervised), manual';

-- +goose Down
ALTER TABLE anomalies
DROP COLUMN IF EXISTS composite_score;

ALTER TABLE streams
DROP COLUMN IF EXISTS composite_threshold,
DROP COLUMN IF EXISTS ncd_weight,
DROP COLUMN IF EXISTS p_value_weight,
DROP COLUMN IF EXISTS compression_weight,
DROP COLUMN IF EXISTS calibration_method;

DROP TABLE IF EXISTS feedback_statistics;
DROP TABLE IF EXISTS stream_calibrations;
DROP TABLE IF EXISTS profile_calibrations;

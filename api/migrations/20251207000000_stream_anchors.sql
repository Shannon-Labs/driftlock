-- +goose Up
-- SHA-140: Anchor Baseline Strategy for Drift Detection
-- Frozen snapshots of normal behavior at calibration time
-- Detects gradual "boiling frog" drift that sliding windows miss

-- Anchor table stores frozen baseline snapshots
CREATE TABLE IF NOT EXISTS stream_anchors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stream_id UUID NOT NULL REFERENCES streams(id) ON DELETE CASCADE,

    -- Anchor data (compressed baseline snapshot)
    anchor_data BYTEA NOT NULL,
    compressor TEXT NOT NULL DEFAULT 'zstd',
    event_count INTEGER NOT NULL,

    -- Anchor metadata
    calibration_completed_at TIMESTAMPTZ NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    -- Baseline statistics for comparison
    baseline_entropy NUMERIC(8,4),
    baseline_compression_ratio NUMERIC(8,4),
    baseline_ncd_self NUMERIC(8,4),  -- NCD of baseline against itself (should be ~0)

    -- Drift detection thresholds
    drift_ncd_threshold NUMERIC(8,4) NOT NULL DEFAULT 0.35,  -- Slightly higher than shock (0.25)

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    superseded_at TIMESTAMPTZ,  -- When this anchor was replaced
    superseded_by UUID REFERENCES stream_anchors(id)
);

-- Only one active anchor per stream (partial unique index)
CREATE UNIQUE INDEX IF NOT EXISTS unique_active_anchor
ON stream_anchors (stream_id)
WHERE is_active = TRUE;

-- Index for fast lookup of active anchor by stream
CREATE INDEX IF NOT EXISTS stream_anchors_active_idx
ON stream_anchors (stream_id)
WHERE is_active = TRUE;

-- Index for audit trail queries
CREATE INDEX IF NOT EXISTS stream_anchors_history_idx
ON stream_anchors (stream_id, created_at DESC);

-- Add anchor-related settings to streams
ALTER TABLE streams
ADD COLUMN IF NOT EXISTS anchor_enabled BOOLEAN NOT NULL DEFAULT TRUE,
ADD COLUMN IF NOT EXISTS drift_ncd_threshold NUMERIC(8,4) NOT NULL DEFAULT 0.35,
ADD COLUMN IF NOT EXISTS anchor_reset_on_drift BOOLEAN NOT NULL DEFAULT FALSE;

-- Documentation
COMMENT ON TABLE stream_anchors IS 'Frozen baseline snapshots for drift detection (SHA-140)';
COMMENT ON COLUMN stream_anchors.anchor_data IS 'Compressed representation of normal events at calibration time';
COMMENT ON COLUMN stream_anchors.baseline_ncd_self IS 'Self-NCD sanity check (should be ~0)';
COMMENT ON COLUMN stream_anchors.drift_ncd_threshold IS 'NCD threshold above which drift is detected';
COMMENT ON COLUMN stream_anchors.is_active IS 'Only one active anchor per stream; old ones are superseded';
COMMENT ON COLUMN streams.anchor_enabled IS 'Whether drift detection against anchor is enabled';
COMMENT ON COLUMN streams.drift_ncd_threshold IS 'Default NCD threshold for drift detection';
COMMENT ON COLUMN streams.anchor_reset_on_drift IS 'Auto-reset anchor when drift detected (use with caution)';

-- +goose Down
ALTER TABLE streams
DROP COLUMN IF EXISTS anchor_enabled,
DROP COLUMN IF EXISTS drift_ncd_threshold,
DROP COLUMN IF EXISTS anchor_reset_on_drift;

DROP INDEX IF EXISTS stream_anchors_history_idx;
DROP INDEX IF EXISTS stream_anchors_active_idx;
DROP TABLE IF EXISTS stream_anchors;

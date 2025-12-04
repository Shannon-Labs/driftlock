-- +goose Up
-- SHA-139: Cold Start Guardrail - Track stream calibration state
-- NCD requires ~50+ events for statistical significance

-- Add calibration tracking columns to streams table
ALTER TABLE streams
ADD COLUMN IF NOT EXISTS events_ingested BIGINT NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS is_calibrated BOOLEAN NOT NULL DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS calibrated_at TIMESTAMPTZ,
ADD COLUMN IF NOT EXISTS min_baseline_size INTEGER NOT NULL DEFAULT 50;

-- Index for quick lookup of uncalibrated streams (for dashboard queries)
CREATE INDEX IF NOT EXISTS streams_calibration_idx
ON streams (tenant_id, is_calibrated)
WHERE is_calibrated = FALSE;

-- Documentation comments
COMMENT ON COLUMN streams.events_ingested IS 'Total events ever ingested to this stream (cumulative, never resets)';
COMMENT ON COLUMN streams.is_calibrated IS 'True when events_ingested >= min_baseline_size';
COMMENT ON COLUMN streams.calibrated_at IS 'Timestamp when calibration completed (events_ingested reached threshold)';
COMMENT ON COLUMN streams.min_baseline_size IS 'Minimum events required before NCD calculations are statistically significant (default 50)';

-- +goose Down
DROP INDEX IF EXISTS streams_calibration_idx;

ALTER TABLE streams
DROP COLUMN IF EXISTS events_ingested,
DROP COLUMN IF EXISTS is_calibrated,
DROP COLUMN IF EXISTS calibrated_at,
DROP COLUMN IF EXISTS min_baseline_size;

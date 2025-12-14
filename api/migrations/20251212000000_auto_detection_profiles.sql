-- +goose Up
-- Expand stream detection profiles and default new streams to auto-detection.

ALTER TABLE streams
    ALTER COLUMN detection_profile SET DEFAULT 'auto';

-- The original schema constrained detection_profile to a small fixed set.
-- Drop the constraint to allow new auto-detected profiles (financial, llm_safety, etc.)
-- and future DB-driven profile names.
ALTER TABLE streams
    DROP CONSTRAINT IF EXISTS streams_detection_profile_check;

COMMENT ON COLUMN streams.detection_profile IS 'Profile name used for CBAD defaults and DB-driven calibration (auto-detected or user-selected)';

-- +goose Down
-- Best-effort rollback: coerce unknown profiles back to balanced and restore original constraint.

UPDATE streams
SET detection_profile = 'balanced'
WHERE detection_profile NOT IN ('sensitive', 'balanced', 'strict', 'custom');

ALTER TABLE streams
    ALTER COLUMN detection_profile SET DEFAULT 'balanced';

ALTER TABLE streams
    ADD CONSTRAINT streams_detection_profile_check
    CHECK (detection_profile IN ('sensitive', 'balanced', 'strict', 'custom'));

COMMENT ON COLUMN streams.detection_profile IS 'Sensitivity preset: sensitive (low thresholds), balanced (default), strict (high thresholds), custom (use tuned values)';


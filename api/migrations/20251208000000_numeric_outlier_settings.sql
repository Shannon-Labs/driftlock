-- +goose Up
-- SHA-143: Numeric outlier detection settings
ALTER TABLE streams ADD COLUMN IF NOT EXISTS numeric_outlier_enabled BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE streams ADD COLUMN IF NOT EXISTS numeric_k_sigma NUMERIC(4,2) NOT NULL DEFAULT 3.0;

-- +goose Down
ALTER TABLE streams DROP COLUMN IF EXISTS numeric_outlier_enabled;
ALTER TABLE streams DROP COLUMN IF EXISTS numeric_k_sigma;

-- +goose Up
-- Add firebase_uid for Firebase Auth integration
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS firebase_uid TEXT UNIQUE;
CREATE INDEX IF NOT EXISTS tenants_firebase_uid_idx ON tenants (firebase_uid) WHERE firebase_uid IS NOT NULL;

-- Add soft delete support for API keys
ALTER TABLE api_keys ADD COLUMN IF NOT EXISTS revoked_at TIMESTAMPTZ;
CREATE INDEX IF NOT EXISTS api_keys_active_idx ON api_keys (tenant_id) WHERE revoked_at IS NULL;

-- Add billing grace period fields
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS grace_period_ends_at TIMESTAMPTZ;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS payment_failure_count INTEGER DEFAULT 0;

-- +goose Down
DROP INDEX IF EXISTS tenants_firebase_uid_idx;
ALTER TABLE tenants DROP COLUMN IF EXISTS firebase_uid;

DROP INDEX IF EXISTS api_keys_active_idx;
ALTER TABLE api_keys DROP COLUMN IF EXISTS revoked_at;

ALTER TABLE tenants DROP COLUMN IF EXISTS grace_period_ends_at;
ALTER TABLE tenants DROP COLUMN IF EXISTS payment_failure_count;

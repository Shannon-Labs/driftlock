-- +goose Up
-- Add expiration for verification tokens (15-minute TTL for security)
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS verification_token_expires_at TIMESTAMPTZ;

-- Create index for efficient cleanup of expired pending tenants
CREATE INDEX IF NOT EXISTS tenants_verification_expires_idx
ON tenants (verification_token_expires_at)
WHERE verification_token IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS tenants_verification_expires_idx;
ALTER TABLE tenants DROP COLUMN IF EXISTS verification_token_expires_at;

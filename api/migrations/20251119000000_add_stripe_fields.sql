-- +goose Up
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS stripe_customer_id TEXT;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS stripe_subscription_id TEXT;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS stripe_status TEXT;

-- +goose Down
ALTER TABLE tenants DROP COLUMN IF EXISTS stripe_status;
ALTER TABLE tenants DROP COLUMN IF EXISTS stripe_subscription_id;
ALTER TABLE tenants DROP COLUMN IF EXISTS stripe_customer_id;

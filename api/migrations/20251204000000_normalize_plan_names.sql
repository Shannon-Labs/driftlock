-- +goose Up
-- Normalize plan names to canonical values: pulse, radar, tensor, orbit

-- Update legacy free tier names to 'pulse'
UPDATE tenants SET plan = 'pulse' WHERE plan IN ('trial', 'pilot', 'starter');

-- Update legacy basic tier names to 'radar'
UPDATE tenants SET plan = 'radar' WHERE plan IN ('basic', 'signal');

-- Update legacy pro tier names to 'tensor'
UPDATE tenants SET plan = 'tensor' WHERE plan IN ('pro', 'lock', 'transistor', 'sentinel', 'growth');

-- Update legacy enterprise tier names to 'orbit'
UPDATE tenants SET plan = 'orbit' WHERE plan = 'enterprise';

-- Update the default plan in the tenants table
ALTER TABLE tenants ALTER COLUMN plan SET DEFAULT 'pulse';

-- Add a constraint to validate plan values (optional, but recommended)
-- Note: We only add this if it doesn't already exist
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'tenants_plan_check'
    ) THEN
        ALTER TABLE tenants ADD CONSTRAINT tenants_plan_check
            CHECK (plan IN ('pulse', 'radar', 'tensor', 'orbit'));
    END IF;
END $$;

-- +goose Down
-- Remove the constraint if it was added
ALTER TABLE tenants DROP CONSTRAINT IF EXISTS tenants_plan_check;

-- Revert the default
ALTER TABLE tenants ALTER COLUMN plan SET DEFAULT 'pilot';

-- Note: We cannot reliably revert the data changes since we don't know
-- which legacy name each tenant originally had. This is a one-way migration.

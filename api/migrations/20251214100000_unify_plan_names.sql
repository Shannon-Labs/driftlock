-- Unify plan names to canonical Free/Pro/Team/Enterprise
-- This migration converts legacy plan names to the new standard

-- Convert legacy plan names to canonical names
UPDATE tenants SET plan = 'free' WHERE plan IN ('trial', 'pilot', 'pulse');
UPDATE tenants SET plan = 'pro' WHERE plan IN ('radar', 'signal', 'starter');
UPDATE tenants SET plan = 'team' WHERE plan IN ('tensor', 'growth', 'lock');
UPDATE tenants SET plan = 'enterprise' WHERE plan IN ('orbit', 'horizon');

-- Update status for free users who were in trial
UPDATE tenants SET status = 'active' WHERE plan = 'free' AND status = 'trialing' AND trial_ends_at < NOW();

-- Add check constraint for valid plan names (optional, for future enforcement)
-- Note: Keeping legacy names valid during transition period
-- ALTER TABLE tenants ADD CONSTRAINT valid_plan_names
--     CHECK (plan IN ('free', 'pro', 'team', 'enterprise',
--                     'trial', 'pilot', 'pulse', 'radar', 'signal', 'starter',
--                     'tensor', 'growth', 'lock', 'orbit', 'horizon'));

-- Log the migration
DO $$
BEGIN
    RAISE NOTICE 'Plan names unified: trial/pilot/pulse -> free, radar/signal/starter -> pro, tensor/growth/lock -> team, orbit/horizon -> enterprise';
END $$;

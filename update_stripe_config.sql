-- Update Stripe product and price mapping in the database
-- This script sets up the correct Stripe Product and Price IDs for DriftLock

-- First, let's check if plan_price_map table exists
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.tables
        WHERE table_schema = 'public'
        AND table_name = 'plan_price_map'
    ) THEN
        -- Create the table if it doesn't exist
        CREATE TABLE public.plan_price_map (
            plan TEXT PRIMARY KEY,
            stripe_product_id TEXT NOT NULL,
            stripe_price_id TEXT NOT NULL,
            price_cents INTEGER NOT NULL,
            included_calls INTEGER NOT NULL,
            overage_rate_per_call NUMERIC(10,6) NOT NULL,
            created_at TIMESTAMPTZ DEFAULT NOW(),
            updated_at TIMESTAMPTZ DEFAULT NOW()
        );

        RAISE NOTICE 'Created plan_price_map table';
    ELSE
        RAISE NOTICE 'plan_price_map table already exists';
    END IF;
END $$;

-- Clear existing data and insert correct Stripe configuration
DELETE FROM public.plan_price_map;

-- Insert Pro Plan configuration
INSERT INTO public.plan_price_map (
    plan,
    stripe_product_id,
    stripe_price_id,
    price_cents,
    included_calls,
    overage_rate_per_call
) VALUES (
    'pro',
    'prod_TJKXbWnB3ExnqJ',
    'price_1SMhsZL4rhSbUSqA51lWvPlQ',
    4900,  -- $49.00 in cents
    50000, -- 50,000 included calls
    0.001  -- $0.001 per overage call
);

-- Insert Enterprise Plan configuration
INSERT INTO public.plan_price_map (
    plan,
    stripe_product_id,
    stripe_price_id,
    price_cents,
    included_calls,
    overage_rate_per_call
) VALUES (
    'enterprise',
    'prod_TJKXEFXBjkcsAB',
    'price_1SMhshL4rhSbUSqAyHfhWUSQ',
    24900, -- $249.00 in cents
    500000, -- 500,000 included calls
    0.0005  -- $0.0005 per overage call
);

-- Add free/developer plan if needed
INSERT INTO public.plan_price_map (
    plan,
    stripe_product_id,
    stripe_price_id,
    price_cents,
    included_calls,
    overage_rate_per_call
) VALUES (
    'developer',
    NULL, -- No product ID for free plan
    NULL, -- No price ID for free plan
    0,    -- Free
    1000, -- 1,000 included calls
    0.0   -- No overage (hard limit)
) ON CONFLICT (plan) DO NOTHING;

-- Verify the data was inserted correctly
SELECT
    plan,
    stripe_product_id,
    stripe_price_id,
    price_cents,
    included_calls,
    overage_rate_per_call
FROM public.plan_price_map
ORDER BY price_cents;

-- Add RLS policies for the plan_price_map table
ALTER TABLE public.plan_price_map ENABLE ROW LEVEL SECURITY;

-- Drop existing policies if they exist
DROP POLICY IF EXISTS "plan_price_map_public_read" ON public.plan_price_map;
DROP POLICY IF EXISTS "plan_price_map_service_role_all" ON public.plan_price_map;

-- Create policies for the plan_price_map table
CREATE POLICY "plan_price_map_public_read" ON public.plan_price_map
FOR SELECT TO anon, authenticated USING (true);

CREATE POLICY "plan_price_map_service_role_all" ON public.plan_price_map
FOR ALL TO service_role USING (true) WITH CHECK (true);

-- Add updated_at trigger
CREATE OR REPLACE FUNCTION public.set_plan_price_map_updated_at()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS set_plan_price_map_updated_at ON public.plan_price_map;
CREATE TRIGGER set_plan_price_map_updated_at
BEFORE UPDATE ON public.plan_price_map
FOR EACH ROW EXECUTE FUNCTION public.set_plan_price_map_updated_at();

RAISE NOTICE '';
RAISE NOTICE '✅ Stripe product configuration updated successfully';
RAISE NOTICE '✅ Pro Plan: $49/month, 50k calls, $0.001 overage';
RAISE NOTICE '✅ Enterprise Plan: $249/month, 500k calls, $0.0005 overage';
RAISE NOTICE '✅ Developer Plan: Free, 1k calls, hard limit';
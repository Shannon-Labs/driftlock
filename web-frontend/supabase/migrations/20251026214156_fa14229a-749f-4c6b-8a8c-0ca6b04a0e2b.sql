-- ============================================
-- DRIFTLOCK PRODUCTION TABLES
-- Plans, pricing, promotions, events, actions
-- ============================================

-- Plans catalog
CREATE TABLE IF NOT EXISTS public.plans (
  code TEXT PRIMARY KEY,
  display_name TEXT NOT NULL,
  base_price_cents INTEGER NOT NULL,
  included_calls BIGINT NOT NULL,
  overage_rate_cents NUMERIC(10,6) NOT NULL,
  features JSONB DEFAULT '{}',
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Stripe price mapping
CREATE TABLE IF NOT EXISTS public.plan_price_map (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  plan_code TEXT NOT NULL REFERENCES public.plans(code),
  currency TEXT NOT NULL DEFAULT 'usd',
  stripe_price_id TEXT NOT NULL,
  stripe_product_id TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(plan_code, currency)
);

-- Promotions
CREATE TABLE IF NOT EXISTS public.promotions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code TEXT UNIQUE NOT NULL,
  stripe_promotion_code TEXT,
  percent_off INTEGER NOT NULL CHECK (percent_off > 0 AND percent_off <= 100),
  applies_to_plans TEXT[] DEFAULT '{}',
  starts_at TIMESTAMPTZ,
  ends_at TIMESTAMPTZ,
  max_redemptions INTEGER,
  times_redeemed INTEGER DEFAULT 0,
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Stripe events log (for idempotency)
CREATE TABLE IF NOT EXISTS public.stripe_events (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  event_id TEXT UNIQUE NOT NULL,
  type TEXT NOT NULL,
  payload JSONB NOT NULL,
  received_at TIMESTAMPTZ DEFAULT NOW()
);

-- Billing actions log
CREATE TABLE IF NOT EXISTS public.billing_actions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  stripe_event_id TEXT REFERENCES public.stripe_events(event_id),
  organization_id UUID REFERENCES public.organizations(id),
  action TEXT NOT NULL,
  result TEXT NOT NULL,
  details JSONB DEFAULT '{}',
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Quota policy per organization
CREATE TABLE IF NOT EXISTS public.quota_policy (
  organization_id UUID PRIMARY KEY REFERENCES public.organizations(id) ON DELETE CASCADE,
  behavior_on_exceed TEXT DEFAULT 'soft_cap' CHECK (behavior_on_exceed IN ('block', 'soft_cap', 'allow')),
  cap_percent INTEGER DEFAULT 120 CHECK (cap_percent >= 100),
  alert_70 BOOLEAN DEFAULT true,
  alert_90 BOOLEAN DEFAULT true,
  alert_100 BOOLEAN DEFAULT true,
  invoice_threshold_cents INTEGER DEFAULT 50000,
  last_alert_sent_at TIMESTAMPTZ,
  last_alert_type TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Dunning states
CREATE TABLE IF NOT EXISTS public.dunning_states (
  organization_id UUID PRIMARY KEY REFERENCES public.organizations(id) ON DELETE CASCADE,
  state TEXT DEFAULT 'ok' CHECK (state IN ('ok', 'grace', 'suspended', 'cancelled')),
  since TIMESTAMPTZ DEFAULT NOW(),
  notes TEXT,
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================
-- RLS POLICIES
-- ============================================

-- Plans (public read)
ALTER TABLE public.plans ENABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS "plans_public_read" ON public.plans;
CREATE POLICY "plans_public_read" ON public.plans
FOR SELECT USING (is_active = true);

DROP POLICY IF EXISTS "plans_service_role_all" ON public.plans;
CREATE POLICY "plans_service_role_all" ON public.plans
FOR ALL TO service_role USING (true) WITH CHECK (true);

-- Plan price map (public read)
ALTER TABLE public.plan_price_map ENABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS "plan_price_map_public_read" ON public.plan_price_map;
CREATE POLICY "plan_price_map_public_read" ON public.plan_price_map
FOR SELECT USING (true);

DROP POLICY IF EXISTS "plan_price_map_service_role_all" ON public.plan_price_map;
CREATE POLICY "plan_price_map_service_role_all" ON public.plan_price_map
FOR ALL TO service_role USING (true) WITH CHECK (true);

-- Promotions (public read active ones)
ALTER TABLE public.promotions ENABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS "promotions_public_read" ON public.promotions;
CREATE POLICY "promotions_public_read" ON public.promotions
FOR SELECT USING (is_active = true AND NOW() BETWEEN COALESCE(starts_at, '-infinity') AND COALESCE(ends_at, 'infinity'));

DROP POLICY IF EXISTS "promotions_service_role_all" ON public.promotions;
CREATE POLICY "promotions_service_role_all" ON public.promotions
FOR ALL TO service_role USING (true) WITH CHECK (true);

-- Stripe events (service role only)
ALTER TABLE public.stripe_events ENABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS "stripe_events_service_role_all" ON public.stripe_events;
CREATE POLICY "stripe_events_service_role_all" ON public.stripe_events
FOR ALL TO service_role USING (true) WITH CHECK (true);

-- Billing actions (org members can read their own)
ALTER TABLE public.billing_actions ENABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS "billing_actions_org_members" ON public.billing_actions;
CREATE POLICY "billing_actions_org_members" ON public.billing_actions
FOR SELECT TO authenticated USING (public.is_org_member(organization_id));

DROP POLICY IF EXISTS "billing_actions_service_role_all" ON public.billing_actions;
CREATE POLICY "billing_actions_service_role_all" ON public.billing_actions
FOR ALL TO service_role USING (true) WITH CHECK (true);

-- Quota policy (org members can read/update their own)
ALTER TABLE public.quota_policy ENABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS "quota_policy_org_members" ON public.quota_policy;
CREATE POLICY "quota_policy_org_members" ON public.quota_policy
FOR ALL TO authenticated USING (public.is_org_member(organization_id));

DROP POLICY IF EXISTS "quota_policy_service_role_all" ON public.quota_policy;
CREATE POLICY "quota_policy_service_role_all" ON public.quota_policy
FOR ALL TO service_role USING (true) WITH CHECK (true);

-- Dunning states (org members can read their own)
ALTER TABLE public.dunning_states ENABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS "dunning_states_org_members" ON public.dunning_states;
CREATE POLICY "dunning_states_org_members" ON public.dunning_states
FOR SELECT TO authenticated USING (public.is_org_member(organization_id));

DROP POLICY IF EXISTS "dunning_states_service_role_all" ON public.dunning_states;
CREATE POLICY "dunning_states_service_role_all" ON public.dunning_states
FOR ALL TO service_role USING (true) WITH CHECK (true);

-- ============================================
-- INDEXES
-- ============================================

CREATE INDEX IF NOT EXISTS idx_stripe_events_event_id ON public.stripe_events(event_id);
CREATE INDEX IF NOT EXISTS idx_stripe_events_type ON public.stripe_events(type);
CREATE INDEX IF NOT EXISTS idx_billing_actions_org_created ON public.billing_actions(organization_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_plan_price_map_plan_code ON public.plan_price_map(plan_code);
CREATE INDEX IF NOT EXISTS idx_promotions_code ON public.promotions(code);

-- ============================================
-- UPDATED VIEWS
-- ============================================

-- Enhanced usage view with overage calculations
DROP VIEW IF EXISTS public.v_current_period_usage;
CREATE VIEW public.v_current_period_usage
WITH (security_invoker=true)
AS
SELECT
  s.organization_id,
  s.plan,
  s.status,
  s.current_period_start,
  s.current_period_end,
  COALESCE(u.total_calls, 0) as total_calls,
  COALESCE(u.included_calls_used, 0) as included_calls_used,
  COALESCE(u.overage_calls, 0) as overage_calls,
  s.included_calls,
  ROUND((COALESCE(u.total_calls, 0)::DECIMAL / NULLIF(s.included_calls, 0)) * 100, 2) as percent_used,
  COALESCE(u.estimated_charges_cents, 0) as estimated_charges_cents,
  ROUND(COALESCE(u.estimated_charges_cents, 0)::DECIMAL / 100, 2) as estimated_overage_usd,
  EXTRACT(DAY FROM (s.current_period_end - NOW())) as days_remaining,
  (COALESCE(u.total_calls, 0)::DECIMAL / NULLIF(s.included_calls, 0)) >= 0.70 as needs_70_alert,
  (COALESCE(u.total_calls, 0)::DECIMAL / NULLIF(s.included_calls, 0)) >= 0.90 as needs_90_alert,
  (COALESCE(u.total_calls, 0)::DECIMAL / NULLIF(s.included_calls, 0)) >= 1.00 as needs_100_alert,
  qp.cap_percent,
  qp.behavior_on_exceed,
  ds.state as dunning_state
FROM public.subscriptions s
LEFT JOIN public.usage_counters u ON u.organization_id = s.organization_id
  AND u.period_start = s.current_period_start
LEFT JOIN public.quota_policy qp ON qp.organization_id = s.organization_id
LEFT JOIN public.dunning_states ds ON ds.organization_id = s.organization_id
WHERE s.status IN ('active', 'trialing', 'incomplete');

-- ============================================
-- SEED DATA
-- ============================================

-- Insert plans
INSERT INTO public.plans (code, display_name, base_price_cents, included_calls, overage_rate_cents, features) VALUES
('developer', 'Developer', 0, 10000, 0, '{"apis": ["stream", "monitor"], "support": "community"}'),
('standard', 'Standard', 4900, 250000, 0.35, '{"apis": ["stream", "monitor"], "support": "priority", "sla": "99.9%", "retention_days": 90}'),
('growth', 'Growth', 24900, 2000000, 0.18, '{"apis": ["stream", "monitor"], "support": "dedicated", "sla": "99.9%", "retention_days": 180, "custom_thresholds": true}')
ON CONFLICT (code) DO UPDATE SET
  display_name = EXCLUDED.display_name,
  base_price_cents = EXCLUDED.base_price_cents,
  included_calls = EXCLUDED.included_calls,
  overage_rate_cents = EXCLUDED.overage_rate_cents,
  features = EXCLUDED.features;

-- Insert promotion (LAUNCH50)
INSERT INTO public.promotions (code, percent_off, applies_to_plans, starts_at, ends_at, max_redemptions) VALUES
('LAUNCH50', 50, ARRAY['standard', 'growth'], NOW(), NOW() + INTERVAL '3 months', 1000)
ON CONFLICT (code) DO UPDATE SET
  percent_off = EXCLUDED.percent_off,
  applies_to_plans = EXCLUDED.applies_to_plans,
  ends_at = EXCLUDED.ends_at;

-- ============================================
-- TRIGGERS
-- ============================================

DROP TRIGGER IF EXISTS set_updated_at ON public.quota_policy;
CREATE TRIGGER set_updated_at
BEFORE UPDATE ON public.quota_policy
FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();

DROP TRIGGER IF EXISTS set_updated_at ON public.dunning_states;
CREATE TRIGGER set_updated_at
BEFORE UPDATE ON public.dunning_states
FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();
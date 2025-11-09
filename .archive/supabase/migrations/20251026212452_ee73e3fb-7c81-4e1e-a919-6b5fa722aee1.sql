-- ============================================
-- DRIFTLOCK MULTI-TENANT BILLING SCHEMA
-- Implements organizations, Stripe integration, 
-- usage-based billing, and RLS isolation
-- ============================================

-- ============================================
-- 1. ORGANIZATIONS & MULTI-TENANCY
-- ============================================

-- Organizations table (core multi-tenancy)
CREATE TABLE IF NOT EXISTS public.organizations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  slug TEXT UNIQUE NOT NULL,
  settings JSONB DEFAULT '{}',
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Organization membership
CREATE TABLE IF NOT EXISTS public.organization_members (
  organization_id UUID REFERENCES public.organizations(id) ON DELETE CASCADE,
  user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
  role TEXT NOT NULL CHECK (role IN ('owner', 'admin', 'member', 'viewer')),
  joined_at TIMESTAMPTZ DEFAULT NOW(),
  PRIMARY KEY (organization_id, user_id)
);

-- ============================================
-- 2. BILLING TABLES
-- ============================================

-- Billing customers (Stripe integration)
CREATE TABLE IF NOT EXISTS public.billing_customers (
  organization_id UUID PRIMARY KEY REFERENCES public.organizations(id) ON DELETE CASCADE,
  stripe_customer_id TEXT UNIQUE NOT NULL,
  billing_email TEXT,
  company_name TEXT,
  tax_id TEXT,
  tax_country TEXT,
  tax_postal_code TEXT,
  default_payment_method_last4 TEXT,
  default_payment_method_brand TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Modify existing subscriptions table to support organization model
ALTER TABLE public.subscriptions 
  ADD COLUMN IF NOT EXISTS organization_id UUID REFERENCES public.organizations(id) ON DELETE CASCADE,
  ADD COLUMN IF NOT EXISTS plan TEXT CHECK (plan IN ('developer', 'standard', 'growth', 'enterprise')),
  ADD COLUMN IF NOT EXISTS status TEXT DEFAULT 'active' CHECK (status IN ('active', 'canceled', 'past_due', 'unpaid', 'incomplete')),
  ADD COLUMN IF NOT EXISTS current_period_start TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS current_period_end TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS price_stream_id TEXT,
  ADD COLUMN IF NOT EXISTS price_monitor_id TEXT,
  ADD COLUMN IF NOT EXISTS included_calls BIGINT DEFAULT 0,
  ADD COLUMN IF NOT EXISTS overage_rate_per_call NUMERIC(10,6) DEFAULT 0,
  ADD COLUMN IF NOT EXISTS has_launch50_promo BOOLEAN DEFAULT FALSE;

-- Usage counters (pooled billing for Stream + Monitor APIs)
CREATE TABLE IF NOT EXISTS public.usage_counters (
  organization_id UUID NOT NULL REFERENCES public.organizations(id) ON DELETE CASCADE,
  period_start TIMESTAMPTZ NOT NULL,
  period_end TIMESTAMPTZ NOT NULL,
  total_calls BIGINT NOT NULL DEFAULT 0,
  included_calls_used BIGINT NOT NULL DEFAULT 0,
  overage_calls BIGINT NOT NULL DEFAULT 0,
  estimated_charges_cents BIGINT NOT NULL DEFAULT 0,
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  PRIMARY KEY (organization_id, period_start)
);

-- Billing events (Stripe webhook tracking)
CREATE TABLE IF NOT EXISTS public.billing_events (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  organization_id UUID REFERENCES public.organizations(id) ON DELETE CASCADE,
  stripe_event_id TEXT UNIQUE,
  event_type TEXT NOT NULL,
  payload JSONB NOT NULL,
  processed_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Invoices mirror (Stripe invoice cache)
CREATE TABLE IF NOT EXISTS public.invoices_mirror (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  organization_id UUID REFERENCES public.organizations(id) ON DELETE CASCADE,
  stripe_invoice_id TEXT UNIQUE NOT NULL,
  status TEXT NOT NULL,
  amount_due_cents BIGINT NOT NULL,
  amount_paid_cents BIGINT NOT NULL,
  hosted_invoice_url TEXT,
  invoice_pdf_url TEXT,
  finalized_at TIMESTAMPTZ,
  paid_at TIMESTAMPTZ,
  period_start TIMESTAMPTZ,
  period_end TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================
-- 3. APPLICATION TABLES
-- ============================================

-- Modify api_keys to support organizations
ALTER TABLE public.api_keys 
  ADD COLUMN IF NOT EXISTS organization_id UUID REFERENCES public.organizations(id) ON DELETE CASCADE,
  ADD COLUMN IF NOT EXISTS key_name TEXT,
  ADD COLUMN IF NOT EXISTS permissions JSONB DEFAULT '{"stream": true, "monitor": true}',
  ADD COLUMN IF NOT EXISTS created_by UUID REFERENCES auth.users(id);

-- Anomaly events (billable events - only these count toward usage)
CREATE TABLE IF NOT EXISTS public.anomaly_events (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  organization_id UUID NOT NULL REFERENCES public.organizations(id) ON DELETE CASCADE,
  anomaly_type TEXT NOT NULL,
  severity TEXT NOT NULL CHECK (severity IN ('critical', 'high', 'medium', 'low')),
  description TEXT,
  raw_data JSONB,
  explanation TEXT,
  detection_timestamp TIMESTAMPTZ DEFAULT NOW(),
  resolved_at TIMESTAMPTZ,
  resolved_by UUID REFERENCES auth.users(id),
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Modify detections table to link to organizations
ALTER TABLE public.detections
  ADD COLUMN IF NOT EXISTS organization_id UUID REFERENCES public.organizations(id) ON DELETE CASCADE;

-- Audit logs (compliance & security)
CREATE TABLE IF NOT EXISTS public.audit_logs (
  id BIGSERIAL PRIMARY KEY,
  organization_id UUID REFERENCES public.organizations(id) ON DELETE CASCADE,
  user_id UUID REFERENCES auth.users(id),
  action TEXT NOT NULL,
  resource_type TEXT NOT NULL,
  resource_id UUID,
  ip_address INET,
  user_agent TEXT,
  details JSONB,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Organization settings (sensitivity controls, dunning behavior)
CREATE TABLE IF NOT EXISTS public.org_settings (
  organization_id UUID PRIMARY KEY REFERENCES public.organizations(id) ON DELETE CASCADE,
  anomaly_sensitivity NUMERIC(3,2) DEFAULT 0.5 CHECK (anomaly_sensitivity BETWEEN 0 AND 1),
  dunning_behavior TEXT DEFAULT 'soft_cap' CHECK (dunning_behavior IN ('block_immediately', 'soft_cap', 'allow_overage')),
  usage_alert_thresholds JSONB DEFAULT '{"70": true, "90": true, "100": true}',
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================
-- 4. RLS HELPER FUNCTIONS
-- ============================================

-- Check if user is member of organization
CREATE OR REPLACE FUNCTION public.is_org_member(org_id UUID)
RETURNS BOOLEAN
LANGUAGE SQL STABLE SECURITY DEFINER
SET search_path = public
AS $$
  SELECT EXISTS (
    SELECT 1
    FROM public.organization_members m
    WHERE m.organization_id = org_id
      AND m.user_id = auth.uid()
  );
$$;

-- Check if user is member or service role
CREATE OR REPLACE FUNCTION public.is_org_member_or_service(org_id UUID)
RETURNS BOOLEAN
LANGUAGE SQL STABLE SECURITY DEFINER
SET search_path = public
AS $$
  SELECT
    auth.role() = 'service_role'
    OR EXISTS (
      SELECT 1
      FROM public.organization_members m
      WHERE m.organization_id = org_id AND m.user_id = auth.uid()
    );
$$;

-- ============================================
-- 5. ENABLE RLS ON ALL TABLES
-- ============================================

ALTER TABLE public.organizations ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.organization_members ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.billing_customers ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.usage_counters ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.billing_events ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.invoices_mirror ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.anomaly_events ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.audit_logs ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.org_settings ENABLE ROW LEVEL SECURITY;

-- ============================================
-- 6. RLS POLICIES
-- ============================================

-- Organizations policies
DROP POLICY IF EXISTS "orgs_members_read" ON public.organizations;
CREATE POLICY "orgs_members_read" ON public.organizations
FOR SELECT TO authenticated USING (public.is_org_member(id));

DROP POLICY IF EXISTS "orgs_service_role_all" ON public.organizations;
CREATE POLICY "orgs_service_role_all" ON public.organizations
FOR ALL TO service_role USING (true) WITH CHECK (true);

-- Organization members policies
DROP POLICY IF EXISTS "org_members_self_org" ON public.organization_members;
CREATE POLICY "org_members_self_org" ON public.organization_members
FOR SELECT TO authenticated USING (public.is_org_member(organization_id));

DROP POLICY IF EXISTS "org_members_service_role_all" ON public.organization_members;
CREATE POLICY "org_members_service_role_all" ON public.organization_members
FOR ALL TO service_role USING (true) WITH CHECK (true);

-- Billing customers policies
DROP POLICY IF EXISTS "billing_customers_org_members" ON public.billing_customers;
CREATE POLICY "billing_customers_org_members" ON public.billing_customers
FOR ALL TO authenticated USING (public.is_org_member(organization_id));

DROP POLICY IF EXISTS "billing_customers_service_role_all" ON public.billing_customers;
CREATE POLICY "billing_customers_service_role_all" ON public.billing_customers
FOR ALL TO service_role USING (true) WITH CHECK (true);

-- Subscriptions policies (update existing)
DROP POLICY IF EXISTS "subscriptions_org_members" ON public.subscriptions;
CREATE POLICY "subscriptions_org_members" ON public.subscriptions
FOR ALL TO authenticated USING (
  organization_id IS NULL AND auth.uid() = user_id 
  OR public.is_org_member(organization_id)
);

-- Usage counters policies
DROP POLICY IF EXISTS "usage_counters_org_members" ON public.usage_counters;
CREATE POLICY "usage_counters_org_members" ON public.usage_counters
FOR ALL TO authenticated USING (public.is_org_member(organization_id));

DROP POLICY IF EXISTS "usage_counters_service_role_all" ON public.usage_counters;
CREATE POLICY "usage_counters_service_role_all" ON public.usage_counters
FOR ALL TO service_role USING (true) WITH CHECK (true);

-- Billing events policies
DROP POLICY IF EXISTS "billing_events_org_members" ON public.billing_events;
CREATE POLICY "billing_events_org_members" ON public.billing_events
FOR SELECT TO authenticated USING (public.is_org_member(organization_id));

DROP POLICY IF EXISTS "billing_events_service_role_all" ON public.billing_events;
CREATE POLICY "billing_events_service_role_all" ON public.billing_events
FOR ALL TO service_role USING (true) WITH CHECK (true);

-- Invoices mirror policies
DROP POLICY IF EXISTS "invoices_mirror_org_members" ON public.invoices_mirror;
CREATE POLICY "invoices_mirror_org_members" ON public.invoices_mirror
FOR SELECT TO authenticated USING (public.is_org_member(organization_id));

DROP POLICY IF EXISTS "invoices_mirror_service_role_all" ON public.invoices_mirror;
CREATE POLICY "invoices_mirror_service_role_all" ON public.invoices_mirror
FOR ALL TO service_role USING (true) WITH CHECK (true);

-- API keys policies (update existing)
DROP POLICY IF EXISTS "api_keys_org_members" ON public.api_keys;
CREATE POLICY "api_keys_org_members" ON public.api_keys
FOR ALL TO authenticated USING (
  organization_id IS NULL AND user_id IN (SELECT id FROM users WHERE firebase_uid = (auth.uid())::text OR id = auth.uid())
  OR public.is_org_member(organization_id)
);

-- Anomaly events policies
DROP POLICY IF EXISTS "anomaly_events_org_members" ON public.anomaly_events;
CREATE POLICY "anomaly_events_org_members" ON public.anomaly_events
FOR ALL TO authenticated USING (public.is_org_member(organization_id));

DROP POLICY IF EXISTS "anomaly_events_service_role_all" ON public.anomaly_events;
CREATE POLICY "anomaly_events_service_role_all" ON public.anomaly_events
FOR ALL TO service_role USING (true) WITH CHECK (true);

-- Detections policies (update existing to support organizations)
DROP POLICY IF EXISTS "detections_org_members" ON public.detections;
CREATE POLICY "detections_org_members" ON public.detections
FOR ALL TO authenticated USING (
  organization_id IS NULL AND user_id IN (SELECT id FROM users WHERE firebase_uid = (auth.uid())::text OR id = auth.uid())
  OR public.is_org_member(organization_id)
);

-- Audit logs policies
DROP POLICY IF EXISTS "audit_logs_org_members" ON public.audit_logs;
CREATE POLICY "audit_logs_org_members" ON public.audit_logs
FOR SELECT TO authenticated USING (public.is_org_member(organization_id));

DROP POLICY IF EXISTS "audit_logs_service_role_all" ON public.audit_logs;
CREATE POLICY "audit_logs_service_role_all" ON public.audit_logs
FOR ALL TO service_role USING (true) WITH CHECK (true);

-- Org settings policies
DROP POLICY IF EXISTS "org_settings_org_members" ON public.org_settings;
CREATE POLICY "org_settings_org_members" ON public.org_settings
FOR ALL TO authenticated USING (public.is_org_member(organization_id));

DROP POLICY IF EXISTS "org_settings_service_role_all" ON public.org_settings;
CREATE POLICY "org_settings_service_role_all" ON public.org_settings
FOR ALL TO service_role USING (true) WITH CHECK (true);

-- ============================================
-- 7. USAGE METERING RPC FUNCTION
-- ============================================

CREATE OR REPLACE FUNCTION public.increment_usage(
  p_org UUID,
  p_period_start TIMESTAMPTZ,
  p_count BIGINT DEFAULT 1
) RETURNS VOID
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = public
AS $$
DECLARE
  v_included_calls BIGINT;
BEGIN
  -- Get included calls for this subscription
  SELECT included_calls INTO v_included_calls
  FROM public.subscriptions
  WHERE organization_id = p_org
    AND current_period_start = p_period_start
  LIMIT 1;

  -- Atomically increment the usage counter
  UPDATE public.usage_counters
  SET
    total_calls = total_calls + p_count,
    included_calls_used = LEAST(total_calls + p_count, COALESCE(v_included_calls, 0)),
    overage_calls = GREATEST(0, (total_calls + p_count) - COALESCE(v_included_calls, 0)),
    updated_at = NOW()
  WHERE organization_id = p_org AND period_start = p_period_start;

  -- If no rows were updated, the counter doesn't exist yet - create it
  IF NOT FOUND THEN
    INSERT INTO public.usage_counters (
      organization_id,
      period_start,
      period_end,
      total_calls,
      included_calls_used,
      overage_calls
    )
    SELECT
      p_org,
      p_period_start,
      s.current_period_end,
      p_count,
      LEAST(p_count, COALESCE(s.included_calls, 0)),
      GREATEST(0, p_count - COALESCE(s.included_calls, 0))
    FROM public.subscriptions s
    WHERE s.organization_id = p_org
      AND s.current_period_start = p_period_start
    LIMIT 1;
  END IF;
END$$;

-- Security: only allow service role to call this
REVOKE ALL ON FUNCTION public.increment_usage(uuid, timestamptz, bigint) FROM PUBLIC;
GRANT EXECUTE ON FUNCTION public.increment_usage(uuid, timestamptz, bigint) TO service_role;

-- ============================================
-- 8. DASHBOARD VIEWS
-- ============================================

-- Current period usage view
CREATE OR REPLACE VIEW public.v_current_period_usage AS
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
  EXTRACT(DAY FROM (s.current_period_end - NOW())) as days_remaining,
  (COALESCE(u.total_calls, 0)::DECIMAL / NULLIF(s.included_calls, 0)) >= 0.70 as needs_70_alert,
  (COALESCE(u.total_calls, 0)::DECIMAL / NULLIF(s.included_calls, 0)) >= 0.90 as needs_90_alert,
  (COALESCE(u.total_calls, 0)::DECIMAL / NULLIF(s.included_calls, 0)) >= 1.00 as needs_100_alert
FROM public.subscriptions s
LEFT JOIN public.usage_counters u ON u.organization_id = s.organization_id
  AND u.period_start = s.current_period_start
WHERE s.status IN ('active', 'trialing', 'incomplete');

-- Anomaly summary view
CREATE OR REPLACE VIEW public.v_org_anomaly_summary AS
SELECT
  ae.organization_id,
  s.current_period_start,
  s.current_period_end,
  COUNT(*) as total_anomalies,
  COUNT(*) FILTER (WHERE ae.severity = 'critical') as critical_anomalies,
  COUNT(*) FILTER (WHERE ae.severity = 'high') as high_anomalies,
  COUNT(*) FILTER (WHERE ae.severity = 'medium') as medium_anomalies,
  COUNT(*) FILTER (WHERE ae.severity = 'low') as low_anomalies,
  MAX(ae.detection_timestamp) as last_detection,
  MIN(ae.detection_timestamp) as first_detection_in_period
FROM public.anomaly_events ae
JOIN public.subscriptions s ON ae.organization_id = s.organization_id
WHERE ae.detection_timestamp >= s.current_period_start
  AND ae.detection_timestamp <= s.current_period_end
GROUP BY ae.organization_id, s.current_period_start, s.current_period_end;

-- ============================================
-- 9. TRIGGERS FOR UPDATED_AT
-- ============================================

CREATE OR REPLACE FUNCTION public.set_updated_at()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS set_updated_at ON public.organizations;
CREATE TRIGGER set_updated_at
BEFORE UPDATE ON public.organizations
FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();

DROP TRIGGER IF EXISTS set_updated_at ON public.billing_customers;
CREATE TRIGGER set_updated_at
BEFORE UPDATE ON public.billing_customers
FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();

DROP TRIGGER IF EXISTS set_updated_at ON public.org_settings;
CREATE TRIGGER set_updated_at
BEFORE UPDATE ON public.org_settings
FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();

-- ============================================
-- 10. INDEXES FOR PERFORMANCE
-- ============================================

CREATE INDEX IF NOT EXISTS idx_org_members_user_id ON public.organization_members(user_id);
CREATE INDEX IF NOT EXISTS idx_org_members_org_id ON public.organization_members(organization_id);
CREATE INDEX IF NOT EXISTS idx_usage_counters_org_period ON public.usage_counters(organization_id, period_start);
CREATE INDEX IF NOT EXISTS idx_billing_events_org_id ON public.billing_events(organization_id);
CREATE INDEX IF NOT EXISTS idx_billing_events_stripe_event_id ON public.billing_events(stripe_event_id);
CREATE INDEX IF NOT EXISTS idx_anomaly_events_org_timestamp ON public.anomaly_events(organization_id, detection_timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_org_created ON public.audit_logs(organization_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_api_keys_org_id ON public.api_keys(organization_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_org_id ON public.subscriptions(organization_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_stripe_id ON public.subscriptions(stripe_subscription_id);
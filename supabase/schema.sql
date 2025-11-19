-- DriftLock Supabase Database Schema
-- This creates all tables required by the Cloudflare Workers API Gateway

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Organizations table - Multi-tenant support
CREATE TABLE IF NOT EXISTS organizations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Users table - User authentication
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    encrypted_password VARCHAR(255),
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    organization_id UUID REFERENCES organizations(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- API Keys table - For Cloudflare Worker authentication
CREATE TABLE IF NOT EXISTS api_keys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    key_hash VARCHAR(255) UNIQUE NOT NULL,  -- SHA-256 hashed API key
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    scopes TEXT[] DEFAULT '{}',
    expires_at TIMESTAMPTZ,
    is_active BOOLEAN NOT NULL DEFAULT true,
    rate_limit_per_minute INTEGER DEFAULT 60,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_used_at TIMESTAMPTZ
);

-- Indexes for api_keys
CREATE INDEX IF NOT EXISTS idx_api_keys_hash ON api_keys(key_hash);
CREATE INDEX IF NOT EXISTS idx_api_keys_org ON api_keys(organization_id);
CREATE INDEX IF NOT EXISTS idx_api_keys_active ON api_keys(is_active) WHERE is_active = true;

-- Billing Customers - Stripe customer mapping
CREATE TABLE IF NOT EXISTS billing_customers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    stripe_customer_id VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index
CREATE INDEX IF NOT EXISTS idx_billing_customers_stripe ON billing_customers(stripe_customer_id);

-- Subscriptions - Stripe subscription tracking
CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    stripe_subscription_id VARCHAR(255) UNIQUE,
    stripe_customer_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    plan VARCHAR(50) NOT NULL,
    current_period_start TIMESTAMPTZ,
    current_period_end TIMESTAMPTZ,
    cancel_at_period_end BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_subscriptions_org ON subscriptions(organization_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_status ON subscriptions(status);
CREATE INDEX IF NOT EXISTS idx_subscriptions_stripe ON subscriptions(stripe_subscription_id);

-- Usage tracking - For billing and metering
CREATE TABLE IF NOT EXISTS usage_tracking (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    event_type VARCHAR(100) NOT NULL,
    count INTEGER NOT NULL DEFAULT 1,
    period_start TIMESTAMPTZ NOT NULL,
    period_end TIMESTAMPTZ NOT NULL,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_usage_tracking_org ON usage_tracking(organization_id);
CREATE INDEX IF NOT EXISTS idx_usage_tracking_period ON usage_tracking(period_start, period_end);

-- Audit logs - API call tracking
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    action VARCHAR(255) NOT NULL,
    resource_type VARCHAR(100),
    ip_address INET,
    user_agent TEXT,
    details JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_audit_logs_org ON audit_logs(organization_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created ON audit_logs(created_at DESC);

-- Invoices mirror - Stripe invoice tracking
CREATE TABLE IF NOT EXISTS invoices_mirror (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    stripe_invoice_id VARCHAR(255) UNIQUE NOT NULL,
    stripe_customer_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    amount_due_cents INTEGER NOT NULL,
    amount_paid_cents INTEGER DEFAULT 0,
    currency VARCHAR(10) NOT NULL DEFAULT 'usd',
    period_start TIMESTAMPTZ,
    period_end TIMESTAMPTZ,
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_invoices_org ON invoices_mirror(organization_id);
CREATE INDEX IF NOT EXISTS idx_invoices_stripe ON invoices_mirror(stripe_invoice_id);

-- Triggers to auto-update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_organizations_updated_at
    BEFORE UPDATE ON organizations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_billing_customers_updated_at
    BEFORE UPDATE ON billing_customers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_subscriptions_updated_at
    BEFORE UPDATE ON subscriptions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert default organization (idempotent)
INSERT INTO organizations (id, name, slug)
VALUES ('00000000-0000-0000-0000-000000000000', 'Default Organization', 'default')
ON CONFLICT (id) DO NOTHING;

-- Row Level Security (RLS) - Enable and create policies
ALTER TABLE organizations ENABLE ROW LEVEL SECURITY;
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE api_keys ENABLE ROW LEVEL SECURITY;
ALTER TABLE billing_customers ENABLE ROW LEVEL SECURITY;
ALTER TABLE subscriptions ENABLE ROW LEVEL SECURITY;
ALTER TABLE usage_tracking ENABLE ROW LEVEL SECURITY;
ALTER TABLE audit_logs ENABLE ROW LEVEL SECURITY;
ALTER TABLE invoices_mirror ENABLE ROW LEVEL SECURITY;

-- Simple policies (adjust for your auth model)
CREATE POLICY "Enable read access for all users" ON organizations FOR SELECT USING (true);
CREATE POLICY "Enable read access for all users" ON users FOR SELECT USING (true);
CREATE POLICY "Enable read access for all users" ON api_keys FOR SELECT USING (true);
CREATE POLICY "Enable read access for all users" ON billing_customers FOR SELECT USING (true);
CREATE POLICY "Enable read access for all users" ON subscriptions FOR SELECT USING (true);
CREATE POLICY "Enable read access for all users" ON usage_tracking FOR SELECT USING (true);
CREATE POLICY "Enable read access for all users" ON audit_logs FOR SELECT USING (true);
CREATE POLICY "Enable read access for all users" ON invoices_mirror FOR SELECT USING (true);

-- Enable inserts/updates for service role (bypasses RLS)
CREATE POLICY "Enable all for service role" ON organizations FOR ALL TO service_role USING (true) WITH CHECK (true);
CREATE POLICY "Enable all for service role" ON users FOR ALL TO service_role USING (true) WITH CHECK (true);
CREATE POLICY "Enable all for service role" ON api_keys FOR ALL TO service_role USING (true) WITH CHECK (true);
CREATE POLICY "Enable all for service role" ON billing_customers FOR ALL TO service_role USING (true) WITH CHECK (true);
CREATE POLICY "Enable all for service role" ON subscriptions FOR ALL TO service_role USING (true) WITH CHECK (true);
CREATE POLICY "Enable all for service role" ON usage_tracking FOR ALL TO service_role USING (true) WITH CHECK (true);
CREATE POLICY "Enable all for service role" ON audit_logs FOR ALL TO service_role USING (true) WITH CHECK (true);
CREATE POLICY "Enable all for service role" ON invoices_mirror FOR ALL TO service_role USING (true) WITH CHECK (true);

-- Comments for documentation
COMMENT ON TABLE organizations IS 'Multi-tenant organization support';
COMMENT ON TABLE users IS 'User accounts with organization membership';
COMMENT ON TABLE api_keys IS 'API keys for programmatic access, hashed for security';
COMMENT ON TABLE billing_customers IS 'Maps Stripe customers to organizations';
COMMENT ON TABLE subscriptions IS 'Stripe subscription status tracking';
COMMENT ON TABLE usage_tracking IS 'Usage metering for billing';
COMMENT ON TABLE audit_logs IS 'Audit trail for API calls and changes';
COMMENT ON TABLE invoices_mirror IS 'Stripe invoice status mirror';

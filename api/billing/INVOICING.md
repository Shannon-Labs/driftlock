# Billing & Subscription Management

## Stripe Integration Plan

### Supported Plans

| Plan | Monthly Cost | Events Included | Rate Limit | Support |
|------|-------------|-----------------|------------|---------|
| Trial | $0 | 10,000 | 60 req/min | Community |
| Starter | $99 | 500,000 | 100 req/min | Email |
| Growth | $499 | 5M | 500 req/min | Email + Slack |
| Enterprise | Custom | Unlimited | Custom | Dedicated |

### Implementation Phases

#### Phase 1: Track Usage (Pre-launch)
- [x] Database schema has `tenant.plan` column
- [ ] Add Stripe customer_id to tenants table
- [ ] Track event counts per tenant/stream
- [ ] Track API request counts
- [ ] Log overages for monitoring

#### Phase 2: Stripe Integration (Post-launch)
- [ ] Create Stripe products and prices
- [ ] Implement checkout flow
- [ ] Add subscription webhooks
- [ ] Handle payment failures
- [ ] Create customer portal

#### Phase 3: Enforcement (Post-launch)
- [ ] Soft limits with warnings
- [ ] Hard limits with upgrade prompts
- [ ] Rate limit based on plan
- [ ] Auto-upgrade on payment

### Database Changes

```sql
CREATE TABLE stripe_customers (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    stripe_customer_id TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE tenants ADD COLUMN stripe_customer_id TEXT;
ALTER TABLE tenants ADD COLUMN plan_started_at TIMESTAMPTZ;
ALTER TABLE tenants ADD COLUMN current_period_end TIMESTAMPTZ;

CREATE TABLE usage_metrics (
    tenant_id UUID NOT NULL,
    stream_id UUID NOT NULL,
    date DATE NOT NULL,
    event_count BIGINT NOT NULL DEFAULT 0,
    api_request_count BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (tenant_id, stream_id, date)
);
```

### Environment Variables
- `STRIPE_SECRET_KEY` - Stripe API key (live or test)
- `STRIPE_WEBHOOK_SECRET` - For webhook validation
- `STRIPE_PUBLISHABLE_KEY` - Frontend checkout

### Monitoring
- Daily usage alerts at 80% and 100% of plan limits
- Failed payment notifications
- Churn risk detection
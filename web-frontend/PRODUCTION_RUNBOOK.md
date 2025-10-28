# Driftlock Production Runbook

## 🎯 System Overview

Driftlock provides anomaly detection through two APIs (Stream + Monitor) with pooled billing - only anomaly detections count toward quota, not data ingestion.

**Architecture:**
- **Frontend**: React + TypeScript
- **Database**: Supabase PostgreSQL with RLS
- **Billing**: Stripe (subscriptions + metered usage)
- **Backend**: Supabase Edge Functions
- **Email**: Resend

---

## 📋 Pre-Launch Checklist

### 1. Stripe Configuration

#### Create Products & Prices
```bash
# In Stripe Dashboard (https://dashboard.stripe.com/test/products)

# Developer Plan (Free)
- Name: Driftlock Developer
- Description: 10k included anomaly detections/month
- Price: $0/month
- Type: Subscription
- Metadata: { tier: "developer", included_calls: 10000 }

# Standard Plan
- Name: Driftlock Standard  
- Description: 250k included anomaly detections/month
- Base Price: $49/month
- Metered Price: $0.0035 per anomaly detection (overage)
- Metadata: { tier: "standard", included_calls: 250000 }

# Growth Plan
- Name: Driftlock Growth
- Description: 2M included anomaly detections/month
- Base Price: $249/month
- Metered Price: $0.0018 per anomaly detection (overage)
- Metadata: { tier: "growth", included_calls: 2000000 }
```

#### Update plan_price_map Table
```sql
-- Run in Supabase SQL Editor
INSERT INTO public.plan_price_map (plan_code, stripe_price_id, stripe_product_id) VALUES
('standard', 'price_xxxxxxxxxxxxx', 'prod_xxxxxxxxxxxxx'),
('growth', 'price_xxxxxxxxxxxxx', 'prod_xxxxxxxxxxxxx')
ON CONFLICT (plan_code, currency) DO UPDATE SET
  stripe_price_id = EXCLUDED.stripe_price_id,
  stripe_product_id = EXCLUDED.stripe_product_id;
```

#### Configure Webhook
```
Webhook URL: https://nfkdeeunyvnntvpvwpwh.supabase.co/functions/v1/stripe-webhook

Events to subscribe:
☑ checkout.session.completed
☑ customer.subscription.created
☑ customer.subscription.updated
☑ customer.subscription.deleted
☑ invoice.paid
☑ invoice.payment_failed

Webhook Secret: Save to STRIPE_WEBHOOK_SECRET in Supabase
```

#### Create Promotion Code
```
Code: LAUNCH50
Discount: 50% off
Duration: 3 months
Applies to: Standard & Growth plans
```

### 2. Email Configuration (Resend)

1. Sign up at https://resend.com
2. Verify domain: https://resend.com/domains
3. Create API key: https://resend.com/api-keys
4. Add to Supabase secrets as `RESEND_API_KEY`
5. Update "from" address in `send-alert-email/index.ts`

### 3. Environment Variables

Ensure all secrets are set in Supabase:
```bash
STRIPE_SECRET_KEY=sk_live_xxxxx
STRIPE_PUBLISHABLE_KEY=pk_live_xxxxx
STRIPE_WEBHOOK_SECRET=whsec_xxxxx
RESEND_API_KEY=re_xxxxx
```

---

## 🚀 Deployment Steps

### 1. Deploy Edge Functions
```bash
# Functions deploy automatically via Lovable
# Verify deployment at:
# https://supabase.com/dashboard/project/nfkdeeunyvnntvpvwpwh/functions
```

### 2. Verify Database Schema
```sql
-- Check all tables exist
SELECT table_name FROM information_schema.tables 
WHERE table_schema = 'public' 
AND table_name IN (
  'organizations', 'subscriptions', 'usage_counters',
  'plans', 'plan_price_map', 'promotions',
  'stripe_events', 'billing_actions', 'quota_policy',
  'dunning_states', 'invoices_mirror'
);

-- Verify plans are seeded
SELECT * FROM public.plans;
```

### 3. Test Webhook
```bash
# Use Stripe CLI to test webhook locally
stripe listen --forward-to https://nfkdeeunyvnntvpvwpwh.supabase.co/functions/v1/stripe-webhook

# Trigger test events
stripe trigger checkout.session.completed
stripe trigger invoice.paid
```

---

## 🧪 Testing Scenarios

### Test 1: New User Signup (Developer Plan)
```
1. User signs up via OAuth
2. Verify organization created
3. Verify subscription record with plan='developer', included_calls=10000
4. Verify usage_counter row exists with 0 calls
5. Verify quota_policy created with defaults
```

### Test 2: Upgrade to Standard Plan
```
1. User clicks "Upgrade" in dashboard
2. Stripe Checkout opens with LAUNCH50 promo applied (50% off)
3. Complete checkout
4. Webhook fires: checkout.session.completed
5. Verify subscription updated to plan='standard', included_calls=250000
6. Verify billing_customer record created
7. Verify usage_counter persists from Developer plan
```

### Test 3: Usage Metering (Anomaly Detection Only)
```
# Test with anomaly=false (should NOT increment)
curl -X POST https://nfkdeeunyvnntvpvwpwh.supabase.co/functions/v1/meter-usage \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"organization_id": "xxx", "anomaly": false, "count": 100}'

# Verify total_calls unchanged

# Test with anomaly=true (should increment)
curl -X POST https://nfkdeeunyvnntvpvwpwh.supabase.co/functions/v1/meter-usage \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"organization_id": "xxx", "anomaly": true, "count": 1}'

# Verify total_calls incremented by 1
```

### Test 4: Usage Alerts
```
# Simulate 70% usage
UPDATE usage_counters 
SET total_calls = 175000 
WHERE organization_id = 'xxx';

# Trigger meter-usage
# Verify alert email sent for 'usage_70'

# Simulate 100% usage
UPDATE usage_counters 
SET total_calls = 250000 
WHERE organization_id = 'xxx';

# Trigger meter-usage
# Verify alert email sent for 'usage_100'
```

### Test 5: Soft Cap (120% limit)
```
# Set usage to 121% of included
UPDATE usage_counters 
SET total_calls = 302500 
WHERE organization_id = 'xxx';

# Trigger meter-usage with behavior_on_exceed='soft_cap'
# Should return 429 status code
```

### Test 6: Payment Failure
```
# Trigger Stripe event
stripe trigger invoice.payment_failed

# Verify:
1. Invoice mirrored to invoices_mirror table
2. dunning_states set to 'grace'
3. Alert email sent to billing contact
```

---

## 📊 Monitoring

### Key Metrics to Watch

#### Database Queries
```sql
-- Active subscriptions by plan
SELECT plan, COUNT(*) 
FROM subscriptions 
WHERE status = 'active' 
GROUP BY plan;

-- Total usage this month
SELECT SUM(total_calls) as total_detections,
       SUM(overage_calls) as total_overage
FROM usage_counters 
WHERE period_start >= date_trunc('month', NOW());

-- Revenue estimate (overage)
SELECT SUM(estimated_charges_cents) / 100 as estimated_revenue_usd
FROM usage_counters
WHERE period_start >= date_trunc('month', NOW());

-- Organizations in dunning
SELECT COUNT(*) FROM dunning_states WHERE state != 'ok';
```

#### Edge Function Logs
```
Monitor at:
- Stripe Webhook: https://supabase.com/dashboard/project/nfkdeeunyvnntvpvwpwh/functions/stripe-webhook/logs
- Meter Usage: https://supabase.com/dashboard/project/nfkdeeunyvnntvpvwpwh/functions/meter-usage/logs
- Alert Emails: https://supabase.com/dashboard/project/nfkdeeunyvnntvpvwpwh/functions/send-alert-email/logs
```

### Alert Conditions
- Stripe webhook failures
- Usage metering errors  
- Email delivery failures
- Dunning state escalations
- High overage charges (>$1000/org)

---

## 🔧 Common Operations

### Manually Issue Invoice
```typescript
// Call from Supabase SQL Editor or Edge Function
const stripe = new Stripe(STRIPE_SECRET_KEY);
await stripe.invoices.create({
  customer: 'cus_xxxxx',
  auto_advance: true,
});
```

### Reset Usage for Testing
```sql
UPDATE usage_counters 
SET total_calls = 0, 
    included_calls_used = 0, 
    overage_calls = 0,
    estimated_charges_cents = 0
WHERE organization_id = 'xxx';
```

### Change Dunning State
```sql
UPDATE dunning_states 
SET state = 'ok', 
    since = NOW(),
    notes = 'Payment resolved manually'
WHERE organization_id = 'xxx';
```

### Adjust Quota Policy
```sql
UPDATE quota_policy
SET behavior_on_exceed = 'soft_cap',
    cap_percent = 150,
    invoice_threshold_cents = 100000
WHERE organization_id = 'xxx';
```

---

## 🚨 Incident Response

### Webhook Not Processing
1. Check Supabase function logs
2. Verify webhook signature matches STRIPE_WEBHOOK_SECRET
3. Check stripe_events table for duplicate event_id
4. Re-send event from Stripe Dashboard

### Payment Failed but User Still Active
1. Check dunning_states table
2. Verify invoice status in Stripe
3. Manually trigger invoice.paid webhook if needed
4. Update dunning state to 'ok'

### Usage Not Incrementing
1. Verify API calls include anomaly=true
2. Check meter-usage function logs
3. Ensure organization_id is correct
4. Verify subscription is active

### Email Alerts Not Sending
1. Check Resend dashboard for delivery status
2. Verify RESEND_API_KEY is set
3. Check send-alert-email function logs
4. Verify billing_email exists in billing_customers

---

## 📚 Resources

- **Stripe Dashboard**: https://dashboard.stripe.com
- **Supabase Dashboard**: https://supabase.com/dashboard/project/nfkdeeunyvnntvpvwpwh
- **Resend Dashboard**: https://resend.com
- **Edge Function Logs**: https://supabase.com/dashboard/project/nfkdeeunyvnntvpvwpwh/functions
- **Database SQL Editor**: https://supabase.com/dashboard/project/nfkdeeunyvnntvpvwpwh/sql

---

## 🔐 Security Checklist

- [x] RLS enabled on all tables
- [x] Service role used for sensitive operations
- [x] Stripe webhook signature verification
- [x] No PII logged in Edge Functions
- [x] API keys hashed in database
- [x] CORS properly configured
- [x] Rate limiting on meter-usage endpoint
- [x] Audit logs for billing actions

---

## 📝 Next Steps

1. **Set up monitoring alerts** (e.g., PagerDuty, Sentry)
2. **Configure backup/disaster recovery**
3. **Load testing** with realistic traffic patterns
4. **Customer portal** integration (Stripe Customer Portal)
5. **Analytics dashboard** for business metrics
6. **Documentation** for API consumers

# Driftlock SaaS Launch Checklist

**Target Launch Date:** [To be set]  
**Current Phase:** Pre-Launch Infrastructure ✅ Ready

## 🎯 Pre-Launch Status

### ✅ Infrastructure Complete
- [x] Cloud Run API deployment configured
- [x] Cloudflare Pages frontend ready
- [x] Supabase PostgreSQL schemas ready
- [x] CI/CD pipeline (`cloudbuild.yaml`)
- [x] Rate limiting & authentication implemented
- [x] Multi-tenant database schema
- [x] Health checks & monitoring
- [x] API proxy functions for frontend
- [x] Domain & SSL configuration ready

### 📋 Remaining Tasks Before Launch

#### 1. User Onboarding (High Priority)
- [ ] **Create public onboarding endpoint** (`/v1/onboard/signup`)
- [ ] **Add email verification flow** (SendGrid or similar)
- [ ] **Create welcome email template**
- [ ] **Implement rate limiting by IP** for signup
- [ ] **Add signup analytics tracking**

**Estimated:** 2-3 hours

#### 2. Self-Service Features (Medium Priority)
- [ ] **Add signup form to landing page**
- [ ] **Create "Get Started" CTA flow**
- [ ] **Build minimal dashboard for new tenants**
- [ ] **Add "Upgrade Plan" buttons**

**Estimated:** 4-6 hours

#### 3. Billing Foundation (Medium Priority)
- [ ] **Create Stripe account and products**
- [ ] **Add Stripe customer tracking to database**
- [ ] **Implement usage metrics collection**
- [ ] **Add overage email alerts**
- [ ] **Create customer portal (Stripe Billing)**

**Estimated:** 3-4 hours

#### 4. Admin & Monitoring (Medium Priority)
- [ ] **Create admin dashboard view**
- [ ] **Set up tenant analytics queries**
- [ ] **Configure alert channels (PagerDuty/Slack)**
- [ ] **Create runbook for common issues**
- [ ] **Set up log aggregation (if needed)**

**Estimated:** 2-3 hours

#### 5. Testing & Validation (High Priority)
- [ ] **End-to-end deployment test** (see below)
- [ ] **Load testing** (1000 req/sec baseline)
- [ ] **Security review** (CORS, auth, rate limiting)
- [ ] **Compliance check** (GDPR data handling)
- [ ] **Chaos testing** (what happens if DB fails)

**Estimated:** 4-5 hours

#### 6. Documentation & Launch Prep (Low Priority)
- [ ] **Create API documentation** (Swagger/OpenAPI)
- [ ] **Write "Getting Started" guide**
- [ ] **Create pricing page**
- [ ] **Set up support email/Discord**
- [ ] **Write launch announcement**

**Estimated:** 3-4 hours

---

## 🚀 Quick Launch Path (MVP Approach)

If you want to launch **this week**, focus on these critical items:

### Week 1: MVP Launch
**Goal:** First 5 beta customers

**Day 1-2: Onboarding**
- [ ] Implement `/v1/onboard/signup` endpoint
- [ ] Manual email verification (you check database, email them)
- [ ] Add simple signup: "Email support@driftlock.net for access"

**Day 3-4: Deployment**
- [ ] Follow [COMPLETE_DEPLOYMENT_PLAN.md](./COMPLETE_DEPLOYMENT_PLAN.md)
- [ ] Test with your own data
- [ ] Create 1-2 demo tenants
- [ ] Verify billing page exists (even if manual invoicing)

**Day 5: Launch**
- [ ] Tweet/Y Combinator post
- [ ] Email 10 potential beta customers
- [ ] Monitor first signups manually
- [ ] Invoice manually via Stripe (no automated billing yet)

**Time to launch:** ~2-3 days of focused work

### Week 2-3: Iterate
- [ ] Automate email verification
- [ ] Add self-service upgrade flow
- [ ] Implement usage tracking
- [ ] Add more monitoring/alerts

### Week 4: Growth
- [ ] Full Stripe integration
- [ ] Customer portal
- [ ] Advanced analytics
- [ ] Scale Cloud Run resources

---

## 🔍 End-to-End Testing Procedure

Run this before launch to verify everything works:

### Step 1: Database
```bash
# 1. Connect to Supabase
gcloud secrets versions access latest --secret=driftlock-db-url | psql

# 2. Verify tables exist
\dt

# 3. Should see: tenants, streams, api_keys, anomalies, etc.
```

### Step 2: API Deployment
```bash
# 1. Deploy to Cloud Run
gcloud builds submit --config=cloudbuild.yaml

# 2. Get service URL
SERVICE_URL=$(gcloud run services describe driftlock-api \
  --region us-central1 \
  --format 'value(status.url)')

# 3. Test health endpoint
curl $SERVICE_URL/healthz | jq
```

**Expected:** All green, database connected

### Step 3: Frontend Deployment
```bash
# 1. Deploy landing page
cd landing-page && npm run build
wrangler pages deploy dist --project-name=driftlock

# 2. Test custom domain
curl https://driftlock.net/healthz
```

**Expected:** Returns same JSON as direct Cloud Run call

### Step 4: Multi-tenant Flow
```bash
# 1. Create test tenant
API_KEY=$(./collector-processor/cmd/driftlock-http/driftlock-http \
  --dev-mode \
  create-tenant \
  --name "Test Corp" \
  --plan trial \
  --key-role admin \
  --json | jq -r .api_key)

# 2. Test authentication
curl -X POST https://driftlock.net/api/v1/detect \
  -H "X-Api-Key: $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"events": [...], "window_size": 100}'

# 3. Check anomalies saved
psql $DATABASE_URL -c "SELECT COUNT(*) FROM anomalies LIMIT 5"
```

### Step 5: Load Test
```bash
# Install k6 if needed
npm install -g k6

# Run load test (adjust duration as needed)
k6 run --vus 10 --duration 30s -e API_URL=$SERVICE_URL -e API_KEY=$API_KEY scripts/load-test.js
```

**Expected:** <500ms p95 latency, <5% error rate

---

## 📊 Success Metrics (First 30 Days)

### Technical
- [ ] 99.9% uptime (max 43 minutes downtime/month)
- [ ] <500ms p95 API latency
- [ ] <1% error rate on `/v1/detect`
- [ ] Zero data loss

### Business
- [ ] 10+ beta signups
- [ ] 5+ active tenants
- [ ] 1+ paying customer
- [ ] <$50 total infra cost

### Product
- [ ] 3+ customers complete onboarding solo
- [ ] Zero critical security issues
- [ ] First export job completed
- [ ] First upgrade from trial to paid

---

## 🚨 Escalation Contacts

**Technical Issues:**
- API down → Check Cloud Run dashboard, rollback if needed
- DB issues → Check Supabase status, contact support if >15 min
- Payment issues → Manually invoice via Stripe while debugging

**Business Issues:**
- Customer complaint → Respond within 2 hours, offer refund if needed
- Churn risk → Personal email from founder
- Feature request → Add to backlog, communicate timeline

---

## ✅ Go/No-Go Decision Criteria

**Launch is GO if:**
- [ ] End-to-end test passes
- [ ] Demo tenant works flawlessly
- [ ] Payment processing tested (even if manual)
- [ ] Monitoring alerts configured
- [ ] Rollback plan documented
- [ ] Support channel ready

**Launch is NO-GO if:**
- [ ] Database connection unreliable
- [ ] API errors >5% in testing
- [ ] No way to contact support
- [ ] No billing process (even manual)
- [ ] Security concerns unresolved

---

**Estimated Time to Launch:** 2-3 days for MVP, 2 weeks for full launch

**Next Step:** Pick "Quick Launch" or "Full Launch" path and start Day 1 tasks.
# Driftlock Full-Featured Launch Roadmap

**Last Updated**: November 2024
**Timeline**: 2-3 Weeks for Full Launch
**Scope**: Everything except full Stripe integration

---

## Current Status Summary

### Completed Features

- **Core API** - Full anomaly detection with Rust/Go FFI
- **Database Schema** - Multi-tenant with all required tables
- **Authentication** - API key-based auth with rate limiting
- **Landing Page** - Marketing site with playground
- **Deployment Config** - Cloud Run, Supabase, Cloudflare ready
- **Documentation** - Comprehensive docs for developers
- **Testing** - API tests and demo scripts
- **CI/CD** - GitHub workflows for builds

### Newly Implemented (This Session)

- **Onboarding System** - Self-service signup endpoint
- **Email Service** - SendGrid integration for notifications
- **Usage Tracking** - Tenant usage metrics and plan limits
- **Admin Dashboard** - Web-based tenant management
- **SignupForm** - Frontend signup component
- **Load Testing** - k6-based performance tests
- **Getting Started Guide** - User onboarding documentation

---

## Phase Breakdown

### Phase 1: Infrastructure Setup (Days 1-2) - COMPLETE

- [x] Supabase database with connection pooling
- [x] Initial schema migration (`20250301120000_initial_schema.sql`)
- [x] Google Cloud project configuration
- [x] Secret Manager setup for credentials
- [x] service.yaml and cloudbuild.yaml configured

**Files**:
- `api/migrations/20250301120000_initial_schema.sql`
- `service.yaml`
- `cloudbuild.yaml`

---

### Phase 2: API Deployment (Days 2-3) - COMPLETE

- [x] Rust core (cbad-core) compilation
- [x] Go HTTP service build
- [x] Docker image creation
- [x] Cloud Run deployment automation
- [x] Health check endpoint

**Commands**:
```bash
gcloud builds submit --config=cloudbuild.yaml
```

---

### Phase 3: Frontend Deployment (Day 3) - COMPLETE

- [x] Cloudflare Pages configuration
- [x] API proxy functions
- [x] Landing page with playground
- [x] Contact form submission

**Commands**:
```bash
cd landing-page && npm run build && wrangler pages deploy dist
```

---

### Phase 4: Onboarding System (Days 4-6) - COMPLETE

- [x] Signup endpoint (`/v1/onboard/signup`)
- [x] Rate limiting (5 signups/hour/IP)
- [x] Email validation
- [x] Duplicate checking
- [x] Auto-tenant creation with trial plan
- [x] Immediate API key return
- [x] SignupForm.vue component
- [x] Onboarding migration (`20250302000000_onboarding.sql`)

**New Files**:
- `collector-processor/cmd/driftlock-http/onboarding.go`
- `api/migrations/20250302000000_onboarding.sql`
- `landing-page/src/components/cta/SignupForm.vue`

**Updated Files**:
- `collector-processor/cmd/driftlock-http/main.go` (added routes)

---

### Phase 5: Email Automation (Days 7-9) - COMPLETE

- [x] SendGrid API integration
- [x] Welcome email template
- [x] Verification email template
- [x] Trial expiration warning template
- [x] Async email sending

**New Files**:
- `collector-processor/cmd/driftlock-http/email.go`

**Environment Variables**:
```bash
SENDGRID_API_KEY=SG.xxx
EMAIL_FROM_ADDRESS=noreply@driftlock.net
EMAIL_FROM_NAME=Driftlock
```

**TODO for Production**:
- [ ] Set up SendGrid account and verify sender
- [ ] Store API key in Secret Manager
- [ ] Implement verification flow endpoint

---

### Phase 6: Usage Tracking (Days 10-11) - COMPLETE

- [x] Usage metrics table in database
- [x] Event/anomaly/request tracking
- [x] Plan limits definition (trial, starter, growth, enterprise)
- [x] Usage summary queries
- [x] Plan limit checking (80%, 100%, 120% thresholds)
- [x] Daily aggregation job structure

**New Files**:
- `collector-processor/cmd/driftlock-http/usage.go`

**TODO for Production**:
- [ ] Add usage tracking call to detectHandler
- [ ] Set up cron job for daily aggregation
- [ ] Implement usage warning emails

---

### Phase 7: Admin Dashboard (Days 12-14) - COMPLETE

- [x] Admin authentication (X-Admin-Key header)
- [x] Tenant list endpoint (`/v1/admin/tenants`)
- [x] Usage metrics endpoint (`/v1/admin/tenants/:id/usage`)
- [x] AdminDashboard.vue with full UI
- [x] Search and filter functionality
- [x] Usage details modal

**New Files**:
- `landing-page/src/views/AdminDashboard.vue`

**Updated Files**:
- `landing-page/src/router/index.ts` (added /admin route)

**Access**:
- URL: `https://driftlock.net/admin`
- Auth: X-Admin-Key header with ADMIN_KEY env var

---

### Phase 8: Stripe Setup (Days 15-16) - PARTIAL

**Completed**:
- [x] Plan definitions in code
- [x] Usage tracking foundation

**TODO**:
- [ ] Create Stripe products (Trial, Starter, Growth, Enterprise)
- [ ] Store Stripe keys in Secret Manager
- [ ] Document manual billing workflow in `api/billing/INVOICING.md`

**NOT implementing in MVP**:
- Checkout flow
- Subscription webhooks
- Customer portal
- Automatic billing

---

### Phase 9: Testing & Validation (Days 17-18) - PARTIAL

**Completed**:
- [x] Load testing script (`scripts/load-test.js`)
- [x] Existing API tests

**TODO**:
- [ ] Run comprehensive test suite
- [ ] Manual testing checklist
- [ ] Run load tests with k6
- [ ] Fix any failing tests

**Commands**:
```bash
# Install k6
brew install k6  # macOS

# Run load test
k6 run scripts/load-test.js

# Run with custom settings
k6 run --vus 10 --duration 30s scripts/load-test.js
```

**Target Metrics**:
- p95 latency < 500ms for /healthz
- p95 latency < 5s for /v1/detect
- Error rate < 1%
- Throughput > 100 req/sec

---

### Phase 10: Launch Preparation (Days 19-21) - PARTIAL

**Completed**:
- [x] GETTING_STARTED.md documentation
- [x] This ROADMAP.md

**TODO**:
- [ ] Update README.md with signup instructions
- [ ] Set up monitoring and alerts
- [ ] Configure backup strategy
- [ ] Security review
- [ ] Cost monitoring

---

## Quick Reference

### New API Endpoints

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/v1/onboard/signup` | POST | None | Create account |
| `/v1/admin/tenants` | GET | Admin | List tenants |
| `/v1/admin/tenants/:id/usage` | GET | Admin | Get usage |

### Environment Variables Needed

```bash
# Required
DATABASE_URL=postgresql://...
DRIFTLOCK_DEV_MODE=true  # or license key

# Optional
SENDGRID_API_KEY=SG.xxx
ADMIN_KEY=your-admin-secret
EMAIL_FROM_ADDRESS=noreply@driftlock.net
APP_URL=https://driftlock.net
```

### Key Commands

```bash
# Build and deploy API
gcloud builds submit --config=cloudbuild.yaml

# Deploy frontend
cd landing-page && npm run build && wrangler pages deploy dist

# Run load tests
k6 run scripts/load-test.js

# View logs
gcloud run services logs read driftlock-api --region us-central1 --limit 50

# Check health
curl https://driftlock.net/api/v1/healthz | jq
```

---

## Remaining Work (Prioritized)

### High Priority (Before Launch)

1. **Run database migration** - Apply `20250302000000_onboarding.sql`
2. **Set ADMIN_KEY** - Configure admin authentication
3. **Test signup flow** - End-to-end testing
4. **Deploy updates** - Push to Cloud Run and Cloudflare
5. **Monitoring setup** - Basic uptime and error alerts

### Medium Priority (Week After Launch)

1. **SendGrid setup** - Enable email notifications
2. **Usage tracking integration** - Connect to detect handler
3. **Manual billing workflow** - Document Stripe process
4. **Load testing** - Run performance benchmarks
5. **Security review** - Complete checklist

### Low Priority (Future Iterations)

1. **Email verification flow** - Full verification system
2. **Plan enforcement** - Hard limits and overage handling
3. **Customer portal** - Self-service plan management
4. **Automated billing** - Full Stripe integration
5. **Advanced analytics** - Usage dashboards and charts

---

## Success Metrics (First 30 Days)

### Technical

- [ ] 99.9% uptime
- [ ] < 500ms p95 API latency
- [ ] < 1% error rate
- [ ] Zero data loss

### Business

- [ ] 20+ signups
- [ ] 10+ verified users
- [ ] 5+ active API users
- [ ] 1+ paying customer (manual)

### Product

- [ ] 5+ users complete onboarding solo
- [ ] 3+ users make multiple API calls
- [ ] Zero critical security issues

---

## Launch Day Checklist

### Pre-Launch

- [ ] Run migrations
- [ ] Deploy API with new endpoints
- [ ] Deploy frontend with signup form
- [ ] Configure ADMIN_KEY
- [ ] Test signup flow
- [ ] Test admin dashboard
- [ ] Enable monitoring

### Launch

- [ ] Soft launch to personal network
- [ ] Post on Twitter/LinkedIn
- [ ] Submit to Hacker News
- [ ] Email beta list
- [ ] Monitor logs

### Post-Launch

- [ ] Daily monitoring
- [ ] Personal emails to new users
- [ ] Collect feedback
- [ ] Fix bugs immediately

---

## File Reference

### New Files Created

```
collector-processor/cmd/driftlock-http/
├── onboarding.go    # Signup endpoint + admin endpoints
├── email.go         # SendGrid email service
└── usage.go         # Usage tracking + plan limits

api/migrations/
└── 20250302000000_onboarding.sql  # Email + usage tables

landing-page/src/
├── components/cta/
│   └── SignupForm.vue    # Signup form component
└── views/
    └── AdminDashboard.vue  # Admin management UI

scripts/
└── load-test.js     # k6 load testing script

docs/
└── GETTING_STARTED.md  # User onboarding guide

ROADMAP.md           # This file
```

### Modified Files

```
collector-processor/cmd/driftlock-http/main.go  # Added routes
landing-page/src/router/index.ts               # Added /admin route
```

---

## Contact & Support

- **Project**: [github.com/Shannon-Labs/driftlock](https://github.com/Shannon-Labs/driftlock)
- **Email**: hunter@shannonlabs.dev
- **Website**: [driftlock.net](https://driftlock.net)

---

*This roadmap is a living document. Update as features are completed and priorities change.*

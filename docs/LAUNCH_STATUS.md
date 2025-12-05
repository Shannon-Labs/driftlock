# Launch Development Guide

> **For AI Assistants:** Start here! This section tracks what's done, what's next, and how to continue development.

## Current Status: ~98% Launch Ready (Deployment In Progress)

**Last Updated:** 2025-12-04
**Target:** Public launch with self-serve signup, working billing, demo/playground, API access

### What's Done
- All backend features implemented and tested
- All frontend features implemented
- Documentation updated
- E2E tests written
- Security audit passed (no SQL injection, all endpoints auth-protected)
- **[2025-12-04] AI Agent Integration Complete:**
  - Fixed Firebase→Cloud Run routing (service name + /api prefix)
  - Complete OpenAPI 3.0 spec with all 20+ endpoints
  - AI agent integration guide (`/docs/ai-agents/INTEGRATION.md`)
  - Use case documentation (`/docs/use-cases/general-anomaly-detection.md`)

### What's In Progress
- Cloud Build deploying new code (build ID: 37eefe5f)
- Final production verification pending

---

## Project Management

**Primary**: [Linear - shannon-labs/driftlock](https://linear.app/shannon-labs/project/driftlock-a8c80503816c/overview)
**Automation Level**: Heavy (auto-triage, daily standups, velocity reports)

### Linear MCP Commands

```
"Show me issues in the driftlock project"
"Create a Linear issue: Bug - API returns 500 on empty payload"
"Mark DRI-123 as Done"
"What issues are blocked or need attention?"
```

### Workflow

1. Check Linear for current priorities
2. Pick an issue and move to "In Progress"
3. Create branch: `git checkout -b DRI-123-description`
4. Code & commit - PRs auto-link to Linear
5. Merge - Issue auto-closes when PR merges

---

## Launch Checklist

### Phase 1: User Onboarding (COMPLETE)

- [x] Database migration for verification flow
- [x] Email verification backend
- [x] Verify endpoint implementation
- [x] API key regeneration/create/revoke endpoints
- [x] SignupForm.vue - pending verification state
- [x] VerifyEmailView.vue - verification landing page
- [x] Router update for `/verify` route

### Phase 2: Stripe Billing (COMPLETE)

- [x] 14-day trial to checkout
- [x] Complete webhook handlers (trial_will_end, payment_failed, payment_succeeded)
- [x] Grace period logic (7-day after payment failure)
- [x] Billing status endpoint
- [x] Frontend billing UI (trial banners, grace period warning)
- [x] Pricing page checkout

### Phase 3: Developer Experience (COMPLETE)

- [x] Anonymous demo endpoint (`POST /v1/demo/detect`)
- [x] Playground demo mode
- [x] Usage dashboard with charts

### Phase 4: Polish & Testing (MOSTLY COMPLETE)

- [x] E2E test: signup → verify → detect → anomaly
- [x] E2E test: trial → checkout → subscription
- [x] Error code reference page
- [ ] Load test Cloud Run deployment (optional)
- [x] Production deployment verification (in progress)

### Phase 5: AI Agent Integration (COMPLETE)

- [x] Firebase routing fix
- [x] Path prefix middleware
- [x] Complete OpenAPI 3.0 spec
- [x] AI agent integration guide
- [x] Use case documentation

---

## Next Priority Tasks

### Immediate (if deployment not verified)
1. Check build status: `gcloud builds list --limit=1 --project=driftlock`
2. Verify production: `curl https://driftlock.net/api/healthz`
3. Test demo endpoint

### Pre-Launch Verification
4. Run E2E tests: `go test ./collector-processor/cmd/driftlock-http/... -v`
5. Verify Stripe webhooks: `stripe listen`
6. Test full flow: Signup → Verify → Detect → Dashboard

### Optional Enhancements
7. Load testing
8. Redis rate limiting for multi-instance

---

## Pricing Tiers

| Tier | Price | Events/Month |
|------|-------|--------------|
| Pulse (Free) | $0 | 10,000 |
| Radar | $15/mo | 500,000 |
| Tensor | $100/mo | 5,000,000 |
| Orbit | $499/mo | Unlimited |

---

## Useful Commands

```bash
# Local development
docker compose up -d
cd collector-processor && go run ./cmd/driftlock-http

# Frontend
cd landing-page && npm run dev

# Database
goose -dir api/migrations postgres "$DATABASE_URL" up

# Testing
curl -X POST http://localhost:8080/api/v1/onboard/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","company_name":"Test Co","plan":"trial"}'

# Stripe webhooks
stripe listen --forward-to localhost:8080/api/v1/billing/webhook
```

---

**Last updated:** 2025-12-04

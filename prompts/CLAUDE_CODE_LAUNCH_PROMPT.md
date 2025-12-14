# Driftlock Launch Sprint - Claude Code Multi-Agent Prompt

**Objective:** Get Driftlock to production launch using parallel agents coordinated via Linear.

---

## Project Context

**Driftlock** is a compression-based anomaly detection (CBAD) platform for OpenTelemetry data. We're at ~98% launch ready. The remaining work is verification, final fixes, and deployment confirmation.

### Tech Stack
- **Backend:** Go 1.22+, PostgreSQL, Redis, Kafka
- **Core Engine:** Rust FFI (cbad-core)
- **Frontend:** Vue 3, TypeScript, Tailwind
- **Billing:** Stripe
- **Auth:** Firebase JWT + API keys
- **AI Analysis:** Gemini 2.5 Flash (just integrated)
- **Hosting:** GCP Cloud Run, Firebase Hosting

### Key URLs
- **Production:** https://driftlock.net (or https://driftlock.web.app)
- **API:** https://driftlock.net/api/v1
- **Linear Project:** https://linear.app/shannon-labs/project/driftlock-a8c80503816c

---

## MCP Servers Available

Use these MCP servers for coordination:

1. **Linear MCP** - Issue tracking and coordination
   - `list_issues` - Get current issues
   - `create_issue` - Create new issues
   - `update_issue` - Update status
   
2. **Supabase MCP** - Database operations (if needed)

3. **Stripe MCP** - Billing verification

---

## Agent Assignments

Spawn these agents in parallel. Each agent should:
1. Check Linear for their assigned issues
2. Create issues for blockers they discover
3. Update issue status as they work
4. Coordinate via Linear comments

### Agent 1: Backend Verification (`dl-backend`)

**Focus:** Verify API is production-ready

```
Tasks:
1. Run: curl https://driftlock.net/api/v1/healthz
2. Run: curl https://driftlock.net/api/v1/readyz
3. Test demo endpoint:
   curl -X POST https://driftlock.net/api/v1/demo/detect \
     -H "Content-Type: application/json" \
     -d '{"events":[{"message":"test event","timestamp":"2025-01-01T00:00:00Z"}],"window_size":50}'
4. Verify Gemini AI integration works (check logs for "AI client initialized")
5. Check Cloud Run logs for errors: gcloud logging read "resource.type=cloud_run_revision" --limit=50

Create Linear issues for any failures found.
```

### Agent 2: Frontend Verification (`dl-frontend`)

**Focus:** Verify landing page and dashboard

```
Tasks:
1. Visit https://driftlock.net - verify page loads
2. Check signup form renders and validates
3. Verify playground/demo mode works
4. Check pricing page links to Stripe checkout
5. Verify mobile responsiveness
6. Check browser console for JS errors

Create Linear issues for any UI bugs.
```

### Agent 3: Database & Migrations (`dl-db`)

**Focus:** Ensure database is ready

```
Tasks:
1. Verify firebase_uid column exists:
   SELECT column_name FROM information_schema.columns 
   WHERE table_name = 'tenants' AND column_name = 'firebase_uid';

2. Check migration status:
   goose -dir api/migrations postgres "$DATABASE_URL" status

3. Verify Stripe price IDs are configured (env vars):
   - STRIPE_PRICE_ID_STARTER
   - STRIPE_PRICE_ID_PRO
   - STRIPE_PRICE_ID_TEAM
   - STRIPE_PRICE_ID_SCALE
   - (optional) STRIPE_PRICE_ID_ENTERPRISE

4. Check tenant table has proper indexes

Create Linear issues for any schema problems.
```

### Agent 4: DevOps & Deployment (`dl-devops`)

**Focus:** Verify infrastructure

```
Tasks:
1. Check Cloud Build status:
   gcloud builds list --limit=5 --project=driftlock

2. Verify Cloud Run service is healthy:
   gcloud run services describe driftlock-api --region=us-central1

3. Check Firebase Hosting deployment:
   firebase hosting:channel:list

4. Verify environment variables are set in Cloud Run:
   - DATABASE_URL
   - GEMINI_API_KEY (NEW - just added)
   - STRIPE_SECRET_KEY
   - SENDGRID_API_KEY

5. If GEMINI_API_KEY not set, add it:
   gcloud run services update driftlock-api \
     --set-env-vars="GEMINI_API_KEY=YOUR_GEMINI_API_KEY,GEMINI_MODEL=gemini-2.5-flash"

Create Linear issues for any infra problems.
```

### Agent 5: E2E Testing (`dl-testing`)

**Focus:** Run full user journey tests

```
Tasks:
1. Run E2E tests:
   cd collector-processor && go test ./cmd/driftlock-http/... -v -run TestE2E

2. Test signup flow manually:
   - POST /api/v1/onboard/signup
   - Verify email sent
   - GET /api/v1/onboard/verify?token=...
   - Verify API key returned

3. Test detection flow:
   - POST /api/v1/detect with API key
   - Verify anomaly detection works
   - Check AI analysis is included in response

4. Test Stripe webhook:
   stripe listen --forward-to https://driftlock.net/api/v1/billing/webhook

Create Linear issues for any test failures.
```

### Agent 6: Documentation (`dl-docs`)

**Focus:** Final docs review

```
Tasks:
1. Update docs/LAUNCH_STATUS.md with current status
2. Verify README.md has correct URLs
3. Check API docs match actual endpoints
4. Update DEPLOYMENT_CHECKLIST.md
5. Add Gemini API key to environment variable docs

Create Linear issues for any doc gaps.
```

---

## Coordination Protocol

1. **Start:** Each agent creates a Linear issue: "Agent X: Starting launch verification"
2. **Progress:** Update issue with findings every 10 minutes
3. **Blockers:** Create new issues with `blocker` label, assign to relevant agent
4. **Complete:** Mark issue as Done, summarize in comment

### Linear Commands

```
# Check all launch issues
"Show me all issues in driftlock project with label 'launch'"

# Create blocker
"Create issue: [BLOCKER] API healthz returning 500 - needs immediate fix"

# Update status
"Update DRI-XXX: Completed frontend verification, no issues found"
```

---

## Environment Variables Needed

```bash
# Already configured
DATABASE_URL=postgres://...
STRIPE_SECRET_KEY=sk_live_...
SENDGRID_API_KEY=SG....

# NEW - Add to Cloud Run
GEMINI_API_KEY=YOUR_GEMINI_API_KEY
GEMINI_MODEL=gemini-2.5-flash
AI_PROVIDER=gemini
```

---

## Success Criteria

Launch is complete when:

1. [ ] `curl https://driftlock.net/api/v1/healthz` returns 200
2. [ ] Demo detection endpoint works without auth
3. [ ] Signup → Verify → API Key flow works
4. [ ] Stripe checkout creates subscription
5. [ ] AI analysis appears in detection responses
6. [ ] No critical errors in Cloud Run logs
7. [ ] All Linear launch issues marked Done

---

## Quick Commands

```bash
# Start local dev
docker compose up -d
cd collector-processor && go run ./cmd/driftlock-http

# Deploy
gcloud builds submit --config cloudbuild.yaml .
firebase deploy --only hosting

# Test production
curl https://driftlock.net/api/v1/healthz
curl -X POST https://driftlock.net/api/v1/demo/detect \
  -H "Content-Type: application/json" \
  -d '{"events":[{"message":"test"}],"window_size":50}'
```

---

## Files to Know

| Purpose | Path |
|---------|------|
| Main API server | `collector-processor/cmd/driftlock-http/main.go` |
| AI clients | `collector-processor/internal/ai/` |
| Frontend | `landing-page/src/` |
| Migrations | `api/migrations/` |
| Deploy config | `cloudbuild.yaml`, `deploy/docker-compose.yml` |
| Launch status | `docs/LAUNCH_STATUS.md` |

---

**GO SHIP IT! 🚀**

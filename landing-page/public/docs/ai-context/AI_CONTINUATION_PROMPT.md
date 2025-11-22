# 🚀 Prompt for Next AI Session (Claude/Cursor)

## Context
You're continuing work on Driftlock's SaaS launch. The backend is **95% complete and deployed to Cloud Run**. Your main job is to **deploy the frontend to Cloudflare Pages** and test the complete signup flow.

## Current Branch
```bash
git checkout saas-launch
```

## What's Already Done ✅
- Backend API deployed to Cloud Run with:
  - `/v1/onboard/signup` - User signup with rate limiting
  - `/v1/onboard/verify` - Email verification flow
  - `/v1/detect` - Core anomaly detection with usage tracking
- SendGrid email service integrated
- Database schema extended for onboarding
- Signup form component created in Vue
- Cloud Build CI/CD pipeline working

## Your Mission 🎯

### Primary Goal (30-45 minutes)
**Deploy the landing page to Cloudflare Pages and test the complete signup flow.**

### Step-by-Step Instructions

#### 1. Get Cloud Run API URL (2 min)
```bash
gcloud run services describe driftlock-api --region=us-central1 --format='value(status.url)'
# Example output: https://driftlock-api-xxxxx-uc.a.run.app
```

#### 2. Build Frontend (5 min)
```bash
cd landing-page
npm install
npm run build
```

#### 3. Deploy to Cloudflare Pages (10 min)
```bash
# Install wrangler if needed
npm install -g wrangler

# Login to Cloudflare (if not already)
wrangler login

# Deploy
wrangler pages deploy dist --project-name=driftlock
```

#### 4. Configure Cloudflare (10 min)
In Cloudflare Dashboard → Pages → driftlock:

**Environment Variables:**
- Add `VITE_API_BASE_URL` = `https://[YOUR-CLOUD-RUN-URL]/api/v1`

**Custom Domain:**
- Add `driftlock.net` 
- Wait for SSL (5-10 min)

#### 5. Test Signup Flow (15 min)
1. Visit `https://driftlock.net`
2. Fill out signup form
3. Check email for verification link
4. Click verification link
5. Check email for API key
6. Test API key:
```bash
curl -X POST [YOUR-CLOUD-RUN-URL]/v1/detect \
  -H "X-Api-Key: dlk_[YOUR-KEY]" \
  -H "Content-Type: application/json" \
  -d '{"events":[{"test":"data"}]}'
```

## If You Have Extra Time 🎨

### Option A: Admin Dashboard (2-3 hours)
Create `landing-page/src/views/AdminDashboard.vue`:
- List all tenants
- View usage metrics
- Search/filter capabilities
- Add to router at `/admin`

Add admin API endpoints in Go:
```go
// collector-processor/cmd/driftlock-http/admin.go
- GET /v1/admin/tenants
- GET /v1/admin/tenants/:id/usage
- PATCH /v1/admin/tenants/:id/status
```

### Option B: Testing & Documentation (1-2 hours)
- Run `./scripts/test-deployment.sh` and fix any issues
- Update `README.md` with getting started guide
- Test on mobile devices
- Check CORS configuration

### Option C: Monitoring Setup (1 hour)
```bash
# Create uptime check
gcloud monitoring uptime-checks create https://driftlock.net/api/v1/healthz \
  --check-interval=5m \
  --display-name="Driftlock API Health"

# Set up alerts for errors/latency
```

## Important Files to Reference 📚

- **`HANDOFF.md`** - Comprehensive status document with all details
- **`full-launch.plan.md`** - Original launch plan (Phases 1-10)
- **`collector-processor/cmd/driftlock-http/main.go`** - API routes
- **`landing-page/src/views/HomeView.vue`** - Landing page with signup form

## Quick Reference Commands 🔧

```bash
# Check API health
curl $(gcloud run services describe driftlock-api --region=us-central1 --format='value(status.url)')/healthz

# View logs
gcloud run services logs read driftlock-api --region=us-central1 --limit=50

# Rebuild API (if needed)
SHORT_SHA=$(git rev-parse --short HEAD) && \
gcloud builds submit --config=cloudbuild.yaml --substitutions=SHORT_SHA=$SHORT_SHA

# Check database
gcloud sql connect driftlock-db --user=driftlock_user
# Password: hunter will provide if needed
```

## Known Issues ⚠️

1. **Cloud Run NOT Publicly Accessible**
   - Org policy blocks `allUsers` access
   - API works, but needs auth token or policy change
   - Use `gcloud run services invoke` for testing

2. **SendGrid Sender Verification**
   - Make sure sender email is verified in SendGrid dashboard
   - Free tier: 100 emails/day

3. **CORS**
   - Already configured for `https://driftlock.net`
   - Uses URL-encoded commas in env var

## Success Criteria ✨

When you're done, these should all work:
- [ ] `https://driftlock.net` loads successfully
- [ ] Signup form submits without errors
- [ ] Verification email arrives
- [ ] Clicking verification link works
- [ ] Welcome email with API key arrives
- [ ] API key works with `/v1/detect` endpoint

## Where to Get Help 🆘

- **Full status:** Read `HANDOFF.md`
- **Launch plan:** Check `full-launch.plan.md`
- **API docs:** See `api/onboarding/ONBOARDING_API.md`
- **Database schema:** Check `api/migrations/20250302000000_onboarding.sql`

## Commit Your Work 💾

When done:
```bash
git add -A
git commit -m "deploy: Frontend to Cloudflare Pages + E2E testing

- Deployed landing page to Cloudflare Pages
- Configured custom domain driftlock.net
- Tested complete signup and verification flow
- Verified API key issuance and usage
[Add any additional changes you made]"

git push origin saas-launch
```

Then create a PR from `saas-launch` → `main` on GitHub.

---

**Good luck! The backend is solid, you've got this! 🚀**

**Estimated time: 45-60 minutes for primary goal**
**Stretch goals: Additional 1-3 hours**


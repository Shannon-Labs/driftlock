# Complete Deployment Plan: Driftlock Production Setup

**Last Updated:** 2025-01-27  
**Status:** Ready for Production Deployment  
**Estimated Time:** 2-3 hours (excluding DNS propagation)

This document provides a **step-by-step guide** to deploy Driftlock to production, covering:
- Supabase database setup
- Google Cloud Run API deployment  
- Cloudflare Pages frontend deployment
- Firebase Hosting (alternative)
- End-to-end connectivity
- Domain configuration

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    User's Browser                           │
└──────────────────────┬──────────────────────────────────────┘
                       │ HTTPS
                       ▼
┌─────────────────────────────────────────────────────────────┐
│         Cloudflare Pages (driftlock.net)                    │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Vue.js Landing Page (Static)                        │   │
│  │  - HomeView.vue                                      │   │
│  │  - PlaygroundShell.vue                               │   │
│  └──────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Cloudflare Pages Functions                          │   │
│  │  - /api/v1/[[path]].ts → Proxies to Cloud Run       │   │
│  │  - /api/v1/contact.ts → Contact form handler         │   │
│  └──────────────────────────────────────────────────────┘   │
└──────────────────────┬──────────────────────────────────────┘
                       │ HTTPS (via API_BACKEND_URL)
                       ▼
┌─────────────────────────────────────────────────────────────┐
│      Google Cloud Run (us-central1)                          │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  driftlock-api Service                               │   │
│  │  - Go HTTP API (driftlock-http)                      │   │
│  │  - Endpoints: /v1/detect, /v1/anomalies, /healthz    │   │
│  │  - Port: 8080                                        │   │
│  └──────────────────────────────────────────────────────┘   │
└──────────────────────┬──────────────────────────────────────┘
                       │ PostgreSQL (SSL)
                       ▼
┌─────────────────────────────────────────────────────────────┐
│              Supabase PostgreSQL                             │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Database: postgres                                  │   │
│  │  - Tables: tenants, streams, anomalies, etc.        │   │
│  │  - Connection: Pooler (port 6543)                   │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

---

## Prerequisites Checklist

Before starting, ensure you have:

- [ ] Google Cloud Project with billing enabled
- [ ] Supabase account (free tier OK for dev)
- [ ] Cloudflare account (free tier OK)
- [ ] Domain name (`driftlock.net`) registered
- [ ] `gcloud` CLI installed and authenticated
- [ ] `npm` and Node.js 18+ installed
- [ ] Docker installed (for local testing)

---

## Phase 1: Supabase Database Setup (15 minutes)

### Step 1.1: Create Supabase Project

1. Go to https://supabase.com and sign in
2. Click "New Project"
3. Fill in:
   - **Name**: `driftlock-production`
   - **Database Password**: Generate a strong password (save it!)
   - **Region**: Choose closest to `us-central1` (e.g., `us-east-1` or `us-west-1`)
   - **Pricing Plan**: Free tier is fine for now
4. Wait 2-3 minutes for provisioning

### Step 1.2: Get Connection Strings

1. In Supabase Dashboard → **Settings** → **Database**
2. Scroll to **Connection string** section
3. Copy **Connection pooling** URL (port 6543) - this is what Cloud Run needs:
   ```
   postgresql://postgres.[PROJECT-REF]:[YOUR-PASSWORD]@aws-0-[REGION].pooler.supabase.com:6543/postgres
   ```
4. **Important**: Add `?sslmode=require` to the end:
   ```
   postgresql://postgres.[PROJECT-REF]:[YOUR-PASSWORD]@aws-0-[REGION].pooler.supabase.com:6543/postgres?sslmode=require
   ```
5. Save this URL - you'll need it for Google Secret Manager

### Step 1.3: Run Database Migrations

**Option A: Using Supabase SQL Editor (Easiest)**

1. Go to Supabase Dashboard → **SQL Editor**
2. Click **New Query**
3. Open `api/migrations/20250301120000_initial_schema.sql` from your repo
4. Copy entire contents
5. Paste into SQL Editor
6. Click **Run** (or press Cmd/Ctrl+Enter)
7. Verify success - you should see "Success. No rows returned"

**Option B: Using psql (For Verification)**

```bash
# Get direct connection string (port 5432, not pooler)
# From Supabase Dashboard → Settings → Database → Connection string → URI

export DATABASE_URL="postgresql://postgres:[PASSWORD]@db.[PROJECT-REF].supabase.co:5432/postgres"

# Run migrations using your Go binary
cd /path/to/driftlock
./bin/driftlock-http migrate up
```

**Verify Tables Created:**

```sql
-- Run in Supabase SQL Editor
SELECT table_name 
FROM information_schema.tables 
WHERE table_schema = 'public' 
ORDER BY table_name;

-- Should show:
-- api_keys
-- anomaly_evidence
-- anomalies
-- export_jobs
-- ingest_batches
-- stream_configs
-- streams
-- tenants
```

---

## Phase 2: Google Cloud Setup (30 minutes)

### Step 2.1: Initialize GCP Project

```bash
# Set your project ID
export PROJECT_ID="your-gcp-project-id"
gcloud config set project $PROJECT_ID

# Enable required APIs
gcloud services enable \
  cloudbuild.googleapis.com \
  run.googleapis.com \
  artifactregistry.googleapis.com \
  secretmanager.googleapis.com \
  cloudresourcemanager.googleapis.com
```

### Step 2.2: Create Secrets in Secret Manager

```bash
# Store Supabase database URL
echo -n "postgresql://postgres.[PROJECT-REF]:[PASSWORD]@aws-0-[REGION].pooler.supabase.com:6543/postgres?sslmode=require" | \
  gcloud secrets create driftlock-db-url --data-file=-

# Store license key (or use dev mode for testing)
echo -n "your-license-key-here" | \
  gcloud secrets create driftlock-license-key --data-file=-

# Grant Cloud Run service account access
PROJECT_NUMBER=$(gcloud projects describe $PROJECT_ID --format="value(projectNumber)")
SERVICE_ACCOUNT="${PROJECT_NUMBER}-compute@developer.gserviceaccount.com"

gcloud secrets add-iam-policy-binding driftlock-db-url \
  --member="serviceAccount:${SERVICE_ACCOUNT}" \
  --role="roles/secretmanager.secretAccessor"

gcloud secrets add-iam-policy-binding driftlock-license-key \
  --member="serviceAccount:${SERVICE_ACCOUNT}" \
  --role="roles/secretmanager.secretAccessor"
```

### Step 2.3: Update Configuration Files

**Update `service.yaml`:**

```bash
# Replace PROJECT_ID placeholder
sed -i.bak "s/PROJECT_ID/$PROJECT_ID/g" service.yaml
```

**Update `cloudbuild.yaml`:**

The file already uses `$PROJECT_ID` substitution variable, which Cloud Build will automatically replace.

### Step 2.4: Deploy API to Cloud Run

**Option A: Using Cloud Build (Recommended for CI/CD)**

```bash
# Submit build (this builds and deploys automatically)
gcloud builds submit --config=cloudbuild.yaml

# Monitor build progress
gcloud builds list --limit=1
```

**Option B: Manual Docker Build**

```bash
# Build Docker image
docker build \
  -f collector-processor/cmd/driftlock-http/Dockerfile \
  -t gcr.io/$PROJECT_ID/driftlock-api:latest \
  --build-arg USE_OPENZL=false \
  .

# Push to Container Registry
docker push gcr.io/$PROJECT_ID/driftlock-api:latest

# Deploy to Cloud Run
gcloud run deploy driftlock-api \
  --image gcr.io/$PROJECT_ID/driftlock-api:latest \
  --region us-central1 \
  --platform managed \
  --allow-unauthenticated \
  --set-env-vars PORT=8080,CORS_ALLOW_ORIGINS=https://driftlock.net,https://www.driftlock.net,LOG_LEVEL=info \
  --set-secrets DATABASE_URL=driftlock-db-url:latest,DRIFTLOCK_LICENSE_KEY=driftlock-license-key:latest \
  --memory 2Gi \
  --cpu 2 \
  --min-instances 1 \
  --max-instances 10 \
  --timeout 300 \
  --concurrency 80
```

### Step 2.5: Get Cloud Run Service URL

```bash
# Get the service URL
SERVICE_URL=$(gcloud run services describe driftlock-api \
  --region us-central1 \
  --format 'value(status.url)')

echo "API URL: $SERVICE_URL"
# Should look like: https://driftlock-api-xxxxx-uc.a.run.app
```

**Save this URL** - you'll need it for Cloudflare Pages configuration.

### Step 2.6: Verify API Deployment

```bash
# Test health endpoint
curl $SERVICE_URL/healthz | jq

# Should return:
# {
#   "success": true,
#   "library_status": "healthy",
#   "database": "connected",
#   "license": { ... }
# }
```

### Step 2.7: Create Initial Tenant (Optional)

```bash
# Create a Cloud Run Job for migrations (if not done via Supabase SQL Editor)
gcloud run jobs create driftlock-migrate \
  --image gcr.io/$PROJECT_ID/driftlock-api:latest \
  --region us-central1 \
  --set-secrets DATABASE_URL=driftlock-db-url:latest \
  --command /usr/local/bin/driftlock-http \
  --args migrate,up \
  --max-retries 1

# Execute migration job
gcloud run jobs execute driftlock-migrate --region us-central1

# Create initial tenant
gcloud run jobs create driftlock-create-tenant \
  --image gcr.io/$PROJECT_ID/driftlock-api:latest \
  --region us-central1 \
  --set-secrets DATABASE_URL=driftlock-db-url:latest,DRIFTLOCK_LICENSE_KEY=driftlock-license-key:latest \
  --command /usr/local/bin/driftlock-http \
  --args create-tenant,--name,Production,--key-role,admin,--json \
  --max-retries 1

# Execute tenant creation
gcloud run jobs execute driftlock-create-tenant --region us-central1
```

---

## Phase 3: Cloudflare Pages Frontend Setup (20 minutes)

### Step 3.1: Install Wrangler CLI

```bash
npm install -g wrangler
wrangler login
```

### Step 3.2: Configure Environment Variables

In Cloudflare Dashboard → Pages → Your Project → Settings → Environment Variables:

**Production Environment:**
- `API_BACKEND_URL`: `https://driftlock-api-xxxxx-uc.a.run.app` (your Cloud Run URL)
- `VITE_API_BASE_URL`: `https://driftlock.net/api/v1` (will proxy through Pages Functions)
- `CRM_WEBHOOK_URL`: (optional) Your CRM webhook URL for contact form submissions

**Preview Environment:**
- `API_BACKEND_URL`: `https://driftlock-api-xxxxx-uc.a.run.app` (same as production)
- `VITE_API_BASE_URL`: `https://preview.driftlock.net/api/v1` (or your preview domain)

### Step 3.3: Deploy Landing Page

**Option A: Connect GitHub Repository (Recommended)**

1. Go to Cloudflare Dashboard → Pages → **Create a project**
2. Connect your GitHub repository
3. Select repository: `Shannon-Labs/driftlock`
4. Configure build:
   - **Framework preset**: Vite
   - **Build command**: `cd landing-page && npm install && npm run build`
   - **Build output directory**: `landing-page/dist`
   - **Root directory**: `/landing-page`
5. Add environment variables (from Step 3.2)
6. Click **Save and Deploy**

**Option B: Manual Deploy via Wrangler**

```bash
cd landing-page

# Install dependencies
npm install

# Build
npm run build

# Deploy
wrangler pages deploy dist --project-name=driftlock
```

### Step 3.4: Configure Custom Domain

1. In Cloudflare Pages → Your Project → **Custom domains**
2. Click **Set up a custom domain**
3. Enter: `driftlock.net`
4. Cloudflare will automatically configure DNS (if domain is managed by Cloudflare)
5. Wait for SSL certificate provisioning (1-5 minutes)

**If domain is NOT managed by Cloudflare:**

At your domain registrar, add:
- **CNAME**: `driftlock.net` → `driftlock.pages.dev`
- **CNAME**: `www.driftlock.net` → `driftlock.net`

### Step 3.5: Verify Pages Functions

The landing page includes Cloudflare Pages Functions:
- `functions/api/v1/[[path]].ts` - Proxies API requests to Cloud Run
- `functions/api/v1/contact.ts` - Handles contact form submissions

These should work automatically once deployed. Test:

```bash
# Test API proxy
curl https://driftlock.net/api/v1/healthz

# Should return the same as direct Cloud Run call
```

---

## Phase 4: Firebase Hosting (Alternative to Cloudflare)

If you prefer Firebase Hosting over Cloudflare Pages:

### Step 4.1: Install Firebase CLI

```bash
npm install -g firebase-tools
firebase login
```

### Step 4.2: Initialize Firebase Project

```bash
cd landing-page
firebase init hosting

# Select:
# - Use existing project (or create new)
# - Public directory: dist
# - Single-page app: Yes
# - Set up automatic builds: No (we'll deploy manually)
```

### Step 4.3: Configure firebase.json

Create/update `firebase.json`:

```json
{
  "hosting": {
    "public": "dist",
    "ignore": [
      "firebase.json",
      "**/.*",
      "**/node_modules/**"
    ],
    "rewrites": [
      {
        "source": "/api/**",
        "function": "api"
      },
      {
        "source": "**",
        "destination": "/index.html"
      }
    ],
    "headers": [
      {
        "source": "/assets/**",
        "headers": [
          {
            "key": "Cache-Control",
            "value": "public, max-age=31536000, immutable"
          }
        ]
      }
    ]
  },
  "functions": [
    {
      "source": "functions",
      "codebase": "default",
      "runtime": "nodejs20"
    }
  ]
}
```

### Step 4.4: Deploy to Firebase

```bash
# Build
npm run build

# Deploy
firebase deploy --only hosting,functions
```

**Note**: Firebase Functions would need to be rewritten to proxy to Cloud Run (similar to Cloudflare Pages Functions).

---

## Phase 5: End-to-End Connectivity (15 minutes)

### Step 5.1: Update CORS Configuration

Ensure Cloud Run API allows requests from your domain:

```bash
# Update Cloud Run service CORS settings
gcloud run services update driftlock-api \
  --region us-central1 \
  --set-env-vars CORS_ALLOW_ORIGINS=https://driftlock.net,https://www.driftlock.net
```

### Step 5.2: Test API Connectivity

```bash
# Test direct Cloud Run access
curl https://driftlock-api-xxxxx-uc.a.run.app/healthz | jq

# Test via Cloudflare Pages proxy
curl https://driftlock.net/api/v1/healthz | jq

# Both should return the same response
```

### Step 5.3: Test Playground

1. Visit `https://driftlock.net`
2. Scroll to playground section (or visit `https://driftlock.net/playground`)
3. Click "Check API Status" - should show "Connected"
4. Upload sample data and test `/v1/detect`

### Step 5.4: Test Contact Form

1. Fill out contact form on landing page
2. Submit
3. Check Cloudflare Pages Function logs:
   - Dashboard → Pages → Your Project → Functions → Logs
4. Verify webhook fires (if `CRM_WEBHOOK_URL` is set)

---

## Phase 6: Domain & DNS Configuration (10 minutes)

### Step 6.1: Cloudflare DNS Setup

If using Cloudflare for DNS:

1. Go to Cloudflare Dashboard → **DNS** → Your Domain
2. Ensure these records exist:
   - **A** or **CNAME**: `driftlock.net` → Cloudflare Pages IP or CNAME
   - **CNAME**: `www.driftlock.net` → `driftlock.net`

### Step 6.2: SSL Certificate

Cloudflare automatically provisions SSL certificates:
- **Universal SSL**: Free, auto-renewing
- Usually active within 5 minutes of domain configuration

Verify:
```bash
curl -I https://driftlock.net
# Should show HTTP/2 200
```

### Step 6.3: Redirect www to Non-www (Optional)

In Cloudflare Dashboard → **Page Rules**:
- **URL**: `www.driftlock.net/*`
- **Setting**: Forwarding URL → 301 Permanent Redirect → `https://driftlock.net/$1`

---

## Phase 7: Monitoring & Verification (10 minutes)

### Step 7.1: Set Up Monitoring

**Google Cloud Monitoring:**

```bash
# View Cloud Run metrics
gcloud monitoring dashboards list

# View logs
gcloud logging read "resource.type=cloud_run_revision AND resource.labels.service_name=driftlock-api" \
  --limit 50 \
  --format json
```

**Cloudflare Analytics:**

- Dashboard → Pages → Your Project → Analytics
- View page views, bandwidth, requests

### Step 7.2: Health Check Endpoints

Create monitoring alerts for:

1. **API Health**: `https://driftlock-api-xxxxx-uc.a.run.app/healthz`
2. **Frontend**: `https://driftlock.net`
3. **API Proxy**: `https://driftlock.net/api/v1/healthz`

### Step 7.3: Performance Testing

```bash
# Test API latency
time curl -s https://driftlock.net/api/v1/healthz

# Test page load
lighthouse https://driftlock.net --view
```

---

## Phase 8: Production Checklist

Before going live, verify:

### Infrastructure
- [ ] Supabase database is running and accessible
- [ ] Cloud Run service is deployed and healthy
- [ ] Cloudflare Pages is deployed
- [ ] Custom domain is configured
- [ ] SSL certificate is active
- [ ] DNS records are correct

### Functionality
- [ ] `/healthz` endpoint returns `success: true`
- [ ] Database connection works (`database: "connected"`)
- [ ] Landing page loads correctly
- [ ] Playground can connect to API
- [ ] `/v1/detect` endpoint works with sample data
- [ ] Contact form submits successfully
- [ ] CORS headers are correct (no browser errors)

### Security
- [ ] Secrets are stored in Secret Manager (not in code)
- [ ] License key is configured (or dev mode is intentional)
- [ ] CORS only allows your domains
- [ ] API requires authentication (`X-Api-Key` header)
- [ ] SSL/TLS is enforced everywhere

### Performance
- [ ] Page load time < 3 seconds
- [ ] API response time < 500ms for `/healthz`
- [ ] API response time < 5s for `/v1/detect` with sample data
- [ ] Cloud Run auto-scaling is configured

---

## Troubleshooting Guide

### API Not Responding

**Check Cloud Run logs:**
```bash
gcloud run services logs read driftlock-api --region us-central1 --limit 50
```

**Common issues:**
- Database connection string incorrect → Check Secret Manager
- License key invalid → Check `/healthz` response
- CORS errors → Verify `CORS_ALLOW_ORIGINS` env var

### Database Connection Issues

**Test connection:**
```bash
# Get connection string from Secret Manager
gcloud secrets versions access latest --secret=driftlock-db-url

# Test with psql
psql "postgresql://postgres.[PROJECT-REF]:[PASSWORD]@aws-0-[REGION].pooler.supabase.com:6543/postgres?sslmode=require"
```

**Common issues:**
- Using direct connection (5432) instead of pooler (6543) → Use pooler for Cloud Run
- Missing `sslmode=require` → Supabase requires SSL
- Wrong region → Ensure pooler region matches Supabase project region

### Frontend Not Loading

**Check Cloudflare Pages logs:**
- Dashboard → Pages → Your Project → Functions → Logs

**Common issues:**
- Build failed → Check build logs in Cloudflare Dashboard
- Environment variables not set → Verify in Pages Settings
- API proxy not working → Check `API_BACKEND_URL` is set correctly

### CORS Errors

**Verify CORS configuration:**
```bash
# Check Cloud Run CORS env var
gcloud run services describe driftlock-api \
  --region us-central1 \
  --format="value(spec.template.spec.containers[0].env)"

# Should include: CORS_ALLOW_ORIGINS=https://driftlock.net,https://www.driftlock.net
```

**Test CORS:**
```bash
curl -X OPTIONS https://driftlock.net/api/v1/detect \
  -H "Origin: https://driftlock.net" \
  -H "Access-Control-Request-Method: POST" \
  -v
```

---

## Cost Estimates

### Supabase
- **Free Tier**: 500MB database, 2GB bandwidth - **$0/month**
- **Pro Tier**: 8GB database, 50GB bandwidth - **$25/month** (recommended for production)

### Google Cloud Run
- **Compute**: ~$0.00002400 per vCPU-second, ~$0.00000250 per GiB-second
- **Requests**: First 2 million requests free, then $0.40 per million
- **Estimated**: **$10-50/month** depending on traffic (with 1 min instance always warm)

### Cloudflare Pages
- **Free Tier**: Unlimited requests, 500 builds/month - **$0/month**
- **Pro Tier**: Advanced features - **$20/month** (optional)

### Total Estimated Cost
- **Development/Testing**: **$0-10/month** (Supabase free + Cloud Run minimal usage)
- **Production (Low Traffic)**: **$35-75/month** (Supabase Pro + Cloud Run with 1 min instance)
- **Production (Medium Traffic)**: **$100-200/month** (with scaling)

---

## Next Steps After Deployment

1. **Set up CI/CD**: Connect GitHub to Cloud Build for automatic deployments
2. **Configure alerts**: Set up Cloud Monitoring alerts for errors and latency
3. **Create backup strategy**: Configure Supabase backups
4. **Set up analytics**: Add Google Analytics or similar
5. **Document API**: Ensure API docs are up to date
6. **Load testing**: Test with expected production load
7. **Security audit**: Review security settings and access controls

---

## Rollback Plan

If something goes wrong:

1. **Revert Cloud Run**: Deploy previous image version
   ```bash
   gcloud run services update driftlock-api \
     --image gcr.io/$PROJECT_ID/driftlock-api:[PREVIOUS_SHA] \
     --region us-central1
   ```

2. **Revert Frontend**: Rollback Cloudflare Pages deployment
   - Dashboard → Pages → Your Project → Deployments → Rollback

3. **Disable Custom Domain**: Remove custom domain temporarily
   - Keep `driftlock.pages.dev` active as fallback

4. **Check Logs**: Review Cloud Run and Cloudflare logs for errors

---

## Support & Resources

- **Supabase Docs**: https://supabase.com/docs
- **Cloud Run Docs**: https://cloud.google.com/run/docs
- **Cloudflare Pages Docs**: https://developers.cloudflare.com/pages
- **Project Issues**: GitHub Issues in `Shannon-Labs/driftlock`

---

**This plan is comprehensive and step-by-step. Follow it sequentially, and you'll have a fully functional production deployment.**


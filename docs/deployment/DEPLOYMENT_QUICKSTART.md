# Deployment Quick Start

**For the complete step-by-step guide, see [COMPLETE_DEPLOYMENT_PLAN.md](./COMPLETE_DEPLOYMENT_PLAN.md)**

## TL;DR - Get Online in 30 Minutes

### 1. Supabase (5 min)
```bash
# Create project at https://supabase.com
# Get pooler URL: postgresql://postgres.[REF]:[PASS]@aws-0-[REGION].pooler.supabase.com:6543/postgres?sslmode=require
# Run migrations via SQL Editor (copy/paste api/migrations/20250301120000_initial_schema.sql)
```

### 2. Google Cloud Run (15 min)
```bash
export PROJECT_ID="your-project-id"
gcloud config set project $PROJECT_ID

# Create secrets
echo -n "YOUR_SUPABASE_URL" | gcloud secrets create driftlock-db-url --data-file=-
echo -n "YOUR_LICENSE_KEY" | gcloud secrets create driftlock-license-key --data-file=-

# Deploy
gcloud builds submit --config=cloudbuild.yaml

# Get URL
gcloud run services describe driftlock-api --region us-central1 --format 'value(status.url)'
```

### 3. Cloudflare Pages (10 min)
```bash
cd landing-page
npm install
npm run build

# Deploy via Cloudflare Dashboard or:
wrangler pages deploy dist --project-name=driftlock

# Set environment variable in Cloudflare Dashboard:
# API_BACKEND_URL = https://your-cloud-run-url.a.run.app
```

### 4. Domain (5 min)
- Cloudflare Dashboard → Pages → Custom domains → Add `driftlock.net`
- Wait for SSL (1-5 min)
- Done!

## Architecture

```
User → Cloudflare Pages (driftlock.net) 
     → Pages Functions (/api/v1/*) 
     → Cloud Run API (us-central1)
     → Supabase PostgreSQL
```

## Key Files

- **Complete Guide**: `docs/COMPLETE_DEPLOYMENT_PLAN.md`
- **Cloud Run Config**: `service.yaml`, `cloudbuild.yaml`
- **Supabase Setup**: `docs/deployment/supabase-setup.md`
- **Database Schema**: `api/migrations/20250301120000_initial_schema.sql`

## Verification

```bash
# Test API
curl https://your-cloud-run-url.a.run.app/healthz

# Test via proxy
curl https://driftlock.net/api/v1/healthz

# Both should return: {"success": true, "database": "connected", ...}
```


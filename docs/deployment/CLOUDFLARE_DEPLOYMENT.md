# Cloudflare Deployment Quick Reference

This is a quick reference for deploying Driftlock to Cloudflare. For detailed instructions, see [docs/deployment/CLOUDFLARE_MIGRATION.md](docs/deployment/CLOUDFLARE_MIGRATION.md).

## Quick Deploy

```bash
# One-command deployment (Workers + Pages)
./scripts/deploy-cloudflare.sh
```

## Manual Deployment

### 1. Deploy Workers

```bash
cd workers
npm install
npm run deploy
```

### 2. Deploy Pages

```bash
cd landing-page
npm install
npm run build
npm run deploy:cloudflare
```

## Required Environment Variables

### Workers (set via `wrangler secret put`)

```bash
wrangler secret put GEMINI_API_KEY
wrangler secret put VITE_FIREBASE_API_KEY
wrangler secret put VITE_FIREBASE_AUTH_DOMAIN
wrangler secret put VITE_FIREBASE_PROJECT_ID
wrangler secret put VITE_FIREBASE_STORAGE_BUCKET
wrangler secret put VITE_FIREBASE_MESSAGING_SENDER_ID
wrangler secret put VITE_FIREBASE_APP_ID
wrangler secret put VITE_FIREBASE_MEASUREMENT_ID
wrangler secret put CLOUD_RUN_API_URL
```

### Pages (set in Cloudflare Dashboard)

Go to Pages → Settings → Environment Variables:

- `VITE_FIREBASE_API_KEY`
- `VITE_FIREBASE_AUTH_DOMAIN`
- `VITE_FIREBASE_PROJECT_ID`
- `VITE_FIREBASE_STORAGE_BUCKET`
- `VITE_FIREBASE_MESSAGING_SENDER_ID`
- `VITE_FIREBASE_APP_ID`
- `VITE_FIREBASE_MEASUREMENT_ID`
- `VITE_API_BASE_URL` (set to `/api` for path-based routing or `https://api.driftlock.net` for subdomain)

## Testing

```bash
# Test Workers locally
cd workers
npm run dev
curl http://localhost:8787/healthz

# Test frontend locally
cd landing-page
VITE_WORKERS_DEV_URL=http://localhost:8787 npm run dev
```

## Architecture

```
Frontend (Cloudflare Pages)
    ↓
API Routes (/api/*)
    ↓
Cloudflare Workers
    ↓
Cloud Run Backend (existing)
```

## What Changed

- ✅ Firebase Functions → Cloudflare Workers
- ✅ Firebase Hosting → Cloudflare Pages
- ✅ Firebase Auth → Still using Firebase Auth (can migrate later)
- ✅ API Proxy → Now handled by Workers

## Next Steps

1. Deploy Workers and Pages
2. Configure DNS routing
3. Test all endpoints
4. Monitor for errors
5. Decommission Firebase Functions (after verification)

See [docs/deployment/CLOUDFLARE_MIGRATION.md](docs/deployment/CLOUDFLARE_MIGRATION.md) for complete details.



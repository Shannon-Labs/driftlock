# Cloudflare Migration Guide

This document describes the migration from Firebase to Cloudflare for the Driftlock SaaS platform.

## Architecture Overview

### Before (Firebase)
- **Frontend**: Firebase Hosting (`landing-page/dist`)
- **Backend Functions**: Firebase Functions (`functions/`)
- **Auth**: Firebase Auth
- **API Proxy**: Firebase Functions proxy to Cloud Run

### After (Cloudflare)
- **Frontend**: Cloudflare Pages (`landing-page/dist`)
- **Backend Functions**: Cloudflare Workers (`workers/`)
- **Auth**: Firebase Auth (maintained for now, can migrate to Cloudflare Workers Auth later)
- **API Proxy**: Cloudflare Workers proxy to Cloud Run

## Components

### 1. Cloudflare Workers (`workers/`)

The Workers replace all Firebase Functions:

- **`/api/v1/onboard/signup`** - User signup proxy
- **`/api/analyze`** - Anomaly analysis with Gemini AI
- **`/api/compliance`** - Compliance report generation
- **`/api/v1/healthz`** - Health check endpoint
- **`/getFirebaseConfig`** - Firebase config for frontend (backward compatibility)
- **`/api/v1/**`** - General API proxy to Cloud Run backend
- **`/webhooks/stripe`** - Stripe webhook proxy

### 2. Cloudflare Pages (`landing-page/`)

The Vue.js frontend is deployed to Cloudflare Pages:
- Static files served from `dist/` directory
- API routes handled by Workers (via routing rules)
- Environment variables for Firebase config

## Setup Instructions

### Prerequisites

1. **Cloudflare Account**
   - Sign up at https://dash.cloudflare.com
   - Add your domain (e.g., `driftlock.net`)

2. **Wrangler CLI**
   ```bash
   npm install -g wrangler
   wrangler login
   ```

3. **Environment Variables**
   - Get Firebase config from Firebase Console
   - Get Gemini API key from Google AI Studio
   - Get Cloud Run API URL

### Step 1: Deploy Cloudflare Workers

```bash
cd workers

# Install dependencies
npm install

# Set environment variables (secrets)
wrangler secret put GEMINI_API_KEY
wrangler secret put VITE_FIREBASE_API_KEY
wrangler secret put VITE_FIREBASE_AUTH_DOMAIN
wrangler secret put VITE_FIREBASE_PROJECT_ID
wrangler secret put VITE_FIREBASE_STORAGE_BUCKET
wrangler secret put VITE_FIREBASE_MESSAGING_SENDER_ID
wrangler secret put VITE_FIREBASE_APP_ID
wrangler secret put VITE_FIREBASE_MEASUREMENT_ID
wrangler secret put CLOUD_RUN_API_URL

# Deploy to production
npm run deploy

# Or deploy to preview
npm run deploy:preview
```

### Step 2: Configure Workers Routes

In Cloudflare Dashboard:
1. Go to Workers & Pages → `driftlock-api`
2. Settings → Triggers → Routes
3. Add routes:
   - `api.driftlock.net/*` (subdomain)
   - `driftlock.net/api/*` (path-based)

### Step 3: Deploy Frontend to Cloudflare Pages

```bash
cd landing-page

# Build the frontend
npm run build

# Deploy to Cloudflare Pages
npm run deploy:cloudflare

# Or connect via Git (recommended):
# 1. Go to Cloudflare Dashboard → Pages
# 2. Connect Git repository
# 3. Set build command: `npm run build`
# 4. Set output directory: `dist`
# 5. Add environment variables (see below)
```

### Step 4: Configure Pages Environment Variables

In Cloudflare Dashboard → Pages → `driftlock` → Settings → Environment Variables:

```
VITE_FIREBASE_API_KEY=<your-key>
VITE_FIREBASE_AUTH_DOMAIN=driftlock.firebaseapp.com
VITE_FIREBASE_PROJECT_ID=driftlock
VITE_FIREBASE_STORAGE_BUCKET=driftlock.appspot.com
VITE_FIREBASE_MESSAGING_SENDER_ID=131489574303
VITE_FIREBASE_APP_ID=1:131489574303:web:e83e3e433912d05a8d61aa
VITE_FIREBASE_MEASUREMENT_ID=G-CXBMVS3G8H
VITE_API_BASE_URL=/api
```

### Step 5: Configure DNS

1. **Subdomain Route** (recommended):
   - Add CNAME: `api.driftlock.net` → `driftlock-api.<account>.workers.dev`
   - Update frontend `VITE_API_BASE_URL` to `https://api.driftlock.net`

2. **Path-based Route** (alternative):
   - Workers route: `driftlock.net/api/*`
   - Frontend uses relative paths: `/api/*`

### Step 6: Update Frontend API Base URL

If using subdomain routing, update `landing-page/src/components/playground/PlaygroundShell.vue`:

```typescript
const apiBase = import.meta.env.VITE_API_BASE_URL || 'https://api.driftlock.net'
```

## Local Development

### Run Workers Locally

```bash
cd workers
npm run dev
# Workers run on http://localhost:8787
```

### Run Frontend Locally

```bash
cd landing-page

# Set environment variable for local Workers
export VITE_WORKERS_DEV_URL=http://localhost:8787

npm run dev
# Frontend runs on http://localhost:5173
```

### Test Integration

1. Frontend at `http://localhost:5173`
2. Workers at `http://localhost:8787`
3. API calls from frontend should proxy through Workers to Cloud Run

## Migration Checklist

- [x] Create Cloudflare Workers to replace Firebase Functions
- [x] Update frontend to use Workers endpoints
- [x] Configure wrangler.toml for Workers and Pages
- [ ] Deploy Workers to Cloudflare
- [ ] Deploy Pages to Cloudflare
- [ ] Configure DNS routing
- [ ] Test all API endpoints
- [ ] Verify Firebase Auth still works
- [ ] Update documentation
- [ ] Monitor for errors
- [ ] Decommission Firebase Functions (after verification)

## Cost Comparison

### Firebase
- Hosting: Free tier (10GB storage, 360MB/day transfer)
- Functions: $0.40 per million invocations
- Bandwidth: $0.12/GB after free tier

### Cloudflare
- Pages: Free (unlimited requests, 500 builds/month)
- Workers: Free (100,000 requests/day), then $5/month for 10M requests
- Bandwidth: Free (unlimited)

**Estimated Savings**: ~$50-200/month depending on traffic

## Troubleshooting

### CORS Errors

Workers include CORS headers. If you see CORS errors:
1. Check that Workers are deployed and accessible
2. Verify routes are configured correctly
3. Check browser console for actual error

### API 404 Errors

1. Verify Workers routes are configured
2. Check that Workers are deployed
3. Test Workers directly: `curl https://api.driftlock.net/healthz`

### Firebase Auth Not Working

1. Verify Firebase config is correct in Workers secrets
2. Check that `/getFirebaseConfig` endpoint returns valid config
3. Verify Firebase project settings allow your domain

### Environment Variables Not Working

1. Check Cloudflare Dashboard → Workers → Settings → Variables
2. Verify variable names match exactly (case-sensitive)
3. Redeploy after adding new variables

## Rollback Plan

If issues occur, you can rollback:

1. **Keep Firebase Functions running** during migration
2. **Update DNS** to point back to Firebase Hosting
3. **Revert frontend** to use Firebase Functions endpoints
4. **Investigate issues** in Cloudflare Workers logs

## Next Steps

1. **Migrate Auth**: Consider Cloudflare Workers Auth or Auth0
2. **Add Rate Limiting**: Use Cloudflare Rate Limiting
3. **Add Analytics**: Use Cloudflare Analytics
4. **Optimize Caching**: Configure Cloudflare Cache Rules
5. **Add DDoS Protection**: Already included with Cloudflare

## Support

- Cloudflare Workers Docs: https://developers.cloudflare.com/workers/
- Cloudflare Pages Docs: https://developers.cloudflare.com/pages/
- Wrangler CLI Docs: https://developers.cloudflare.com/workers/wrangler/




# API Deployment Guide

## Current Situation

The API is **not deployed yet**. The Go API server needs to be deployed somewhere, then we can route `/api/v1` requests to it.

## Deployment Options

### Option 1: Deploy Go API + Cloudflare Pages Function Proxy (Recommended)

**Step 1: Deploy Go API to a hosting service**

Choose one:
- **Render** (easiest): https://render.com
- **Railway**: https://railway.app  
- **Fly.io**: https://fly.io
- **Google Cloud Run**: https://cloud.google.com/run

**Step 2: Configure Cloudflare Pages Function**

1. Set the backend URL in Cloudflare Pages environment variables:
   ```bash
   wrangler pages secret put API_BACKEND_URL --project-name=driftlock
   # Enter your backend URL: https://your-app.onrender.com
   ```

2. The Pages Function at `landing-page/functions/api/v1/[[path]].ts` will automatically proxy requests.

**Step 3: Deploy**

```bash
cd landing-page
npm run build
npx wrangler pages deploy dist --project-name=driftlock
```

### Option 2: Deploy Go API Directly (Simpler, but separate domain)

**Deploy to Render:**

1. Go to https://render.com
2. Create new "Web Service"
3. Connect your GitHub repo
4. Set:
   - **Build Command**: `docker build -t driftlock-http -f collector-processor/cmd/driftlock-http/Dockerfile .`
   - **Start Command**: `docker run -p 8080:8080 -e PORT=8080 -e CORS_ALLOW_ORIGINS=https://driftlock.net driftlock-http`
   - **Environment**: `CORS_ALLOW_ORIGINS=https://driftlock.net,https://www.driftlock.net`

5. Update frontend to point to Render URL:
   ```typescript
   // landing-page/src/views/PlaygroundView.vue
   const apiBase = import.meta.env.VITE_API_BASE_URL || 'https://your-app.onrender.com'
   ```

### Option 3: Deploy Go API as Cloudflare Worker (Advanced)

Cloudflare Workers now supports Go via WASM, but requires compiling Go to WASM which is complex with CGO dependencies.

## Quick Start: Render Deployment

1. **Create Render account**: https://render.com

2. **Create new Web Service**:
   - Connect GitHub repo
   - Select `collector-processor/cmd/driftlock-http/Dockerfile`
   - Set environment variables:
     ```
     PORT=8080
     CORS_ALLOW_ORIGINS=https://driftlock.net,https://www.driftlock.net
     ```

3. **Get your Render URL**: `https://driftlock-api.onrender.com`

4. **Update Cloudflare Pages Function**:
   ```bash
   wrangler pages secret put API_BACKEND_URL --project-name=driftlock
   # Enter: https://driftlock-api.onrender.com
   ```

5. **Redeploy Pages**:
   ```bash
   cd landing-page
   npm run build
   npx wrangler pages deploy dist --project-name=driftlock
   ```

## Testing

After deployment, test:
```bash
curl https://driftlock.net/api/v1/healthz
# Should return JSON, not HTML

curl https://driftlock.net/api/v1/detect \
  -X POST \
  -H "Content-Type: application/json" \
  -d '[{"amount": 100}]'
```

## Current Status

- ✅ Pages Function proxy code created: `landing-page/functions/api/v1/[[path]].ts`
- ⏳ Need to deploy Go API backend
- ⏳ Need to set `API_BACKEND_URL` environment variable
- ⏳ Need to redeploy Pages


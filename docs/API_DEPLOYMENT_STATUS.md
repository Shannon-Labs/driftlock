# API Deployment Status Investigation

## Current Situation

**Problem**: The API endpoint `https://driftlock.net/api/v1/healthz` returns HTML (the SPA index page) instead of JSON, indicating the API is not actually deployed or routed correctly.

## Findings

1. **No Cloudflare Pages Functions**: The `landing-page` directory has no `functions/` folder, so there are no Pages Functions handling `/api/v1` routes.

2. **No Cloudflare Worker**: There's no active Cloudflare Worker deployed for the API (based on wrangler commands).

3. **API Code Exists**: The Go API server code exists at `collector-processor/cmd/driftlock-http/main.go` but it's not deployed to Cloudflare.

4. **CORS Secret Set**: We set `CORS_ALLOW_ORIGINS` secret on Cloudflare Pages, but this won't help if the API isn't running on Cloudflare.

## Next Steps

### Option 1: Deploy API as Cloudflare Pages Function
Create a Pages Function to handle `/api/v1` routes:
- Create `landing-page/functions/api/v1/[[path]].ts` 
- Proxy requests to the actual API backend (wherever it's deployed)
- Or implement the API logic directly in the Function

### Option 2: Deploy API as Cloudflare Worker
- Create a new Cloudflare Worker project
- Deploy the Go API (if possible) or create a TypeScript Worker that proxies to the backend
- Configure routing: `api.driftlock.net/*` → Worker

### Option 3: Deploy API Elsewhere and Route Through Cloudflare
- Deploy Go API to Render, Cloud Run, Railway, Fly.io, etc.
- Use Cloudflare Workers/Pages Functions as a reverse proxy
- Or configure DNS to point `api.driftlock.net` directly to the API service

### Option 4: Use Cloudflare Workers with Go Runtime (Experimental)
- Cloudflare Workers now supports Go via WASM
- Could potentially compile the Go API to WASM and deploy as a Worker

## Questions to Answer

1. **Where is the API currently deployed?** (Vercel, Render, Cloud Run, etc.)
2. **Do you want to deploy it on Cloudflare?** (Workers or Pages Functions)
3. **What's the actual API URL?** (Maybe it's `api.driftlock.net` or a different subdomain?)

## Recommendation

Since you're already using Cloudflare Pages for the frontend, the cleanest solution would be:

1. **Deploy API as Cloudflare Pages Function** (easiest, same project)
   - Create `landing-page/functions/api/v1/[[path]].ts`
   - Proxy to wherever the Go API is currently deployed
   - Or implement API logic directly in TypeScript

2. **Or deploy API as separate Cloudflare Worker** (more scalable)
   - Create new Worker project
   - Route `api.driftlock.net` → Worker
   - Deploy Go API or proxy to backend

Let me know which approach you prefer and I can help implement it!


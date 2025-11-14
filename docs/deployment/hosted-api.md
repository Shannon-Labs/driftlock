## Hosted API Deployment (Render and Cloud Run)

This guide shows how to deploy the Driftlock HTTP API container with compression-based anomaly detection to a managed platform.

Images referenced below are published by the GitHub Actions workflow:
- ghcr.io/OWNER/driftlock-driftlock-http:latest
- ghcr.io/OWNER/driftlock-driftlock-http:openzl

Replace OWNER with your GitHub org or username.

### Render

1) Create a new Web Service
- Select “Deploy an existing image”
- Image URL: ghcr.io/OWNER/driftlock-driftlock-http:latest
- Runtime: Docker
- Region: closest to your users
- Port: 8080

2) Environment
- Add PORT=8080
- Optional: set LOG_LEVEL=info

3) Authentication (pilot)
- For a simple API key gate, configure a reverse proxy (e.g., Render’s native header rules or put Cloudflare in front).

4) Health check
- Path: /healthz (expects {"ok":true})

### Google Cloud Run

1) Enable Cloud Run and Artifact Registry
```bash
gcloud services enable run.googleapis.com artifactregistry.googleapis.com
```

2) Deploy the image
```bash
gcloud run deploy driftlock-http \
  --image ghcr.io/OWNER/driftlock-driftlock-http:latest \
  --region us-central1 \
  --platform managed \
  --port 8080 \
  --allow-unauthenticated
```

3) Verify
```bash
curl -s https://YOUR-SERVICE-URL/healthz
```

### Request examples

NDJSON:
```bash
curl -s -X POST "https://YOUR-SERVICE-URL/v1/detect?format=ndjson&baseline=400&window=1&algo=zstd" \
  -H "Content-Type: application/json" \
  --data-binary @test-data/normal-transactions.jsonl | jq .
```

JSON array (autodetected):
```bash
curl -s -X POST "https://YOUR-SERVICE-URL/v1/detect" \
  -H "Content-Type: application/json" \
  --data-binary @test-data/test-demo.json | jq .
```

### CORS Configuration

The API supports CORS via the `CORS_ALLOW_ORIGINS` environment variable. You can specify multiple origins separated by commas.

#### Cloudflare Workers/Pages Functions Deployment

**Option 1: Using Cloudflare Dashboard**
1. Go to Cloudflare Dashboard → Workers & Pages → Your API Project
2. Navigate to **Settings** → **Variables**
3. Add or update the variable:
   - **Variable name**: `CORS_ALLOW_ORIGINS`
   - **Value**: `https://driftlock.net,https://www.driftlock.net`
   - **Environment**: Production (and Preview if needed)
4. Save changes (Workers update immediately, Pages Functions require redeploy)

**Option 2: Using Wrangler CLI**
```bash
# Set as secret (recommended for sensitive values)
wrangler secret put CORS_ALLOW_ORIGINS
# Enter: https://driftlock.net,https://www.driftlock.net

# Or add to wrangler.toml for non-secret values
[env.production.vars]
CORS_ALLOW_ORIGINS = "https://driftlock.net,https://www.driftlock.net"
```

#### Other Platforms (Render, Cloud Run, etc.)

Set the `CORS_ALLOW_ORIGINS` environment variable to your frontend origin(s):
- Single origin: `https://driftlock.net`
- Multiple origins: `https://driftlock.net,https://www.driftlock.net`

The API responds to preflight `OPTIONS` requests and returns `Access-Control-Allow-*` headers for allowed origins.




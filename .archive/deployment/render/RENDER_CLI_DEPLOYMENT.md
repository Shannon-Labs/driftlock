# Deploy API to Render via CLI

## Prerequisites

✅ Render CLI installed: `brew install render`

## Step 1: Login to Render

```bash
render login
```

This will open your browser to authenticate. After logging in, you'll be able to use the CLI.

## Step 2: Deploy Using render.yaml

Since we already have `render.yaml` configured, you can deploy directly:

```bash
# From project root
render deploy
```

Or if you want to specify the blueprint file:

```bash
render deploy --blueprint render.yaml
```

## Step 3: Check Deployment Status

```bash
# List all services
render services list

# View logs
render logs driftlock-api

# Check service status
render services show driftlock-api
```

## Step 4: Get Your API URL

After deployment completes, Render will provide a URL like:
- `https://driftlock-api.onrender.com`

You can also get it via CLI:

```bash
render services show driftlock-api | grep "URL"
```

## Step 5: Connect Cloudflare to Render

Once you have the Render URL:

```bash
# Set the backend URL in Cloudflare Pages
wrangler pages secret put API_BACKEND_URL --project-name=driftlock
# Enter your Render URL when prompted

# Redeploy Pages
cd landing-page
npm run build
npx wrangler pages deploy dist --project-name=driftlock
```

## Alternative: Create Service Manually via CLI

If you prefer to create the service manually:

```bash
# Create a new web service
render services create web \
  --name driftlock-api \
  --dockerfilePath ./collector-processor/cmd/driftlock-http/Dockerfile \
  --dockerContext . \
  --env PORT=8080 \
  --env CORS_ALLOW_ORIGINS="https://driftlock.net,https://www.driftlock.net" \
  --env LOG_LEVEL=info \
  --healthCheckPath /healthz \
  --plan starter
```

## Useful CLI Commands

```bash
# View service details
render services show driftlock-api

# View live logs
render logs driftlock-api --follow

# Restart service
render services restart driftlock-api

# Update environment variables
render services update driftlock-api --env CORS_ALLOW_ORIGINS="https://driftlock.net,https://www.driftlock.net"

# Delete service (if needed)
render services delete driftlock-api
```

## Troubleshooting

### Authentication Issues
```bash
# Re-authenticate
render logout
render login
```

### Deployment Fails
- Check Dockerfile path is correct
- Verify you're in the project root directory
- Check build logs: `render logs driftlock-api`

### Service Not Starting
- Check health check path: `/healthz`
- Verify PORT environment variable
- Check service logs: `render logs driftlock-api --follow`


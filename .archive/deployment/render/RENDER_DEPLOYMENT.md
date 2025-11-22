# Deploy API to Render

## Quick Setup

### Option 1: Using render.yaml (Recommended)

1. **Push render.yaml to your repo** (already created in project root)

2. **Go to Render Dashboard**: https://dashboard.render.com

3. **Create New Blueprint**:
   - Click "New" → "Blueprint"
   - Connect your GitHub repository
   - Render will detect `render.yaml` automatically
   - Click "Apply"

4. **Get your API URL**: 
   - Render will provide: `https://driftlock-api.onrender.com` (or similar)
   - Copy this URL

5. **Update Cloudflare Pages Function**:
   ```bash
   wrangler pages secret put API_BACKEND_URL --project-name=driftlock
   # Enter your Render URL when prompted
   ```

6. **Redeploy Pages**:
   ```bash
   cd landing-page
   npm run build
   npx wrangler pages deploy dist --project-name=driftlock
   ```

### Option 2: Manual Setup

1. **Go to Render Dashboard**: https://dashboard.render.com

2. **Create New Web Service**:
   - Click "New" → "Web Service"
   - Connect your GitHub repository

3. **Configure Service**:
   - **Name**: `driftlock-api`
   - **Environment**: `Docker`
   - **Dockerfile Path**: `collector-processor/cmd/driftlock-http/Dockerfile`
   - **Docker Context**: `.` (root directory)

4. **Set Environment Variables**:
   ```
   PORT=8080
   CORS_ALLOW_ORIGINS=https://driftlock.net,https://www.driftlock.net
   LOG_LEVEL=info
   ```

5. **Set Health Check**:
   - **Health Check Path**: `/healthz`

6. **Choose Plan**:
   - **Starter** (Free tier) - good for testing
   - **Standard** ($7/month) - better for production

7. **Deploy**:
   - Click "Create Web Service"
   - Wait for build to complete (~5-10 minutes)
   - Get your URL: `https://driftlock-api.onrender.com`

## After Deployment

1. **Test the API**:
   ```bash
   curl https://driftlock-api.onrender.com/healthz
   # Should return JSON with success: true
   ```

2. **Set Backend URL in Cloudflare**:
   ```bash
   wrangler pages secret put API_BACKEND_URL --project-name=driftlock
   # Enter: https://driftlock-api.onrender.com
   ```

3. **Redeploy Pages**:
   ```bash
   cd landing-page
   npm run build
   npx wrangler pages deploy dist --project-name=driftlock
   ```

4. **Test Full Stack**:
   ```bash
   curl https://driftlock.net/api/v1/healthz
   # Should return JSON (not HTML)
   ```

## Render Free Tier Limitations

- **Spins down after 15 minutes of inactivity** (cold start ~30 seconds)
- **512MB RAM limit**
- **750 hours/month** (enough for always-on if you upgrade)

For production, consider upgrading to Standard plan ($7/month) for:
- Always-on (no spin-down)
- More RAM
- Better performance

## Troubleshooting

### Build Fails
- Check Dockerfile path is correct
- Ensure Rust and Go toolchains are available
- Check build logs in Render dashboard

### API Not Responding
- Check health check path is `/healthz`
- Verify PORT environment variable is set
- Check Render service logs

### CORS Errors
- Verify `CORS_ALLOW_ORIGINS` includes your frontend domain
- Check that Cloudflare Pages Function is proxying correctly


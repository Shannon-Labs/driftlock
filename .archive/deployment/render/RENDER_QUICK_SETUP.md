# Render Quick Setup Guide

## Current Configuration (What to Enter)

### Basic Settings
- **Name**: `driftlock-api` (or `driftlock`)
- **Language**: `Docker` ✅
- **Branch**: `main` ✅
- **Region**: `Oregon (US West)` ✅ (or your preferred region)

### Docker Settings
- **Dockerfile Path**: `./collector-processor/cmd/driftlock-http/Dockerfile` ⚠️ **IMPORTANT**
- **Root Directory**: Leave empty (or set to `.`)

### Instance Type
- **Free** ($0/month) - Good for testing
  - ⚠️ Spins down after 15 min inactivity (30 sec cold start)
- **Starter** ($9/month) - Recommended for production
  - Always-on, no spin-down

### Environment Variables
Click "Add Environment Variable" and add these:

1. **PORT**
   - Value: `8080`

2. **CORS_ALLOW_ORIGINS**
   - Value: `https://driftlock.net,https://www.driftlock.net`

3. **LOG_LEVEL**
   - Value: `info`

### Advanced Settings
- **Health Check Path**: `/healthz`
- **Auto-Deploy**: Enabled (default)

## After Deployment

1. **Get your API URL**: Render will show something like:
   - `https://driftlock-api.onrender.com`

2. **Test the API**:
   ```bash
   curl https://driftlock-api.onrender.com/healthz
   ```

3. **Connect Cloudflare**:
   ```bash
   wrangler pages secret put API_BACKEND_URL --project-name=driftlock
   # Enter your Render URL when prompted
   
   cd landing-page
   npm run build
   npx wrangler pages deploy dist --project-name=driftlock
   ```

4. **Test Full Stack**:
   ```bash
   curl https://driftlock.net/api/v1/healthz
   ```

## Common Issues

### Build Fails
- ✅ Check Dockerfile Path is: `./collector-processor/cmd/driftlock-http/Dockerfile`
- ✅ Ensure Root Directory is empty or `.`
- ✅ Check build logs in Render dashboard

### Service Won't Start
- ✅ Verify PORT=8080 is set
- ✅ Check Health Check Path is `/healthz`
- ✅ View logs: `render logs driftlock-api` (if using CLI)

### CORS Errors
- ✅ Verify CORS_ALLOW_ORIGINS includes your domain
- ✅ Check Cloudflare Pages Function is configured


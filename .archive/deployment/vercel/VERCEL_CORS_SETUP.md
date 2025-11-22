# Cloudflare API CORS Configuration Guide

## Quick Setup

To allow requests from `driftlock.net` to your API at `https://driftlock.net/api/v1`:

### Option 1: Using Cloudflare Dashboard (Recommended)

#### Step 1: Access Cloudflare Dashboard
1. Go to [Cloudflare Dashboard](https://dash.cloudflare.com)
2. Navigate to **Workers & Pages**
3. Find your API Worker/Pages project (the one handling `/api/v1` routes)
4. Click on the project name

#### Step 2: Add/Update CORS Environment Variable
1. Go to **Settings** → **Variables**
2. Scroll to **Environment Variables** section
3. Click **Add variable** (or edit existing `CORS_ALLOW_ORIGINS` if it exists)
4. Set the following:
   - **Variable name**: `CORS_ALLOW_ORIGINS`
   - **Value**: `https://driftlock.net,https://www.driftlock.net`
   - **Environment**: Select **Production** (and **Preview** if needed)
5. Click **Save**

#### Step 3: Redeploy (if needed)
- If using Cloudflare Workers: Changes take effect automatically on next request (or trigger a redeploy)
- If using Pages Functions: Changes take effect on next deployment

### Option 2: Using Wrangler CLI

```bash
# Set the CORS environment variable
wrangler secret put CORS_ALLOW_ORIGINS

# When prompted, enter:
# https://driftlock.net,https://www.driftlock.net
```

**Note**: If `CORS_ALLOW_ORIGINS` is not sensitive, you can also add it to `wrangler.toml`:

```toml
[env.production.vars]
CORS_ALLOW_ORIGINS = "https://driftlock.net,https://www.driftlock.net"
```

### Step 4: Verify
Test CORS from your browser console on `https://driftlock.net/playground`:

```javascript
fetch('https://driftlock.net/api/v1/healthz', {
  method: 'GET',
  headers: { 'Origin': 'https://driftlock.net' }
})
.then(r => r.json())
.then(console.log)
.catch(console.error)
```

You should see a successful response without CORS errors.

## How It Works

The API code (`collector-processor/cmd/driftlock-http/main.go`) reads `CORS_ALLOW_ORIGINS` and:
- Splits comma-separated values
- Checks incoming `Origin` header against allowed origins
- Returns appropriate `Access-Control-Allow-Origin` header

## Troubleshooting

### CORS Still Not Working?
1. **Check environment variable is set**: Verify in Cloudflare dashboard that `CORS_ALLOW_ORIGINS` exists
2. **Verify deployment**: If using Workers, changes are immediate. If using Pages Functions, ensure latest deployment is active
3. **Check origin match**: The origin must match exactly (including `https://` and no trailing slash)
4. **Browser cache**: Try hard refresh (Cmd+Shift+R / Ctrl+Shift+R) or incognito mode
5. **Check Worker/Pages Function**: Ensure your API code is reading `CORS_ALLOW_ORIGINS` from environment variables

### Testing CORS Manually
```bash
# Test preflight OPTIONS request
curl -X OPTIONS https://driftlock.net/api/v1/detect \
  -H "Origin: https://driftlock.net" \
  -H "Access-Control-Request-Method: POST" \
  -v

# Should return headers like:
# Access-Control-Allow-Origin: https://driftlock.net
# Access-Control-Allow-Methods: POST, GET, OPTIONS
```

## Current Configuration

- **API Endpoint**: `https://driftlock.net/api/v1`
- **Frontend**: `https://driftlock.net`
- **Required CORS Value**: `https://driftlock.net,https://www.driftlock.net`


# DriftLock Cloudflare Deployment Guide 🚀

This guide will walk you through deploying DriftLock to Cloudflare for production.

## Architecture

```
┌─────────────────────────────────────────┐
│          Cloudflare CDN                 │
│  ┌──────────────┐  ┌─────────────────┐ │
│  │   Web App    │  │   API Server    │ │
│  │  (Pages)     │  │   (Workers)     │ │
│  │              │  │                 │ │
│  │  React SPA   │  │   Hono API      │ │
│  └──────┬───────┘  └────────┬────────┘ │
│         │                    │          │
│         └─────────┬──────────┘          │
│                   │                      │
└───────────────────┼──────────────────────┘
                    │
            ┌───────▼────────┐
            │   Supabase     │
            │                │
            │ - Database     │
            │ - Auth         │
            │ - Edge Funcs   │
            │ - Realtime     │
            └────────────────┘
```

## Prerequisites

- [ ] Cloudflare account (free tier works)
- [ ] Cloudflare CLI installed
- [ ] Supabase project configured
- [ ] Stripe account configured

### Install Cloudflare CLI

**Option 1: Via Wrangler (Recommended)**
```bash
npm install -g wrangler
```

**Option 2: Direct installation**
```bash
curl -L https://github.com/cloudflare/cloudflare-cli/releases/latest/download/cf.tgz | tar -xz
sudo cp cf /usr/local/bin
```

**Verify installation:**
```bash
wrangler --version
cf --version
```

## Deployment Overview

We'll deploy:
1. **Web Frontend** → Cloudflare Pages (React app)
2. **API Server** → Cloudflare Workers (TypeScript/Hono)
3. **Database** → Keep using Supabase (already configured)

## Step 1: Deploy API to Cloudflare Workers

### 1.1 Navigate to API Worker directory
```bash
cd /Volumes/VIXinSSD/driftlock/cloudflare-api-worker
```

### 1.2 Install dependencies
```bash
npm install
```

### 1.3 Login to Cloudflare
```bash
wrangler login
```

### 1.4 Set secrets
```bash
# Get these from your .env file
wrangler secret put SUPABASE_SERVICE_ROLE_KEY
# Paste the value from SUPABASE_SERVICE_ROLE_KEY in /Volumes/VIXinSSD/driftlock/.env

wrangler secret put STRIPE_WEBHOOK_SECRET
# Get this from Stripe Dashboard after creating webhook
```

### 1.5 Deploy to staging first
```bash
wrangler deploy --env staging
```

### 1.6 Test staging deployment
```bash
curl https://driftlock-api-staging.YOUR_ACCOUNT.workers.dev/health
```

### 1.7 Deploy to production
```bash
wrangler deploy --env production
```

### 1.8 Update Stripe Webhook

After successful deployment, update your Stripe webhook URL:
- **Old URL:** `https://nfkdeeunyvnntvpvwpwh.supabase.co/functions/v1/stripe-webhook`
- **New URL:** `https://driftlock-api.YOUR_SUBDOMAIN.workers.dev/stripe-webhook`

Get your Worker URL:
```bash
wrangler route list
```

Update webhook in Stripe Dashboard:
1. Go to https://dashboard.stripe.com/test/webhooks
2. Click on your webhook
3. Update "Endpoint URL"
4. Save

## Step 2: Deploy Web Frontend to Cloudflare Pages

### 2.1 Navigate to web-frontend
```bash
cd /Volumes/VIXinSSD/driftlock/web-frontend
```

### 2.2 Install dependencies
```bash
npm install
```

### 2.3 Build the project
```bash
npm run build
```

### 2.4 Login to Cloudflare
```bash
cf pages login
```

### 2.5 Create Pages project
```bash
# Option 1: Interactive (recommended for first time)
cf pages project create driftlock-web-frontend --commit-source control

# Option 2: Direct deployment
cf pages deploy . --project-name driftlock-web-frontend
```

### 2.6 Configure environment variables

In the Cloudflare Pages dashboard:
1. Go to your project → Settings → Environment variables
2. Add the following:

**Production variables:**
```
VITE_SUPABASE_PROJECT_ID = nfkdeeunyvnntvpvwpwh
VITE_SUPABASE_PUBLISHABLE_KEY = eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Im5ma2RlZXVueXZubnR2cHZ3cHdoIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NTg1ODU2NDUsImV4cCI6MjA3NDE2MTY0NX0.nRjQZJG5h66OgvQs8z9dmpQKw3nNHTTjhiRwdt48YGo
VITE_SUPABASE_URL = https://nfkdeeunyvnntvpvwpwh.supabase.co
```

**Staging variables:**
```
VITE_SUPABASE_PROJECT_ID = nfkdeeunyvnntvpvwpwh
VITE_SUPABASE_PUBLISHABLE_KEY = eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Im5ma2RlZXVueXZubnR2cHZ3cHdoIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NTg1ODU2NDUsImV4cCI6MjA3NDE2MTY0NX0.nRjQZJG5h66OgvQs8z9dmpQKw3nNHTTjhiRwdt48YGo
VITE_SUPABASE_URL = https://nfkdeeunyvnntvpvwpwh.supabase.co
```

### 2.7 Deploy
```bash
cf pages deploy . --project-name driftlock-web-frontend
```

### 2.8 Get your deployment URL
```bash
cf pages domain list driftlock-web-frontend
```

## Step 3: Configure Custom Domain (Optional)

### 3.1 Add custom domain to Pages
```bash
# For web frontend
cf pages domain add driftlock-web-frontend driftlock.com
cf pages domain add driftlock-web-frontend www.driftlock.com

# For API Worker
# Add to wrangler.toml:
# [[routes]]
# pattern = "api.driftlock.com/*"
# zone_name = "driftlock.com"

# Then deploy:
wrangler deploy --env production
```

### 3.2 DNS Configuration

In Cloudflare DNS dashboard:
1. Add CNAME record for `driftlock.com` → `driftlock-web-frontend.pages.dev`
2. Add CNAME record for `www.driftlock.com` → `driftlock-web-frontend.pages.dev`
3. Add CNAME record for `api.driftlock.com` → `driftlock-api.YOUR_ACCOUNT.workers.dev`

### 3.3 SSL Certificate

Cloudflare automatically provides SSL certificates. Ensure:
- SSL/TLS encryption mode: "Full (strict)"
- Always Use HTTPS: ON
- HTTP Strict Transport Security (HSTS): Enabled

## Step 4: Update Supabase Settings

If using custom domains, update Supabase:

1. Go to Supabase Dashboard → Authentication → Settings
2. Update "Site URL" to your custom domain
3. Add redirect URLs for any OAuth providers

## Step 5: Test the Deployment

### 5.1 Test Web Frontend
```bash
# Visit your Pages URL
open https://driftlock-web-frontend.pages.dev

# Check if page loads
curl -I https://driftlock-web-frontend.pages.dev
```

### 5.2 Test API Worker
```bash
# Health check
curl https://driftlock-api.YOUR_SUBDOMAIN.workers.dev/health

# Get API info
curl https://driftlock-api.YOUR_SUBDOMAIN.workers.dev/

# Test anomaly creation (replace ORG_ID with actual ID)
curl -X POST https://driftlock-api.YOUR_SUBDOMAIN.workers.dev/anomalies \
  -H "Content-Type: application/json" \
  -d '{"organization_id": "ORG_ID", "event_type": "test", "severity": "medium"}'
```

### 5.3 Test Stripe Integration
```bash
# Trigger test webhook
stripe listen --forward-to https://driftlock-api.YOUR_SUBDOMAIN.workers.dev/stripe-webhook
stripe trigger checkout.session.completed
```

## Step 6: Update Environment in Frontend

### 6.1 Update API Endpoint

If you changed the API endpoint, update the frontend code:

**File:** `web-frontend/src/lib/supabase.ts`

```typescript
// Update this to your Worker URL
const supabaseUrl = 'https://driftlock-api.YOUR_SUBDOMAIN.workers.dev'
```

Or create a new environment variable:
```typescript
const supabaseUrl = import.meta.env.VITE_API_URL || 'https://nfkdeeunyvnntvpvwpwh.supabase.co'
```

### 6.2 Rebuild and redeploy
```bash
cd /Volumes/VIXinSSD/driftlock/web-frontend
npm run build
cf pages deploy . --project-name driftlock-web-frontend
```

## Automated Deployment Scripts

We created deployment scripts to make this easier:

### Deploy API
```bash
cd /Volumes/VIXinSSD/driftlock/cloudflare-api-worker
bash deploy.sh
```

### Deploy Frontend
```bash
cd /Volumes/VIXinSSD/driftlock/web-frontend
bash deploy-pages.sh
```

## Monitoring

### 5.1 Cloudflare Analytics

**Pages Analytics:**
- Dashboard → Cloudflare Pages → Your Project → Analytics

**Worker Metrics:**
- Dashboard → Workers & Pages → Your Worker → Metrics

### 5.2 Application Monitoring

**Supabase:**
- Dashboard → Logs
- https://supabase.com/dashboard/project/nfkdeeunyvnntvpvwpwh/logs

**Stripe:**
- Dashboard → Webhooks
- Monitor delivery success rate

### 5.3 Error Tracking

Consider adding:
- Sentry for error tracking
- LogRocket for session replay
- Datadog for APM

## Security Checklist

- [ ] All secrets stored in Cloudflare, not in code
- [ ] HTTPS enforced (SSL/TLS: Full strict)
- [ ] CORS properly configured
- [ ] RLS enabled in Supabase
- [ ] Webhook signature verification enabled
- [ ] Rate limiting configured (Cloudflare)
- [ ] Security headers added (CSP, HSTS, etc.)

## Performance Optimization

### Cloudflare Cache
```javascript
// Add to _headers file in web-frontend/public/
/*
  Cache-Control: public, max-age=31536000, immutable
*/
```

### Edge Functions
- Cloudflare Workers run at the edge (low latency)
- API responses cached via Cloudflare CDN
- Static assets cached for 1 year

## Troubleshooting

### Build Fails

**Frontend:**
```bash
# Check Node version
node --version  # Should be 18+

# Clear cache
rm -rf node_modules package-lock.json
npm install
npm run build
```

**API Worker:**
```bash
# Check TypeScript errors
npm run typecheck

# Check Wrangler config
wrangler.toml --validate
```

### API Returns 500 Error

**Check logs:**
```bash
# Worker logs
wrangler tail driftlock-api --env production

# Edge function logs
# https://supabase.com/dashboard/project/nfkdeeunyvnntvpvwpwh/functions
```

### CORS Errors

**Check:**
1. CORS_ORIGIN environment variable in Worker
2. Exact domain match in Supabase Auth settings
3. Headers in API responses

### Stripe Webhook Failing

1. Check webhook URL is correct
2. Verify STRIPE_WEBHOOK_SECRET matches Stripe dashboard
3. Test with Stripe CLI:
   ```bash
   stripe listen --forward-to YOUR_WEBHOOK_URL
   stripe trigger checkout.session.completed
   ```

## Production Checklist

- [ ] Deployed to Cloudflare Workers (API)
- [ ] Deployed to Cloudflare Pages (Frontend)
- [ ] SSL certificates active
- [ ] Environment variables configured
- [ ] Secrets stored securely
- [ ] Stripe webhook URL updated
- [ ] Custom domain configured (optional)
- [ ] DNS configured
- [ ] Monitoring set up
- [ ] Backup strategy in place
- [ ] Load tested
- [ ] Security audit completed

## Cost Estimation

**Free Tier Includes:**
- 100,000 Workers requests/day
- 500 build minutes/month (Pages)
- Unlimited personal projects
- Unlimited bandwidth

**Paid Tier ($20/month):**
- 10M Workers requests
- 10,000 build minutes
- Priority support

**Stripe Fees:**
- 2.9% + 30¢ per successful payment
- Volume discounts available

## Support

- [Cloudflare Docs](https://developers.cloudflare.com/)
- [Supabase Docs](https://supabase.com/docs)
- [Stripe Docs](https://stripe.com/docs)
- [Hono Docs](https://hono.dev/)

## Summary

✅ Cloudflare Workers deployed for API server
✅ Cloudflare Pages deployed for web frontend
✅ Supabase continues as managed database
✅ Stripe webhook integrated
✅ Custom domains configured
✅ SSL certificates active
✅ Production monitoring enabled

**Your DriftLock platform is now live on Cloudflare!** 🎉

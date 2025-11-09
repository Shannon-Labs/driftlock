# DriftLock Cloudflare Deployment Guide

This directory contains configuration files and scripts for deploying DriftLock to Cloudflare.

## Prerequisites

1. **Cloudflare Account** - Sign up at https://cloudflare.com
2. **Cloudflare CLI** - Install `wrangler` and `cf` CLI tools
3. **Supabase Project** - Already configured with credentials
4. **Stripe Account** - Already configured with products

## Architecture

```
┌─────────────────┐
│  Cloudflare     │
│  Pages          │  ← Web Frontend (React)
│  (driftlock-web)│
└────────┬────────┘
         │
         │ API Calls
         ▼
┌─────────────────┐
│  Cloudflare     │
│  Workers        │  ← API Server (TypeScript)
│  (driftlock-api)│
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Supabase      │  ← Database + Auth + Functions
│  (Managed)      │
└─────────────────┘
```

## Deployment Options

### Option 1: Pages + Workers (Recommended)

**Web-Frontend** → Cloudflare Pages  
**API Server** → Cloudflare Workers  
**Database** → Supabase (managed)

**Steps:**
1. Copy `pages.toml` to `/Volumes/VIXinSSD/driftlock/web-frontend/`
2. Update environment variables in pages.toml
3. Run `cf pages deploy`

### Option 2: Full Cloudflare

**Web-Frontend** → Cloudflare Pages  
**API Server** → Cloudflare Workers  
**Database** → Cloudflare D1

**Steps:**
1. Use `wrangler.toml` as template
2. Create D1 database
3. Run migrations
4. Deploy both services

### Option 3: Hybrid

**Web-Frontend** → Cloudflare Pages  
**API Server** → Your own VPS/cloud provider  
**Database** → Supabase

## Quick Start (Option 1)

### 1. Deploy Web-Frontend to Pages

```bash
cd /Volumes/VIXinSSD/driftlock/web-frontend

# Create Pages project
cf pages project create driftlock-web-frontend

# Deploy
cf pages deploy . --project-name driftlock-web-frontend

# Set environment variables
cf pages env var VITE_SUPABASE_PROJECT_ID nfkdeeunyvnntvpvwpwh
cf pages env var VITE_SUPABASE_PUBLISHABLE_KEY eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
cf pages env var VITE_SUPABASE_URL https://nfkdeeunyvnntvpvwpwh.supabase.co
```

### 2. Deploy API to Workers

```bash
cd /Volumes/VIXinSSD/driftlock/cloudflare-deployment

# Login to Cloudflare
wrangler login

# Deploy
wrangler deploy
```

### 3. Update Stripe Webhook

After deployment, update your Stripe webhook URL:
```
https://driftlock-api.YOUR_SUBDOMAIN.workers.dev/stripe-webhook
```

## Environment Variables

### Web-Frontend (Pages)
```bash
VITE_SUPABASE_PROJECT_ID=nfkdeeunyvnntvpvwpwh
VITE_SUPABASE_PUBLISHABLE_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
VITE_SUPABASE_URL=https://nfkdeeunyvnntvpvwpwh.supabase.co
```

### API Worker
```bash
SUPABASE_PROJECT_ID=nfkdeeunyvnntvpvwpwh
SUPABASE_URL=https://nfkdeeunyvnntvpvwpwh.supabase.co
SUPABASE_SERVICE_ROLE_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
ENVIRONMENT=production
```

## Configuration Files

- `pages.toml` - Cloudflare Pages configuration for web-frontend
- `wrangler.toml` - Cloudflare Workers configuration for API
- `deploy.sh` - Automated deployment script

## Testing Deployment

### Test Web-Frontend
```bash
# Get Pages URL
cf pages domain list driftlock-web-frontend

# Visit in browser
open https://driftlock-web-frontend.pages.dev
```

### Test API Worker
```bash
# Get Worker URL
wrangler route list

# Test health endpoint
curl https://driftlock-api.YOUR_SUBDOMAIN.workers.dev/health
```

### Test Stripe Integration
```bash
# Trigger test webhook
stripe trigger checkout.session.completed

# Check Stripe logs
stripe logs tail
```

## Custom Domains

### Web-Frontend
```bash
# Add custom domain
cf pages domain add driftlock-web-frontend www.driftlock.com

# Or use Pages.dev subdomain
# https://driftlock-web-frontend.pages.dev
```

### API Worker
```bash
# Add route to wrangler.toml
[[routes]]
pattern = "api.driftlock.com/*"
zone_name = "driftlock.com"
```

## Troubleshooting

### Build Fails
- Check Node.js version (use 18)
- Verify all dependencies are installed
- Check for TypeScript errors

### Worker Deployment Fails
- Verify wrangler.toml configuration
- Check environment variables are set
- Ensure D1 database is created (if using)

### API Not Responding
- Check Worker logs: `wrangler tail`
- Verify route configuration
- Check CORS settings

### Stripe Webhook Failing
- Verify webhook URL in Stripe dashboard
- Check webhook secret matches environment
- Test with Stripe CLI: `stripe listen --forward-to YOUR_WEBHOOK_URL`

## Production Checklist

- [ ] All environment variables set
- [ ] Custom domains configured
- [ ] SSL certificates active
- [ ] Stripe webhook URL updated
- [ ] Database migrations applied
- [ ] Edge Functions deployed
- [ ] Monitoring enabled
- [ ] Error tracking configured
- [ ] Backup strategy in place

## Monitoring

### Cloudflare Analytics
- Pages: Automatic analytics
- Workers: Built-in metrics

### Application Monitoring
- Add Sentry or similar for error tracking
- Monitor Supabase metrics
- Track Stripe webhook success rate

## Security

- Always use environment variables for secrets
- Enable RLS in Supabase
- Use HTTPS only
- Configure CORS properly
- Regular security audits

## Support

- Cloudflare Docs: https://developers.cloudflare.com/
- Supabase Docs: https://supabase.com/docs
- Stripe Docs: https://stripe.com/docs

#!/bin/bash
# Cloudflare Deployment Script for DriftLock

set -e

echo "🚀 Starting DriftLock deployment to Cloudflare..."

# Deploy Web-Frontend to Cloudflare Pages
echo ""
echo "📦 Deploying web-frontend to Cloudflare Pages..."
cd /Volumes/VIXinSSD/driftlock/web-frontend

# Set environment variables for Pages
cf pages project create driftlock-web-frontend || true
cf pages deploy . --project-name driftlock-web-frontend \
  --prod \
  --env production \
  --variable VITE_SUPABASE_PROJECT_ID:$VITE_SUPABASE_PROJECT_ID \
  --variable VITE_SUPABASE_PUBLISHABLE_KEY:$VITE_SUPABASE_PUBLISHABLE_KEY \
  --variable VITE_SUPABASE_URL:$VITE_SUPABASE_URL

# Deploy API Worker
echo ""
echo "⚡ Deploying API to Cloudflare Workers..."
cd /Volumes/VIXinSSD/driftlock/cloudflare-deployment

wrangler deploy driftlock-api --name driftlock-api

echo ""
echo "✅ Deployment complete!"
echo ""
echo "📊 Web-Frontend: https://driftlock-web-frontend.pages.dev"
echo "🔧 API Worker: https://driftlock-api.your-subdomain.workers.dev"
echo ""
echo "🔗 Update your Stripe webhook URL to:"
echo "   https://driftlock-api.your-subdomain.workers.dev/stripe-webhook"

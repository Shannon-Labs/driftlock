#!/bin/bash

echo "========================================="
echo "DriftLock API - Cloudflare Workers Deployment"
echo "========================================="
echo ""

# Check if wrangler is installed
if ! command -v wrangler &> /dev/null; then
    echo "❌ Wrangler CLI not found. Installing..."
    npm install -g wrangler
fi

echo "✅ Wrangler version: $(wrangler --version)"
echo ""

# Login check
echo "Checking Cloudflare authentication..."
if ! wrangler whoami &> /dev/null; then
    echo "❌ Not logged in to Cloudflare"
    echo "Please run: wrangler login"
    exit 1
fi
echo "✅ Logged in to Cloudflare"
echo ""

# Install dependencies
echo "Installing dependencies..."
npm install
echo ""

# Set secrets
echo "========================================="
echo "Setting secrets..."
echo "========================================="
echo ""
echo "Setting SUPABASE_SERVICE_ROLE_KEY..."
echo "Please get this from: https://supabase.com/dashboard/project/nfkdeeunyvnntvpvwpwh/settings/api"
echo ""

# Read the service role key from the .env file
SERVICE_ROLE_KEY=$(grep SUPABASE_SERVICE_ROLE_KEY /Volumes/VIXinSSD/driftlock/.env | cut -d '=' -f2)
if [ -z "$SERVICE_ROLE_KEY" ]; then
    echo "❌ SUPABASE_SERVICE_ROLE_KEY not found in .env file"
    echo "Please run: wrangler secret put SUPABASE_SERVICE_ROLE_KEY"
else
    echo $SERVICE_ROLE_KEY | wrangler secret put SUPABASE_SERVICE_ROLE_KEY
fi

echo ""
echo "Setting STRIPE_WEBHOOK_SECRET..."
echo "Please get this from: https://dashboard.stripe.com/test/webhooks"
echo "⚠️  You'll need to create a webhook endpoint first: https://driftlock-api.YOUR_SUBDOMAIN.workers.dev/stripe-webhook"
echo ""

# Deploy to staging first
echo "========================================="
echo "Deploying to STAGING environment..."
echo "========================================="
wrangler deploy --env staging

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Successfully deployed to STAGING"
    STAGING_URL=$(wrangler deployments list | head -1 | awk '{print $2}')
    echo "📍 Staging URL: $STAGING_URL"
    echo ""

    # Ask if user wants to deploy to production
    echo "========================================="
    read -p "Deploy to PRODUCTION? (y/N): " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo ""
        echo "========================================="
        echo "Deploying to PRODUCTION environment..."
        echo "========================================="
        wrangler deploy --env production

        if [ $? -eq 0 ]; then
            echo ""
            echo "✅ Successfully deployed to PRODUCTION"
            echo ""
            echo "========================================="
            echo "🎉 DEPLOYMENT COMPLETE!"
            echo "========================================="
            echo ""
            echo "Next steps:"
            echo "1. Update Stripe webhook URL:"
            echo "   Old: https://nfkdeeunyvnntvpvwpwh.supabase.co/functions/v1/stripe-webhook"
            echo "   New: https://driftlock-api.YOUR_SUBDOMAIN.workers.dev/stripe-webhook"
            echo ""
            echo "2. Update web-frontend API endpoint if needed"
            echo ""
            echo "3. Test the deployment:"
            echo "   curl https://driftlock-api.YOUR_SUBDOMAIN.workers.dev/health"
            echo ""
        else
            echo "❌ Production deployment failed"
            exit 1
        fi
    fi
else
    echo "❌ Staging deployment failed"
    exit 1
fi

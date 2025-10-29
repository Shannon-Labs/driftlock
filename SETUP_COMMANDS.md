# Commands to Set Up Secrets for DriftLock API Gateway

## Set All Required Secrets

Run these commands in your terminal:

```bash
cd /Volumes/VIXinSSD/driftlock/cloudflare-api-worker

# Set Supabase service role key
wrangler secret put SUPABASE_SERVICE_ROLE_KEY

# Set JWT secret
# Optional, if the worker validates JWTs (not required by default)
wrangler secret put JWT_SECRET

# Set Stripe webhook secret (use the one from your Stripe dashboard)
wrangler secret put STRIPE_WEBHOOK_SECRET
# When prompted, enter: whsec_DHBt7I8WKbWVb1RK0mGAg51hcRLZ7CSC

# If needed, set Supabase URL as secret instead of in vars
# wrangler secret put SUPABASE_URL
```

## Update wrangler.toml for Production

When you're ready to deploy to production, update the `ENV_GO_BACKEND_URL` in wrangler.toml:

```toml
[vars]
SUPABASE_URL = "https://your-project.supabase.co"
# Optional proxy target if routing to Go backend behind the worker
GO_BACKEND_URL = "https://your-backend-url.com"
```

## Deploy the Updated Worker

After setting all secrets:

```bash
wrangler deploy
```

## Update Stripe Dashboard

After deployment, if you want to change your webhook URL to point to your Cloudflare Worker:

1. Go to https://dashboard.stripe.com/webhooks
2. Click on your current webhook endpoint
3. Update the URL to: `https://driftlock-api.YOUR-USERNAME.workers.dev/stripe-webhook`
4. Save the changes

## Verify Everything Works

1. Check that your worker is deployed: `curl https://driftlock-api.YOUR-USERNAME.workers.dev/health`
2. Monitor your webhook deliveries in the Stripe dashboard
3. Test subscription events to ensure they're processed correctly

Your API Gateway is now fully configured to handle both API requests and Stripe webhooks! 🎉

#!/bin/bash

echo "Testing the DriftLock API Gateway setup..."

# Test health endpoint of the Go backend directly
echo "Testing Go backend health endpoint at http://localhost:8081/health..."
response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8081/health 2>/dev/null)
if [ "$response" = "200" ]; then
    echo "✅ Go backend is running and accessible on port 8081"
else
    echo "❌ Go backend is not responding on port 8081 (HTTP $response)"
fi

echo ""
echo "Cloudflare Workers API Gateway is running on port 8787"
echo "To test the gateway, make requests to: http://localhost:8787/api/v1/..."
echo "The gateway will proxy to your backend at http://localhost:8081"
echo ""
echo "When you deploy the worker to Cloudflare, it will be available at: https://[your-worker].workers.dev/"
echo ""
echo "Next steps:"
echo "1. Set your secrets with: wrangler secret put ENV_SUPABASE_SERVICE_ROLE_KEY"
echo "2. Set JWT secret with: wrangler secret put ENV_JWT_SECRET"
echo "3. When deploying to production, update ENV_GO_BACKEND_URL in wrangler.toml to your production backend URL"
echo "4. Deploy with: wrangler deploy"
echo ""
echo "For Stripe webhook in your Stripe dashboard, subscribe to these events:"
echo "- customer.subscription.created, updated, deleted"
echo "- invoice.payment_succeeded, payment_failed"
echo "- checkout.session.completed"
echo "- customer.subscription.paused, resumed (optional)"
echo "- customer.subscription.trial_will_end (optional)"
echo ""
echo "Webhook endpoint will be: https://[your-worker-subdomain].your-username.workers.dev/"
echo "Remember to save the webhook signing secret in your Stripe dashboard and Cloudflare secrets"
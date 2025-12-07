#!/bin/bash
# Simple test of Stripe checkout endpoint
# Usage: ./test-stripe-checkout.sh <email> <plan>

set -e

EMAIL="${1:-test-$(date +%s)@driftlock.net}"
PLAN="${2:-radar}"

echo "=== Stripe Checkout Endpoint Test ==="
echo "Email: $EMAIL"
echo "Plan: $PLAN"
echo ""

# First, let's use an existing verified user to test the checkout endpoint
# We'll use the API key we know works from email verification
echo "Testing checkout endpoint with authenticated user..."
echo ""

# For this test, we need an API key. Let's check if we have any in the database
# or use a test endpoint if available
echo "Note: This test requires an authenticated user (with API key) to create a checkout session"
echo "The flow is:"
echo "1. User signs up and verifies email"
echo "2. User logs in and gets API key"
echo "3. User calls /v1/billing/checkout with API key to create Stripe session"
echo "4. User is redirected to Stripe to complete payment"
echo ""

# Let's test if the checkout endpoint exists and responds correctly to unauthenticated requests
echo "Testing checkout endpoint (unauthenticated)..."
CHECKOUT_RESPONSE=$(curl -s -X POST "https://driftlock.net/api/v1/billing/checkout" \
  -H "Content-Type: application/json" \
  -d "{\"plan\":\"$PLAN\"}")

echo "Response (expecting error about missing authentication):"
echo "$CHECKOUT_RESPONSE" | jq '.' 2>/dev/null || echo "$CHECKOUT_RESPONSE"
echo ""

if echo "$CHECKOUT_RESPONSE" | grep -q "unauthorized\|authentication\|missing"; then
    echo "✅ Checkout endpoint exists and properly requires authentication"
else
    echo "⚠️  Unexpected response from checkout endpoint"
fi

echo ""
echo "=== Expected Checkout Flow (with authenticated user) ==="
echo "Request:"
echo "curl -X POST 'https://driftlock.net/api/v1/billing/checkout' \\"
echo "  -H 'Authorization: Bearer <API_KEY>' \\"
echo "  -H 'Content-Type: application/json' \\"
echo "  -d '{\"plan\":\"radar\"}'"
echo ""
echo "Expected Response:"
echo '{
  "success": true,
  "checkout_url": "https://checkout.stripe.com/pay/cs_test_...",
  "stripe_session_id": "cs_test_..."
}'
echo ""
echo "Note: The checkout_url will redirect to Stripe's hosted checkout page"
echo "After payment, Stripe will send a webhook to update the subscription"
#!/bin/bash
# Test Stripe checkout flow
# Usage: ./test-stripe.sh <firebase_auth_token> [plan]
#
# This script tests the Stripe checkout flow:
# 1. Creates a checkout session with auth token
# 2. Returns Stripe checkout URL
#
# Prerequisites:
# - Must have a verified Driftlock account
# - Need a valid Firebase auth token
# - STRIPE_SECRET_KEY must be configured in production
#
# To get a Firebase auth token:
# 1. Sign in to driftlock.net
# 2. Open browser dev tools > Application > IndexedDB > firebaseLocalStorage
# 3. Find the access token

set -e

API_URL="${API_URL:-https://driftlock.net}"
TOKEN="${1:-}"
PLAN="${2:-radar}"

if [ -z "$TOKEN" ]; then
  echo "Error: Firebase auth token required"
  echo ""
  echo "Usage: $0 <firebase_auth_token> [plan]"
  echo ""
  echo "Plans: radar ($15/mo), tensor ($100/mo), orbit ($499/mo)"
  echo ""
  echo "To get your Firebase token:"
  echo "1. Sign in at https://driftlock.net"
  echo "2. Open browser DevTools (F12)"
  echo "3. Go to Application > IndexedDB > firebaseLocalStorageDb"
  echo "4. Find the 'firebase:authUser:...' entry"
  echo "5. Copy the 'stsTokenManager.accessToken' value"
  exit 1
fi

echo "Testing Stripe checkout flow..."
echo "API URL: $API_URL"
echo "Plan: $PLAN"
echo ""

# Create checkout session
echo "Creating checkout session..."
RESPONSE=$(curl -s -X POST "$API_URL/api/v1/billing/checkout" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"plan\":\"$PLAN\"}")

echo "Response: $RESPONSE"
echo ""

# Extract checkout URL
if echo "$RESPONSE" | grep -q '"checkout_url"'; then
  CHECKOUT_URL=$(echo "$RESPONSE" | grep -o '"checkout_url":"[^"]*"' | cut -d'"' -f4)
  echo "Checkout URL: $CHECKOUT_URL"
  echo ""
  echo "Next steps:"
  echo "1. Open the checkout URL in your browser"
  echo "2. Complete payment with Stripe test card: 4242 4242 4242 4242"
  echo "3. Verify subscription in Stripe dashboard"
  echo "4. Check webhook delivery in Stripe dashboard > Developers > Webhooks"
else
  echo "Failed to create checkout session. Check response above."
  exit 1
fi

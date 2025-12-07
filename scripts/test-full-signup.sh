#!/bin/bash
# Test complete user signup flow
# Usage: ./test-full-signup.sh <your-email>

set -e

EMAIL="${1:-test-$(date +%s)@driftlock.net}"

echo "=== Testing Complete User Signup Flow ==="
echo "Email: $EMAIL"
echo "Timestamp: $(date)"
echo ""

# Step 1: Sign up
echo "Step 1: Creating account..."
echo "Running: curl -X POST https://driftlock.net/api/v1/onboard/signup"
echo "Data: {\"email\":\"$EMAIL\",\"company_name\":\"Test Company $(date +%s)\",\"plan\":\"trial\"}"
echo ""

SIGNUP_RESPONSE=$(curl -s -X POST https://driftlock.net/api/v1/onboard/signup \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"company_name\":\"Test Company $(date +%s)\",\"plan\":\"trial\"}")

echo "Signup Response:"
echo "$SIGNUP_RESPONSE" | jq '.' 2>/dev/null || echo "$SIGNUP_RESPONSE"

# Extract tenant info if available
TENANT_ID=$(echo "$SIGNUP_RESPONSE" | jq -r '.tenant.id // empty' 2>/dev/null)
if [ -n "$TENANT_ID" ] && [ "$TENANT_ID" != "empty" ]; then
    echo ""
    echo "✅ Account created!"
    echo "Tenant ID: $TENANT_ID"

    # The verification link is sent via email
    # We can't programmatically retrieve it without checking email
    # So let's show the user what to expect

    echo ""
    echo "✅ Account created!"
    echo "Tenant ID: $TENANT_ID"
    echo ""

    # Check if email was marked as sent
    if echo "$SIGNUP_RESPONSE" | grep -q '"verification_sent":true'; then
        echo ""
        echo "✅ Verification email sent!"
        echo "📧 Check your inbox at $EMAIL"
        echo ""
        echo "Step 3: Manual Verification Required"
        echo "1. Open the verification email"
        echo "2. Click the verification link"
        echo "3. The link will look like: https://driftlock.net/api/v1/onboard/verify?token=..."
        echo "4. After clicking, you'll receive your API key in the response"
        echo ""
        echo "Expected verification response format:"
        echo '{
          "success": true,
          "message": "Email verified successfully!",
          "api_key": "dlk_<uuid>.<secret>",
          "pending_verification": false
        }'
        echo ""

        # Extract API key from the signup response in case it's included
        API_KEY=$(echo "$SIGNUP_RESPONSE" | jq -r '.api_key // "Not provided in signup"')

        if [ "$API_KEY" != "Not provided in signup" ]; then
            echo "Note: Some implementations include the API key in the signup response"
            echo "If this is the case, your API key would be: $API_KEY"
        fi
    else
        echo ""
        echo "❌ Verification email not marked as sent"
    fi
else
    echo ""
    echo "❌ Failed to create account"
    echo "Response: $SIGNUP_RESPONSE"
fi
#!/bin/bash
# Test full checkout flow including signup → verification → upgrade → payment
# Usage: ./test-full-checkout.sh <email> <plan>
# Plans: pilot, radar, tensor, orbit

set -e

EMAIL="${1:-test-$(date +%s)@driftlock.net}"
PLAN="${2:-radar}"  # Default to radar tier ($100/mo)
# PASSWORD is not needed for signup

echo "=== Full Checkout Flow Test ==="
echo "Email: $EMAIL"
echo "Plan: $PLAN"
echo "Timestamp: $(date)"
echo ""

# Validate plan
case "$PLAN" in
    pilot|radar|tensor|orbit)
        echo "✅ Valid plan: $PLAN"
        ;;
    *)
        echo "❌ Invalid plan. Options: pilot, radar, tensor, orbit"
        exit 1
        ;;
esac

# Step 1: Create account
echo "Step 1: Creating account..."
echo ""

SIGNUP_RESPONSE=$(curl -s -X POST https://driftlock.net/api/v1/onboard/signup \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"company_name\":\"Test Checkout $(date +%s)\",\"plan\":\"trial\"}")

TENANT_ID=$(echo "$SIGNUP_RESPONSE" | jq -r '.tenant.id // empty' 2>/dev/null)

if [ -z "$TENANT_ID" ] || [ "$TENANT_ID" = "empty" ]; then
    echo "❌ Failed to create account"
    exit 1
fi

echo "✅ Account created (Tenant ID: $TENANT_ID)"
echo ""

# Step 2: Wait for verification (in production)
echo "Step 2: Waiting for email verification..."
echo ""

echo "⏰ Monitoring for email verification (checking every 5 seconds for up to 2 minutes)..."
for i in {1..24}; do
    sleep 5
    echo "Check #$i: $(date)"

    STATUS_RESPONSE=$(curl -s -X GET "https://driftlock.net/api/v1/onboard/status?tenant_id=$TENANT_ID" \
      -H "Content-Type: application/json")

    if echo "$STATUS_RESPONSE" | jq -e '.pending_verification' | grep -q "false"; then
        echo ""
        echo "✅ Email verified!"

        # Get API key
        API_KEY=$(echo "$STATUS_RESPONSE" | jq -r '.api_key // "Not found"' 2>/dev/null)

        if [ "$API_KEY" != "Not found" ] && [ -n "$API_KEY" ]; then
            echo "✅ API Key obtained: ${API_KEY:0:16}..."

            # Step 3: Test API key works
            echo ""
            echo "Step 3: Testing API key..."
            DETECT_RESPONSE=$(curl -s -X POST https://driftlock.net/api/v1/detect \
              -H "Authorization: Bearer $API_KEY" \
              -H "Content-Type: application/json" \
              -d '{"events": ["test event 1", "test event 2", "anomaly detected"]}')

            if echo "$DETECT_RESPONSE" | grep -q "success"; then
                echo "✅ API key working"

                # Step 4: Get billing status
                echo ""
                echo "Step 4: Checking billing status..."
                BILLING_RESPONSE=$(curl -s -X GET "https://driftlock.net/api/v1/me/billing" \
                  -H "Authorization: Bearer $API_KEY" \
                  -H "Content-Type: application/json")

                echo "Current billing status:"
                echo "$BILLING_RESPONSE" | jq '.' 2>/dev/null

                # Step 5: Create checkout session
                echo ""
                echo "Step 5: Creating $PLAN checkout session..."
                CHECKOUT_RESPONSE=$(curl -s -X POST "https://driftlock.net/api/v1/billing/checkout" \
                  -H "Authorization: Bearer $API_KEY" \
                  -H "Content-Type: application/json" \
                  -d "{\"plan\":\"$PLAN\"}")

                echo "Checkout session created:"
                echo "$CHECKOUT_RESPONSE" | jq '.' 2>/dev/null

                # Extract checkout URL
                CHECKOUT_URL=$(echo "$CHECKOUT_RESPONSE" | jq -r '.checkout_url // empty' 2/dev/null)
                STRIPE_SESSION_ID=$(echo "$CHECKOUT_RESPONSE" | jq -r '.stripe_session_id // empty' 2>/dev/null)

                if [ -n "$CHECKOUT_URL" ] && [ "$CHECKOUT_URL" != "empty" ]; then
                    echo ""
                    echo "✅ Checkout session created!"
                    echo ""
                    echo "🔗 Checkout URL: $CHECKOUT_URL"
                    echo ""
                    echo "Step 6: Manual Stripe Checkout Required"
                    echo "Please:"
                    echo "1. Open the checkout URL in your browser"
                    echo "2. Login with test card: 4242 4242 4242 4242"
                    echo "   - Name: Test User"
                    "   - Email: $EMAIL"
                    "   - Billing address: Any US address"
                    echo "3. Complete payment"
                    echo ""
                    echo "After payment, the system will:"
                    echo "- Create a Stripe subscription"
                    echo "- Update your tenant record"
                    echo "- Send a confirmation email"
                    echo ""

                    # Monitor for subscription activation
                    echo "Monitoring for subscription activation (checking every 5 seconds for up to 2 minutes)..."
                    for j in {1..24}; do
                        sleep 5
                        echo "Check #$j: $(date)"

                        BILLING_UPDATE=$(curl -s -X GET "https://driftlock.net/api/v1/me/billing" \
                          -H "Authorization: Bearer $API_KEY" \
                          -H "Content-Type: application/json")

                        # Check if subscription is active
                        if echo "$BILLING_UPDATE" | jq -e '.plan' | grep -q "$PLAN"; then
                            echo ""
                            echo "🎉 Subscription activated!"
                            echo ""
                            echo "Plan: $PLAN"
                            echo "Billing response:"
                            echo "$BILLING_UPDATE" | jq '. | {plan, status, stripe_customer_id}' 2>/dev/null

                            # Step 7: Verify API access with paid plan
                            echo ""
                            echo "Step 7: Verifying paid plan API access..."
                            PAID_DETECT_RESPONSE=$(curl -s -X POST https://driftlock.net/api/v1/detect \
                              -H "Authorization: Bearer $API_KEY" \
                              -H "Content-Type: application/json" \
                              -d '{"events": ["paid plan test 1", "paid plan test 2", "another event"]}')

                            if echo "$PAID_DETECT_RESPONSE" | grep -q "success"; then
                                echo "✅ Paid plan API access working!"
                                echo ""
                                echo "=== CHECKOUT FLOW COMPLETE ==="
                                echo "✅ Account created and verified"
                                echo "✅ API key generated"
                                echo "✅ Trial plan active"
                                echo "✅ Checkout session created"
                                echo "✅ Payment completed"
                                echo "✅ $PLAN subscription active"
                                echo "✅ Paid plan API access working"
                                echo ""
                                echo "Test successful! 🚀"
                                exit 0
                            else
                                echo "⚠️  Subscription active but API test failed"
                                echo "Response: $PAID_DETECT_RESPONSE"
                            fi
                            break
                        else
                            echo "  Still waiting for subscription activation..."
                        fi
                    done

                    echo ""
                    echo "⏰ Timeout - subscription not activated after 2 minutes"
                    echo "Please check:"
                    echo "- Stripe dashboard for payment status"
                    echo "- Driftlock billing status: https://driftlock.net/api/v1/me/billing"
                else
                    echo "❌ Failed to create checkout session"
                    echo "Response: $CHECKOUT_RESPONSE"
                fi
            else
                echo "❌ API key test failed"
                echo "Response: $DETECT_RESPONSE"
            fi
        else
            echo "⚠️  Account verified but no API key found"
        fi
        break
    else
        echo "  Still waiting for verification..."
    fi
done

echo ""
echo "⏰ Timeout - email not verified after 2 minutes"
echo "Please check your email (including spam folder) and try again"
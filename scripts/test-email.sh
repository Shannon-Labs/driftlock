#!/bin/bash
# Test signup flow - triggers verification email
# Usage: ./test-email.sh [email@example.com]
#
# This script tests the email verification flow:
# 1. Signs up a new user
# 2. Checks Cloud Run logs for email sending
#
# Prerequisites:
# - Production API must be deployed
# - SENDGRID_API_KEY must be configured in production

set -e

API_URL="${API_URL:-https://driftlock.net}"
EMAIL="${1:-test-$(date +%s)@example.com}"
COMPANY="CLI Test $(date +%Y%m%d-%H%M%S)"

echo "Testing email flow..."
echo "API URL: $API_URL"
echo "Email: $EMAIL"
echo "Company: $COMPANY"
echo ""

# Step 1: Sign up
echo "Step 1: Signing up..."
RESPONSE=$(curl -s -X POST "$API_URL/api/v1/onboard/signup" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"company_name\":\"$COMPANY\",\"plan\":\"trial\"}")

echo "Response: $RESPONSE"
echo ""

# Check for success
if echo "$RESPONSE" | grep -q '"success":true'; then
  echo "Signup successful! Check your inbox for verification email."
  echo ""
  echo "Step 2: Check Cloud Run logs for email activity:"
  echo "gcloud logging read 'resource.type=cloud_run_revision AND textPayload:verification' --limit=20 --format='value(textPayload)'"
else
  echo "Signup may have failed. Check response above."
  exit 1
fi

echo ""
echo "Manual verification steps:"
echo "1. Check email inbox for verification link"
echo "2. Click link to verify account"
echo "3. Note the API key from the response"

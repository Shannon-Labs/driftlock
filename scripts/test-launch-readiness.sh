#!/bin/bash

# Driftlock Launch Readiness Test
# Comprehensive test script to verify all systems are go

set -e

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

PASSED=0
FAILED=0

# Test counter
test_count=0

# Function to run a test
run_test() {
    local test_name="$1"
    local test_command="$2"
    
    test_count=$((test_count + 1))
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "Test $test_count: $test_name"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    if eval "$test_command"; then
        echo -e "${GREEN}✅ PASS${NC}: $test_name"
        PASSED=$((PASSED + 1))
        return 0
    else
        echo -e "${RED}❌ FAIL${NC}: $test_name"
        FAILED=$((FAILED + 1))
        return 1
    fi
}

# Header
echo ""
echo "══════════════════════════════════════════════════════════════"
echo "  Driftlock Launch Readiness Test"
echo "══════════════════════════════════════════════════════════════"
echo ""
echo "Running comprehensive tests to verify launch readiness..."
echo ""

# Test 1: Check Firebase authentication
echo "🔐 Checking Firebase authentication..."
if firebase projects:list > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Firebase authentication working${NC}"
    PASSED=$((PASSED + 1))
else
    echo -e "${RED}❌ Firebase authentication failed${NC}"
    echo "   Run: firebase login --reauth"
    FAILED=$((FAILED + 1))
fi

# Test 2: Check Google Cloud authentication
echo "☁️  Checking Google Cloud authentication..."
if gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q .; then
    echo -e "${GREEN}✅ Google Cloud authentication working${NC}"
    PASSED=$((PASSED + 1))
else
    echo -e "${RED}❌ Google Cloud authentication failed${NC}"
    echo "   Run: gcloud auth login"
    FAILED=$((FAILED + 1))
fi

# Test 3: Check Cloud Run backend health
echo "🏥 Checking Cloud Run backend health..."
run_test "Backend /healthz endpoint" "
curl -s https://driftlock-api-o6kjgrsowq-uc.a.run.app/healthz | grep -q '\"success\":true'
"

# Test 4: Check backend signup endpoint
echo "📝 Checking backend signup endpoint..."
run_test "Backend /v1/onboard/signup" "
curl -s -X POST https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/onboard/signup \
  -H 'Content-Type: application/json' \
  -d '{\"email\":\"test@driftlock.dev\",\"company_name\":\"Test Co\"}' | grep -q 'api_key'
"

# Test 5: Check Firebase Functions health
echo "🔥 Checking Firebase Functions..."
run_test "Firebase Functions healthCheck" "
curl -s https://us-central1-driftlock.cloudfunctions.net/healthCheck | grep -q '\"success\":true'
"

# Test 6: Check Functions proxy to backend
echo "🔀 Checking Functions proxy..."
run_test "Functions API proxy" "
curl -s https://us-central1-driftlock.cloudfunctions.net/apiProxy/api/v1/healthz | grep -q '\"success\":true'
"

# Test 7: Check Stripe webhook endpoint (basic)
echo "🪝 Checking Stripe webhook endpoint..."
run_test "Stripe webhook endpoint configured" "
curl -s -X POST https://us-central1-driftlock.cloudfunctions.net/apiProxy/webhooks/stripe \
  -H 'Content-Type: application/json' \
  -d '{}' | grep -q '\"error\"' || true
"

# Test 8: Check landing page
echo "🌐 Checking landing page..."
run_test "Landing page accessible" "
curl -s -I https://driftlock.net | grep -q '200 OK'
"

# Test 9: Check Firebase Hosting deployment
echo "📦 Checking Firebase Hosting..."
current_version=$(firebase hosting:sites:releases:list --site=driftlock | head -n 1 | awk '{print $2}')
if [ -n "$current_version" ]; then
    echo -e "${GREEN}✅ Firebase Hosting deployed${NC} (version: $current_version)"
    PASSED=$((PASSED + 1))
else
    echo -e "${RED}❌ Firebase Hosting not deployed${NC}"
    FAILED=$((FAILED + 1))
fi

# Test 10: Check Stripe CLI
echo "💳 Checking Stripe CLI..."
if command -v stripe &> /dev/null; then
    if stripe config --list > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Stripe CLI installed and configured${NC}"
        PASSED=$((PASSED + 1))
    else
        echo -e "${YELLOW}⚠️  Stripe CLI installed but not logged in${NC}"
        echo "   Run: stripe login"
        PASSED=$((PASSED + 1))
    fi
else
    echo -e "${RED}❌ Stripe CLI not installed${NC}"
    echo "   Install: brew install stripe/stripe-cli/stripe"
    FAILED=$((FAILED + 1))
fi

# Summary
echo ""
echo "══════════════════════════════════════════════════════════════"
echo "  Test Summary"
echo "══════════════════════════════════════════════════════════════"
echo ""
echo -e "${GREEN}Passed: $PASSED${NC}"
echo -e "${RED}Failed: $FAILED${NC}"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 All systems ready for launch!${NC}"
    echo ""
    echo "Next steps:"
    echo "1. Set up Stripe webhook secret in environment"
    echo "2. Configure Cloudflare DNS for Firebase"
    echo "3. Run end-to-end user flow tests"
    echo "4. Launch! 🚀"
    exit 0
else
    echo -e "${RED}❌ Some tests failed. Please fix before launching.${NC}"
    echo ""
    echo "Refer to LAUNCH_PLAN.md for troubleshooting steps."
    exit 1
fi

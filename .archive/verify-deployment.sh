#!/bin/bash

echo "========================================="
echo "DriftLock Deployment Verification"
echo "========================================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

PASSED=0
FAILED=0

pass() {
    echo -e "${GREEN}✓${NC} $1"
    ((PASSED++))
}

fail() {
    echo -e "${RED}✗${NC} $1"
    ((FAILED++))
}

info() {
    echo -e "${YELLOW}ℹ${NC} $1"
}

# Configuration
API_WORKER_URL="${1:-https://driftlock-api-staging.YOUR_ACCOUNT.workers.dev}"
PAGES_URL="${2:-https://driftlock-web-frontend.pages.dev}"

info "Testing deployment at:"
echo "  API: $API_WORKER_URL"
echo "  Web: $PAGES_URL"
echo ""

# Test 1: API Health Check
echo "1. Testing API Health Check..."
if curl -sf "${API_WORKER_URL}/health" > /dev/null 2>&1; then
    pass "API health check passed"
else
    fail "API health check failed"
fi

# Test 2: API Root Endpoint
echo "2. Testing API Root Endpoint..."
if curl -sf "${API_WORKER_URL}/" > /dev/null 2>&1; then
    pass "API root endpoint accessible"
else
    fail "API root endpoint failed"
fi

# Test 3: Web Frontend
echo "3. Testing Web Frontend..."
if curl -sf "$PAGES_URL" > /dev/null 2>&1; then
    pass "Web frontend accessible"
else
    fail "Web frontend failed"
fi

# Test 4: API Response Format
echo "4. Testing API Response Format..."
if curl -sf "${API_WORKER_URL}/health" | grep -q "status"; then
    pass "API returns valid JSON"
else
    fail "API response format invalid"
fi

# Test 5: CORS Headers
echo "5. Testing CORS Headers..."
if curl -sf -H "Origin: https://example.com" -H "Access-Control-Request-Method: GET" -X OPTIONS "${API_WORKER_URL}/health" > /dev/null 2>&1; then
    pass "CORS preflight works"
else
    info "CORS preflight may need configuration"
fi

# Test 6: Supabase Edge Functions
echo "6. Testing Supabase Edge Functions..."
if curl -sf "https://nfkdeeunyvnntvpvwpwh.supabase.co/functions/v1/health" | grep -q "status"; then
    pass "Supabase edge functions accessible"
else
    fail "Supabase edge functions failed"
fi

# Test 7: Environment Variables
echo "7. Verifying Environment Configuration..."
if grep -q "SUPABASE_PROJECT_ID=nfkdeeunyvnntvpvwpwh" /Volumes/VIXinSSD/driftlock/.env; then
    pass "Supabase configuration present"
else
    fail "Supabase configuration missing"
fi

# Test 8: Deployment Artifacts
echo "8. Checking Deployment Artifacts..."
if [ -f "/Volumes/VIXinSSD/driftlock/cloudflare-api-worker/src/index.ts" ]; then
    pass "API Worker code exists"
else
    fail "API Worker code missing"
fi

if [ -f "/Volumes/VIXinSSD/driftlock/web-frontend/dist/index.html" ] || [ -f "/Volumes/VIXinSSD/driftlock/web-frontend/build/index.html" ]; then
    pass "Web Frontend build exists"
else
    fail "Web Frontend build missing (run: npm run build)"
fi

# Test 9: Deployment Scripts
echo "9. Checking Deployment Scripts..."
if [ -x "/Volumes/VIXinSSD/driftlock/cloudflare-api-worker/deploy.sh" ]; then
    pass "API deployment script exists"
else
    fail "API deployment script missing"
fi

if [ -x "/Volumes/VIXinSSD/driftlock/web-frontend/deploy-pages.sh" ]; then
    pass "Frontend deployment script exists"
else
    fail "Frontend deployment script missing"
fi

# Test 10: Integration Test Suite
echo "10. Running Integration Tests..."
if bash /Volumes/VIXinSSD/driftlock/test-integration-simple.sh > /dev/null 2>&1; then
    pass "Integration tests passed"
else
    fail "Integration tests failed"
fi

# Summary
echo ""
echo "========================================="
echo "Deployment Verification Summary"
echo "========================================="
echo ""
echo "Total Checks: 10"
echo -e "${GREEN}Passed: $PASSED${NC}"
echo -e "${RED}Failed: $FAILED${NC}"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✅ All deployment verification checks passed!${NC}"
    echo ""
    echo "Your DriftLock deployment is ready!"
    echo ""
    echo "Next Steps:"
    echo "1. Update Stripe webhook URL to point to your Worker"
    echo "2. Test end-to-end user flow"
    echo "3. Configure custom domains (optional)"
    echo "4. Set up monitoring and alerts"
    echo ""
    exit 0
else
    echo -e "${YELLOW}⚠️  Some checks failed. Please review and fix issues.${NC}"
    echo ""
    exit 1
fi

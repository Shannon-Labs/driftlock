#!/bin/bash

echo "========================================="
echo "DriftLock Integration Test Suite"
echo "========================================="
echo ""

# Test counters
PASSED=0
FAILED=0

# Helper functions
pass() {
    echo -e "\033[0;32m✓\033[0m $1"
    PASSED=$((PASSED + 1))
}

fail() {
    echo -e "\033[0;31m✗\033[0m $1"
    FAILED=$((FAILED + 1))
}

section() {
    echo ""
    echo "========================================="
    echo "$1"
    echo "========================================="
}

# Section 1: File Structure Validation
section "1. File Structure Validation"

if [ -d "web-frontend/src" ]; then
    pass "web-frontend/src directory exists"
else
    fail "web-frontend/src directory missing"
fi

if [ -d "web-frontend/supabase/functions" ]; then
    pass "web-frontend/supabase/functions directory exists"
else
    fail "web-frontend/supabase/functions directory missing"
fi

if [ -d "web-frontend/supabase/migrations" ]; then
    pass "web-frontend/supabase/migrations directory exists"
else
    fail "web-frontend/supabase/migrations directory missing"
fi

# Check edge functions
FUNCTIONS=("health" "meter-usage" "send-alert-email" "stripe-webhook")
for func in "${FUNCTIONS[@]}"; do
    if [ -d "web-frontend/supabase/functions/$func" ]; then
        pass "Edge function '$func' exists"
    else
        fail "Edge function '$func' missing"
    fi
done

# Section 2: Configuration Files
section "2. Configuration Files"

if [ -f ".env" ]; then
    pass ".env file exists"

    if grep -q "SUPABASE_PROJECT_ID" .env && grep -q "SUPABASE_ANON_KEY" .env && grep -q "SUPABASE_SERVICE_ROLE_KEY" .env; then
        pass "Supabase configuration present in .env"
    else
        fail "Supabase configuration incomplete in .env"
    fi
else
    fail ".env file missing"
fi

if [ -f "web-frontend/.env" ]; then
    pass "web-frontend/.env file exists"
else
    fail "web-frontend/.env file missing"
fi

if [ -f "docker-compose.yml" ]; then
    pass "docker-compose.yml exists"

    if grep -q "web-frontend:" docker-compose.yml; then
        pass "web-frontend service defined in docker-compose.yml"
    else
        fail "web-frontend service missing in docker-compose.yml"
    fi

    if grep -q "api:" docker-compose.yml; then
        pass "api service defined in docker-compose.yml"
    else
        fail "api service missing in docker-compose.yml"
    fi
else
    fail "docker-compose.yml missing"
fi

if [ -f "web-frontend/package.json" ]; then
    pass "web-frontend/package.json exists"

    if grep -q "@supabase" web-frontend/package.json; then
        pass "Supabase dependency present (Stripe handled via edge functions)"
    else
        fail "Missing Supabase dependency"
    fi
else
    fail "web-frontend/package.json missing"
fi

# Section 3: API Server Integration
section "3. API Server Integration"

if [ -f "api-server/internal/supabase/client.go" ]; then
    pass "Supabase client exists in API server"

    if grep -q "CreateAnomaly" api-server/internal/supabase/client.go; then
        pass "CreateAnomaly function present"
    else
        fail "CreateAnomaly function missing"
    fi

    if grep -q "CreateUsageRecord" api-server/internal/supabase/client.go; then
        pass "CreateUsageRecord function present"
    else
        fail "CreateUsageRecord function missing"
    fi

    if grep -q "NotifyWebhook\|SendWebhook" api-server/internal/supabase/client.go; then
        pass "NotifyWebhook/SendWebhook function present"
    else
        fail "NotifyWebhook/SendWebhook function missing"
    fi
else
    fail "Supabase client missing from API server"
fi

# Section 4: Edge Function Source Files
section "4. Edge Function Source Files"

for func in "${FUNCTIONS[@]}"; do
    if [ -f "web-frontend/supabase/functions/$func/index.ts" ]; then
        pass "Edge function '$func/index.ts' exists"
    else
        fail "Edge function '$func/index.ts' missing"
    fi
done

# Section 5: Database Migration Files
section "5. Database Migration Files"

MIGRATION_COUNT=$(ls -1 web-frontend/supabase/migrations/*.sql 2>/dev/null | wc -l)
if [ "$MIGRATION_COUNT" -gt 0 ]; then
    pass "Found $MIGRATION_COUNT migration files"
else
    fail "No migration files found"
fi

# Section 6: Documentation
section "6. Documentation"

if [ -f "web-frontend/IMPLEMENTATION_COMPLETE.md" ]; then
    pass "web-frontend/IMPLEMENTATION_COMPLETE.md exists"
else
    fail "web-frontend/IMPLEMENTATION_COMPLETE.md missing"
fi

if [ -f "web-frontend/PRODUCTION_RUNBOOK.md" ]; then
    pass "web-frontend/PRODUCTION_RUNBOOK.md exists"
else
    fail "web-frontend/PRODUCTION_RUNBOOK.md missing"
fi

# Summary
section "Test Summary"

TOTAL=$((PASSED + FAILED))
echo ""
echo "Total Tests: $TOTAL"
echo -e "\033[0;32mPassed: $PASSED\033[0m"
echo -e "\033[0;31mFailed: $FAILED\033[0m"

if [ $FAILED -eq 0 ]; then
    echo ""
    echo -e "\033[0;32mAll integration tests passed! ✅\033[0m"
    echo ""
    echo "Next steps:"
    echo "1. Run: docker-compose -f docker-compose.yml up -d"
    echo "2. Visit: http://localhost:3000 (web-frontend)"
    echo "3. Visit: http://localhost:8080/healthz (API)"
else
    echo ""
    echo -e "\033[1;33mSome tests failed. Please review and fix issues. ⚠️\033[0m"
fi

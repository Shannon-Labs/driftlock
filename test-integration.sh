#!/bin/bash
set -e

echo "========================================="
echo "DriftLock Integration Test Suite"
echo "========================================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counters
PASSED=0
FAILED=0
TOTAL=0

# Helper functions
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

section() {
    echo ""
    echo "========================================="
    echo "$1"
    echo "========================================="
}

# Section 1: File Structure Validation
section "1. File Structure Validation"

TOTAL=$((TOTAL + 1))
if [ -d "web-frontend/src" ]; then
    pass "web-frontend/src directory exists"
else
    fail "web-frontend/src directory missing"
fi

TOTAL=$((TOTAL + 1))
if [ -d "web-frontend/supabase/functions" ]; then
    pass "web-frontend/supabase/functions directory exists"
else
    fail "web-frontend/supabase/functions directory missing"
fi

TOTAL=$((TOTAL + 1))
if [ -d "web-frontend/supabase/migrations" ]; then
    pass "web-frontend/supabase/migrations directory exists"
else
    fail "web-frontend/supabase/migrations directory missing"
fi

# Check edge functions
FUNCTIONS=("health" "meter-usage" "send-alert-email" "stripe-webhook")
for func in "${FUNCTIONS[@]}"; do
    TOTAL=$((TOTAL + 1))
    if [ -d "web-frontend/supabase/functions/$func" ]; then
        pass "Edge function '$func' exists"
    else
        fail "Edge function '$func' missing"
    fi
done

# Section 2: Configuration Files
section "2. Configuration Files"

TOTAL=$((TOTAL + 1))
if [ -f ".env" ]; then
    pass ".env file exists"

    # Check required Supabase config
    TOTAL=$((TOTAL + 1))
    if grep -q "SUPABASE_PROJECT_ID" .env && grep -q "SUPABASE_ANON_KEY" .env && grep -q "SUPABASE_SERVICE_ROLE_KEY" .env; then
        pass "Supabase configuration present in .env"
    else
        fail "Supabase configuration incomplete in .env"
    fi
else
    fail ".env file missing"
fi

TOTAL=$((TOTAL + 1))
if [ -f "web-frontend/.env" ]; then
    pass "web-frontend/.env file exists"
else
    fail "web-frontend/.env file missing"
fi

TOTAL=$((TOTAL + 1))
if [ -f "docker-compose.yml" ]; then
    pass "docker-compose.yml exists"

    # Check for required services
    TOTAL=$((TOTAL + 1))
    if grep -q "web-frontend:" docker-compose.yml; then
        pass "web-frontend service defined in docker-compose.yml"
    else
        fail "web-frontend service missing in docker-compose.yml"
    fi

    TOTAL=$((TOTAL + 1))
    if grep -q "api:" docker-compose.yml; then
        pass "api service defined in docker-compose.yml"
    else
        fail "api service missing in docker-compose.yml"
    fi
else
    fail "docker-compose.yml missing"
fi

TOTAL=$((TOTAL + 1))
if [ -f "web-frontend/package.json" ]; then
    pass "web-frontend/package.json exists"

    # Check key dependencies
    TOTAL=$((TOTAL + 1))
    if grep -q "@supabase" web-frontend/package.json && grep -q "@stripe" web-frontend/package.json; then
        pass "Supabase and Stripe dependencies present"
    else
        fail "Missing Supabase or Stripe dependencies"
    fi
else
    fail "web-frontend/package.json missing"
fi

# Section 3: API Server Integration
section "3. API Server Integration"

TOTAL=$((TOTAL + 1))
if [ -f "api-server/internal/supabase/client.go" ]; then
    pass "Supabase client exists in API server"

    # Check for key functions
    TOTAL=$((TOTAL + 1))
    if grep -q "CreateAnomaly" api-server/internal/supabase/client.go; then
        pass "CreateAnomaly function present"
    else
        fail "CreateAnomaly function missing"
    fi

    TOTAL=$((TOTAL + 1))
    if grep -q "CreateUsageRecord" api-server/internal/supabase/client.go; then
        pass "CreateUsageRecord function present"
    else
        fail "CreateUsageRecord function missing"
    fi

    TOTAL=$((TOTAL + 1))
    if grep -q "SendWebhook" api-server/internal/supabase/client.go; then
        pass "SendWebhook function present"
    else
        fail "SendWebhook function missing"
    fi
else
    fail "Supabase client missing from API server"
fi

# Section 4: Build Tests
section "4. Build Tests"

info "Testing web-frontend build..."
if cd web-frontend && npm run build > /dev/null 2>&1; then
    pass "web-frontend builds successfully"
    ((TOTAL++))
else
    fail "web-frontend build failed"
    ((TOTAL++))
fi
cd ..

info "Testing API server build..."
if cd api-server && go build -o /tmp/driftlock-api ./cmd/api-server > /dev/null 2>&1; then
    pass "API server builds successfully"
    ((TOTAL++))
else
    fail "API server build failed"
    ((TOTAL++))
fi
cd ..

# Section 5: Docker Compose Configuration
section "5. Docker Compose Configuration"

TOTAL=$((TOTAL + 1))
if docker-compose -f docker-compose.yml config > /dev/null 2>&1; then
    pass "docker-compose.yml is valid"
else
    fail "docker-compose.yml has errors"
fi

# Section 6: Supabase Edge Function Files
section "6. Edge Function Source Files"

for func in "${FUNCTIONS[@]}"; do
    TOTAL=$((TOTAL + 1))
    if [ -f "web-frontend/supabase/functions/$func/index.ts" ]; then
        pass "Edge function '$func/index.ts' exists"
    else
        fail "Edge function '$func/index.ts' missing"
    fi
done

# Section 7: Database Migration Files
section "7. Database Migration Files"

MIGRATION_COUNT=$(ls -1 web-frontend/supabase/migrations/*.sql 2>/dev/null | wc -l)
TOTAL=$((TOTAL + 1))
if [ "$MIGRATION_COUNT" -gt 0 ]; then
    pass "Found $MIGRATION_COUNT migration files"
else
    fail "No migration files found"
fi

# Section 8: Documentation
section "8. Documentation"

DOCS=("web-frontend/IMPLEMENTATION_COMPLETE.md" "web-frontend/PRODUCTION_RUNBOOK.md")
for doc in "${DOCS[@]}"; do
    TOTAL=$((TOTAL + 1))
    if [ -f "$doc" ]; then
        pass "Documentation '$doc' exists"
    else
        fail "Documentation '$doc' missing"
    fi
done

# Summary
section "Test Summary"

TOTAL=$((TOTAL + PASSED + FAILED))
echo ""
echo "Total Tests: $TOTAL"
echo -e "${GREEN}Passed: $PASSED${NC}"
echo -e "${RED}Failed: $FAILED${NC}"

if [ $FAILED -eq 0 ]; then
    echo ""
    echo -e "${GREEN}All integration tests passed! ✅${NC}"
    echo ""
    echo "Next steps:"
    echo "1. Run: docker-compose -f docker-compose.yml up -d"
    echo "2. Visit: http://localhost:3000 (web-frontend)"
    echo "3. Visit: http://localhost:8080/healthz (API)"
    exit 0
else
    echo ""
    echo -e "${YELLOW}Some tests failed. Please review and fix issues. ⚠️${NC}"
    exit 1
fi

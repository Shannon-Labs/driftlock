#!/bin/bash

# DriftLock API Integration Test Script
# Tests all major endpoints and validates functionality

set -e

# Configuration
API_BASE_URL="${API_BASE_URL:-http://localhost:8080}"
TEST_EMAIL="test-$(date +%s)@driftlock.test"
TEST_PASSWORD="TestPassword123!"

echo "========================================"
echo "DriftLock API Integration Test"
echo "========================================"
echo "API URL: $API_BASE_URL"
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Helper function to test endpoint
test_endpoint() {
    local name=$1
    local method=$2
    local endpoint=$3
    local data=$4
    local expected_status=$5
    local headers=$6

    echo -n "Testing $name... "

    if [ -n "$headers" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$API_BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -H "$headers" \
            -d "$data" 2>&1)
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$API_BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data" 2>&1)
    fi

    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)

    if [ "$http_code" -eq "$expected_status" ]; then
        echo -e "${GREEN}✓ PASS${NC} (HTTP $http_code)"
        TESTS_PASSED=$((TESTS_PASSED + 1))
        echo "$body"
        return 0
    else
        echo -e "${RED}✗ FAIL${NC} (Expected HTTP $expected_status, got $http_code)"
        echo "Response: $body"
        TESTS_FAILED=$((TESTS_FAILED + 1))
        return 1
    fi
}

echo "========== Health Checks =========="
test_endpoint "Health Check" "GET" "/health" "" 200

echo ""
echo "========== Authentication Tests =========="

# Register user
echo "Registering test user..."
register_response=$(curl -s -X POST "$API_BASE_URL/api/v1/auth/register" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\",\"name\":\"Test User\"}")

echo "Register response: $register_response"

# Login user
echo "Logging in test user..."
login_response=$(curl -s -X POST "$API_BASE_URL/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\"}")

echo "Login response: $login_response"

# Extract JWT token
JWT_TOKEN=$(echo "$login_response" | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$JWT_TOKEN" ]; then
    echo -e "${RED}Failed to get JWT token. Cannot continue with authenticated tests.${NC}"
    echo "Login response was: $login_response"
else
    echo -e "${GREEN}✓ Successfully obtained JWT token${NC}"
    echo "Token: ${JWT_TOKEN:0:20}..."

    echo ""
    echo "========== User Management Tests =========="
    test_endpoint "Get User Profile" "GET" "/api/v1/user" "" 200 "Authorization: Bearer $JWT_TOKEN"

    echo ""
    echo "========== Dashboard Tests =========="
    test_endpoint "Get Dashboard Stats" "GET" "/api/v1/dashboard/stats" "" 200 "Authorization: Bearer $JWT_TOKEN"
    test_endpoint "Get Recent Anomalies" "GET" "/api/v1/dashboard/recent" "" 200 "Authorization: Bearer $JWT_TOKEN"

    echo ""
    echo "========== Event Ingestion Tests =========="

    # Ingest test event
    event_data='[{
        "timestamp": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'",
        "stream_type": "logs",
        "data": "Test log message for anomaly detection",
        "metadata": {
            "source": "test-script",
            "severity": "info"
        }
    }]'

    test_endpoint "Ingest Events" "POST" "/api/v1/events/ingest" "$event_data" 200 "Authorization: Bearer $JWT_TOKEN"

    echo ""
    echo "========== Anomaly Tests =========="
    test_endpoint "List Anomalies" "GET" "/api/v1/anomalies" "" 200 "Authorization: Bearer $JWT_TOKEN"
    test_endpoint "Get Events" "GET" "/api/v1/events" "" 200 "Authorization: Bearer $JWT_TOKEN"

    echo ""
    echo "========== Billing Tests =========="
    test_endpoint "List Billing Plans" "GET" "/api/v1/billing/plans" "" 200 "Authorization: Bearer $JWT_TOKEN"
    test_endpoint "Get Current Subscription" "GET" "/api/v1/billing/subscription" "" 200 "Authorization: Bearer $JWT_TOKEN"
    test_endpoint "Get Usage" "GET" "/api/v1/billing/usage" "" 200 "Authorization: Bearer $JWT_TOKEN"

    echo ""
    echo "========== Onboarding Tests =========="
    test_endpoint "Get Onboarding Progress" "GET" "/api/v1/onboarding/progress" "" 200 "Authorization: Bearer $JWT_TOKEN"
    test_endpoint "Get Onboarding Resources" "GET" "/api/v1/onboarding/resources" "" 200 "Authorization: Bearer $JWT_TOKEN"
fi

echo ""
echo "========================================"
echo "Test Summary"
echo "========================================"
echo -e "Tests Passed: ${GREEN}$TESTS_PASSED${NC}"
echo -e "Tests Failed: ${RED}$TESTS_FAILED${NC}"
echo "========================================"

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed.${NC}"
    exit 1
fi

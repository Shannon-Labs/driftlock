#!/bin/bash
# Test Driftlock API endpoints with various inputs

set -e

echo "🧪 Testing Driftlock API endpoints..."
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Track results
FAILED=0
PASSED=0

API_URL="${API_URL:-http://localhost:8080}"

# Function to test API endpoint
test_endpoint() {
    local name=$1
    local method=$2
    local endpoint=$3
    local data=$4
    local expected_status=${5:-200}
    
    echo -n "Testing $name... "
    
    if [ -n "$data" ]; then
        RESPONSE=$(curl -s -w "\n%{http_code}" -X "$method" \
            -H "Content-Type: application/json" \
            --data-binary "$data" \
            "$API_URL$endpoint" 2>&1)
    else
        RESPONSE=$(curl -s -w "\n%{http_code}" -X "$method" \
            "$API_URL$endpoint" 2>&1)
    fi
    
    HTTP_CODE=$(echo "$RESPONSE" | tail -1)
    BODY=$(echo "$RESPONSE" | sed '$d')
    
    if [ "$HTTP_CODE" -eq "$expected_status" ]; then
        echo -e "${GREEN}✓ PASSED${NC} (HTTP $HTTP_CODE)"
        PASSED=$((PASSED + 1))
        return 0
    else
        echo -e "${RED}✗ FAILED${NC} (HTTP $HTTP_CODE, expected $expected_status)"
        echo "  Response: $BODY"
        FAILED=$((FAILED + 1))
        return 1
    fi
}

# Ensure service is running
echo "=== Checking API availability ==="
if ! curl -sf "$API_URL/healthz" > /dev/null 2>&1; then
    echo -e "${RED}API is not available at $API_URL${NC}"
    echo "Please start the service: docker compose up -d driftlock-http"
    exit 1
fi
echo -e "${GREEN}API is available${NC}"
echo ""

# Test health endpoint
echo "=== Testing /healthz endpoint ==="
test_endpoint "GET /healthz" "GET" "/healthz" "" 200

# Verify health response structure
HEALTH_RESPONSE=$(curl -s "$API_URL/healthz")
if echo "$HEALTH_RESPONSE" | jq -e '.success == true' > /dev/null 2>&1; then
    echo -e "  ${GREEN}✓ Health response structure valid${NC}"
    PASSED=$((PASSED + 1))
else
    echo -e "  ${RED}✗ Health response structure invalid${NC}"
    FAILED=$((FAILED + 1))
fi

# Test metrics endpoint
echo ""
echo "=== Testing /metrics endpoint ==="
test_endpoint "GET /metrics" "GET" "/metrics" "" 200

# Test /v1/detect with NDJSON
echo ""
echo "=== Testing /v1/detect endpoint ==="

# Test with small NDJSON
TEST_DATA='{"timestamp":"2025-01-01T00:00:00Z","value":100}
{"timestamp":"2025-01-01T00:00:01Z","value":101}
{"timestamp":"2025-01-01T00:00:02Z","value":102}'

test_endpoint "POST /v1/detect (NDJSON)" "POST" "/v1/detect?format=ndjson" "$TEST_DATA" 200

# Verify response structure
DETECT_RESPONSE=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    --data-binary "$TEST_DATA" \
    "$API_URL/v1/detect?format=ndjson")

if echo "$DETECT_RESPONSE" | jq -e '.success == true' > /dev/null 2>&1 && \
   echo "$DETECT_RESPONSE" | jq -e '.request_id' > /dev/null 2>&1; then
    echo -e "  ${GREEN}✓ Response structure valid${NC}"
    PASSED=$((PASSED + 1))
else
    echo -e "  ${RED}✗ Response structure invalid${NC}"
    FAILED=$((FAILED + 1))
fi

# Test with JSON array
TEST_JSON='[{"timestamp":"2025-01-01T00:00:00Z","value":100},{"timestamp":"2025-01-01T00:00:01Z","value":101}]'
test_endpoint "POST /v1/detect (JSON array)" "POST" "/v1/detect?format=json" "$TEST_JSON" 200

# Test auto-detection (no format parameter)
test_endpoint "POST /v1/detect (auto-detect NDJSON)" "POST" "/v1/detect" "$TEST_DATA" 200

# Test with query parameters
test_endpoint "POST /v1/detect (with params)" "POST" "/v1/detect?baseline=10&window=1&hop=1&algo=zstd" "$TEST_DATA" 200

# Test error handling
echo ""
echo "=== Testing error handling ==="

# Invalid JSON
test_endpoint "POST /v1/detect (invalid JSON)" "POST" "/v1/detect" "invalid json" 400

# Empty body
test_endpoint "POST /v1/detect (empty body)" "POST" "/v1/detect" "" 400

# Test security headers
echo ""
echo "=== Testing security headers ==="
HEADERS=$(curl -s -I "$API_URL/healthz")
if echo "$HEADERS" | grep -q "X-Frame-Options: DENY"; then
    echo -e "  ${GREEN}✓ X-Frame-Options header present${NC}"
    PASSED=$((PASSED + 1))
else
    echo -e "  ${RED}✗ X-Frame-Options header missing${NC}"
    FAILED=$((FAILED + 1))
fi

if echo "$HEADERS" | grep -iq "X-XSS-Protection"; then
    echo -e "  ${GREEN}✓ X-XSS-Protection header present${NC}"
    PASSED=$((PASSED + 1))
else
    echo -e "  ${RED}✗ X-XSS-Protection header missing${NC}"
    FAILED=$((FAILED + 1))
fi

# Summary
echo ""
echo "=== API Test Summary ==="
echo -e "${GREEN}Passed: $PASSED${NC}"
if [ $FAILED -gt 0 ]; then
    echo -e "${RED}Failed: $FAILED${NC}"
    exit 1
else
    echo -e "${GREEN}All API tests passed!${NC}"
    exit 0
fi

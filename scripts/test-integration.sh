#!/bin/bash
# Full integration test for Driftlock

set -e

echo "🔗 Running full integration tests..."
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

API_URL="${API_URL:-http://localhost:8080}"
TEST_DATA_DIR="${TEST_DATA_DIR:-test-data}"

# Track results
FAILED=0
PASSED=0

# Function to test anomaly detection with test data
test_anomaly_detection() {
    local file=$1
    local expected_min=$2
    local expected_max=$3
    local description=$4
    
    echo -n "Testing $description... "
    
    if [ ! -f "$file" ]; then
        echo -e "${YELLOW}⚠ SKIPPED (file not found: $file)${NC}"
        return 0
    fi
    
    RESPONSE=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        --data-binary "@$file" \
        "$API_URL/v1/detect?format=ndjson")
    
    ANOMALY_COUNT=$(echo "$RESPONSE" | jq -r '.anomaly_count // 0')
    
    if [ "$ANOMALY_COUNT" -ge "$expected_min" ] && [ "$ANOMALY_COUNT" -le "$expected_max" ]; then
        echo -e "${GREEN}✓ PASSED${NC} (found $ANOMALY_COUNT anomalies, expected $expected_min-$expected_max)"
        PASSED=$((PASSED + 1))
        return 0
    else
        echo -e "${RED}✗ FAILED${NC} (found $ANOMALY_COUNT anomalies, expected $expected_min-$expected_max)"
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

# Test with normal transactions (should find < 5 anomalies)
echo "=== Testing with normal transactions ==="
test_anomaly_detection "$TEST_DATA_DIR/normal-transactions.jsonl" 0 5 "normal transactions"

# Test with anomalous transactions (should find > 80 anomalies)
echo ""
echo "=== Testing with anomalous transactions ==="
test_anomaly_detection "$TEST_DATA_DIR/anomalous-transactions.jsonl" 80 100 "anomalous transactions"

# Test with mixed transactions (should find 45-55 anomalies)
echo ""
echo "=== Testing with mixed transactions ==="
test_anomaly_detection "$TEST_DATA_DIR/mixed-transactions.jsonl" 45 55 "mixed transactions"

# Test Prometheus metrics increment
echo ""
echo "=== Testing Prometheus metrics ==="
INITIAL_COUNT=$(curl -s "$API_URL/metrics" | grep "driftlock_http_requests_total" | awk '{print $2}' | head -1 || echo "0")

# Make a request
curl -s -X POST \
    -H "Content-Type: application/json" \
    --data-binary '{"test":"data"}' \
    "$API_URL/v1/detect?format=ndjson" > /dev/null

sleep 1

NEW_COUNT=$(curl -s "$API_URL/metrics" | grep "driftlock_http_requests_total" | awk '{print $2}' | head -1 || echo "0")

if [ "$NEW_COUNT" -gt "$INITIAL_COUNT" ]; then
    echo -e "  ${GREEN}✓ Metrics incremented correctly${NC}"
    PASSED=$((PASSED + 1))
else
    echo -e "  ${RED}✗ Metrics did not increment${NC}"
    FAILED=$((FAILED + 1))
fi

# Test response structure
echo ""
echo "=== Testing response structure ==="
RESPONSE=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    --data-binary '{"test":"data"}' \
    "$API_URL/v1/detect?format=ndjson")

REQUIRED_FIELDS=("success" "request_id" "total_events" "anomaly_count")
for field in "${REQUIRED_FIELDS[@]}"; do
    if echo "$RESPONSE" | jq -e ".$field" > /dev/null 2>&1; then
        echo -e "  ${GREEN}✓ Field '$field' present${NC}"
        PASSED=$((PASSED + 1))
    else
        echo -e "  ${RED}✗ Field '$field' missing${NC}"
        FAILED=$((FAILED + 1))
    fi
done

# Summary
echo ""
echo "=== Integration Test Summary ==="
echo -e "${GREEN}Passed: $PASSED${NC}"
if [ $FAILED -gt 0 ]; then
    echo -e "${RED}Failed: $FAILED${NC}"
    exit 1
else
    echo -e "${GREEN}All integration tests passed!${NC}"
    exit 0
fi


#!/bin/bash
# Test Docker Compose service startup and health checks

set -e

echo "🚀 Testing Driftlock service startup and health..."
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Track results
FAILED=0
PASSED=0

# Function to test service health
test_health() {
    local service=$1
    local url=$2
    local max_wait=${3:-30}
    
    echo -n "Testing $service health... "
    
    local elapsed=0
    while [ $elapsed -lt $max_wait ]; do
        if curl -sf "$url" > /dev/null 2>&1; then
            echo -e "${GREEN}✓ HEALTHY${NC}"
            PASSED=$((PASSED + 1))
            return 0
        fi
        sleep 1
        elapsed=$((elapsed + 1))
    done
    
    echo -e "${RED}✗ UNHEALTHY (timeout after ${max_wait}s)${NC}"
    FAILED=$((FAILED + 1))
    return 1
}

# Start HTTP API server
echo "=== Starting driftlock-http service ==="
docker compose up -d driftlock-http
sleep 5

# Test health endpoint
test_health "driftlock-http" "http://localhost:8080/healthz" 30

# Verify health response structure
echo -n "Verifying health response structure... "
HEALTH_RESPONSE=$(curl -s http://localhost:8080/healthz)
if echo "$HEALTH_RESPONSE" | jq -e '.success == true' > /dev/null 2>&1 && \
   echo "$HEALTH_RESPONSE" | jq -e '.library_status' > /dev/null 2>&1; then
    echo -e "${GREEN}✓ PASSED${NC}"
    PASSED=$((PASSED + 1))
else
    echo -e "${RED}✗ FAILED${NC}"
    echo "  Response: $HEALTH_RESPONSE"
    FAILED=$((FAILED + 1))
fi

# Test metrics endpoint
echo -n "Testing /metrics endpoint... "
if curl -sf http://localhost:8080/metrics > /dev/null 2>&1; then
    echo -e "${GREEN}✓ PASSED${NC}"
    PASSED=$((PASSED + 1))
    
    # Verify Prometheus format
    METRICS=$(curl -s http://localhost:8080/metrics)
    if echo "$METRICS" | grep -q "driftlock_http_requests_total"; then
        echo -e "  ${GREEN}✓ Prometheus metrics found${NC}"
    else
        echo -e "  ${YELLOW}⚠ Prometheus metrics not found${NC}"
    fi
else
    echo -e "${RED}✗ FAILED${NC}"
    FAILED=$((FAILED + 1))
fi

# Test graceful shutdown
echo -n "Testing graceful shutdown... "
docker compose stop driftlock-http
sleep 2
if ! curl -sf http://localhost:8080/healthz > /dev/null 2>&1; then
    echo -e "${GREEN}✓ PASSED${NC}"
    PASSED=$((PASSED + 1))
else
    echo -e "${RED}✗ FAILED (service still responding)${NC}"
    FAILED=$((FAILED + 1))
fi

# Summary
echo ""
echo "=== Service Test Summary ==="
echo -e "${GREEN}Passed: $PASSED${NC}"
if [ $FAILED -gt 0 ]; then
    echo -e "${RED}Failed: $FAILED${NC}"
    exit 1
else
    echo -e "${GREEN}All service tests passed!${NC}"
    exit 0
fi


#!/bin/bash
# Test Kafka integration for Driftlock

set -e

echo "📨 Testing Kafka integration..."
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Track results
FAILED=0
PASSED=0

# Check if Kafka profile is enabled
echo "=== Checking Kafka services ==="

# Check if Kafka is running
if ! docker compose ps kafka | grep -q "Up"; then
    echo -e "${YELLOW}⚠ Kafka services not running${NC}"
    echo "Starting Kafka services..."
    docker compose --profile kafka up -d kafka zookeeper
    echo "Waiting for Kafka to be ready..."
    sleep 10
fi

# Test Kafka broker connectivity
echo -n "Testing Kafka broker connectivity... "
if docker compose exec -T kafka kafka-broker-api-versions --bootstrap-server localhost:9092 > /dev/null 2>&1; then
    echo -e "${GREEN}✓ PASSED${NC}"
    PASSED=$((PASSED + 1))
else
    echo -e "${RED}✗ FAILED${NC}"
    FAILED=$((FAILED + 1))
fi

# Test Zookeeper connectivity
echo -n "Testing Zookeeper connectivity... "
if docker compose exec -T zookeeper zkServer.sh status > /dev/null 2>&1; then
    echo -e "${GREEN}✓ PASSED${NC}"
    PASSED=$((PASSED + 1))
else
    echo -e "${YELLOW}⚠ Zookeeper status check failed (may still be starting)${NC}"
fi

# Check if collector is running
echo ""
echo "=== Checking collector service ==="
if docker compose ps driftlock-collector | grep -q "Up"; then
    echo -e "${GREEN}✓ Collector is running${NC}"
    PASSED=$((PASSED + 1))
    
    # Check collector logs for Kafka connection
    echo -n "Checking collector logs for Kafka connection... "
    if docker compose logs driftlock-collector 2>&1 | grep -q "Kafka publisher initialized"; then
        echo -e "${GREEN}✓ PASSED${NC}"
        PASSED=$((PASSED + 1))
    else
        echo -e "${YELLOW}⚠ Kafka publisher initialization not found in logs${NC}"
    fi
else
    echo -e "${YELLOW}⚠ Collector is not running${NC}"
    echo "Start with: docker compose --profile kafka up -d driftlock-collector"
fi

# Test topic creation (if kafka-topics.sh is available)
echo ""
echo "=== Testing Kafka topic operations ==="
if docker compose exec -T kafka kafka-topics.sh --bootstrap-server localhost:9092 --list > /dev/null 2>&1; then
    echo -e "  ${GREEN}✓ Kafka topics command works${NC}"
    PASSED=$((PASSED + 1))
    
    # List existing topics
    TOPICS=$(docker compose exec -T kafka kafka-topics.sh --bootstrap-server localhost:9092 --list 2>/dev/null || echo "")
    if [ -n "$TOPICS" ]; then
        echo "  Existing topics: $TOPICS"
    fi
else
    echo -e "  ${YELLOW}⚠ Kafka topics command not available${NC}"
fi

# Summary
echo ""
echo "=== Kafka Test Summary ==="
echo -e "${GREEN}Passed: $PASSED${NC}"
if [ $FAILED -gt 0 ]; then
    echo -e "${RED}Failed: $FAILED${NC}"
    echo ""
    echo "Note: Some tests may require the collector to be configured with Kafka enabled."
    exit 1
else
    echo -e "${GREEN}All Kafka tests passed!${NC}"
    exit 0
fi


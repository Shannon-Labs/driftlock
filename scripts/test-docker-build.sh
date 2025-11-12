#!/bin/bash
# Test Docker builds for all Driftlock services

set -e

echo "🐳 Testing Docker builds for Driftlock services..."
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Track results
FAILED=0
PASSED=0

# Function to test a build
test_build() {
    local name=$1
    local dockerfile=$2
    local build_args=$3
    
    echo -n "Building $name... "
    
    if docker build -t "driftlock-$name:test" $build_args -f "$dockerfile" . > /tmp/docker-build-$name.log 2>&1; then
        echo -e "${GREEN}✓ PASSED${NC}"
        PASSED=$((PASSED + 1))
        
        # Get image size
        SIZE=$(docker images "driftlock-$name:test" --format "{{.Size}}" | head -1)
        echo "  Image size: $SIZE"
        
        # Check if size is reasonable (< 500MB)
        SIZE_MB=$(echo $SIZE | sed 's/[^0-9.]//g')
        if (( $(echo "$SIZE_MB > 500" | bc -l) )); then
            echo -e "  ${YELLOW}⚠ Warning: Image size exceeds 500MB${NC}"
        fi
    else
        echo -e "${RED}✗ FAILED${NC}"
        echo "  See /tmp/docker-build-$name.log for details"
        FAILED=$((FAILED + 1))
    fi
    echo ""
}

# Test driftlock-http (default, no OpenZL)
echo "=== Testing driftlock-http (USE_OPENZL=false) ==="
test_build "http" "collector-processor/cmd/driftlock-http/Dockerfile" "--build-arg USE_OPENZL=false"

# Test driftlock-http (with OpenZL)
echo "=== Testing driftlock-http (USE_OPENZL=true) ==="
test_build "http-openzl" "collector-processor/cmd/driftlock-http/Dockerfile" "--build-arg USE_OPENZL=true"

# Test driftlock-collector (default, no OpenZL)
echo "=== Testing driftlock-collector (USE_OPENZL=false) ==="
test_build "collector" "collector-processor/cmd/driftlock-collector/Dockerfile" "--build-arg USE_OPENZL=false"

# Test driftlock-collector (with OpenZL)
echo "=== Testing driftlock-collector (USE_OPENZL=true) ==="
test_build "collector-openzl" "collector-processor/cmd/driftlock-collector/Dockerfile" "--build-arg USE_OPENZL=true"

# Summary
echo "=== Build Test Summary ==="
echo -e "${GREEN}Passed: $PASSED${NC}"
if [ $FAILED -gt 0 ]; then
    echo -e "${RED}Failed: $FAILED${NC}"
    exit 1
else
    echo -e "${GREEN}All builds passed!${NC}"
    exit 0
fi


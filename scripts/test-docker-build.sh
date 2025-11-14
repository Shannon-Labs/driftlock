#!/bin/bash
# Test Docker builds for all Driftlock services

set -euo pipefail

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

detect_openzl() {
    if [ -n "${OPENZL_LIB_DIR:-}" ] && [ -f "$OPENZL_LIB_DIR/libopenzl.a" ]; then
        return 0
    fi
    if [ -d "openzl" ] && ls openzl/libopenzl.a >/dev/null 2>&1; then
        return 0
    fi
    if ls cbad-core/openzl/libopenzl.a >/dev/null 2>&1; then
        return 0
    fi
    return 1
}

ENABLE_OPENZL_BUILD=${ENABLE_OPENZL_BUILD:-auto}
OPENZL_TESTS_ENABLED=0
case "$ENABLE_OPENZL_BUILD" in
    true|TRUE|1)
        OPENZL_TESTS_ENABLED=1
        ;;
    false|FALSE|0)
        OPENZL_TESTS_ENABLED=0
        ;;
    *)
        if detect_openzl; then
            OPENZL_TESTS_ENABLED=1
        fi
        ;;
esac

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
        [ -n "$SIZE" ] && echo "  Image size: $SIZE"
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
if [ "$OPENZL_TESTS_ENABLED" -eq 1 ]; then
    echo "=== Testing driftlock-http (USE_OPENZL=true) ==="
    test_build "http-openzl" "collector-processor/cmd/driftlock-http/Dockerfile" "--build-arg USE_OPENZL=true"
else
    echo "=== Skipping driftlock-http OpenZL build (libraries not detected) ==="
fi

# Test driftlock-collector (default, no OpenZL)
echo "=== Testing driftlock-collector (USE_OPENZL=false) ==="
test_build "collector" "collector-processor/cmd/driftlock-collector/Dockerfile" "--build-arg USE_OPENZL=false"

# Test driftlock-collector (with OpenZL)
if [ "$OPENZL_TESTS_ENABLED" -eq 1 ]; then
    echo "=== Testing driftlock-collector (USE_OPENZL=true) ==="
    test_build "collector-openzl" "collector-processor/cmd/driftlock-collector/Dockerfile" "--build-arg USE_OPENZL=true"
else
    echo "=== Skipping driftlock-collector OpenZL build (libraries not detected) ==="
fi

if [ "$OPENZL_TESTS_ENABLED" -eq 0 ]; then
    echo "Note: Set ENABLE_OPENZL_BUILD=true and provide libopenzl.a to test OpenZL images."
fi

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

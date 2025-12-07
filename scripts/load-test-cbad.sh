#!/bin/bash
# Load test script for CBAD detection performance
# Usage: ./scripts/load-test-cbad.sh [concurrent_streams] [events_per_stream]

set -e

CONCURRENT_STREAMS="${1:-10}"
EVENTS_PER_STREAM="${2:-1000}"

echo "=== CBAD Load Test ==="
echo "Concurrent streams: $CONCURRENT_STREAMS"
echo "Events per stream: $EVENTS_PER_STREAM"
echo "Timestamp: $(date)"
echo ""

# Build test binary if needed
echo "Building load tester..."
cd /Volumes/VIXinSSD/driftlock/collector-processor
GOFLAGS="-mod=mod" go build -o bin/load-tester -ldflags="-s -w" ./scripts/cbad_load_test.go

# Run load test
echo ""
echo "Running load test..."
echo "./bin/load-tester -streams $CONCURRENT_STREAMS -events $EVENTS_PER_STREAM"
echo ""

./bin/load-tester -streams "$CONCURRENT_STREAMS" -events "$EVENTS_PER_STREAM"

echo ""
echo "=== Load Test Complete ==="
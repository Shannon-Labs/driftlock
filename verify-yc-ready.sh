#!/bin/bash
set -e

echo "🔍 Verifying Driftlock YC Demo..."

# Build
echo "Building Rust core..."
if [ ! -f cbad-core/target/release/libcbad_core.a ]; then
  (cd cbad-core && timeout 120s cargo build --release) || { echo "❌ Rust core build failed"; exit 1; }
fi

echo "Building Go demo..."
go build -o driftlock-demo cmd/demo/main.go || { echo "❌ Go build failed"; exit 1; }

# Run and time
echo "Running demo..."
timeout 30s ./driftlock-demo test-data/financial-demo.json > verify.log 2>&1 || {
  echo "❌ Demo run failed or timed out"; exit 1;
}

# Check output
if [ ! -f demo-output.html ]; then
  echo "❌ demo-output.html missing"
  exit 1
fi

ANOMALIES=$(grep -c "ANOMALY DETECTED" demo-output.html || echo "0")
if [ "$ANOMALIES" -lt 10 ] || [ "$ANOMALIES" -gt 30 ]; then
  echo "❌ Anomaly count $ANOMALIES outside 10-30 range"
  exit 1
fi

echo "✅ Demo ready: $ANOMALIES anomalies, <30s runtime"
echo "✅ YC partners can run: ./verify-yc-ready.sh"

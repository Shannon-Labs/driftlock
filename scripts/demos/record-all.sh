#!/bin/bash
# Record all demo videos using VHS
# Usage: ./scripts/demos/record-all.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

cd "$PROJECT_ROOT"

echo "=================================="
echo "Recording Driftlock Demo Videos"
echo "=================================="
echo ""

# Ensure VHS is available
if ! command -v vhs &> /dev/null; then
    echo "Error: VHS not installed. Run: brew install vhs"
    exit 1
fi

# Pre-build to avoid recording compilation
echo "Pre-building examples in release mode..."
cargo build --examples --release
echo ""

# Record crypto_stream demo
echo "Recording: crypto_stream.tape"
echo "  Output: crypto_stream.gif, crypto_stream.mp4"
vhs "$SCRIPT_DIR/crypto_stream.tape"
echo ""

# Record kaggle_benchmarks demo
echo "Recording: kaggle_benchmarks.tape"
echo "  Output: kaggle_benchmarks.gif, kaggle_benchmarks.mp4"
vhs "$SCRIPT_DIR/kaggle_benchmarks.tape"
echo ""

echo "=================================="
echo "Done! Output files:"
echo "=================================="
ls -lh "$SCRIPT_DIR"/*.gif "$SCRIPT_DIR"/*.mp4 2>/dev/null || echo "  (files in scripts/demos/)"

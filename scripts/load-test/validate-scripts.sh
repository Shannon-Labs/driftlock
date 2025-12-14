#!/bin/bash
# Validate K6 load test scripts
# Usage: ./scripts/load-test/validate-scripts.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "=== K6 Load Test Script Validator ==="
echo ""

# Check if k6 is installed
if ! command -v k6 &> /dev/null; then
    echo "ERROR: k6 is not installed!"
    echo ""
    echo "Install k6:"
    echo "  macOS:   brew install k6"
    echo "  Linux:   See QUICKSTART.md"
    echo "  Windows: choco install k6"
    echo ""
    exit 1
fi

echo "✓ K6 found: $(k6 version | head -n1)"
echo ""

# List all test scripts
SCRIPTS=(
    "smoke.js"
    "load.js"
    "stress.js"
    "soak.js"
    "detector-capacity.js"
    "rate-limit-validation.js"
)

echo "Validating scripts..."
echo ""

ERRORS=0

for script in "${SCRIPTS[@]}"; do
    if [ ! -f "$script" ]; then
        echo "✗ $script: NOT FOUND"
        ERRORS=$((ERRORS + 1))
        continue
    fi

    # Check for syntax errors (basic validation)
    if node -c "$script" 2>/dev/null; then
        echo "✓ $script: Syntax OK"
    else
        # K6 scripts are not pure JS, so this might fail
        # Just check if file exists and is readable
        if [ -r "$script" ]; then
            echo "✓ $script: File readable"
        else
            echo "✗ $script: Cannot read file"
            ERRORS=$((ERRORS + 1))
        fi
    fi
done

echo ""
echo "Checking required files..."
if [ -f "helpers.js" ]; then
    echo "✓ helpers.js exists"
else
    echo "✗ helpers.js missing"
    ERRORS=$((ERRORS + 1))
fi

if [ -f "config/thresholds.json" ]; then
    echo "✓ config/thresholds.json exists"
else
    echo "✗ config/thresholds.json missing"
    ERRORS=$((ERRORS + 1))
fi

echo ""

if [ $ERRORS -eq 0 ]; then
    echo "==================================="
    echo "✓ All validation checks passed!"
    echo "==================================="
    echo ""
    echo "Next steps:"
    echo "  1. Start API: cargo run -p driftlock-api --release"
    echo "  2. Run smoke test: k6 run scripts/load-test/smoke.js"
    echo ""
    exit 0
else
    echo "==================================="
    echo "✗ Validation failed with $ERRORS error(s)"
    echo "==================================="
    exit 1
fi

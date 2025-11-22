#!/bin/bash
# Quick start script for 4-hour crypto test
# Automatically creates API key if needed using CLI

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "🚀 Driftlock 4-Hour Crypto Test Setup"
echo ""

# Load .env if it exists
if [ -f "$REPO_ROOT/.env" ]; then
    echo "📥 Loading .env file..."
    set -a
    source "$REPO_ROOT/.env"
    set +a
fi

# Check for API key
if [ -z "$DRIFTLOCK_API_KEY" ]; then
    echo "⚠️  No API key found. Creating one automatically..."
    echo ""
    
    # Try to create one using Cloud Run Jobs (works from anywhere)
    echo "   Attempting to create via Cloud Run Job..."
    if "$SCRIPT_DIR/create-test-api-key-cloudrun.sh"; then
        # Reload .env to get the new key
        if [ -f "$REPO_ROOT/.env" ]; then
            set -a
            source "$REPO_ROOT/.env"
            set +a
        fi
        echo "✅ API key created and loaded!"
    else
        echo ""
        echo "❌ Failed to create API key automatically."
        echo ""
        echo "Please either:"
        echo "  1. Run manually: ./scripts/create-test-api-key-cloudrun.sh"
        echo "  2. Or sign up at https://driftlock.web.app and set:"
        echo "     export DRIFTLOCK_API_KEY='dlk_...'"
        exit 1
    fi
fi

export DRIFTLOCK_API_URL="${DRIFTLOCK_API_URL:-https://driftlock.web.app/api/v1}"

echo "✅ API Key: ${DRIFTLOCK_API_KEY:0:20}..."
echo "✅ API URL: $DRIFTLOCK_API_URL"
echo ""
echo "📊 Starting 4-hour crypto test..."
echo "   This will stream live Binance data and detect anomalies"
echo "   Logs will be saved to: logs/crypto-api-test-*.log"
echo ""

# Run the 4-hour test script
"$SCRIPT_DIR/run_crypto_test_4h.sh"

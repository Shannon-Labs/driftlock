#!/bin/bash

# Driftlock Stripe Webhook Forwarding Script
# Forwards Stripe webhooks to local Firebase Functions emulator

set -e

echo "🚀 Starting Stripe Webhook Forwarding for Driftlock..."
echo ""
echo "This script helps you test Stripe webhooks locally by forwarding them"
echo "from the Stripe CLI to your local Firebase Functions emulator."
echo ""

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if stripe CLI is installed
if ! command -v stripe &> /dev/null; then
    echo -e "${RED}❌ Stripe CLI not found${NC}"
    echo "Please install Stripe CLI first:"
    echo "  macOS: brew install stripe/stripe-cli/stripe"
    echo "  Linux: See https://stripe.com/docs/stripe-cli"
    exit 1
fi

echo -e "${GREEN}✅ Stripe CLI is installed${NC}"

# Check if user is logged in to Stripe
if ! stripe config --list &> /dev/null; then
    echo -e "${YELLOW}⚠️  Not logged into Stripe${NC}"
    echo "Please run: stripe login"
    exit 1
fi

echo -e "${GREEN}✅ Logged into Stripe${NC}"

# Find Firebase Functions emulator port
FIREBASE_CONFIG_FILE="firebase.json"
DEFAULT_PORT=5001
FUNCTIONS_PORT=$DEFAULT_PORT

if [[ -f "$FIREBASE_CONFIG_FILE" ]]; then
    # Try to extract port from firebase.json
    EMU_PORT=$(grep -o '"functionsPort":[[:space:]]*[0-9]*' "$FIREBASE_CONFIG_FILE" | grep -o '[0-9]*' || echo "")
    if [[ -n "$EMU_PORT" ]]; then
        FUNCTIONS_PORT=$EMU_PORT
    fi
fi

echo ""
echo "🔧 Configuration:"
echo "  Functions Emulator Port: $FUNCTIONS_PORT"
echo "  Webhook Endpoint: /webhooks/stripe"
echo ""

# Check if emulator is running
if ! curl -s http://localhost:$FUNCTIONS_PORT > /dev/null 2>&1; then
    echo -e "${RED}❌ Firebase Functions emulator not running${NC}"
    echo ""
    echo "Please start the emulator first:"
    echo "  cd functions"
    echo "  npm run build"
    echo "  firebase emulators:start --only functions"
    echo ""
    echo "Then run this script in another terminal."
    exit 1
fi

echo -e "${GREEN}✅ Firebase Functions emulator is running${NC}"
echo ""

# Create the forwarding URL
FORWARD_URL="http://localhost:${FUNCTIONS_PORT}/driftlock/us-central1/apiProxy/webhooks/stripe"

echo "🔗 Webhook forwarding URL: $FORWARD_URL"
echo ""
echo -e "${YELLOW}Starting webhook listener...${NC}"
echo ""
echo "💡 This will forward all Stripe events to your local emulator."
echo "💡 Keep this terminal open while testing."
echo "💡 In another terminal, trigger events with: stripe trigger <event>"
echo ""
echo "Common events to test:"
echo "  stripe trigger customer.subscription.created"
echo "  stripe trigger customer.subscription.deleted"
echo "  stripe trigger invoice.payment_succeeded"
echo "  stripe trigger invoice.payment_failed"
echo "  stripe trigger checkout.session.completed"
echo ""
echo "Press Ctrl+C to stop"
echo ""
echo "=================================================="
echo ""

# Start stripe listen with forwarding
stripe listen --forward-to "$FORWARD_URL" --format JSON

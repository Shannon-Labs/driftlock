#!/bin/bash
# Test AI Explainability for Driftlock
# Usage: ./test-ai-explainability.sh <API_KEY>
#
# This script tests that AI-enhanced explanations are generated for anomalies.
# Prerequisites:
#   - A paid account (Radar, Tensor, or Orbit tier)
#   - OR trial account with Z.AI configured
#   - API is running with AI_PROVIDER and AI_API_KEY configured

set -e

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

API_URL="${API_URL:-http://localhost:8080}"
API_KEY="${1:-$DRIFTLOCK_API_KEY}"

if [ -z "$API_KEY" ]; then
    echo -e "${RED}Error: API key required${NC}"
    echo "Usage: $0 <API_KEY>"
    echo "Or set DRIFTLOCK_API_KEY environment variable"
    exit 1
fi

echo -e "${CYAN}=== Driftlock AI Explainability Test ===${NC}"
echo ""
echo "API URL: $API_URL"
echo "API Key: ${API_KEY:0:12}..."
echo ""

# Step 1: Check API health
echo -e "${YELLOW}Step 1: Checking API availability...${NC}"
if ! curl -sf "$API_URL/healthz" > /dev/null 2>&1; then
    echo -e "${RED}API is not available at $API_URL${NC}"
    exit 1
fi
echo -e "${GREEN}API is available${NC}"
echo ""

# Step 2: Check billing status to verify account type
echo -e "${YELLOW}Step 2: Checking billing status...${NC}"
BILLING_RESPONSE=$(curl -s -H "Authorization: Bearer $API_KEY" "$API_URL/api/v1/me/billing" 2>&1)
PLAN=$(echo "$BILLING_RESPONSE" | jq -r '.plan // "unknown"')
STATUS=$(echo "$BILLING_RESPONSE" | jq -r '.status // "unknown"')

echo "Plan: $PLAN"
echo "Status: $STATUS"

if [ "$PLAN" = "pulse" ] || [ "$PLAN" = "unknown" ]; then
    echo -e "${YELLOW}Warning: Free tier (pulse) may have limited AI features${NC}"
    echo "For full AI explanations, upgrade to Radar or higher"
fi
echo ""

# Step 3: Check AI usage endpoint (verify AI is configured)
echo -e "${YELLOW}Step 3: Checking AI configuration...${NC}"
AI_USAGE=$(curl -s -H "Authorization: Bearer $API_KEY" "$API_URL/api/v1/me/usage/ai" 2>&1)
AI_MODEL=$(echo "$AI_USAGE" | jq -r '.model_type // "none"')
AI_CALLS=$(echo "$AI_USAGE" | jq -r '.calls_used // 0')
AI_LIMIT=$(echo "$AI_USAGE" | jq -r '.calls_limit // 0')

echo "AI Model: $AI_MODEL"
echo "AI Calls: $AI_CALLS / $AI_LIMIT"

if [ "$AI_MODEL" = "none" ] && [ "$PLAN" != "pulse" ]; then
    echo -e "${YELLOW}Warning: No AI model configured. Check AI_PROVIDER environment variable.${NC}"
fi
echo ""

# Step 4: Send test events designed to trigger anomalies
echo -e "${YELLOW}Step 4: Sending test events to trigger anomaly detection...${NC}"

# Create baseline data (normal pattern)
BASELINE_EVENTS=""
for i in $(seq 1 50); do
    BASELINE_EVENTS+='{"timestamp":"2025-01-01T00:'$(printf "%02d" $((i/60)))':'$(printf "%02d" $((i%60)))'Z","cpu":45,"memory":60,"latency":100,"status":"ok"}'$'\n'
done

# Create anomalous data (significant deviation)
ANOMALY_EVENTS=""
for i in $(seq 1 10); do
    ANOMALY_EVENTS+='{"timestamp":"2025-01-01T01:'$(printf "%02d" $i)':00Z","cpu":95,"memory":95,"latency":5000,"status":"critical","error_rate":0.85}'$'\n'
done

# Combine for full payload
TEST_EVENTS="${BASELINE_EVENTS}${ANOMALY_EVENTS}"

# Send detection request
DETECT_RESPONSE=$(curl -s -X POST \
    -H "Authorization: Bearer $API_KEY" \
    -H "Content-Type: application/json" \
    --data-binary "$TEST_EVENTS" \
    "$API_URL/api/v1/detect?format=ndjson&stream=ai-test-stream" 2>&1)

# Parse response
SUCCESS=$(echo "$DETECT_RESPONSE" | jq -r '.success // false')
BATCH_ID=$(echo "$DETECT_RESPONSE" | jq -r '.batch_id // "unknown"')
ANOMALY_COUNT=$(echo "$DETECT_RESPONSE" | jq -r '.anomaly_count // 0')
TOTAL_EVENTS=$(echo "$DETECT_RESPONSE" | jq -r '.total_events // 0')

echo "Success: $SUCCESS"
echo "Batch ID: $BATCH_ID"
echo "Total Events: $TOTAL_EVENTS"
echo "Anomalies Detected: $ANOMALY_COUNT"

if [ "$ANOMALY_COUNT" -eq 0 ]; then
    echo -e "${YELLOW}No anomalies detected. Try adjusting detection thresholds.${NC}"
    echo ""
    echo "Response:"
    echo "$DETECT_RESPONSE" | jq .
    exit 0
fi
echo ""

# Step 5: Check anomaly explanations
echo -e "${YELLOW}Step 5: Checking anomaly explanations...${NC}"

# Get the first anomaly from the response
FIRST_ANOMALY=$(echo "$DETECT_RESPONSE" | jq -r '.anomalies[0] // empty')
if [ -n "$FIRST_ANOMALY" ]; then
    ANOMALY_ID=$(echo "$FIRST_ANOMALY" | jq -r '.id // "unknown"')
    EXPLANATION=$(echo "$FIRST_ANOMALY" | jq -r '.explanation // "none"')
    NCD=$(echo "$FIRST_ANOMALY" | jq -r '.metrics.ncd // "N/A"')
    CONFIDENCE=$(echo "$FIRST_ANOMALY" | jq -r '.metrics.confidence_level // "N/A"')

    echo "Anomaly ID: $ANOMALY_ID"
    echo "NCD Score: $NCD"
    echo "Confidence: $CONFIDENCE"
    echo ""
    echo -e "${CYAN}Explanation:${NC}"
    echo "$EXPLANATION"
fi
echo ""

# Step 6: Wait for async AI analysis and re-check
echo -e "${YELLOW}Step 6: Waiting for async AI analysis (5 seconds)...${NC}"
sleep 5

# Query anomalies endpoint to get updated explanation
echo "Fetching updated anomaly from API..."
ANOMALY_DETAIL=$(curl -s -H "Authorization: Bearer $API_KEY" \
    "$API_URL/api/v1/anomalies?limit=1&stream=ai-test-stream" 2>&1)

UPDATED_EXPLANATION=$(echo "$ANOMALY_DETAIL" | jq -r '.anomalies[0].explanation // "none"')

echo ""
echo -e "${CYAN}Updated Explanation (after AI analysis):${NC}"
echo "$UPDATED_EXPLANATION"
echo ""

# Step 7: Verify AI explanation quality
echo -e "${YELLOW}Step 7: Verifying AI explanation quality...${NC}"

if [ "$UPDATED_EXPLANATION" = "none" ] || [ -z "$UPDATED_EXPLANATION" ] || [ "$UPDATED_EXPLANATION" = "null" ]; then
    echo -e "${YELLOW}No AI explanation found${NC}"
    echo ""
    echo "Possible reasons:"
    echo "  1. AI is disabled for this plan"
    echo "  2. Anomaly score below analysis threshold"
    echo "  3. AI_PROVIDER or AI_API_KEY not configured"
    echo "  4. Rate limits exceeded"
    echo ""
    echo "Check server logs for more details."
    exit 1
fi

# Check explanation length and content
EXPLANATION_LENGTH=${#UPDATED_EXPLANATION}
if [ $EXPLANATION_LENGTH -gt 50 ]; then
    echo -e "${GREEN}AI explanation present (${EXPLANATION_LENGTH} chars)${NC}"

    # Check for key explanation elements
    HAS_ANALYSIS=false
    if echo "$UPDATED_EXPLANATION" | grep -qi "anomaly\|unusual\|deviation\|spike\|pattern"; then
        HAS_ANALYSIS=true
    fi

    if [ "$HAS_ANALYSIS" = true ]; then
        echo -e "${GREEN}Explanation contains analytical content${NC}"
    else
        echo -e "${YELLOW}Explanation may be generic (check content)${NC}"
    fi
else
    echo -e "${YELLOW}Explanation seems short (${EXPLANATION_LENGTH} chars)${NC}"
fi
echo ""

# Summary
echo -e "${CYAN}=== Test Summary ===${NC}"
echo "Plan: $PLAN"
echo "AI Model: $AI_MODEL"
echo "Anomalies Detected: $ANOMALY_COUNT"
echo "AI Explanation: $([ ${#UPDATED_EXPLANATION} -gt 50 ] && echo 'Present' || echo 'Limited/Missing')"
echo ""

if [ ${#UPDATED_EXPLANATION} -gt 50 ]; then
    echo -e "${GREEN}AI Explainability Test: PASSED${NC}"
    exit 0
else
    echo -e "${YELLOW}AI Explainability Test: PARTIAL${NC}"
    echo "AI explanation was limited. Check configuration and logs."
    exit 0
fi

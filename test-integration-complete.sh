#!/usr/bin/env bash
# Complete end-to-end integration test for Driftlock

set -euo pipefail

echo "🧪 Testing Complete Driftlock Integration"
echo "========================================="

SUPABASE_URL_DEFAULT="https://nfkdeeunyvnntvpvwpwh.supabase.co"
SUPABASE_URL="${SUPABASE_BASE_URL:-$SUPABASE_URL_DEFAULT}"

# 1. Test Go API Health
echo "1. Testing Go API..."
if ! curl -sf http://localhost:8080/healthz >/dev/null; then
  echo "❌ Go API not running at http://localhost:8080/healthz"
else
  echo "✅ Go API is healthy"
fi

# 2. Test Supabase Connection
echo "2. Testing Supabase Edge Functions..."
if ! curl -sf "${SUPABASE_URL}/functions/v1/health" >/dev/null; then
  echo "❌ Supabase health check failed at ${SUPABASE_URL}/functions/v1/health"
else
  echo "✅ Supabase edge function health endpoint reachable"
fi

# 3. Create Test Anomaly via Go API
echo "3. Creating test anomaly via Go API..."
RESPONSE=$(curl -s -X POST http://localhost:8080/v1/events \
  -H "Content-Type: application/json" \
  -d '{
    "organization_id": "test-org-123",
    "event_type": "log",
    "data": {
      "message": "Critical error detected",
      "level": "ERROR",
      "source": "api-gateway"
    }
  }' || true)
echo "Response: $RESPONSE"

# 4. Verify Anomaly Appears in Supabase
if command -v jq >/dev/null 2>&1; then
  echo "4. Checking Supabase for anomaly..."
  sleep 2
  curl -s "${SUPABASE_URL}/rest/v1/anomalies?select=*&limit=1" \
    -H "apikey: ${SUPABASE_ANON_KEY:-}" | jq .
else
  echo "4. Skipping Supabase anomaly check (jq not installed)"
fi

# 5. Check Usage Metering
if command -v jq >/dev/null 2>&1; then
  echo "5. Verifying usage metering..."
  curl -s "${SUPABASE_URL}/rest/v1/usage_counters?organization_id=eq.test-org-123" \
    -H "apikey: ${SUPABASE_ANON_KEY:-}" | jq .
else
  echo "5. Skipping usage metering check (jq not installed)"
fi

echo "✅ Integration test complete"


# cURL Examples

Complete collection of cURL commands for interacting with Driftlock's REST API.

## Prerequisites

Set your API key as an environment variable:
```bash
export DRIFTLOCK_API_KEY="your_api_key_here"
export DRIFTLOCK_API="https://api.driftlock.net"
```

## Detection

### Basic Detection
```bash
curl -X POST "$DRIFTLOCK_API/v1/detect" \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -d '{
    "stream_id": "default",
    "events": [
      {
        "timestamp": "2025-01-01T10:00:00Z",
        "type": "log",
        "body": {"message": "Normal login", "user": "alice"}
      },
      {
        "timestamp": "2025-01-01T10:01:00Z",
        "type": "log",
        "body": {"message": "Normal login", "user": "bob"}
      },
      {
        "timestamp": "2025-01-01T10:02:00Z",
        "type": "log",
        "body": {"message": "SQL INJECTION ATTEMPT!", "user": "hacker"}
      }
    ]
  }'
```

### Detection with Config Override
```bash
curl -X POST "$DRIFTLOCK_API/v1/detect" \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -d '{
    "stream_id": "production",
    "events": [...],
    "config_override": {
      "ncd_threshold": 0.4,
      "p_value_threshold": 0.01,
      "compressor": "lz4",
      "baseline_size": 500
    }
  }'
```

### Metric Monitoring
```bash
curl -X POST "$DRIFTLOCK_API/v1/detect" \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -d '{
    "stream_id": "metrics",
    "events": [
      {
        "timestamp": "2025-01-01T10:00:00Z",
        "type": "metric",
        "body": {"cpu": 45, "memory": 2048, "latency": 120},
        "attributes": {"host": "server-01"}
      },
      {
        "timestamp": "2025-01-01T10:01:00Z",
        "type": "metric",
        "body": {"cpu": 48, "memory": 2100, "latency": 125},
        "attributes": {"host": "server-01"}
      },
      {
        "timestamp": "2025-01-01T10:02:00Z",
        "type": "metric",
        "body": {"cpu": 95, "memory": 7800, "latency": 3500},
        "attributes": {"host": "server-01"}
      }
    ]
  }'
```

## Anomalies

### List All Anomalies
```bash
curl "$DRIFTLOCK_API/v1/anomalies" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY"
```

### List with Pagination
```bash
curl "$DRIFTLOCK_API/v1/anomalies?limit=50" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY"
```

### Filter by Stream
```bash
curl "$DRIFTLOCK_API/v1/anomalies?stream_id=production-logs" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY"
```

### Filter by NCD Threshold
```bash
curl "$DRIFTLOCK_API/v1/anomalies?min_ncd=0.5" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY"
```

### Filter by Date Range
```bash
curl "$DRIFTLOCK_API/v1/anomalies?since=2025-01-01T00:00:00Z&until=2025-01-31T23:59:59Z" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY"
```

### Combined Filters
```bash
curl "$DRIFTLOCK_API/v1/anomalies?stream_id=production&min_ncd=0.5&status=new&limit=100" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY"
```

### Get Anomaly Details
```bash
curl "$DRIFTLOCK_API/v1/anomalies/anom_abc123" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY"
```

## Exports

### Export Anomalies (JSON)
```bash
curl -X POST "$DRIFTLOCK_API/v1/anomalies/export" \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -d '{
    "format": "json",
    "filters": {
      "stream_id": "production",
      "status": "new"
    },
    "delivery": {
      "type": "sync"
    }
  }'
```

### Export Single Anomaly
```bash
curl -X POST "$DRIFTLOCK_API/v1/anomalies/anom_abc123/export" \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -d '{
    "format": "markdown"
  }'
```

## Dashboard (Firebase Auth Required)

### List API Keys
First, get your Firebase ID token from the browser:
```javascript
// In browser console after logging in
const user = firebase.auth().currentUser;
const token = await user.getIdToken();
console.log(token);
```

Then use it:
```bash
export FIREBASE_TOKEN="your_firebase_id_token"

curl "$DRIFTLOCK_API/v1/me/keys" \
  -H "Authorization: Bearer $FIREBASE_TOKEN"
```

### Get Usage Statistics
```bash
curl "$DRIFTLOCK_API/v1/me/usage" \
  -H "Authorization: Bearer $FIREBASE_TOKEN"
```

## Billing

### Create Checkout Session
```bash
curl -X POST "$DRIFTLOCK_API/v1/billing/checkout" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "plan": "starter"
  }'
```

### Get Customer Portal Link
```bash
curl -X POST "$DRIFTLOCK_API/v1/billing/portal" \
  -H "Authorization: Bearer $FIREBASE_TOKEN" \
  -H "Content-Type: application/json"
```

## Health Check

### Basic Health Check
```bash
curl "$DRIFTLOCK_API/healthz"
```

### Detailed Health Check  
```bash
curl "$DRIFTLOCK_API/healthz?full=1" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY"
```

## Debugging

### Pretty Print JSON Response
```bash
curl "$DRIFTLOCK_API/v1/anomalies" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" | jq '.'
```

### Show Response Headers
```bash
curl -v "$DRIFTLOCK_API/v1/anomalies" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY"
```

### Save Response to File
```bash
curl "$DRIFTLOCK_API/v1/anomalies" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -o anomalies.json
```

### Check Rate Limits
```bash
curl -I "$DRIFTLOCK_API/v1/anomalies" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" | grep -i ratelimit
```

## Error Handling

### Retry on Rate Limit
```bash
#!/bin/bash
response=$(curl -s -w "\\n%{http_code}" "$DRIFTLOCK_API/v1/detect" \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -d '{...}')

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" == "429" ]; then
  retry_after=$(echo "$body" | jq -r '.error.retry_after_seconds')
  echo "Rate limited. Sleeping $retry_after seconds..."
  sleep "$retry_after"
  # Retry request...
fi
```

## Complete Example Script

Save this as `detect.sh`:

```bash
#!/bin/bash
set -e

# Configuration
API_KEY="${DRIFTLOCK_API_KEY:?'Set DRIFTLOCK_API_KEY environment variable'}"
API_URL="https://api.driftlock.net"

# Run detection
echo "Running detection..."
response=$(curl -s -X POST "$API_URL/v1/detect" \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: $API_KEY" \
  -d '{
    "stream_id": "default",
    "events": [
      {"timestamp": "2025-01-01T10:00:00Z", "type": "metric", "body": {"value": 100}},
      {"timestamp": "2025-01-01T10:01:00Z", "type": "metric", "body": {"value": 105}},
      {"timestamp": "2025-01-01T10:02:00Z", "type": "metric", "body": {"value": 950}}
    ]
  }')

# Parse response
echo "$response" | jq '.'

# Extract anomaly count
anomaly_count=$(echo "$response" | jq '.anomaly_count')
echo "\\nDetected $anomaly_count anomalies"

# List anomaly IDs
echo "\\nAnomaly IDs:"
echo "$response" | jq -r '.anomalies[].id'
```

Make it executable and run:
```bash
chmod +x detect.sh
./detect.sh
```

## Next Steps

- **[Python Examples](./python-examples.md)** - Python client implementation
- **[Node.js Examples](./node-examples.md)** - JavaScript/TypeScript examples
- **[API Reference](../rest-api.md)** - Complete API documentation

---

**Tip**: Use [jq](https://stedolan.github.io/jq/) for easier JSON parsing in bash scripts!

# cURL Examples

Copy/paste commands for common Driftlock API calls.

## Prerequisites

```bash
export DRIFTLOCK_API_URL="https://api.driftlock.net/v1"
export DRIFTLOCK_API_KEY="your_api_key_here"
```

## Detection

### Basic detection
```bash
curl -X POST "$DRIFTLOCK_API_URL/detect" \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -d '{
    "stream_id": "default",
    "events": [
      {"timestamp": "2025-01-01T10:00:00Z", "type": "metric", "body": {"latency_ms": 120}},
      {"timestamp": "2025-01-01T10:01:00Z", "type": "metric", "body": {"latency_ms": 950}}  /* anomaly */
    ]
  }'
```

### Detection with overrides
```bash
curl -X POST "$DRIFTLOCK_API_URL/detect" \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -d '{
    "events": [{"type": "log", "body": {"message": "SQL injection attempt"}}],
    "config_override": {
      "ncd_threshold": 0.35,
      "p_value_threshold": 0.02,
      "compressor": "lz4"
    }
  }'
```

### Demo detection (no auth)
```bash
curl -X POST "$DRIFTLOCK_API_URL/demo/detect" \
  -H "Content-Type: application/json" \
  -d '{
    "events": [
      {"body": {"cpu": 45, "memory": 2048}},
      {"body": {"cpu": 99, "memory": 8000}}  /* anomaly */
    ]
  }'
```

## Anomalies

### List anomalies
```bash
curl "$DRIFTLOCK_API_URL/anomalies" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY"
```

### Filter by stream and NCD
```bash
curl "$DRIFTLOCK_API_URL/anomalies?stream_id=prod-logs&min_ncd=0.5" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY"
```

### Date range and pagination
```bash
curl "$DRIFTLOCK_API_URL/anomalies?since=2025-01-01T00:00:00Z&until=2025-01-02T00:00:00Z&limit=20" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY"

# Next page using the returned token
curl "$DRIFTLOCK_API_URL/anomalies?page_token=eyJ..." \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY"
```

## Rate limits quick check
```bash
curl -I "$DRIFTLOCK_API_URL/anomalies" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" | grep -i ratelimit
```

**Next:** [Python examples](./python-examples.md)

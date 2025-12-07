# POST /v1/demo/detect

Try anomaly detection instantly without signing up. This endpoint is rate-limited but requires no authentication.

## Endpoint

```http
POST https://api.driftlock.net/v1/demo/detect
```

## Authentication

**None required** - This is a public demo endpoint.

## Rate Limits

| Limit | Value |
|-------|-------|
| Requests per minute | 10 per IP |
| Events per request | 50 max |

## Request

### Headers

```http
Content-Type: application/json
```

### Body Schema

```json
{
  "events": [
    {
      "timestamp": "2025-01-01T10:00:00Z",
      "type": "log|metric|trace|llm",
      "body": {}
    }
  ],
  "config_override": {
    "ncd_threshold": 0.3,
    "compressor": "zstd"
  }
}
```

### Config Override Options

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `baseline_size` | integer | 30 | Number of events in baseline (demo max: 30) |
| `window_size` | integer | 10 | Size of sliding detection window (demo max: 10) |
| `ncd_threshold` | float | 0.3 | Minimum NCD to flag anomaly (0-1) |
| `p_value_threshold` | float | 0.05 | Maximum p-value for significance |
| `compressor` | string | "zstd" | Compression algorithm (see below) |

### Compression Algorithms

| Algorithm | Speed | Use Case |
|-----------|-------|----------|
| `zstd` | Fast | **Default** - Best balance of speed and compression |
| `lz4` | Fastest | High-throughput streaming (10-20x faster) |
| `zlab` | Medium | Deterministic zlib-based compression |
| `gzip` | Slow | Universal compatibility |

**Recommendation**: Use `zstd` (default) for most cases. Use `lz4` for high-volume real-time streaming where speed is critical.

### Parameters

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `events` | array | **Required** | Array of events to analyze (1-50 events) |
| `events[].timestamp` | string (ISO 8601) | Recommended | Event timestamp |
| `events[].type` | string | Recommended | Event type: `log`, `metric`, `trace`, or `llm` |
| `events[].body` | object | **Required** | Event payload (any JSON object) |
| `config_override` | object | Optional | Override detection parameters |

## Response

### Success (HTTP 200)

```json
{
  "success": true,
  "total_events": 4,
  "anomaly_count": 1,
  "processing_time": "125ms",
  "compression_algo": "zstd",
  "anomalies": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "index": 3,
      "metrics": {
        "ncd": 0.72,
        "compression_ratio": 1.41,
        "entropy_change": 0.13,
        "p_value": 0.004,
        "confidence": 0.996
      },
      "event": {"latency": 950},
      "why": "Significant latency spike detected",
      "detected": true
    }
  ],
  "request_id": "req_abc123",
  "demo": {
    "message": "This is a demo response. Sign up for full access with persistence, history, and evidence bundles.",
    "remaining_calls": 9,
    "limit_per_minute": 10,
    "max_events_per_request": 50,
    "signup_url": "https://driftlock.net/#signup"
  }
}
```

## Examples

### Basic Example

```bash
curl -X POST https://api.driftlock.net/v1/demo/detect \
  -H "Content-Type: application/json" \
  -d '{
    "events": [
      {"body": {"latency": 120}},
      {"body": {"latency": 125}},
      {"body": {"latency": 118}},
      {"body": {"latency": 950}}
    ]
  }'
```

### With LZ4 for High Throughput

```bash
curl -X POST https://api.driftlock.net/v1/demo/detect \
  -H "Content-Type: application/json" \
  -d '{
    "events": [...],
    "config_override": {
      "compressor": "lz4"
    }
  }'
```

### Python Example

```python
import requests

response = requests.post(
    "https://api.driftlock.net/v1/demo/detect",
    json={
        "events": [
            {"body": {"cpu": 45, "memory": 2048}},
            {"body": {"cpu": 47, "memory": 2100}},
            {"body": {"cpu": 44, "memory": 2050}},
            {"body": {"cpu": 99, "memory": 8000}},  # Anomaly
        ]
    }
)

data = response.json()
print(f"Found {data['anomaly_count']} anomalies")
print(f"Remaining demo calls: {data['demo']['remaining_calls']}")

for anomaly in data['anomalies']:
    print(f"  - Index {anomaly['index']}: {anomaly['why']}")
```

## Error Responses

### 400 Bad Request - Too Many Events

```json
{
  "error": {
    "code": "invalid_argument",
    "message": "demo limited to 50 events per request (got 100). Sign up for unlimited access",
    "request_id": "req_abc123"
  }
}
```

### 429 Rate Limit Exceeded

```json
{
  "error": {
    "code": "rate_limit_exceeded",
    "message": "Demo rate limit exceeded. Sign up for unlimited access.",
    "request_id": "req_abc123",
    "retry_after_seconds": 60
  }
}
```

## Demo vs Full API

| Feature | Demo | Full API |
|---------|------|----------|
| Authentication | None | API key required |
| Events per request | 50 | 256 |
| Requests per minute | 10 | 60-1000+ |
| Anomaly persistence | No | Yes |
| Evidence bundles | No | Yes |
| Stream management | No | Yes |

## Ready for More?

1. **[Sign up for free](https://driftlock.net/#signup)** - Get 10,000 events/month free
2. **[Get your API key](../getting-started/authentication.md)** - From your dashboard
3. **[Use the full API](./detect.md)** - With persistence, history, and compliance features

---

**Next**: [Full /v1/detect endpoint →](./detect.md)

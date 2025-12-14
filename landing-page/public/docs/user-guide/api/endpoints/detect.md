# POST /v1/detect

Run synchronous anomaly detection on a batch of events. Returns anomalies inline.

## Endpoint

```
POST https://api.driftlock.net/v1/detect
```

## Authentication

`X-Api-Key` header with a `stream` or `admin` role key.

### Headers
```
Content-Type: application/json
X-Api-Key: YOUR_API_KEY
```

## Request body

```json
{
  "stream_id": "string",
  "events": [
    {
      "timestamp": "2025-01-01T10:00:00Z",
      "type": "log|metric|trace|llm",
      "body": {},
      "attributes": {},
      "idempotency_key": "string",
      "sequence": 1234567890
    }
  ],
  "config_override": {
    "baseline_size": 400,
    "window_size": 50,
    "hop_size": 10,
    "ncd_threshold": 0.3,
    "p_value_threshold": 0.05,
    "compressor": "zstd"
  }
}
```

### Fields

| Field | Required | Description |
| --- | --- | --- |
| `stream_id` | Optional | Stream to use; defaults to key-bound stream if omitted |
| `events` | **Yes** | Array of events (1–256) |
| `events[].timestamp` | Recommended | ISO-8601 timestamp |
| `events[].type` | Recommended | `log`, `metric`, `trace`, or `llm` |
| `events[].body` | **Yes** | JSON payload to evaluate |
| `events[].attributes` | Optional | Additional metadata |
| `events[].idempotency_key` | Optional | Deduplicate repeated sends |
| `events[].sequence` | Optional | Ordering hint |
| `config_override` | Optional | Override detection thresholds for this call |

### Config overrides

| Field | Default | Description |
| --- | --- | --- |
| `baseline_size` | 400 | Events used to build baseline |
| `window_size` | 50 | Sliding window size |
| `hop_size` | 10 | Window slide step |
| `ncd_threshold` | 0.30 | Minimum NCD to flag anomaly |
| `p_value_threshold` | 0.05 | Maximum p-value to accept |
| `compressor` | `zstd` | `zstd`, `lz4`, `gzip`, `openzl` |

## Response

```json
{
  "success": true,
  "batch_id": "batch_abc123",
  "stream_id": "c4e6f7a8-1234-5678-90ab-cdef12345678",
  "total_events": 4,
  "anomaly_count": 1,
  "processing_time": "125ms",
  "compression_algo": "zstd",
  "anomalies": [
    {
      "id": "anom_xyz789",
      "index": 3,
      "metrics": {
        "ncd": 0.72,
        "compression_ratio": 1.41,
        "entropy_change": 0.13,
        "p_value": 0.004,
        "confidence": 0.996
      },
      "event": {"body": {"latency": 950}},
      "why": "Significant latency spike detected",
      "detected": true
    }
  ],
  "request_id": "req_abc123"
}
```

### Response fields

| Field | Description |
| --- | --- |
| `anomaly_count` | Number of anomalies detected |
| `processing_time` | Server-side processing time |
| `compression_algo` | Compressor used (may fallback) |
| `anomalies[]` | Array of anomalies with metrics and explanation |
| `request_id` | Use when contacting support |

## Limits

| Limit | Value |
| --- | --- |
| Events per request | 256 |
| Payload size | 10 MB |
| Request rate | Plan-based (see rate limits) |

## Examples

### Basic detection (cURL)

```bash
curl -X POST https://api.driftlock.net/v1/detect \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -d '{
    "stream_id": "default",
    "events": [
      {"timestamp": "2025-01-01T10:03:00Z", "type": "metric", "body": {"latency_ms": 950}}
    ]
  }'
```

### With overrides (Python)

```python
import os
import requests

api_key = os.environ["DRIFTLOCK_API_KEY"]
resp = requests.post(
    "https://api.driftlock.net/v1/detect",
    headers={"X-Api-Key": api_key, "Content-Type": "application/json"},
    json={
        "stream_id": "production-logs",
        "events": [
            {"body": {"message": "SQL injection attempt"}, "type": "log"},
            {"body": {"message": "normal"}, "type": "log"}
        ],
        "config_override": {
            "ncd_threshold": 0.35,
            "p_value_threshold": 0.02,
            "compressor": "lz4"
        }
    },
    timeout=15,
)
resp.raise_for_status()
print(resp.json())
```

## Error responses

- **400 invalid_argument** — missing `events`, invalid JSON, or >256 events
- **401 unauthorized** — missing/invalid API key
- **429 rate_limit_exceeded** — respect `retry_after_seconds` and batch requests

## Related

- [Demo endpoint (no auth)](./demo.md)
- [List anomalies](./anomalies.md)
- [Error codes](../errors.md)

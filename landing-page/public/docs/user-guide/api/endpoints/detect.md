# POST /v1/detect

Run synchronous anomaly detection on a batch of events. This endpoint processes events immediately and returns detected anomalies in the response.

## Endpoint

```
POST https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/detect
```

## Authentication

Requires API key with `stream` or `admin` role:
```
X-Api-Key: YOUR_API_KEY
```

## Request

### Headers
```
Content-Type: application/json
X-Api-Key: YOUR_API_KEY
```

### Body Schema

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
    "ncd_threshold": 0.3,
    "p_value_threshold": 0.05,
    "compressor": "zstd"
  }
}
```

### Parameters

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `stream_id` | string | Optional | Stream to use for detection. If omitted, uses default stream or API key's bound stream |
| `events` | array | **Required** | Array of events to analyze (1-256 events) |
| `events[].timestamp` | string (ISO 8601) | **Required** | Event timestamp |
| `events[].type` | string | **Required** | Event type: `log`, `metric`, `trace`, or `llm` |
| `events[].body` | object | **Required** | Event payload (any JSON object) |
| `events[].attributes` | object | Optional | Additional metadata |
| `events[].idempotency_key` | string | Optional | Unique key for deduplication |
| `events[].sequence` | integer | Optional | Sequence number for ordering |
| `config_override` | object | Optional | Override default detection parameters |

### Config Override Options

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `baseline_size` | integer | 400 | Number of events in baseline |
| `window_size` | integer | 50 | Size of sliding detection window |
| `hop_size` | integer | 10 | Window slide amount |
| `ncd_threshold` | float | 0.3 | Minimum NCD to flag anomaly (0-1) |
| `p_value_threshold` | float | 0.05 | Maximum p-value for significance |
| `permutation_count` | integer | 1000 | Number of permutations for p-value |
| `compressor` | string | "zstd" | Compression algorithm: `zstd`, `lz4`, `gzip`, `openzl` |

## Response

### Success (HTTP 200)

```json
{
  "success": true,
  "batch_id": "batch_abc123",
  "stream_id": "c4e6f7a8-1234-5678-90ab-cdef12345678",
  "total_events": 100,
  "anomaly_count": 3,
  "processing_time": "245ms",
  "compression_algo": "zstd",
  "fallback_from_algo": "",
  "anomalies": [
    {
      "id": "anom_xyz789",
      "index": 42,
      "metrics": {
        "ncd": 0.72,
        "compression_ratio": 1.41,
        "entropy_change": 0.13,
        "p_value": 0.004,
        "confidence": 0.996,
        "baseline_compressed_size": 2048,
        "window_compressed_size": 2890,
        "joined_compressed_size": 4120
      },
      "event": {
        "timestamp": "2025-01-01T10:42:00Z",
        "type": "metric",
        "body": {"latency": 950}
      },
      "why": "Significant latency spike detected compared to baseline",
      "detected": true
    }
  ],
  "request_id": "req_abc123"
}
```

### Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `success` | boolean | Whether detection completed successfully |
| `batch_id` | string | Unique ID for this detection batch |
| `stream_id` | string | UUID of the stream used |
| `total_events` | integer | Number of events processed |
| `anomaly_count` | integer | Number of anomalies detected |
| `processing_time` | string | Time taken to process |
| `compression_algo` | string | Compression algorithm used |
| `fallback_from_algo` | string | If fallback occurred (e.g., openzl→zstd) |
| `anomalies` | array | Detected anomalies with details |
| `request_id` | string | Unique request identifier |

### Anomaly Object

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique anomaly ID |
| `index` | integer | Position in the events array |
| `metrics.ncd` | float | Normalized Compression Distance (0-1) |
| `metrics.compression_ratio` | float | Compression ratio |
| `metrics.entropy_change` | float | Entropy delta from baseline |
| `metrics.p_value` | float | Statistical significance (0-1) |
| `metrics.confidence` | float | Confidence score (1 - p_value) |
| `event` | object | The anomalous event |
| `why` | string | Plain English explanation |
| `detected` | boolean | Always true for anomalies |

## Examples

### Basic Detection

```bash
curl -X POST https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/detect \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: YOUR_API_KEY" \
  -d '{
    "stream_id": "default",
    "events": [
      {
        "timestamp": "2025-01-01T10:00:00Z",
        "type": "log",
        "body": {"message": "User login successful", "user_id": "123"}
      },
      {
        "timestamp": "2025-01-01T10:01:00Z",
        "type": "log",
        "body": {"message": "SQL injection attempt detected!", "user_id": "456"}
      }
    ]
  }'
```

### With Config Override

```bash
curl -X POST https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/detect \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: YOUR_API_KEY" \
  -d '{
    "stream_id": "production-logs",
    "events": [...],
    "config_override": {
      "ncd_threshold": 0.4,
      "p_value_threshold": 0.01,
      "compressor": "lz4"
    }
  }'
```

### Metric Data

```bash
curl -X POST https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/detect \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: YOUR_API_KEY" \
  -d '{
    "stream_id": "metrics",
    "events": [
      {
        "timestamp": "2025-01-01T10:00:00Z",
        "type": "metric",
        "body": {
          "cpu_usage": 45.2,
          "memory_mb": 2048,
          "latency_ms": 120
        },
        "attributes": {
          "host": "server-01",
          "region": "us-east-1"
        }
      }
    ]
  }'
```

## Error Responses

### 400 Bad Request - Invalid Input
```json
{
  "error": {
    "code": "invalid_argument",
    "message": "events array is required",
    "request_id": "req_abc123"
  }
}
```

### 401 Unauthorized - Invalid API Key
```json
{
  "error": {
    "code": "unauthorized",
    "message": "Invalid or missing API key",
    "request_id": "req_abc123"
  }
}
```

### 429 Rate Limit Exceeded
```json
{
  "error": {
    "code": "rate_limit_exceeded",
    "message": "Rate limit exceeded",
    "request_id": "req_abc123",
    "retry_after_seconds": 30
  }
}
```

## Limits

| Limit | Developer | Starter | Pro |
|-------|-----------|---------|-----|
| Events per request | 256 | 256 | 256 |
| Request size | 10 MB | 10 MB | 50 MB |
| Requests per minute | 60 | 300 | 1000+ |

## Best Practices

### Batch Events
```javascript
// Good: Batch multiple events
const events = collectEvents(); // 50-100 events
await detect(events);

// Avoid: One event per request
for (const event of events) {
  await detect([event]); // Too many requests!
}
```

### Use Idempotency Keys
```javascript
const event = {
  timestamp: new Date().toISOString(),
  type: "log",
  body: {...},
  idempotency_key: `evt_${Date.now()}_${randomId}` // Prevent duplicates
};
```

### Handle Errors
```python
import requests
from requests.adapters import HTTPAdapter
from requests.packages.urllib3.util.retry import Retry

session = requests.Session()
retry = Retry(
    total=3,
    status_forcelist=[429, 500, 502, 503, 504],
    backoff_factor=1
)
adapter = HTTPAdapter(max_retries=retry)
session.mount('https://', adapter)

response = session.post(
    'https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/detect',
    headers={'X-Api-Key': api_key},
    json=payload
)
```

## Code Examples

- **[Python](../examples/python-examples.md#detect)** - Python client implementation
- **[Node.js](../examples/node-examples.md#detect)** - JavaScript/TypeScript examples
- **[cURL](../examples/curl-examples.md#detect)** - Complete cURL examples

## Related Endpoints

- **[GET /v1/anomalies](./anomalies.md)** - Query detected anomalies
- **[GET /v1/anomalies/{id}](./anomaly-detail.md)** - Get anomaly details
- **POST /v1/streams/{id}/events** - Async event ingestion

---

**Next**: [Query anomalies →](./anomalies.md)

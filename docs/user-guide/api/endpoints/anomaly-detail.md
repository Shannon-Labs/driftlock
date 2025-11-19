# GET /v1/anomalies/{id}

Retrieve detailed information about a specific anomaly, including full event data, baseline comparison, and evidence.

## Endpoint

```
GET https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/anomalies/{id}
```

## Authentication

Requires API key with `stream` or `admin` role:
```
X-Api-Key: YOUR_API_KEY
```

## Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | string | Yes | The unique ID of the anomaly (e.g., `anom_xyz789`) |

## Response

### Success (HTTP 200)

```json
{
  "id": "anom_xyz789",
  "stream_id": "c4e6f7a8-...",
  "timestamp": "2025-01-01T10:42:00Z",
  "detected_at": "2025-01-01T10:42:05Z",
  "metrics": {
    "ncd": 0.72,
    "p_value": 0.004,
    "confidence": 0.996,
    "compression_ratio": 1.41,
    "entropy_change": 0.13,
    "baseline_size": 400,
    "window_size": 50
  },
  "event": {
    "timestamp": "2025-01-01T10:42:00Z",
    "type": "metric",
    "body": {"latency": 950},
    "attributes": {"host": "server-01"}
  },
  "explanation": {
    "summary": "Significant latency spike detected",
    "details": "Event compressed poorly relative to baseline (NCD 0.72). Entropy increased by 13%.",
    "similar_events": []
  },
  "evidence": {
    "baseline_sample": [...],
    "window_sample": [...]
  },
  "status": "new"
}
```

### Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique anomaly ID |
| `stream_id` | string | ID of the stream this anomaly belongs to |
| `metrics` | object | Detailed detection metrics |
| `event` | object | The full anomalous event data |
| `explanation` | object | Human-readable explanation and details |
| `evidence` | object | Samples of baseline and window data for comparison |
| `status` | string | Current status (`new`, `viewed`, `archived`) |

## Examples

### Get Anomaly Details
```bash
curl "https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/anomalies/anom_xyz789" \
  -H "X-Api-Key: YOUR_API_KEY"
```

## Error Responses

### 404 Not Found
```json
{
  "error": {
    "code": "not_found",
    "message": "Anomaly not found"
  }
}
```

### 403 Forbidden
```json
{
  "error": {
    "code": "forbidden",
    "message": "You do not have permission to view this anomaly"
  }
}
```

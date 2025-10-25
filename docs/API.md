# Driftlock API Documentation

Version: 1.0.0

Base URL: `http://localhost:8080/v1`

All endpoints require authentication via Bearer token unless specified otherwise.

## Authentication

Include API key in Authorization header:
```
Authorization: Bearer YOUR_API_KEY
```

## Endpoints

### Health & Status

#### `GET /healthz`
Health check endpoint (no auth required).

**Response**: `200 OK`
```
ok
```

#### `GET /readyz`
Readiness check - verifies database connectivity (no auth required).

**Response**: `200 OK`
```
ready
```

#### `GET /v1/version`
Get API version (no auth required).

**Response**: `200 OK`
```json
{
  "version": "1.0.0"
}
```

### Anomalies

#### `GET /v1/anomalies`
List anomalies with filtering and pagination.

**Query Parameters**:
- `stream_type` (optional): Filter by stream type (`logs`, `metrics`, `traces`, `llm`)
- `status` (optional): Filter by status (`pending`, `acknowledged`, `dismissed`, `investigating`)
- `min_ncd_score` (optional): Minimum NCD score (0-1)
- `max_p_value` (optional): Maximum p-value (0-1)
- `start_time` (optional): ISO 8601 timestamp
- `end_time` (optional): ISO 8601 timestamp
- `only_significant` (optional): `true` to show only statistically significant anomalies
- `limit` (optional): Page size (default: 50)
- `offset` (optional): Pagination offset (default: 0)

**Response**: `200 OK`
```json
{
  "anomalies": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "timestamp": "2025-10-25T10:30:00Z",
      "stream_type": "logs",
      "ncd_score": 0.45,
      "p_value": 0.023,
      "status": "pending",
      "glass_box_explanation": "Anomaly detected: NCD=0.450, p=0.023, compression ratio dropped from 3.20x to 1.85x due to data structure changes",
      "compression_baseline": 3.20,
      "compression_window": 1.85,
      "compression_combined": 2.52,
      "compression_ratio_change": -42.2,
      "confidence_level": 0.977,
      "is_statistically_significant": true,
      "tags": ["logs", "severity:high"],
      "created_at": "2025-10-25T10:30:05Z",
      "updated_at": "2025-10-25T10:30:05Z"
    }
  ],
  "total": 142,
  "limit": 50,
  "offset": 0,
  "has_more": true
}
```

#### `GET /v1/anomalies/:id`
Get a specific anomaly by ID.

**Response**: `200 OK`
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "timestamp": "2025-10-25T10:30:00Z",
  "stream_type": "logs",
  "ncd_score": 0.45,
  "p_value": 0.023,
  "status": "pending",
  "glass_box_explanation": "Anomaly detected...",
  "detailed_explanation": "CBAD Analysis: NCD=0.450 (statistically significant), confidence=97.7%, baseline entropy=5.42 bits/byte, window entropy=4.21 bits/byte, compression change=-42.2%",
  "compression_baseline": 3.20,
  "compression_window": 1.85,
  "compression_combined": 2.52,
  "compression_ratio_change": -42.2,
  "baseline_entropy": 5.42,
  "window_entropy": 4.21,
  "entropy_change": -1.21,
  "confidence_level": 0.977,
  "is_statistically_significant": true,
  "baseline_data": { /* JSONB payload */ },
  "window_data": { /* JSONB payload */ },
  "metadata": {
    "baseline_entropy": 5.42,
    "window_entropy": 4.21,
    "processing_time_ms": 23
  },
  "tags": ["logs", "severity:high"],
  "created_at": "2025-10-25T10:30:05Z",
  "updated_at": "2025-10-25T10:30:05Z"
}
```

#### `PATCH /v1/anomalies/:id/status`
Update anomaly status (acknowledge, dismiss, investigate).

**Request Body**:
```json
{
  "status": "acknowledged",
  "notes": "Investigated - known issue with log rotation"
}
```

**Response**: `200 OK`
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "status": "acknowledged",
  "acknowledged_by": "admin",
  "acknowledged_at": "2025-10-25T11:00:00Z",
  "notes": "Investigated - known issue with log rotation",
  ...
}
```

#### `GET /v1/anomalies/:id/export`
Export anomaly as evidence bundle (JSON).

**Response**: `200 OK`
```json
{
  "anomaly": { /* full anomaly object */ },
  "exported_at": "2025-10-25T11:30:00Z",
  "exported_by": "admin",
  "version": "1.0.0",
  "signature": "a3f5c8b2...",
  "additional_metadata": {
    "export_format": "driftlock-evidence-v1",
    "compliance_frameworks": ["DORA", "NIS2", "AI Act"]
  }
}
```

### Configuration

#### `GET /v1/config`
Get current detection configuration.

**Response**: `200 OK`
```json
{
  "id": 1,
  "ncd_threshold": 0.3,
  "p_value_threshold": 0.05,
  "baseline_size": 100,
  "window_size": 50,
  "hop_size": 10,
  "stream_overrides": {
    "logs": {
      "ncd_threshold": 0.25
    }
  },
  "is_active": true,
  "created_at": "2025-10-25T00:00:00Z",
  "updated_at": "2025-10-25T00:00:00Z"
}
```

#### `PATCH /v1/config`
Update detection configuration (admin only).

**Request Body**:
```json
{
  "ncd_threshold": 0.25,
  "p_value_threshold": 0.01,
  "baseline_size": 150,
  "notes": "Increased sensitivity for log anomalies"
}
```

**Response**: `200 OK`
```json
{
  "id": 1,
  "ncd_threshold": 0.25,
  "p_value_threshold": 0.01,
  "baseline_size": 150,
  ...
}
```

### Analytics

#### `GET /v1/analytics/summary`
Get statistical summary of anomalies.

**Query Parameters**:
- `start_time` (optional): ISO 8601 timestamp
- `end_time` (optional): ISO 8601 timestamp

**Response**: `200 OK`
```json
{
  "total_anomalies": 1523,
  "anomalies_by_stream_type": {
    "logs": 842,
    "metrics": 421,
    "traces": 260
  },
  "anomalies_by_status": {
    "pending": 143,
    "acknowledged": 1201,
    "dismissed": 179
  },
  "significant_anomalies": 1312,
  "average_ncd_score": 0.38,
  "average_p_value": 0.028,
  "average_compression_change": -31.5,
  "time_range": {
    "start": "2025-10-01T00:00:00Z",
    "end": "2025-10-25T23:59:59Z"
  }
}
```

#### `GET /v1/analytics/compression-timeline`
Get compression ratio timeline.

**Query Parameters**:
- `stream_type` (optional): Filter by stream type
- `start_time` (optional): ISO 8601 timestamp
- `end_time` (optional): ISO 8601 timestamp

**Response**: `200 OK`
```json
{
  "timeline": [
    {
      "timestamp": "2025-10-25T10:00:00Z",
      "baseline_ratio": 3.2,
      "window_ratio": 1.85,
      "compression_change": -42.2
    }
  ],
  "total": 1000
}
```

#### `GET /v1/analytics/ncd-heatmap`
Get NCD score heatmap by stream type and hour.

**Response**: `200 OK`
```json
{
  "heatmap": [
    {
      "stream_type": "logs",
      "hour": 14,
      "average_ncd": 0.42,
      "count": 23
    }
  ],
  "start_time": "2025-10-25T00:00:00Z",
  "end_time": "2025-10-25T23:59:59Z"
}
```

### Real-time Streaming

#### `GET /v1/stream/anomalies`
Server-Sent Events (SSE) endpoint for real-time anomaly notifications.

**Response**: `text/event-stream`

Events:
- `connected` - Initial connection established
- `anomaly` - New anomaly detected
- `heartbeat` - Keepalive (every 15 seconds)

```
event: connected
data: {"client_id": "client-1729850400"}

event: anomaly
data: {"id":"123e4567...","stream_type":"logs",...}

event: heartbeat
data: {"timestamp":"2025-10-25T10:30:15Z"}
```

### Performance Metrics

#### `GET /v1/metrics/performance`
Get real-time performance metrics.

**Response**: `200 OK`
```json
{
  "sse_connections": 42,
  "database": {
    "open": 15,
    "in_use": 8,
    "idle": 7
  }
}
```

## Error Responses

All errors follow RFC 7807 Problem Details format:

```json
{
  "type": "https://docs.driftlock.io/errors/not-found",
  "title": "Anomaly Not Found",
  "status": 404,
  "detail": "Anomaly with ID 123e4567... does not exist",
  "instance": "/v1/anomalies/123e4567..."
}
```

**Status Codes**:
- `200` - Success
- `201` - Created
- `400` - Bad Request (invalid input)
- `401` - Unauthorized (missing/invalid API key)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found
- `429` - Too Many Requests (rate limited)
- `500` - Internal Server Error
- `503` - Service Unavailable

## Rate Limiting

Default limits:
- **100 requests/minute** per API key
- **1000 concurrent SSE connections** per server

Rate limit headers:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1729850460
```

## Pagination

List endpoints support pagination:
- `limit` - Number of items per page (max: 1000, default: 50)
- `offset` - Number of items to skip

Response includes:
```json
{
  "total": 1523,
  "limit": 50,
  "offset": 100,
  "has_more": true
}
```

## Filtering

Timestamps use ISO 8601 format:
```
2025-10-25T10:30:00Z
2025-10-25T10:30:00.123Z
2025-10-25T10:30:00+02:00
```

Boolean parameters accept `true` or `false` (case-insensitive).

## Client Examples

### cURL

```bash
# List recent anomalies
curl -H "Authorization: Bearer YOUR_API_KEY" \
  "http://localhost:8080/v1/anomalies?limit=10&only_significant=true"

# Update anomaly status
curl -X PATCH \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"status":"acknowledged","notes":"Reviewed"}' \
  "http://localhost:8080/v1/anomalies/123e4567.../status"

# Stream real-time anomalies
curl -N -H "Authorization: Bearer YOUR_API_KEY" \
  "http://localhost:8080/v1/stream/anomalies"
```

### JavaScript

```javascript
// List anomalies
const response = await fetch('http://localhost:8080/v1/anomalies', {
  headers: {
    'Authorization': 'Bearer YOUR_API_KEY'
  }
});
const data = await response.json();

// SSE stream
const eventSource = new EventSource('http://localhost:8080/v1/stream/anomalies');
eventSource.addEventListener('anomaly', (event) => {
  const anomaly = JSON.parse(event.data);
  console.log('New anomaly:', anomaly);
});
```

### Python

```python
import requests

# List anomalies
headers = {'Authorization': 'Bearer YOUR_API_KEY'}
response = requests.get('http://localhost:8080/v1/anomalies', headers=headers)
anomalies = response.json()

# SSE stream
import sseclient
messages = sseclient.SSEClient('http://localhost:8080/v1/stream/anomalies', headers=headers)
for msg in messages:
    if msg.event == 'anomaly':
        print(f"New anomaly: {msg.data}")
```

## Versioning

API uses URL versioning (`/v1/`). Breaking changes will increment the major version (`/v2/`).

Current version: **v1** (stable)

# Driftlock API Documentation

## Overview

The Driftlock API provides RESTful endpoints for managing anomaly detection in OpenTelemetry data streams. All endpoints return JSON responses and use standard HTTP status codes.

**Base URL**: `http://localhost:8080`
**API Version**: `v1`

## Authentication

API authentication is handled via API keys or JWT tokens.

```bash
curl -H "Authorization: Bearer YOUR_API_KEY" http://localhost:8080/v1/anomalies
```

## Rate Limiting

- **Default**: 100 requests per minute per IP
- **Burst**: Up to 200 requests
- When rate limited, you'll receive a `429 Too Many Requests` response

## Endpoints

### Health & Status

#### `GET /healthz` - Health check
#### `GET /readyz` - Readiness check  
#### `GET /v1/version` - Get API version

### Event Ingestion

#### `POST /v1/events` - Ingest OpenTelemetry events

**Request**:
```json
{
  "timestamp": "2025-10-25T12:00:00Z",
  "value": 42.5,
  "metadata": {"service": "api-gateway"}
}
```

**Response**: `200 OK`

### Anomalies

#### `GET /v1/anomalies` - List anomalies

**Query Parameters**:
- `limit` (int): Max results (1-1000, default 50)
- `offset` (int): Pagination offset
- `stream_type` (string): Filter by type (logs, metrics, traces, llm)
- `status` (string): Filter by status
- `min_ncd_score` (float): Min NCD score (0-1)
- `max_p_value` (float): Max p-value (0-1)
- `only_significant` (bool): Only significant anomalies

**Response**: `200 OK`
```json
{
  "anomalies": [...],
  "total": 142,
  "limit": 50,
  "offset": 0,
  "has_more": true
}
```

#### `GET /v1/anomalies/:id` - Get specific anomaly
#### `POST /v1/anomalies` - Create anomaly (internal)
#### `PATCH /v1/anomalies/:id/status` - Update anomaly status

**Request**:
```json
{
  "status": "acknowledged",
  "notes": "Investigated"
}
```

### Real-Time Streaming

#### `GET /v1/stream/anomalies` - SSE stream for real-time anomalies

**Event Types**: connected, anomaly, heartbeat

**Example**:
```javascript
const eventSource = new EventSource('/v1/stream/anomalies');
eventSource.addEventListener('anomaly', (e) => {
  const anomaly = JSON.parse(e.data);
  console.log('New anomaly:', anomaly);
});
```

### Export

#### `GET /v1/anomalies/:id/export` - Export anomaly evidence

## HTTP Status Codes

| Code | Meaning |
|------|---------|
| 200 | OK |
| 201 | Created |
| 400 | Bad Request |
| 401 | Unauthorized |
| 404 | Not Found |
| 429 | Rate Limit Exceeded |
| 500 | Internal Server Error |

For full API specification, see the complete documentation.

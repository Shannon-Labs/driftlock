# Streaming API Guide

Driftlock supports real-time anomaly notifications via Server-Sent Events (SSE).

## Endpoint

`GET /v1/stream/anomalies`

### Authentication

Requires `X-Api-Key` header or `api_key` query parameter (for browser `EventSource`).

### Response Format

The stream returns standard SSE events.

- **Event: `ping`**
  - Data: Timestamp (ms)
  - Sent every 30 seconds to keep connection alive.

- **Event: `anomaly`**
  - Data: JSON object of the anomaly.

#### Anomaly Object

```json
{
  "id": "uuid",
  "stream_id": "uuid",
  "ncd": 0.45,
  "confidence": 0.98,
  "explanation": "Unusual pattern in tool usage...",
  "detected_at": "2025-03-01T12:00:00Z",
  "metrics": { ... }
}
```

## Client Example (JavaScript)

```javascript
const apiKey = "your-api-key";
const url = `https://api.driftlock.net/v1/stream/anomalies?api_key=${apiKey}`;

const evtSource = new EventSource(url);

evtSource.addEventListener("anomaly", (e) => {
  const anomaly = JSON.parse(e.data);
  console.log("New Anomaly:", anomaly);
});

evtSource.onopen = () => {
  console.log("Connected to Driftlock Stream");
};
```


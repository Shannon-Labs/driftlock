# API Monitoring

Monitor your microservices for performance degradation and unusual traffic patterns.

## Scenario

You have a critical API service. You want to know if:
- Latency spikes for specific endpoints.
- Error rates increase.
- The request payload structure changes (e.g., clients sending deprecated fields).

## Implementation

### 1. Middleware Setup

The easiest way is to use our [Express](../integrations/express.md) or [Go](../sdks/go.md) middleware.

```javascript
// Express.js
app.use(driftlock({
  streamId: 'api-gateway',
  includeBody: true, // Monitor payload structure
  includeHeaders: ['user-agent', 'content-type'] // Monitor client types
}));
```

### 2. What to Monitor

We recommend monitoring the following metrics as part of the event body:

```json
{
  "method": "POST",
  "path": "/v1/checkout",
  "status": 200,
  "duration_ms": 145,
  "request_size": 1024,
  "response_size": 512,
  "client_id": "client_abc"
}
```

### 3. Detecting Anomalies

Driftlock will automatically learn the "normal" behavior for your API.

#### Latency Anomalies
If `/v1/checkout` usually takes 100-200ms, and suddenly a request takes 800ms, Driftlock will flag it. Unlike static thresholds (e.g., "alert if > 500ms"), Driftlock adapts. If your API is always slow (e.g., 1000ms), it won't alert unless it becomes *slower* or *faster*.

#### Status Code Anomalies
If your API is 99.9% `200 OK`, a sudden burst of `500` or `401` errors represents a high compression distance from the baseline and will be flagged.

#### Payload Drift
If you deploy a breaking change and clients start sending invalid JSON or missing fields, Driftlock will detect the structural change in the request body.

## Alert Routing

You can route alerts based on the `path` or `service` to different teams.

1. **Create separate streams** for each service: `service-auth`, `service-payment`.
2. **Configure Webhooks** for each stream to different Slack channels.

## Comparison with Datadog/New Relic

| Feature | Traditional APM | Driftlock |
|---------|-----------------|-----------|
| **Setup** | Complex agents, heavy config | Simple middleware / API call |
| **Thresholds** | Manual (e.g., "Latency > 500ms") | Automatic (Learns baseline) |
| **Context** | Metrics only (usually) | Full payload analysis |
| **Cost** | High ($$$) | Low ($) |

Driftlock complements your existing APM by providing **explainable anomalies** ("Why is this request slow?") rather than just charts.

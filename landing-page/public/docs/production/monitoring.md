# Monitoring Driftlock

Monitor the health and performance of your Driftlock integration.

## Dashboard Metrics

The [Driftlock Dashboard](https://driftlock.web.app/dashboard) provides real-time charts for:
- **Request Volume**: Total events processed over time.
- **Anomaly Rate**: Percentage of events flagged as anomalous.
- **Latency**: API response times.
- **Errors**: 4xx and 5xx error rates.

## Prometheus Integration

You can export Driftlock metrics to your own Prometheus instance.

### Node.js

If you are using `prom-client`, the Driftlock SDK can automatically register metrics.

```javascript
const client = new DriftlockClient({ ... });
const register = require('prom-client').register;

client.enablePrometheus(register);
```

Exposed metrics:
- `driftlock_requests_total{status="success|error"}`
- `driftlock_anomalies_detected_total`
- `driftlock_request_duration_seconds`

## Health Checks

Implement a health check in your application to verify connectivity to Driftlock.

```javascript
app.get('/health/driftlock', async (req, res) => {
  try {
    // Lightweight check (e.g., verify key)
    await client.verify(); 
    res.status(200).send('OK');
  } catch (err) {
    res.status(503).send('Driftlock Unavailable');
  }
});
```

## Alerting

Set up alerts to be notified when things go wrong.

### Recommended Alerts

1. **High Error Rate**: > 5% of requests failing.
2. **High Latency**: P99 > 500ms for 5 minutes.
3. **Anomaly Spike**: > 100 anomalies detected in 1 minute (could indicate a system-wide issue).

You can configure these alerts in the Dashboard under **Settings > Alerts**.

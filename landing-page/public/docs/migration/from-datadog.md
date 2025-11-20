# Migrating from DataDog

Switching from DataDog's Watchdog or Anomaly Detection to Driftlock? Here's how to make the transition smooth.

## Conceptual Differences

| Feature | DataDog | Driftlock |
|---------|---------|-----------|
| **Philosophy** | Metric-based (Time Series) | Event-based (Compression) |
| **Setup** | Agent installation | API / Middleware |
| **Detection** | Statistical thresholds (ARIMA, etc.) | Normalized Compression Distance (NCD) |
| **Data** | Aggregated metrics | Raw events (JSON) |

## Step 1: Replace the Agent

DataDog relies on the `datadog-agent` running on your host. Driftlock is lighter and runs within your application or as a sidecar.

### Before (DataDog)

```yaml
# datadog.yaml
api_key: ...
logs_enabled: true
```

### After (Driftlock)

Install the SDK in your application.

```javascript
// Node.js
const client = new DriftlockClient({ apiKey: '...' });
```

## Step 2: Porting Monitors

In DataDog, you define monitors like `avg(last_5m):avg:system.cpu.idle{host:host0} < 10`.

In Driftlock, you don't define thresholds. You stream the data, and Driftlock learns what "normal" CPU usage looks like.

### Before (DataDog Monitor)
- **Query**: `avg(last_5m):avg:my.metric > 100`
- **Alert**: "Metric is too high"

### After (Driftlock)
- **Code**: `client.detect({ streamId: 'my-metric', events: [{ value: 105 }] })`
- **Result**: Driftlock alerts if `105` is anomalous compared to recent history.

## Step 3: Log Management

If you use DataDog Logs, you can pipe them to Driftlock for anomaly detection.

### Before
- Logs sent to DataDog via Agent.
- "Log Anomaly Detection" enabled in UI.

### After
- Pipe logs to Driftlock CLI or API.
- [Log Analysis Guide](../examples/log-analysis.md)

## Cost Comparison

- **DataDog**: Charges by host, per million log events, and per custom metric. Costs can spiral with high cardinality.
- **Driftlock**: Simple event-based pricing. No per-host fees. High cardinality is free (it actually helps the model!).

## FAQ

**Q: Can I use both?**
A: Yes! Many customers use DataDog for dashboards/infrastructure monitoring and Driftlock for specialized anomaly detection on critical business events.

**Q: Do I need to backfill data?**
A: No. Driftlock learns quickly. Send ~50-100 events to establish a baseline.

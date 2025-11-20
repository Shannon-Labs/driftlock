# Migrating from New Relic

Transitioning from New Relic Applied Intelligence to Driftlock? Here is what you need to know.

## Key Differences

| Feature | New Relic | Driftlock |
|---------|-----------|-----------|
| **Focus** | APM & Infrastructure | Business Logic & Event Anomalies |
| **Configuration** | Heavy UI configuration | Code-first configuration |
| **Alerting** | Threshold-based + AI | Purely behavior-based |

## Migration Strategy

### 1. Identify Critical Flows

New Relic is great for "Is the server up?". Driftlock is great for "Is the server doing the right thing?".

Identify the business-critical flows you want to monitor:
- Checkout process
- User authentication
- Payment processing
- Critical background jobs

### 2. Instrument with Driftlock SDK

Replace New Relic custom events with Driftlock events.

**New Relic (Python):**
```python
newrelic.agent.record_custom_event('Purchase', {'amount': 100, 'item': 'sku-123'})
```

**Driftlock (Python):**
```python
await client.detect(
    stream_id='purchases', 
    events=[{'type': 'purchase', 'body': {'amount': 100, 'item': 'sku-123'}}]
)
```

### 3. Alerting

New Relic Alerts are often set up in the UI. With Driftlock, you configure [Webhooks](../tools/webhooks.md) to send alerts to your incident management system (PagerDuty, OpsGenie).

## Advantages of Switching

1. **No Training Period**: New Relic often requires weeks of data to train its AI models. Driftlock's compression-based approach works with as few as 50 events.
2. **Explainability**: New Relic might say "Anomaly Detected". Driftlock says "Anomaly Detected: 'amount' is significantly higher than usual".
3. **Cost**: Stop paying for high-cardinality custom metrics.

## Coexistence

You can run Driftlock alongside New Relic. Use New Relic for stack traces and performance profiling, and use Driftlock for detecting logical anomalies in your data.

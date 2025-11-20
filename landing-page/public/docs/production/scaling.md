# Scaling Driftlock

Learn how to scale your Driftlock integration to handle millions of events.

## Rate Limits

Driftlock enforces rate limits to ensure fair usage.

| Plan | Rate Limit | Burst |
|------|------------|-------|
| Developer | 60 req/min | 10 |
| Starter | 1,000 req/min | 50 |
| Pro | Custom | Custom |

### Handling Rate Limits

If you exceed your limit, the API will return `429 Too Many Requests`. Your application should handle this gracefully by:
1. **Logging the error**.
2. **Dropping the event** (if real-time is critical).
3. **Retrying with backoff** (if data completeness is critical).

## Multi-Tenancy

If you are a B2B SaaS platform, you likely have multiple customers (tenants).

### Strategy 1: One Stream per Tenant (Recommended)

Create a separate `stream_id` for each tenant (e.g., `tenant-123`, `tenant-456`).
- **Pros**: Isolated baselines. An anomaly in Tenant A won't affect Tenant B.
- **Cons**: Managing many stream IDs.

```javascript
driftlock.detect({
  streamId: `tenant-${tenantId}`,
  events: [...]
});
```

### Strategy 2: Shared Stream

Use a single stream for all tenants.
- **Pros**: Simple setup.
- **Cons**: No isolation. A spike in one tenant might look normal if the overall volume is high.

## High Volume Architecture

For very high volumes (> 10k events/sec), we recommend a decoupled architecture.

1. **App Servers**: Push events to a message queue (Kafka, SQS, RabbitMQ).
2. **Worker Service**: Consumes events, batches them (e.g., 500 events/batch), and sends to Driftlock.

This ensures that API latency or rate limits never impact your main application flow.

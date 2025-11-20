# Troubleshooting

Solutions to common issues when using Driftlock.

## API Errors

### 401 Unauthorized
- **Cause**: Invalid or missing API key.
- **Solution**: Check your `X-Api-Key` header. Verify the key in the dashboard.

### 403 Forbidden
- **Cause**: API key does not have permission for this action or stream.
- **Solution**: Check if you are using a Restricted Key that doesn't allow `write` access.

### 429 Too Many Requests
- **Cause**: You exceeded your plan's rate limit.
- **Solution**: Implement backoff/retry logic or upgrade your plan.

### 500 Internal Server Error
- **Cause**: Something went wrong on our end.
- **Solution**: Check the [Status Page](https://status.driftlock.io) and retry later.

## Detection Issues

### "No anomalies detected" (False Negative)
- **Cause**: The event is too similar to the baseline.
- **Solution**:
    1. Check if your baseline data is diverse enough.
    2. Wait for more events to accumulate (we need ~50 events to form a stable baseline).

### "Everything is an anomaly" (False Positive)
- **Cause**: The stream has no consistent pattern (high entropy).
- **Solution**:
    1. Ensure you are sending consistent data structures.
    2. Check if you are mixing different types of events in the same `stream_id`.

## SDK Issues

### Timeout Errors
- **Cause**: Network latency or Driftlock API is slow.
- **Solution**: Increase the timeout in your client configuration (default is usually 10s).

```javascript
const client = new DriftlockClient({ timeout: 30000 }); // 30 seconds
```

## Support

If you're still stuck, contact us:
- **Email**: [support@driftlock.io](mailto:support@driftlock.io)
- **GitHub**: [Open an Issue](https://github.com/Shannon-Labs/driftlock/issues)

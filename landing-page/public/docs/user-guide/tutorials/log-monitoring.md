# Tutorial: Log Monitoring

Detect suspicious log entries (e.g., SQL injection attempts) with Driftlock.

## Prerequisites
- API key
- Stream ID (use `logs-production` or `default`)

## Sample payload

```json
{
  "stream_id": "logs-production",
  "events": [
    {"body": {"level": "info", "message": "Login success", "user": "alice"}},
    {"body": {"level": "info", "message": "Login success", "user": "bob"}},
    {"body": {"level": "info", "message": "Password reset start", "user": "carol"}},
    {"body": {"level": "warn", "message": "SQL injection attempt detected", "path": "/login?id=' OR 1=1--"}}
  ]
}
```

## Run detection

```bash
curl -X POST https://api.driftlock.net/v1/detect \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -d @payload.json
```

Expected: the 4th event flags as anomalous with a high NCD and a “SQL injection attempt” explanation.

## Next steps
- Bind a `logs-production` key for least-privilege access.
- Add attributes (service, region) to improve explainability.
- Set profile to `balanced` or `strict` based on alert fatigue.

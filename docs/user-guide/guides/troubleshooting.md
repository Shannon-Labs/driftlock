# Troubleshooting

Common issues and solutions when using Driftlock's anomaly detection API.

## "No anomalies detected"

Most common causes and fixes:

| Symptom | Cause | Solution |
| --- | --- | --- |
| `anomaly_count: 0` always | Not enough events | Send 450+ events to build baseline |
| `detection_ready: false` | Stream still calibrating | Check `detection_events_needed` field |
| `status: calibrating` | Baseline incomplete (< 50 events) | Keep sending events |
| Anomalies expected but not detected | Threshold too strict | Switch to `sensitive` profile |

### Check calibration status

```bash
curl "https://api.driftlock.net/v1/detect" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"stream_id": "my-stream", "events": [{"msg": "test"}]}'
```

Response includes:
- `status`: "calibrating" or "ready" (50 event threshold)
- `detection_ready`: true/false (varies by profile - see table below)
- `calibration.events_ingested`: events received so far
- `calibration.detection_events_needed`: events until detection active

### Detection thresholds by profile

| Profile | Baseline | Window | Detection Ready |
| --- | --- | --- | --- |
| sensitive | 200 | 30 | 230 events |
| balanced (default) | 400 | 50 | 450 events |
| strict | 800 | 100 | 900 events |

Thresholds also adjust automatically with:
- **Adaptive windowing** - based on stream characteristics
- **Auto-tuning** - learns from feedback
- **Custom settings** - per-stream configuration

## JSON format errors

| Error | Cause | Fix |
| --- | --- | --- |
| `400 Bad Request` | Invalid JSON | Validate with `jq` before sending |
| `events must be array` | Events not in array | Wrap: `{"events": [...]}` |
| `empty events array` | No events sent | Include at least 1 event |

### Valid event formats

```json
// Good - JSON object
{"events": [{"message": "user logged in", "level": "info"}]}

// Good - nested objects work
{"events": [{"body": {"cpu": 45.2, "memory": 1024}}]}

// Bad - plain strings
{"events": ["error occurred"]}

// Fix for plain text - wrap it
{"events": [{"text": "error occurred"}]}
```

## Performance issues

| Symptom | Likely cause | Solution |
| --- | --- | --- |
| Slow responses (> 2s) | Large batch size | Reduce to 100-500 events/request |
| Timeouts | Network or rate limit | Check `Retry-After` header |
| High latency on first request | Cold start | Stream stays warm after first request |

### Optimal batch sizes

| Use case | Recommended batch | Reason |
| --- | --- | --- |
| Real-time alerting | 10-50 events | Low latency |
| Bulk processing | 500-1000 events | Throughput |
| Demo/testing | 10-100 events | Quick iteration |

### Rate limits

| Tier | Requests/min | Events/request |
| --- | --- | --- |
| Demo (no auth) | 10 | 200 |
| Radar | 100 | 1,000 |
| Tensor | 500 | 5,000 |
| Orbit | 2,000 | 10,000 |

## Authentication errors

| Error | Code | Fix |
| --- | --- | --- |
| `missing api key` | 401 | Add `X-Api-Key` header |
| `invalid api key` | 401 | Check key in dashboard |
| `subscription inactive` | 403 | Verify billing status |
| `rate limit exceeded` | 429 | Wait or upgrade plan |

### Get a new API key

1. Go to [driftlock.net](https://driftlock.net)
2. Sign in → Dashboard → API Keys
3. Click "Create Key" and copy the value

## Debug mode

Add `X-Request-ID` header to track requests:

```bash
curl "https://api.driftlock.net/v1/detect" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -H "X-Request-ID: debug-$(date +%s)" \
  -H "Content-Type: application/json" \
  -d '{"events": [...]}'
```

Response includes your `request_id` for support inquiries.

## Still stuck?

1. Check `calibration.message` in the response - it explains current status
2. Try the demo endpoint first to validate your payload format
3. Contact support at support@driftlock.net with your `request_id`

**Related:** [Detection Profiles](./detection-profiles.md), [Auto-Tuning](./auto-tuning.md)

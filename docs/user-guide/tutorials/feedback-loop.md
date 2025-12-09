# Tutorial: Use the Feedback Loop

Submit feedback on anomalies to improve sensitivity via auto-tuning.

## Prerequisites

- API key
- Auto-tuning enabled on your stream (see [Auto-Tuning](../guides/auto-tuning.md))
- Recent anomalies to review

## 1) Enable auto-tune (if needed)

```bash
curl -X PATCH "https://api.driftlock.net/v1/streams/logs-production/profile" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"auto_tune_enabled": true}'
```

## 2) List anomalies to review

```bash
curl "https://api.driftlock.net/v1/anomalies?stream_id=logs-production&limit=10" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY"
```

## 3) Submit feedback

```bash
# False positive
curl -X POST "https://api.driftlock.net/v1/anomalies/{id}/feedback" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"feedback_type": "false_positive"}'

# Confirmed anomaly
curl -X POST "https://api.driftlock.net/v1/anomalies/{id}/feedback" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"feedback_type": "confirmed", "reason": "Production incident"}'

# Dismissed (neutral)
curl -X POST "https://api.driftlock.net/v1/anomalies/{id}/feedback" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"feedback_type": "dismissed"}'
```

## 4) Check tuning status

```bash
curl "https://api.driftlock.net/v1/streams/logs-production/tuning" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY"
```

Look for `tune_history` entries and feedback statistics.

## Tips

- Provide honest feedback regularly; aim for 20+ samples before judging results.
- If tuning drifts, switch back to a preset profile (e.g., `balanced`) and re-enable auto-tune.

**Next:** [Detection Profiles](../guides/detection-profiles.md)

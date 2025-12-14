# Auto-Tuning

Auto-tuning adjusts detection thresholds using your feedback to reduce false positives while keeping recall high.

## How it works

```
Detect anomaly → You review → Submit feedback → Thresholds adjust (if enabled)
```

- Feedback types: `false_positive`, `confirmed`, `dismissed`.
- Requires sufficient feedback (20+ samples) before adjusting.
- Adjustments are bounded and cooled down to avoid thrashing.

## Enable/disable

```bash
curl -X PATCH "https://api.driftlock.net/v1/streams/{stream_id}/profile" \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -d '{"auto_tune_enabled": true}'
```

When auto-tune updates thresholds, the stream moves to the `custom` profile.

## Submit feedback

```bash
# False positive
curl -X POST "https://api.driftlock.net/v1/anomalies/{id}/feedback" \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -d '{"feedback_type": "false_positive"}'

# Confirmed
curl -X POST "https://api.driftlock.net/v1/anomalies/{id}/feedback" \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -d '{"feedback_type": "confirmed", "reason": "Production incident"}'

# Dismissed (neutral)
curl -X POST "https://api.driftlock.net/v1/anomalies/{id}/feedback" \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -d '{"feedback_type": "dismissed"}'
```

## View tuning history

```bash
curl "https://api.driftlock.net/v1/streams/{stream_id}/tuning" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY"
```

Returns feedback stats and the last adjustments applied.

## Guidelines

- Provide consistent, honest feedback; avoid marking everything as confirmed or FP.
- Expect changes after meaningful samples (20+). There is a cooldown between adjustments.
- Switch back to a preset profile anytime to reset thresholds.

**Related:** [Detection Profiles](./detection-profiles.md), [Feedback Loop tutorial](../tutorials/feedback-loop.md)

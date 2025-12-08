# Tutorial: Using the Feedback Loop

This tutorial demonstrates how to submit feedback on anomalies and watch auto-tuning improve your detection over time.

## Prerequisites

- A Driftlock API key
- Auto-tuning enabled on your stream (see [Profiles Tutorial](./profiles-tutorial.md))
- Some detected anomalies to review

## Step 1: Enable Auto-Tuning

First, ensure auto-tuning is enabled:

```bash
curl -X PATCH https://driftlock.net/api/v1/streams/logs-production/profile \
  -H "X-Api-Key: dlk_your_key_here" \
  -H "Content-Type: application/json" \
  -d '{"auto_tune_enabled": true}'
```

## Step 2: Get Recent Anomalies

List anomalies to review:

```bash
curl "https://driftlock.net/api/v1/anomalies?stream_id=logs-production&limit=10" \
  -H "X-Api-Key: dlk_your_key_here"
```

Response:

```json
{
  "anomalies": [
    {
      "id": "anom_abc123",
      "stream_id": "logs-production",
      "ncd_score": 0.42,
      "pvalue": 0.02,
      "confidence": 0.85,
      "event_summary": "{\"level\":\"error\",\"message\":\"Connection timeout\"}",
      "created_at": "2025-12-07T10:30:00Z"
    },
    {
      "id": "anom_def456",
      "stream_id": "logs-production",
      "ncd_score": 0.31,
      "pvalue": 0.04,
      "confidence": 0.72,
      "event_summary": "{\"level\":\"warn\",\"message\":\"High memory usage\"}",
      "created_at": "2025-12-07T10:25:00Z"
    }
  ]
}
```

## Step 3: Review and Submit Feedback

### Mark a False Positive

If `anom_def456` wasn't actually anomalous (maybe high memory is normal):

```bash
curl -X POST https://driftlock.net/api/v1/anomalies/anom_def456/feedback \
  -H "X-Api-Key: dlk_your_key_here" \
  -H "Content-Type: application/json" \
  -d '{"feedback_type": "false_positive"}'
```

Response:

```json
{
  "status": "feedback_recorded",
  "anomaly_id": "anom_def456",
  "feedback_type": "false_positive"
}
```

### Confirm a Real Anomaly

If `anom_abc123` was a real issue:

```bash
curl -X POST https://driftlock.net/api/v1/anomalies/anom_abc123/feedback \
  -H "X-Api-Key: dlk_your_key_here" \
  -H "Content-Type: application/json" \
  -d '{"feedback_type": "confirmed"}'
```

## Step 4: Continue Providing Feedback

The algorithm needs **20+ feedback samples** before making adjustments. Continue reviewing and submitting feedback as anomalies come in.

A good practice:
- Check anomalies daily
- Mark obvious false positives immediately
- Confirm real issues when you investigate them
- Use `dismissed` for edge cases that are technically anomalous but not actionable

## Step 5: Check Tuning Status

After providing feedback, check if adjustments have been made:

```bash
curl https://driftlock.net/api/v1/streams/logs-production/tuning \
  -H "X-Api-Key: dlk_your_key_here"
```

Response:

```json
{
  "stream_id": "logs-production",
  "feedback_stats": {
    "total_feedback": 25,
    "false_positives": 8,
    "confirmed": 15,
    "dismissed": 2,
    "false_positive_rate": 0.32
  },
  "tune_history": [
    {
      "tune_type": "ncd",
      "old_value": 0.30,
      "new_value": 0.38,
      "reason": "high_false_positive_rate",
      "confidence": 0.68,
      "created_at": "2025-12-07T11:00:00Z"
    }
  ]
}
```

In this example:
- 32% false positive rate (above 15% target)
- NCD threshold increased from 0.30 to 0.38
- Reason: `high_false_positive_rate`

## Step 6: Observe Improvement

After tuning, you should see fewer false positives. The algorithm will continue adjusting every hour (cooldown period) as you provide more feedback.

## Understanding Tune History

| Field | Meaning |
|-------|---------|
| `tune_type` | What was adjusted (`ncd` or `pvalue`) |
| `old_value` | Previous threshold |
| `new_value` | New threshold |
| `reason` | Why adjustment was made |
| `confidence` | Algorithm confidence (higher = more certain) |

### Common Reasons

- **high_false_positive_rate**: FP rate >15%, thresholds increased
- **low_detection_rate**: FP rate <2.5%, thresholds may decrease
- **fp_ncd_boundary**: False positives cluster near threshold

## Tips for Better Results

1. **Be honest**: Don't mark real anomalies as false positives just to reduce alerts
2. **Be consistent**: Review anomalies regularly
3. **Give it time**: 20+ samples needed, and 1-hour cooldown between adjustments
4. **Check the math**: High confidence scores indicate reliable adjustments

## Resetting Tuning

If auto-tuning goes wrong, reset by switching to a preset profile:

```bash
curl -X PATCH https://driftlock.net/api/v1/streams/logs-production/profile \
  -H "X-Api-Key: dlk_your_key_here" \
  -H "Content-Type: application/json" \
  -d '{"profile": "balanced"}'
```

This resets to profile defaults. You can then re-enable auto-tuning to start fresh.

## What's Next?

- [Auto-Tuning Guide](../guides/auto-tuning.md) - Algorithm details
- [Detection Profiles Guide](../guides/detection-profiles.md) - Understanding profiles

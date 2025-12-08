# Auto-Tuning

Auto-tuning automatically adjusts detection thresholds based on your feedback. When you mark anomalies as false positives or confirm them, Driftlock learns and refines its sensitivity for your stream.

## How It Works

```
Detect Anomaly → Review in Dashboard → Submit Feedback → Auto-Adjust Thresholds
```

1. **Detection**: An anomaly is detected and flagged
2. **Review**: You examine the anomaly and determine if it's real
3. **Feedback**: You mark it as `false_positive`, `confirmed`, or `dismissed`
4. **Learning**: The algorithm analyzes feedback patterns
5. **Adjustment**: Thresholds are adjusted to reduce future false positives

## Enabling Auto-Tune

Auto-tune is **enabled by default** for new streams. If you need to toggle it:

```bash
curl -X PATCH https://driftlock.net/api/v1/streams/{stream_id}/profile \
  -H "X-Api-Key: dlk_..." \
  -H "Content-Type: application/json" \
  -d '{"auto_tune_enabled": true}'
```

## Submitting Feedback

### Mark as False Positive

When an anomaly isn't actually anomalous:

```bash
curl -X POST https://driftlock.net/api/v1/anomalies/{anomaly_id}/feedback \
  -H "X-Api-Key: dlk_..." \
  -H "Content-Type: application/json" \
  -d '{"feedback_type": "false_positive"}'
```

### Confirm Anomaly

When an anomaly is correctly detected:

```bash
curl -X POST https://driftlock.net/api/v1/anomalies/{anomaly_id}/feedback \
  -H "X-Api-Key: dlk_..." \
  -H "Content-Type: application/json" \
  -d '{"feedback_type": "confirmed"}'
```

### Dismiss (Neutral)

When an anomaly is technically correct but not actionable:

```bash
curl -X POST https://driftlock.net/api/v1/anomalies/{anomaly_id}/feedback \
  -H "X-Api-Key: dlk_..." \
  -H "Content-Type: application/json" \
  -d '{"feedback_type": "dismissed"}'
```

## Algorithm Details

The algorithm learns from your feedback to reduce false positives while maintaining detection quality. Key points:

- Requires **20+ feedback samples** before adjusting
- Targets a **5% false positive rate**
- Has a **1-hour cooldown** between adjustments

<details>
<summary><strong>Show technical details</strong> (advanced)</summary>

### Target False Positive Rate

Auto-tune targets a **5% false positive rate**. This means:
- If your FP rate is above 15%: thresholds increase (less sensitive)
- If your FP rate is below 2.5% with confirmed anomalies: thresholds may decrease (more sensitive)

### Minimum Feedback Requirement

The algorithm requires **20+ feedback samples** before making adjustments. This prevents premature tuning based on limited data.

### Cooldown Period

After an adjustment, there's a **1-hour cooldown** before the next adjustment. This prevents thrashing from rapid feedback submissions.

### Adjustment Bounds

- **Learning rate**: 15% of the calculated delta
- **Max single adjustment**: 25% of current threshold
- **NCD bounds**: 0.10 to 0.80
- **P-value bounds**: 0.001 to 0.20

</details>

## Viewing Tuning History

Check what adjustments have been made:

```bash
curl https://driftlock.net/api/v1/streams/{stream_id}/tuning \
  -H "X-Api-Key: dlk_..."
```

Response:

```json
{
  "stream_id": "abc123",
  "feedback_stats": {
    "total_feedback": 45,
    "false_positives": 8,
    "confirmed": 32,
    "dismissed": 5,
    "false_positive_rate": 0.178
  },
  "tune_history": [
    {
      "tune_type": "ncd",
      "old_value": 0.30,
      "new_value": 0.35,
      "reason": "high_false_positive_rate",
      "confidence": 0.822,
      "created_at": "2025-12-07T15:30:00Z"
    }
  ]
}
```

### Tune Reasons

| Reason | Description |
|--------|-------------|
| `high_false_positive_rate` | FP rate >15%, increased thresholds |
| `low_detection_rate` | FP rate <2.5% with confirmed anomalies, decreased thresholds |
| `fp_ncd_boundary` | False positives cluster near threshold, adjusted boundary |
| `insufficient_feedback` | Not enough samples to tune |

## Profile Switching

When auto-tune adjusts thresholds, your stream automatically switches to the `custom` profile. This preserves your tuned values separate from the preset profiles.

You can switch back to a preset profile at any time, which will reset to the profile's defaults:

```bash
curl -X PATCH https://driftlock.net/api/v1/streams/{stream_id}/profile \
  -H "X-Api-Key: dlk_..." \
  -H "Content-Type: application/json" \
  -d '{"profile": "balanced"}'
```

## Best Practices

1. **Be consistent**: Provide feedback regularly for accurate tuning
2. **Be honest**: Mark true false positives, don't dismiss everything
3. **Give it time**: Wait for 20+ feedback samples before expecting changes
4. **Monitor history**: Check tuning history to understand adjustments
5. **Reset if needed**: Switch to a preset profile to reset tuning

## Next Steps

- [Detection Profiles Guide](./detection-profiles.md) - Understanding profile presets
- [Feedback Loop Tutorial](../tutorials/feedback-loop.md) - Hands-on feedback workflow

# Detection Profiles

Detection profiles provide pre-configured sensitivity presets for anomaly detection. Each profile balances detection rate against false positive rate for different use cases.

## Available Profiles

| Profile | NCD Threshold | P-Value | Baseline | Window | Best For |
|---------|--------------|---------|----------|--------|----------|
| **sensitive** | 0.20 | 0.10 | 200 | 30 | Critical systems, high detection rate |
| **balanced** | 0.30 | 0.05 | 400 | 50 | General use (default) |
| **strict** | 0.45 | 0.01 | 800 | 100 | Noisy data, minimize false positives |
| **custom** | user-defined | user-defined | varies | varies | Auto-tuned or manually configured |

## Choosing a Profile

### Sensitive

Use when:
- Missing an anomaly is costly (security, fraud detection)
- You have clean, well-structured data
- You can tolerate some false positives
- Systems are mission-critical

```bash
curl -X PATCH https://driftlock.net/api/v1/streams/{stream_id}/profile \
  -H "X-Api-Key: dlk_..." \
  -H "Content-Type: application/json" \
  -d '{"profile": "sensitive"}'
```

### Balanced (Default)

Use when:
- You want reasonable detection with manageable false positives
- Data has moderate noise levels
- You're starting out and want sensible defaults

This is the default profile for all new streams.

### Strict

Use when:
- False positives disrupt workflows (on-call alerts)
- Data is inherently noisy or variable
- You only want high-confidence anomalies
- Compliance requires documented certainty

```bash
curl -X PATCH https://driftlock.net/api/v1/streams/{stream_id}/profile \
  -H "X-Api-Key: dlk_..." \
  -H "Content-Type: application/json" \
  -d '{"profile": "strict"}'
```

### Custom

The `custom` profile uses thresholds you specify or that were set by auto-tuning. When auto-tune adjusts your thresholds, your stream automatically switches to the `custom` profile.

## API Reference

### Get Current Profile

```bash
curl https://driftlock.net/api/v1/streams/{stream_id}/profile \
  -H "X-Api-Key: dlk_..."
```

Response:

```json
{
  "stream_id": "abc123",
  "profile": "balanced",
  "auto_tune_enabled": false,
  "adaptive_window_enabled": false,
  "current_thresholds": {
    "ncd_threshold": 0.30,
    "pvalue_threshold": 0.05,
    "baseline_size": 400,
    "window_size": 50
  }
}
```

### Update Profile

```bash
curl -X PATCH https://driftlock.net/api/v1/streams/{stream_id}/profile \
  -H "X-Api-Key: dlk_..." \
  -H "Content-Type: application/json" \
  -d '{
    "profile": "sensitive",
    "auto_tune_enabled": true,
    "adaptive_window_enabled": true
  }'
```

### List All Profiles

```bash
curl https://driftlock.net/api/v1/profiles
```

Response:

```json
{
  "profiles": {
    "sensitive": {
      "name": "sensitive",
      "description": "Lower thresholds, more anomalies reported. Best for critical systems requiring high detection rates.",
      "ncd_threshold": 0.20,
      "pvalue_threshold": 0.10,
      "baseline_size": 200,
      "window_size": 30
    },
    "balanced": { ... },
    "strict": { ... },
    "custom": { ... }
  }
}
```

## Understanding Thresholds

### NCD Threshold

The Normalized Compression Distance threshold determines how different an event must be from the baseline to be flagged.

- **Lower values (0.20)**: More sensitive, flags smaller deviations
- **Higher values (0.45)**: Less sensitive, requires larger deviations

### P-Value Threshold

The statistical significance threshold for the permutation test.

- **Higher values (0.10)**: More permissive, easier to flag
- **Lower values (0.01)**: Stricter, requires stronger statistical evidence

## Next Steps

- [Auto-Tuning Guide](./auto-tuning.md) - Learn how feedback adjusts your thresholds
- [Adaptive Windowing Guide](./adaptive-windowing.md) - Automatic window sizing
- [Profiles Tutorial](../tutorials/profiles-tutorial.md) - Hands-on profile selection

# Detection Profiles

Preset sensitivity levels that tune how aggressively Driftlock flags anomalies.

## Profiles

| Profile | NCD threshold | P-value | Baseline | Window | Best for |
| --- | --- | --- | --- | --- | --- |
| `sensitive` | 0.20 | 0.10 | 200 | 30 | Critical systems; catch everything, more noise |
| `balanced` (default) | 0.30 | 0.05 | 400 | 50 | General workloads |
| `strict` | 0.45 | 0.01 | 800 | 100 | Noisy data; minimize false positives |
| `custom` | user-defined | user-defined | varies | varies | Auto-tuned or manually set |

Switching a profile immediately changes thresholds for the stream.

## Choose a profile

- **Sensitive:** security/fraud, low tolerance for missed anomalies, clean data.
- **Balanced:** start here; good trade-off for mixed workloads.
- **Strict:** noisy telemetry or alert fatigue; only high-confidence anomalies.

## API examples

### Get current profile
```bash
curl "https://api.driftlock.net/v1/streams/{stream_id}/profile" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY"
```

### Update profile
```bash
curl -X PATCH "https://api.driftlock.net/v1/streams/{stream_id}/profile" \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -d '{
    "profile": "sensitive",
    "auto_tune_enabled": true,
    "adaptive_window_enabled": true
  }'
```

### List available profiles
```bash
curl "https://api.driftlock.net/v1/profiles" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY"
```

## Tips

- Start with **balanced**, collect feedback, then move to **strict** or **sensitive** as needed.
- If you toggle auto-tuning, the stream switches to `custom` when thresholds are learned.
- Pair strict profiles with alerting channels; pair sensitive profiles with feedback loops to reduce noise.

**Related:** [Auto-Tuning](./auto-tuning.md), [Adaptive Windowing](./adaptive-windowing.md)

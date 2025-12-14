# Adaptive Windowing

Adaptive windowing automatically sizes your baseline and detection windows based on stream characteristics, so you don't have to tune them manually.

## Why adaptive windowing?

- Too small baseline → unstable detection, false positives
- Too large baseline → slow response to real changes
- Too small window → noisy
- Too large window → missed anomalies

Adaptive sizing considers event frequency, event size, pattern diversity, and entropy to pick balanced values.

## Toggle

```bash
curl -X PATCH "https://api.driftlock.net/v1/streams/{stream_id}/profile" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"adaptive_window_enabled": true}'
```

## Sizing bounds

| Parameter | Minimum | Maximum |
|-----------|---------|---------|
| Baseline | 100 events | 2,000 events |
| Window | 10 events | 200 events |
| Memory per stream | — | 50 MB |

## View computed sizes

```bash
curl "https://api.driftlock.net/v1/streams/{stream_id}/profile" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY"
```

Example response includes current window/baseline sizes and stream statistics.

## When to use manual sizing

- Compliance requires fixed windows
- Strong periodic patterns where fixed windows align better
- Load testing/benchmarks with controlled parameters

Disable if needed:

```bash
curl -X PATCH "https://api.driftlock.net/v1/streams/{stream_id}/profile" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"adaptive_window_enabled": false}'
```

## Pairing with profiles

Adaptive sizing works with any detection profile:
- **sensitive + adaptive:** aggressive detection with optimized windows
- **strict + adaptive:** conservative detection with optimized windows
- **custom + adaptive:** tuned thresholds plus optimized windows

**Related:** [Detection Profiles](./detection-profiles.md), [Auto-Tuning](./auto-tuning.md)

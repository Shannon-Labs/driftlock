# Adaptive Windowing

Adaptive windowing automatically sizes your baseline and detection windows based on stream characteristics. Instead of manually configuring window sizes, Driftlock analyzes your data patterns and computes optimal sizes.

## Why Adaptive Windowing?

Manual window sizing requires understanding your data:
- Too small baseline: unstable detection, false positives
- Too large baseline: slow response to real changes
- Too small window: noisy detection
- Too large window: missed anomalies

Adaptive windowing solves this by analyzing:
- **Event frequency**: How often events arrive
- **Event size**: How large your payloads are
- **Pattern diversity**: How varied your data patterns are
- **Entropy**: How random/structured your data is

## Enabling Adaptive Windowing

Adaptive windowing is **enabled by default** for new streams. If you need to toggle it:

```bash
curl -X PATCH https://driftlock.net/api/v1/streams/{stream_id}/profile \
  -H "X-Api-Key: dlk_..." \
  -H "Content-Type: application/json" \
  -d '{"adaptive_window_enabled": true}'
```

## Sizing Bounds

| Parameter | Minimum | Maximum |
|-----------|---------|---------|
| Baseline | 100 events | 2,000 events |
| Window | 10 events | 200 events |
| Memory per stream | - | 50 MB |

## How Sizing Works

The algorithm analyzes your stream characteristics and computes optimal sizes automatically. You don't need to understand the details—it just works.

<details>
<summary><strong>Show algorithm details</strong> (advanced)</summary>

### Factor 1: Event Frequency

Higher event rates suggest larger baselines for stability:

```
baseline = events_per_hour × 0.1
```

A stream receiving 1,000 events/hour gets a baseline around 100 events.

### Factor 2: Memory Constraints

Larger events require smaller windows to stay within memory budget:

```
max_events = (50 MB) / avg_event_size
baseline = max_events / 4  (reserve room for windows)
```

### Factor 3: Pattern Diversity

More varied patterns need larger baselines to capture the full range:

```
diversity_multiplier = 1.0 + (diversity_score × 0.5)
```

Highly diverse data gets up to 1.5× larger baselines.

### Factor 4: Entropy

High-entropy data (more random) needs more samples:

```
entropy_multiplier = 1.0 + (entropy - 6.0) × 0.1
```

Data with entropy above 6 bits/byte gets proportionally larger baselines.

### Final Calculation

The algorithm combines all factors:
1. Computes baseline from each factor
2. Averages frequency, complexity, and entropy baselines
3. Takes minimum with memory constraint
4. Window = baseline × 0.125 (1/8th ratio)

</details>

## Viewing Computed Sizes

Check your stream's adaptive sizes:

```bash
curl https://driftlock.net/api/v1/streams/{stream_id}/profile \
  -H "X-Api-Key: dlk_..."
```

Response includes computed sizes:

```json
{
  "stream_id": "abc123",
  "adaptive_window_enabled": true,
  "current_thresholds": {
    "baseline_size": 450,
    "window_size": 56
  },
  "stream_statistics": {
    "avg_events_per_hour": 2400,
    "avg_event_size_bytes": 512,
    "pattern_diversity_score": 0.35,
    "avg_baseline_entropy": 5.8
  }
}
```

## When to Use Manual Sizing

Disable adaptive windowing when:
- You have specific compliance requirements for window sizes
- Your data has known periodic patterns (use fixed windows aligned to periods)
- Testing or benchmarking with controlled parameters

To disable:

```bash
curl -X PATCH https://driftlock.net/api/v1/streams/{stream_id}/profile \
  -H "X-Api-Key: dlk_..." \
  -H "Content-Type: application/json" \
  -d '{"adaptive_window_enabled": false}'
```

## Combining with Profiles

Adaptive windowing works with any detection profile:
- **sensitive** + adaptive: Aggressive detection with optimized windows
- **strict** + adaptive: Conservative detection with optimized windows
- **custom** + adaptive: Your tuned thresholds with optimized windows

## Learning Period

Adaptive sizing improves over time:
- First batch: Uses profile defaults
- Subsequent batches: Computes characteristics, averages with history
- After ~10 batches: Stable, refined window sizes

## Next Steps

- [Detection Profiles Guide](./detection-profiles.md) - Understanding profile presets
- [Auto-Tuning Guide](./auto-tuning.md) - Threshold adjustment from feedback

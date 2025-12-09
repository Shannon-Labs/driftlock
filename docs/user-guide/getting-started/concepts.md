# Core Concepts

Understand how Driftlock detects anomalies using compression-based anomaly detection (CBAD).

## How Driftlock works (high level)

1) **Baseline** — First ~400 normal events are compressed to form what “normal” looks like.
2) **Window** — New events arrive in a sliding window (default 50, hop 10).
3) **Compression test** — We compare how well `baseline + window` compresses vs. baseline alone.
4) **Metrics** — We compute Normalized Compression Distance (NCD), entropy deltas, and compression ratios.
5) **Significance** — Permutation testing produces a p-value; we flag when NCD is high **and** p-value is below the threshold.
6) **Profiles/auto-tune** — Sensitivity presets or feedback-driven tuning adjust thresholds automatically.

## Baseline

- Built from your earliest normal events (default: 400).
- Should represent steady-state traffic; noisy or abnormal data here will reduce accuracy.
- Baselines are deterministic—same input + config yields the same baseline every time.

## Windows

- Sliding window of recent events (default size 50, hop 10).
- Compared against the baseline; each hop yields a detection decision.
- Adaptive windowing can resize automatically when enabled (see [Adaptive Windowing](../guides/adaptive-windowing.md)).

## Metrics we compute

- **NCD (Normalized Compression Distance):** Distance between baseline and window (0–1). Higher means more novel.
- **p-value:** Probability the observed NCD is due to chance (lower is stronger evidence).
- **confidence:** `1 - p_value` (how sure we are an anomaly is real).
- **compression ratio / entropy change:** Additional explainability signals.

### Default thresholds (balanced profile)

- `ncd_threshold`: **0.30**
- `p_value_threshold`: **0.05**
- `baseline_size`: **400**
- `window_size`: **50**
- `compressor`: **zstd** (use `lz4` for maximum speed)

See [Detection Profiles](../guides/detection-profiles.md) for `sensitive`, `balanced`, and `strict` presets.

## Sensitivity profiles

- **sensitive:** Lower NCD threshold, higher p-value threshold (catches more, more noise).
- **balanced (default):** Good trade-off for most workloads.
- **strict:** Higher NCD, lower p-value (fewer false positives, might miss subtle drift).
- **custom:** Set per-stream overrides or values learned from auto-tuning.

Switching profiles changes thresholds immediately; anomalies become more or less likely based on your tolerance for noise.

## Auto-tuning & feedback

- Submitting feedback (`false_positive`, `confirmed`) can steer thresholds when enabled.
- Feedback influences the `custom` profile for that stream.
- Use it after an initial period to lock in the right sensitivity. See [Auto-Tuning](../guides/auto-tuning.md).

## Explainability

Each anomaly includes:
- NCD, p-value, confidence, compression/entropy deltas.
- The original event and a short “why” message.
- When available: baseline/window evidence snapshots for audits.

## Determinism & compliance

- Runs are deterministic for the same input/config (seeded), enabling reproducible evidence.
- P-values provide statistical backing for alerts—important for audits (DORA/NIS2/AI Act).

## Best practices

- Keep baseline clean: avoid known incidents or load tests during baseline formation.
- Separate streams: don’t mix unrelated data sources (e.g., logs vs metrics) in one stream.
- Batch efficiently: 50–200 events per call keeps throughput high and minimizes rate limits.
- Use idempotency keys when replaying data.
- Start with **balanced**, then tune via profiles or `config_override` only when needed.

## Learn more

- [Detection Profiles](../guides/detection-profiles.md)
- [Auto-Tuning](../guides/auto-tuning.md)
- [Adaptive Windowing](../guides/adaptive-windowing.md)
- [API Reference](../api/rest-api.md)

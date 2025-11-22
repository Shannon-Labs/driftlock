# Core Concepts

Understanding how Driftlock detects anomalies will help you get the most out of the platform. This guide explains the key concepts and algorithms.

## What is Anomaly Detection?

Anomaly detection identifies unusual patterns in data that don't conform to expected behavior. Driftlock uses **compression-based anomaly detection (CBAD)**, which doesn't require training data or pre-defined patterns.

## How Driftlock Works

### 1. Compression-Based Detection

At the core of Driftlock is a simple but powerful idea: **normal data compresses well, anomalous data doesn't**.

When you send events to Driftlock:
1. We build a **baseline** from your initial events (typically first ~400 events)
2. New events are compared to this baseline using compression algorithms
3. Events that don't compress well relative to the baseline are flagged as anomalies

### 2. Normalized Compression Distance (NCD)

**NCD** measures how different a new event is from your baseline data.

```
NCD = (compressed(baseline + event) - min(compressed(baseline), compressed(event)))
      / max(compressed(baseline), compressed(event))
```

- **NCD ≈ 0**: Event is very similar to baseline (normal)
- **NCD ≈ 1**: Event is very different from baseline (anomalous)
- **Threshold**: By default, NCD > 0.3 is considered anomalous

### 3. Statistical Validation

To avoid false positives, Driftlock uses **permutation testing**:

1. We randomly shuffle your baseline data 1,000 times
2. Calculate NCD for each permutation
3. Compute a **p-value**: what's the probability of seeing this NCD by chance?

**p-value < 0.05** means the anomaly is statistically significant (less than 5% chance it's random).

### 4. Confidence Score

Combines NCD and p-value into a single metric:

```
confidence = 1 - p_value
```

- **0.95+**: Very confident this is an anomaly
- **0.90-0.95**: Likely an anomaly
- **< 0.90**: Low confidence, might be noise

## Key Terminology

### Baseline
The "normal" dataset that Driftlock learns from. By default:
- First **400 events** form the baseline
- Must be representative of normal behavior
- Updated periodically as new normal data arrives

### Window
A sliding set of recent events being analyzed:
- Default size: **50 events**
- Compared against the baseline
- Slides forward by **hop size** (default: 10 events)

### Stream
A logical grouping of related events:
- Each stream has its own baseline and configuration
- Examples: "production-logs", "payment-metrics", "iot-sensors"
- Isolates different data sources from each other

### Tenant
Your organization or account:
- All your streams, API keys, and usage belong to your tenant
- Multi-tenant isolation ensures your data stays private

### Compression Algorithms
Driftlock supports multiple compressors:
- **zstd** (default): Fast, good compression ratio
- **lz4**: Extremely fast, lower compression ratio
- **gzip**: Widely compatible, moderate speed
- **openzl** (optional): Advanced, best for structured data

Different compressors work better for different data types. You can override per-stream or per-detect call.

## How Driftlock is Different

### vs. Machine Learning Models
- ❌ ML: Requires training data
- ✅ Driftlock: Works immediately with no training

- ❌ ML: Black box, hard to explain
- ✅ Driftlock: Deterministic, explainable results

- ❌ ML: Needs labeled anomalies
- ✅ Driftlock: Unsupervised, no labels needed

### vs. Rule-Based Systems
- ❌ Rules: Need to define thresholds manually
- ✅ Driftlock: Automatically learns what's normal

- ❌ Rules: Miss novel anomalies
- ✅ Driftlock: Catches any deviation from baseline

- ❌ Rules: Hard to maintain as data evolves
- ✅ Driftlock: Adapts to changing baselines

## Determinism & Reproducibility

Driftlock's detection is **deterministic**:
- Same events + same configuration = same results every time
- Uses seeded random number generation
- Critical for compliance and debugging

Seed is derived from: `tenant_id + stream_id + config.seed`

This means you can re-run detection on historical data and get identical results.

## Explainability

When an anomaly is detected, Driftlock provides:

1. **Metrics**: NCD, p-value, confidence, compression ratios
2. **Evidence**: Baseline vs. window snapshots
3. **Explanation**: Plain English description (optional via Gemini integration)

Example explanation:
> "Latency spike detected: event compressed poorly relative to baseline. Entropy increased by 13%, indicating novel pattern. P-value 0.004 confirms statistical significance."

## Tuning Parameters

You can adjust detection sensitivity per stream:

| Parameter | Default | Description |
|-----------|---------|-------------|
| `baseline_size` | 400 | Events in baseline |
| `window_size` | 50 | Events in sliding window |
| `ncd_threshold` | 0.3 | Minimum NCD to flag anomaly |
| `p_value_threshold` | 0.05 | Maximum p-value for significance |
| `compressor` | "zstd" | Compression algorithm |

**More sensitive** (catch more anomalies, more false positives):
- Lower `ncd_threshold` (e.g., 0.2)
- Higher `p_value_threshold` (e.g., 0.10)

**Less sensitive** (fewer false positives, might miss subtle anomalies):
- Higher `ncd_threshold` (e.g., 0.4)
- Lower `p_value_threshold` (e.g., 0.01)

## Best Practices

### Baseline Quality
- Ensure baseline contains only **normal events**
- Minimum 200-400 events for statistical reliability
- If baseline is contaminated with anomalies, detection degrades

### Stream Separation
- Create separate streams for different event types
- Don't mix logs and metrics in the same stream
- Use stream-specific configurations for optimal detection

### Event Structure
- More structured data compresses better
- JSON is ideal (consistent schema)
- Include timestamps for temporal analysis

### Performance
- Batch events where possible (up to 256 events per request)
- Use async ingestion (`/v1/streams/{id}/events`) for high-throughput
- Monitor rate limits (`X-RateLimit-*` headers)

## Learn More

- **[API Reference](../api/rest-api.md)**: Full endpoint documentation
- **[Algorithms Deep Dive](../reference/algorithms.md)**: Mathematical details
- **[Tutorials](../tutorials/)**: Step-by-step guides
- **[Research Paper](https://arxiv.org/pdf/cs/0111054.pdf)**: NCD theory

---

**Next**: [Set up authentication and manage API keys →](./authentication.md)

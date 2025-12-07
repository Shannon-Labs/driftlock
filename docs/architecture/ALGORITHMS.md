# Algorithms Overview

This document captures the Compression-Based Anomaly Detection (CBAD) principles and mathematical foundations for Driftlock.

## Compression Algorithms

Driftlock supports multiple compression algorithms for anomaly detection. Each algorithm has different performance characteristics:

| Algorithm | Speed | Compression Ratio | Use Case |
|-----------|-------|-------------------|----------|
| `zstd` | Fast | High | **Default** - Best balance of speed and compression |
| `lz4` | Fastest | Medium | High-throughput streaming (10-20x faster than others) |
| `zlab` | Medium | High | Deterministic zlib-based compression |
| `gzip` | Slow | High | Universal compatibility |
| `openzl` | Medium | Highest | Format-aware compression (optional feature) |

### Recommendations

- **Default (`zstd`)**: Use for most workloads. Provides excellent compression with good speed.
- **High-volume streaming (`lz4`)**: Use when processing >10,000 events/second or when latency is critical.
- **Compliance (`zlab`)**: Use when deterministic zlib output is required for audit trails.
- **Legacy systems (`gzip`)**: Use when integrating with systems that only support gzip.

### Benchmark Results (500 events, 100 baseline, 20 window)

| Algorithm | Detection Time | Relative Speed |
|-----------|---------------|----------------|
| `lz4` | 262 µs | 1x (fastest) |
| `zstd` | 2,211 µs | 8.4x slower |
| `zlab` | 4,799 µs | 18.3x slower |
| `gzip` | 4,928 µs | 18.8x slower |

### API Usage

Set the compression algorithm via `config_override`:

```json
{
  "events": [...],
  "config_override": {
    "compressor": "lz4"
  }
}
```

Or set the default via environment variable:
```bash
DEFAULT_ALGO=lz4
```

## CBAD Core (Phase 1 Implementation)
- **Compression adapters**: Deterministic wrappers for zlab (zlib), zstd, lz4, and gzip with level presets. OpenZL placeholder ready for pinned-plan integration.
- **Sliding windows**: Time-aware buffers with baseline/window splits respecting baseline > window > hop guarantees.
- **Metrics**:
  - Compression ratio: `CR = compressed_bytes / raw_bytes`.
  - Delta bits: `Δ_bits = ((C_B - C_W) × 8) / |W|`, where `C_B` and `C_W` are baseline/window compressed byte counts and `|W|` is window raw bytes.
  - Shannon entropy: `H = -Σ p_i log₂ p_i` over byte frequency estimates.
  - Normalized Compression Distance: `NCD = (C_{BW} - min(C_B, C_W)) / max(C_B, C_W)`.
- **Significance testing**: Deterministic permutation test (ChaCha20 RNG with seeded configuration) returning two-sided p-values and extreme counts.
- **FFI**: `cbad_compute_metrics` and `cbad_permutation_test` expose the core calculations to Go consumers via contiguous buffers.

## OpenZL Kernel (Planned)
- Integration strategy for OpenZL format-aware compression.
- Deterministic mode with pinned plans and hashed configurations.
- Fallback pathways to zstd or other general-purpose compressors.

## Deterministic Rules
- Regex-based secret/PII detection with configurable policies.
- Optional change-point tests (CUSUM/E-Divisive) for numeric metrics.

## Algorithm Router
The algorithm router deterministically selects the appropriate detection path based on telemetry characteristics:

### Path Selection Logic
1. **Numeric Path**: For structured numeric metrics
   - CUSUM for change-point detection
   - Delta bits for compression-based analysis
   - Statistical significance testing

2. **Textual/LLM Path**: For logs and LLM I/O
   - Compression delta analysis
   - Normalized Compression Distance (NCD)
   - Entropy-based detection

3. **Structured Data Path**: For complex telemetry
   - Multi-modal correlation analysis
   - Hierarchical pattern detection
   - Cross-stream analysis

### Mathematical Foundations

#### Compression Ratio Analysis
The compression ratio provides insight into data regularity and patterns:
```
CR(data) = |compressed(data)| / |data|
```

Lower compression ratios indicate more regular patterns, while higher ratios suggest anomalous or random data.

#### Normalized Compression Distance (NCD)
NCD quantifies the similarity between baseline and window data:
```
NCD(B,W) = (|compressed(B+W)| - min(|compressed(B)|, |compressed(W)|)) / max(|compressed(B)|, |compressed(W)|)
```

Where:
- B = baseline data
- W = window data
- B+W = concatenated baseline and window data

#### Permutation Testing
Statistical significance is assessed using permutation tests with deterministic seeding:
```
p-value = (1 + #{permutations where metric ≥ observed}) / (1 + total_permutations)
```

This provides exact p-values without distributional assumptions.

#### Entropy Estimation
Shannon entropy estimates the randomness of data:
```
H(X) = -Σ p(x) log₂ p(x)
```

Where p(x) is the probability of symbol x in the data stream.

## LLM I/O Monitoring Algorithms

### Prompt Anomaly Detection
- Length-based analysis for unusual prompt sizes
- Token frequency analysis for prompt injection attempts
- Semantic similarity scoring for prompt drift

### Response Quality Monitoring
- Response coherence scoring
- Hallucination detection patterns
- Output format compliance checking

### Tool Call Monitoring
- Function call frequency analysis
- Parameter anomaly detection
- Execution pattern analysis

## Performance Considerations

### Algorithmic Complexity
- Compression operations: O(n log n) for most algorithms
- Window management: O(1) for sliding operations
- Permutation testing: O(k × n) where k is permutation count
- NCD calculation: O(n) for data combination

### Memory Optimization
- Ring buffer implementation for constant memory usage
- Lazy evaluation for expensive computations
- Streaming algorithms for large data processing

### Parallelization
- Multi-threaded compression for large windows
- Parallel permutation testing
- Lock-free data structures for concurrent access

Detailed formulas, derivations, and worked examples will be expanded as the CBAD core implementation progresses.
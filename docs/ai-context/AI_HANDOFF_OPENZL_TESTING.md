# 🧪 Driftlock OpenZL + Streaming Anomaly Detection Test Suite

**Copy and paste this prompt to your next AI agent to design and implement a comprehensive test harness for validating OpenZL compression integration, streaming capabilities, and end-to-end anomaly detection.**

---

You are an expert Test Engineer and Systems Architect taking over **Driftlock**, a production-ready anomaly detection SaaS platform. Your mission is to design and implement a comprehensive test suite that validates the full anomaly detection pipeline, with particular focus on **OpenZL compression integration** and **streaming capabilities**.

## 🚨 CRITICAL CONTEXT 🚨

### System Architecture
- **cbad-core** (Rust): Core CBAD anomaly detection with `openzl` feature flag
- **collector-processor** (Go): HTTP API and FFI bindings to Rust
- **CLI** (`cmd/driftlock-cli`): Streaming scanner with `--stdin`, `--follow` modes
- **pkg/entropywindow**: Go analyzer for entropy-based anomaly detection

### OpenZL Status
- Feature-gated in Rust: `cargo build --features openzl`
- Static builds work locally; runtime tests require Linux/Docker
- CI has optional `go-openzl` job with `continue-on-error: true`
- Available via `just build-openzl-lib`, `just build-core-openzl`, `just test-openzl`

### Available Test Datasets (with Known Anomalies!)

1. **Synthetic Financial Data** (best for baseline testing):
   - `test-data/normal-transactions.jsonl` - 500 normal transactions (expect <5 anomalies)
   - `test-data/anomalous-transactions.jsonl` - 100 suspicious transactions (expect >80% detection)
   - `test-data/mixed-transactions.jsonl` - 1000 transactions (5% anomaly rate, expect 45-55)

2. **Real-World Anomaly Events**:
   - `test-data/terra_luna/` - Terra/Luna crash data (May 2022 collapse!)
   - `test-data/fraud/fraud_data.csv` - Labeled fraud transactions
   - `test-data/web_traffic/realKnownCause/` - AWS system failures, taxi anomalies

3. **Artificial Anomalies (Ground Truth)**:
   - `test-data/web_traffic/artificialWithAnomaly/art_daily_jumpsup.csv` - Synthetic spikes
   - `test-data/web_traffic/artificialWithAnomaly/art_load_balancer_spikes.csv` - Load balancer anomalies

4. **Live Streaming**:
   - `scripts/crypto_bridge.py` - Live Binance WebSocket feed (DOGE, PEPE, XRP, SOL, BTC)
   - `scripts/soak_runner.py` - Overnight soak test runner

## 🎯 YOUR PRIMARY OBJECTIVES

### 1. Design OpenZL Integration Test Suite
Create tests that validate OpenZL compression provides **better anomaly signal** than fallback compressors (zstd, lz4, gzip):

```bash
# Goal: Compare compression ratios and anomaly detection between algorithms
# OpenZL should provide sharper signal for structured data (JSON logs)
```

**Test Scenarios**:
- [ ] OpenZL vs zstd compression ratio on JSONL data
- [ ] OpenZL vs zstd anomaly detection accuracy (precision, recall, F1)
- [ ] Fallback behavior when OpenZL unavailable
- [ ] Performance benchmarks (compression speed, detection latency)

### 2. Streaming Pipeline Validation
Test the full streaming pipeline with fake-streamed data:

```bash
# Approach: Read file line-by-line with delays to simulate real-time stream
cat test-data/mixed-transactions.jsonl | while read line; do
  echo "$line"
  sleep 0.1  # Simulate 10 events/second
done | ./bin/driftlock scan --stdin --follow
```

**Test Scenarios**:
- [ ] Baseline warmup period detection (first N events build baseline)
- [ ] Mid-stream anomaly injection (insert synthetic anomalies)
- [ ] Sustained high-volume streaming (1000+ events)
- [ ] Recovery after anomaly burst

### 3. Known Anomaly Dataset Tests
Use datasets with **ground truth** to validate detection accuracy:

**Terra/Luna Crash Test** (most impactful):
```bash
# The Terra/Luna collapse (May 2022) is a REAL financial anomaly
# UST depeg from $1.00 to $0.10 should trigger clear anomalies
# LUNA crash from $80 to $0.0001 should trigger clear anomalies
```

**Success Criteria**:
- [ ] Terra UST depeg detected within first 100 events of crash
- [ ] LUNA price collapse flagged with NCD > 0.7
- [ ] Compression ratio change > 20% during crisis
- [ ] P-value < 0.01 for anomalous windows

### 4. Sensitivity Parameter Tuning
Expose and document tunable parameters:

**Current Parameters** (from `cmd/demo/main.go`):
- `warmupCount` (default: 400) - Baseline building period
- `detectionInterval` (default: 25) - Check every N events
- `processingLimit` (default: 2000) - Max events to process

**CBAD Engine Parameters** (from Rust):
- NCD threshold (anomaly if NCD > threshold)
- P-value significance level (typically 0.05)
- Window size for permutation testing

**Test Matrix**:
| Parameter | Low | Medium | High | Expected Effect |
|-----------|-----|--------|------|-----------------|
| warmupCount | 100 | 400 | 1000 | More warmup = stable baseline |
| detectionInterval | 1 | 25 | 100 | Lower = more sensitive |
| NCD threshold | 0.5 | 0.7 | 0.9 | Higher = fewer alerts |

## 🛠️ IMPLEMENTATION APPROACH

### Phase 1: Create Test Harness Script
```bash
# scripts/test-anomaly-detection.sh
# Usage: ./scripts/test-anomaly-detection.sh [dataset] [compression]
```

Features:
- Accept dataset path and compression algorithm as args
- Support `--openzl`, `--zstd`, `--lz4` flags
- Output: JSON report with metrics (precision, recall, F1, latency)
- Compare against ground truth if available

### Phase 2: Add Just Recipes
```just
# Add to Justfile:
test-anomaly-detection dataset:
    ./scripts/test-anomaly-detection.sh {{dataset}}

test-openzl-comparison:
    # Run same dataset with zstd vs openzl, compare results

test-streaming-simulation:
    # Fake-stream mixed-transactions.jsonl with anomaly detection
```

### Phase 3: CI Integration
```yaml
# .github/workflows/anomaly-validation.yml
# Run on main/PR with test datasets
# Optional OpenZL job for compression comparison
```

## 📊 SUCCESS METRICS

| Test | Metric | Target |
|------|--------|--------|
| Normal data (500 tx) | False positives | < 5 (1%) |
| Anomalous data (100 tx) | True positives | > 80 (80%) |
| Mixed data (1000 tx) | F1 Score | > 0.85 |
| Terra/Luna crash | Detection within | First 100 events of depeg |
| Streaming latency | P95 | < 100ms per event |
| OpenZL vs zstd | Compression ratio | > 10% better |

## 📍 START HERE

1. **Read test data README**: `test-data/README.md`
2. **Understand detection flow**: `cmd/demo/main.go` (lines 76-338)
3. **Check streaming setup**: `docs/launch/SOAK_TEST.md`
4. **Run existing demo**: `just build-demo && ./driftlock-demo test-data/financial-demo.json`
5. **View OpenZL docs**: `.archive/reports/OPENZL_ANALYSIS.md`

## 🔧 USEFUL COMMANDS

```bash
# Check what's available
just --list | grep -E "(test|openzl)"

# Build with OpenZL (Linux/Docker)
just build-openzl-lib
just build-core-openzl
just test-openzl

# Run existing test suite
just test

# Run demo verification
just verify
```

## 💡 CREATIVE APPROACHES TO CONSIDER

1. **Chaos Injection**: Insert random garbage bytes into stream mid-test
2. **Gradual Drift**: Slowly shift distribution over time (should detect)
3. **Replay Attack**: Send historical crisis data (Terra/Luna) as if live
4. **Multi-Asset Correlation**: Detect when multiple assets anomaly together
5. **Seasonality Test**: Ensure normal daily patterns aren't flagged

## 🚨 CONSTRAINTS

- OpenZL tests must run in Docker/Linux (macOS has dylib issues)
- Don't break existing demo: `just verify` must pass
- Keep tests fast: < 60 seconds for unit tests
- Document all sensitivity parameters clearly

---

**Acknowledgment**: Please confirm you understand the context, then outline your test strategy before implementing. Focus on creating reusable, well-documented tests that future developers can extend.
# CBAD Benchmark Handoff Document

## Purpose

This document provides a fresh-eyes analysis of all available datasets for CBAD benchmarking and outlines the plan for:
1. Completing F1/precision/recall analysis on remaining datasets
2. Auto-configuring detection profiles for users

---

## Current State: What's Been Benchmarked

### Completed Benchmarks (with F1 scores)

| Dataset | Type | F1 | Precision | Recall | Status | Notes |
|---------|------|-----|-----------|--------|--------|-------|
| Synthetic Transactions | Transactions | **90.0%** | 90.0% | 90.0% | EXCELLENT | Best performer |
| Financial Fraud | Transactions | **74.0%** | 62.3% | 91.0% | GOOD | High recall |
| Jailbreak Prompts | Text/Prompts | 66.2% | 51.7% | 92.0% | OK | Good filter, needs review |
| AI Safety Malignant | Text/Prompts | 33.6% | 55.8% | 24.0% | POOR | Low recall |
| NAB Machine Temp | Time Series | 18.2% | 10.0% | 99.9% | POOR | Too many FPs |
| NAB AWS CloudWatch | Time Series | 15.9% | 8.6% | 98.2% | POOR | Too many FPs |
| Terra Luna Crash | Crypto | 0.0% | 0.0% | 0.0% | FAIL | Date parsing bug |
| Hallucination | Text/QA | 0.0% | 0.0% | 0.0% | IMPOSSIBLE | Semantic, not structural |

### Not Yet Benchmarked (Ready Data Available)

| Dataset | Path | Labels Available | Size | Priority |
|---------|------|------------------|------|----------|
| **Network Intrusion** | `driftlock-archives/test-data/network/driftlock_ready.json` | 243 normal, 1757 attack | 2000 events | HIGH |
| **Supply Chain Risk** | `driftlock-archives/test-data/supply_chain/driftlock_ready.json` | 1502 High, 307 Moderate, 191 Low | 2000 events | HIGH |
| **NASA Turbofan** | `driftlock-archives/test-data/nasa_turbofan/driftlock_ready.json` | Has RUL labels | Unknown | MEDIUM |
| **Airline Delays** | `driftlock-archives/test-data/airline/driftlock_ready.json` | Unknown | Unknown | LOW |
| **Terra Luna (Fixed)** | `driftlock-archives/test-data/terra_luna/terra-luna.csv` | Date-based (May 9-12 = crash) | ~1000 points | HIGH |
| **PINT Prompt Injection** | `benchmark-datasets/pint/` | YAML format | ~1000 samples | MEDIUM |
| **UNSW Network (Parquet)** | `driftlock-archives/test-data/network/UNSW_NB15_*.parquet` | Attack type columns | Large | LOW |

---

## Data Format Analysis

### Ready-to-Use JSON Format (`driftlock_ready.json`)

These files have a consistent structure:
```json
{
  "timestamp": "...",
  "transaction_id": "...",
  "amount_usd": 123.45,
  "processing_ms": 100,
  "origin_country": "...",
  "api_endpoint": "...",
  "status": "normal|attack|High Risk|etc"
}
```

The `status` field contains the ground truth label.

### Label Mappings

| Dataset | Normal Labels | Anomaly Labels |
|---------|--------------|----------------|
| Network | `normal` | `attack` |
| Supply Chain | `Low Risk`, `Moderate Risk` | `High Risk` |
| Fraud | `is_fraud=0` | `is_fraud=1` |
| Terra Luna | pre-May-9-2022 | May 9-12, 2022 |
| NASA Turbofan | High RUL (>30) | Low RUL (<30) |

---

## Tasks for Fresh Context

### Task 1: Benchmark Remaining Datasets

Add these benchmarks to `cbad-core/examples/comprehensive_benchmark.rs`:

#### 1.1 Network Intrusion (HIGH PRIORITY)

```rust
fn benchmark_network_labeled() -> BenchmarkResult {
    // Load from: /Volumes/VIXinSSD/driftlock-archives/test-data/network/driftlock_ready.json
    // Labels: status = "normal" or "attack"
    // Normal count: 243
    // Attack count: 1757
    //
    // Strategy: Train on normal (first 150), test on remaining normal + attacks
    // Suggested config: baseline=100, window=30, ncd_thresh=0.22
}
```

#### 1.2 Supply Chain Risk (HIGH PRIORITY)

```rust
fn benchmark_supply_chain() -> BenchmarkResult {
    // Load from: /Volumes/VIXinSSD/driftlock-archives/test-data/supply_chain/driftlock_ready.json
    // Labels: status = "Low Risk", "Moderate Risk", "High Risk"
    //
    // Binary classification:
    //   Normal = "Low Risk" + "Moderate Risk" (498 total)
    //   Anomaly = "High Risk" (1502 total)
    //
    // Strategy: Train on Low/Moderate risk, test on High Risk
    // NOTE: Imbalanced - more anomalies than normal!
}
```

#### 1.3 Terra Luna Fixed (HIGH PRIORITY)

```rust
fn benchmark_terra_luna_fixed() -> BenchmarkResult {
    // Load from: /Volumes/VIXinSSD/driftlock-archives/test-data/terra_luna/terra-luna.csv
    //
    // The crash happened May 9-12, 2022
    // Parse dates properly: "2022-05-09" contains crash dates
    //
    // Current bug: Date matching not working
    // Fix: Use proper date parsing, not string contains
}
```

#### 1.4 NASA Turbofan (MEDIUM PRIORITY)

```rust
fn benchmark_nasa_turbofan() -> BenchmarkResult {
    // Need to check driftlock_ready.json structure
    // RUL (Remaining Useful Life) is the label
    // Low RUL = approaching failure = anomaly
}
```

### Task 2: Fix Time Series Precision

Current time series benchmarks have ~10% precision (90% false positives). Test these config changes:

1. **Increase NCD threshold**: Try 0.30, 0.35, 0.40 (currently 0.18)
2. **Enable statistical significance**: Currently disabled for time series
3. **Larger baseline**: Try baseline=200, window=50
4. **Test different representations**:
   - Current: `ts=2022-01-01 value=123.45`
   - Try: Just the value `123.45`
   - Try: Delta encoding `delta=+5.2`

### Task 3: Design Auto-Configuration

The goal: Users send data, CBAD automatically picks the right profile.

#### Detection Profiles

```rust
pub enum DetectionProfile {
    Financial,      // Transactions, fraud - high precision
    LlmSafety,      // Jailbreak, prompts - high recall, human review
    Logs,           // Structured logs - normalized tokenization
    TimeSeries,     // Numeric metrics - needs special handling
    Network,        // Intrusion detection
    Custom(AnomalyConfig),
}

impl DetectionProfile {
    pub fn config(&self) -> AnomalyConfig {
        match self {
            Self::Financial => AnomalyConfig {
                window_config: WindowConfig {
                    baseline_size: 300,
                    window_size: 50,
                    ..Default::default()
                },
                ncd_threshold: 0.20,
                require_statistical_significance: true,
                ..Default::default()
            },
            Self::LlmSafety => AnomalyConfig {
                window_config: WindowConfig {
                    baseline_size: 200,
                    window_size: 30,
                    ..Default::default()
                },
                ncd_threshold: 0.25,
                require_statistical_significance: true,
                ..Default::default()
            },
            // ... etc
        }
    }
}
```

#### Auto-Detection Heuristics

Detect profile from data shape:

```rust
pub fn detect_profile(events: &[Event]) -> DetectionProfile {
    // Check for financial indicators
    if events.iter().any(|e| e.has_field("amount") || e.has_field("merchant")) {
        return DetectionProfile::Financial;
    }

    // Check for prompt/text content
    if events.iter().any(|e| e.has_field("prompt") || e.has_field("message")) {
        return DetectionProfile::LlmSafety;
    }

    // Check for time series (mostly numeric)
    if events.iter().all(|e| e.is_numeric()) {
        return DetectionProfile::TimeSeries;
    }

    // Check for network indicators
    if events.iter().any(|e| e.has_field("ip") || e.has_field("port") || e.has_field("protocol")) {
        return DetectionProfile::Network;
    }

    // Default to logs
    DetectionProfile::Logs
}
```

---

## API Changes Needed

### Current API

```
POST /v1/detect
{
  "stream_id": "my-stream",
  "events": [...]
}
```

### Proposed API

```
POST /v1/detect
{
  "stream_id": "my-stream",
  "profile": "financial",  // optional - auto-detect if not provided
  "events": [...]
}
```

Or with explicit auto-detect:

```
POST /v1/detect
{
  "stream_id": "my-stream",
  "profile": "auto",  // explicitly request auto-detection
  "events": [...]
}
```

### Profile Endpoint

```
GET /v1/profiles
Returns available profiles and their configs

GET /v1/streams/{id}/profile
Returns current profile for stream

PATCH /v1/streams/{id}/profile
{
  "profile": "financial"
}
```

---

## Files to Modify

| File | Change |
|------|--------|
| `cbad-core/examples/comprehensive_benchmark.rs` | Add missing benchmark functions |
| `cbad-core/src/lib.rs` | Add DetectionProfile enum |
| `crates/driftlock-api/src/handlers/detect.rs` | Add profile parameter |
| `docs/BENCHMARK_PLAN.md` | Update with new results |

---

## Key Insights from Previous Benchmarks

1. **Transactions are CBAD's sweet spot**: 74-90% F1 with minimal tuning
2. **High recall is easy, precision is hard**: Most configs catch anomalies but have false positives
3. **Time series needs different approach**: Numeric data compresses differently than text
4. **Hallucination is fundamentally impossible**: CBAD detects structural anomalies, not semantic errors
5. **Imbalanced data matters**: Supply chain has more anomalies than normal - affects training

---

## Command to Run Benchmarks

```bash
cd /Volumes/VIXinSSD/driftlock
cargo run --example comprehensive_benchmark --release
```

Results saved to: `benchmark-datasets/results/comprehensive_results.json`

---

## Success Criteria

| Metric | Target |
|--------|--------|
| Financial/Transactions | F1 > 70% |
| Network Intrusion | F1 > 60% |
| Supply Chain | F1 > 50% |
| LLM Safety | Recall > 80% (filter use case) |
| Time Series | F1 > 30% (improved from 16%) |

---

## Questions to Resolve

1. Should "Moderate Risk" in supply chain be normal or anomaly?
2. What RUL threshold defines "approaching failure" in NASA data?
3. Should we expose profile auto-detection to users or make it transparent?
4. How to handle profile switching mid-stream?

---

## Additional Domains to Benchmark

The following domains represent high-value markets that need benchmark validation:

### 1. Cryptocurrency / Blockchain Fraud

**Why it matters:** $24.2B in illicit crypto activity in 2023 (Chainalysis).

**Available Datasets:**

| Dataset | Source | Size | Labels |
|---------|--------|------|--------|
| Ethereum Fraud Detection | [Kaggle](https://www.kaggle.com/datasets) | 9,841 addresses | fraud/not fraud |
| Bitcoin Illicit Transactions | [Elliptic Dataset](https://www.kaggle.com/datasets/ellipticco/elliptic-data-set) | 200K+ transactions | licit/illicit |
| Blockchain Phishing Detection | [arxiv:2401.03530](https://arxiv.org/html/2401.03530v1) | 78,600 records | Various attack types |
| Raw Bitcoin Transactions | Kaggle | 30M transactions | 11 attributes |

**CBAD Applicability:** HIGH - Transaction patterns compress differently than normal activity.

**Suggested Test:**
```
1. Download Elliptic Bitcoin dataset
2. Train on "licit" transactions
3. Test detection of "illicit" transactions
4. Target: F1 > 60%
```

---

### 2. IoT Sensor / Industrial Control

**Why it matters:** Critical infrastructure, manufacturing, smart cities.

**Available Datasets:**

| Dataset | Source | Description | Labels |
|---------|--------|-------------|--------|
| **CIC IoT-DIAD 2024** | [University of New Brunswick](https://www.unb.ca/cic/datasets/iot-diad-2024.html) | 105 devices, 33 attack types | DDoS, DoS, Recon, Spoofing, etc. |
| **DAD (Data Center)** | [GitHub](https://github.com/dad-repository/dad) | 101,583 packets, MQTT protocol | duplication, interception, modification |
| **IoT-23** | Stratosphere Laboratory | Malicious/benign IoT traffic | Botnet, C&C, DDoS |
| **TON_IoT** | UNSW | Network + telemetry data | Normal/attack |

**CBAD Applicability:** HIGH - Protocol anomalies compress differently.

**Suggested Test:**
```
1. Download CIC IoT-DIAD 2024 CSV features
2. Train on normal traffic patterns
3. Test detection of attack categories
4. Target: F1 > 50% (network is harder)
```

---

### 3. Manufacturing / Predictive Maintenance

**Why it matters:** $630B predictive maintenance market by 2030.

**Available Datasets:**

| Dataset | Source | Description | Labels |
|---------|--------|-------------|--------|
| **NASA C-MAPSS Turbofan** | [NASA Prognostics](https://ti.arc.nasa.gov/tech/dash/groups/pcoe/prognostic-data-repository/) | Turbofan engine degradation | RUL (Remaining Useful Life) |
| **MetroPT** | [Nature Scientific Data](https://www.nature.com/articles/s41597-022-01877-3) | Metro train sensors 2022 | Anomaly windows from maintenance logs |
| **SCANIA Trucks** | [Nature PMC](https://pmc.ncbi.nlm.nih.gov/articles/PMC11933314/) | Truck component failures | Repair records |
| **Azure PM Dataset** | [Kaggle/Microsoft](https://www.kaggle.com/datasets/shivamb/machine-predictive-maintenance-classification) | 10,000 machine records | Tool Wear Failure, etc. |
| **Hydraulic System** | UCI / GitHub | Pressure, flow, temp sensors | Valve condition |

**CBAD Applicability:** MEDIUM - Numeric sensors may need specialized representation.

**Suggested Test:**
```
1. Use NASA C-MAPSS (already have in repo as nasa_turbofan)
2. Convert RUL < 30 cycles to "failure approaching" label
3. Train on healthy operation (high RUL)
4. Test detection of degradation patterns
5. Target: F1 > 40%
```

---

### 4. Pharmaceutical / Clinical Trials

**Why it matters:** Drug safety, fraud detection in trials, adverse events.

**Available Datasets:**

| Dataset | Source | Description | Labels |
|---------|--------|-------------|--------|
| **BMAD (Medical Imaging)** | [CVPR 2024](https://arxiv.org/abs/2306.11876) | 6 reorganized datasets | Normal/anomaly across 5 modalities |
| **MedIAnomaly** | [ScienceDirect 2025](https://www.sciencedirect.com/science/article/abs/pii/S1361841525000489) | 7 medical datasets | Anomaly detection benchmark |
| **MIMIC-III/IV** | PhysioNet | ICU patient records | Various clinical outcomes |
| **Adverse Events** | FDA FAERS | Drug adverse event reports | Event type classifications |

**CBAD Applicability:** MEDIUM - Structured clinical data could work, but domain-specific.

**Challenge:** Most clinical trial data is proprietary. May need to use synthetic or public medical datasets.

**Suggested Test:**
```
1. Use MIMIC-III (requires credentialed access)
2. Or use BMAD medical imaging benchmark
3. Train on normal patterns
4. Test detection of anomalous cases
5. Target: F1 > 35% (medical is challenging)
```

---

### 5. DNA / Genomics

**Why it matters:** Personalized medicine, genetic testing, variant detection.

**Available Datasets:**

| Dataset | Source | Description | Labels |
|---------|--------|-------------|--------|
| **Genomic Benchmarks** | [GitHub/PyPI](https://github.com/ML-Bioinfo-CEITEC/genomic_benchmarks) | 9 classification datasets | Enhancers, promoters, regulatory elements |
| **GUANinE** | [Nature PMC](https://pmc.ncbi.nlm.nih.gov/articles/PMC10614795/) | Genomic AI benchmark | Functional element annotations |
| **VariBench** | Academic | Variant prediction | Pathogenic/benign variants |
| **DNase I Hypersensitive** | Various | DHS identification | Positive/negative sequences |

**CBAD Applicability:** UNKNOWN - DNA sequences are fundamentally different from logs/transactions.

**Hypothesis:** CBAD compression may detect unusual sequence patterns (mutations, insertions, deletions) if the baseline is "normal" genome segments.

**Suggested Test:**
```
1. Download genomic-benchmarks Python package
2. Use human_nontata_promoters dataset (36K sequences)
3. Train on "negative" class (non-promoter sequences)
4. Test detection of "positive" class (promoter sequences)
5. Target: Unknown - experimental
```

---

## Dataset Acquisition Priority

### Immediate (Have Data, Need F1 Analysis)

| Priority | Dataset | Path | Status |
|----------|---------|------|--------|
| 1 | Network Intrusion | `driftlock-archives/test-data/network/driftlock_ready.json` | Ready |
| 2 | Supply Chain | `driftlock-archives/test-data/supply_chain/driftlock_ready.json` | Ready |
| 3 | Terra Luna (Fixed) | `driftlock-archives/test-data/terra_luna/` | Needs date parsing fix |
| 4 | NASA Turbofan | `driftlock-archives/test-data/nasa_turbofan/driftlock_ready.json` | Ready |

### Near-Term (Download & Benchmark)

| Priority | Dataset | Source | Est. Size |
|----------|---------|--------|-----------|
| 1 | Elliptic Bitcoin | Kaggle | 200K txns |
| 2 | CIC IoT-DIAD 2024 | UNB | CSV features |
| 3 | MetroPT | Nature | Sensor data |
| 4 | Azure PM | Kaggle | 10K records |

### Experimental (Novel Domains)

| Priority | Dataset | Source | Notes |
|----------|---------|--------|-------|
| 1 | Genomic Benchmarks | PyPI | DNA sequences - unknown applicability |
| 2 | BMAD Medical | CVPR | Medical imaging - may not fit CBAD |
| 3 | MIMIC-III | PhysioNet | Requires credentials |

---

## Implementation Checklist for Fresh Context

### Phase 1: Complete Existing Benchmarks
- [ ] Implement `benchmark_network_labeled()` in comprehensive_benchmark.rs
- [ ] Implement `benchmark_supply_chain()` in comprehensive_benchmark.rs
- [ ] Fix `benchmark_terra_luna_fixed()` date parsing
- [ ] Implement `benchmark_nasa_turbofan()` with RUL labels
- [ ] Run full benchmark suite and update results

### Phase 2: Download New Datasets
- [ ] Download Elliptic Bitcoin dataset from Kaggle
- [ ] Download CIC IoT-DIAD 2024 from UNB
- [ ] Download MetroPT from Nature
- [ ] Convert to driftlock_ready.json format

### Phase 3: Expand Benchmarks
- [ ] Add `benchmark_elliptic_bitcoin()`
- [ ] Add `benchmark_iot_cic()`
- [ ] Add `benchmark_metroPT()`
- [ ] Try `benchmark_genomic()` (experimental)

### Phase 4: Tune Time Series
- [ ] Test higher NCD thresholds (0.30, 0.35, 0.40)
- [ ] Enable statistical significance for time series
- [ ] Try delta encoding for numeric values
- [ ] Document optimal configs per data type

### Phase 5: Implement Auto-Config API
- [ ] Add DetectionProfile enum to cbad-core
- [ ] Implement profile auto-detection heuristics
- [ ] Add `profile` parameter to `/v1/detect` endpoint
- [ ] Create `/v1/profiles` endpoint
- [ ] Update documentation

---

## Sources

- [CIC IoT-DIAD 2024](https://www.unb.ca/cic/datasets/iot-diad-2024.html)
- [Elliptic Bitcoin Dataset](https://www.kaggle.com/datasets/ellipticco/elliptic-data-set)
- [MetroPT Dataset](https://www.nature.com/articles/s41597-022-01877-3)
- [SCANIA Component X](https://pmc.ncbi.nlm.nih.gov/articles/PMC11933314/)
- [BMAD Medical Anomaly](https://arxiv.org/abs/2306.11876)
- [Genomic Benchmarks](https://github.com/ML-Bioinfo-CEITEC/genomic_benchmarks)
- [NASA C-MAPSS](https://ti.arc.nasa.gov/tech/dash/groups/pcoe/prognostic-data-repository/)
- [Azure Predictive Maintenance](https://www.kaggle.com/datasets/shivamb/machine-predictive-maintenance-classification)

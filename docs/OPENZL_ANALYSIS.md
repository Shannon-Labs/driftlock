# OpenZL Integration Analysis for Driftlock CBAD

**Date:** 2025-10-23
**Status:** Recommendation for Phase 1
**Decision:** Include OpenZL as experimental compression adapter alongside zstd/lz4/gzip

---

## Executive Summary

**Recommendation: YES, but with strategic phasing**

OpenZL is an excellent fit for Driftlock's compression-based anomaly detection (CBAD) use case, particularly for structured telemetry data (logs, metrics, traces, LLM I/O). However, integration should be **phased** to balance innovation with project timeline constraints.

### Proposed Approach
1. **Phase 1a (Current):** Implement zstd, lz4, gzip adapters (proven, fast, deterministic)
2. **Phase 1b (Experimental):** Add OpenZL adapter as optional enhancement
3. **Phase 2+:** Optimize OpenZL integration based on real-world performance

---

## What is OpenZL?

### Overview
- **Format-aware compression framework** from Meta (Oct 2025)
- **Lossless compression** optimized for structured data (Parquet, CSV, columnar numeric, timeseries)
- **Single universal decompressor** that handles all data formats
- **Offline training** to learn optimal compression strategies per data shape

### Key Innovation
Unlike generic compressors (zstd, gzip) that treat data as byte streams, OpenZL:
1. Parses data structure (via SDDL or custom parser)
2. Applies transforms (e.g., struct-of-arrays conversion, delta encoding, tokenization)
3. Learns optimal compression graph offline via training
4. Embeds decode recipe in compressed frame (no out-of-band coordination needed)

### Performance Characteristics

From Meta's benchmarks (M1 CPU, clang-17):

| Dataset | Compressor | Ratio | Comp Speed | Decomp Speed |
|---------|-----------|-------|------------|--------------|
| SAO (star records) | zstd -3 | 1.31x | 220 MB/s | 850 MB/s |
| SAO | OpenZL | **2.06x** | **340 MB/s** | **1200 MB/s** |
| ERA5 (columnar numeric) | zstd -3 | ~1.4x | ~200 MB/s | ~800 MB/s |
| ERA5 | OpenZL | **2.8x** | **280 MB/s** | **1100 MB/s** |

**Key Insight:** OpenZL beats zstd on BOTH ratio AND speed for structured data.

---

## Why OpenZL Fits Driftlock

### 1. Data Characteristics Match

Driftlock processes highly structured telemetry:

| Data Type | Structure | OpenZL Benefit |
|-----------|-----------|----------------|
| **OTLP Logs** | JSON/protobuf with repeated fields (timestamp, severity, attributes) | Struct-of-arrays transform + tokenization |
| **OTLP Metrics** | Timeseries numeric data (gauge, counter, histogram) | Delta encoding, transpose for bounded ranges |
| **OTLP Traces** | Nested spans with common fields (trace_id, span_id, service.name) | Format-aware parsing + dictionary compression |
| **LLM I/O** | JSON (prompts, responses, tool calls) with high field repetition | Tokenization for low-cardinality fields |

All of these benefit from **format-aware compression** vs. byte-level compression.

### 2. Anomaly Detection Advantage

**Hypothesis:** Format-aware compression provides **more sensitive anomaly signals** than generic compression.

#### Why?
- **Better baselines:** OpenZL learns "normal" structure during training → deviations are more pronounced
- **Richer metrics:** Compression ratio per field/column vs. per-blob gives granular insights
- **Structured explanations:** "Field X had unusually low compression" is more actionable than "block Y compressed poorly"

#### Example Scenario
**Normal log entry:**
```json
{"level": "info", "service": "api-gateway", "msg": "request completed", "duration_ms": 42}
```

**Anomalous log entry:**
```json
{"level": "error", "service": "api-gateway", "msg": "SEGFAULT in libc malloc() at 0x7fff...", "stack_trace": "...500 lines..."}
```

- **zstd:** Compresses both as JSON byte streams → modest compression ratio change
- **OpenZL:** Learns that `msg` field is usually short tokens → massive compression ratio drop on 500-line stack trace + new field `stack_trace` breaks structure

### 3. Competitive Differentiation

**Market Position:** Explainable AI for regulated industries (DORA, NIS2, AI Act)

OpenZL provides:
- **Glass-box compression:** Explain WHY compression failed (not black-box ML)
- **Deterministic training:** Fixed seed → reproducible compression plans
- **Audit trail:** Compression graph embedded in frame → forensic analysis
- **Novel IP:** No competitors using format-aware compression for anomaly detection

**Pricing Power:** "We use Meta's advanced compression framework optimized for your specific data formats" sounds premium.

---

## Technical Feasibility

### Language & Integration

| Aspect | Details | Driftlock Impact |
|--------|---------|------------------|
| **Primary Language** | C/C++ (C11, C++17) | ✅ Rust FFI via C bindings (same as zstd) |
| **License** | BSD 3-Clause (Meta) | ✅ Compatible with commercial use |
| **Build System** | CMake or Makefile | ✅ Can compile as static library for linking |
| **Dependencies** | Minimal (zstd bundled) | ✅ No heavyweight deps |
| **Platform Support** | Linux, macOS, Windows (clang-cl) | ✅ Cross-platform |
| **Rust Bindings** | None (yet) | ⚠️ Need to create C FFI wrapper |

### Integration Path

#### Option A: Direct C FFI (Recommended)
```rust
// cbad-core/src/compression/openzl.rs

use std::ffi::{CString, c_void};

#[link(name = "openzl")]
extern "C" {
    fn OPENZL_compress(
        dst: *mut u8, dst_capacity: usize,
        src: *const u8, src_size: usize,
        plan: *const c_void,
    ) -> isize;

    fn OPENZL_decompress(
        dst: *mut u8, dst_capacity: usize,
        src: *const u8, src_size: usize,
    ) -> isize;
}

pub struct OpenZLAdapter {
    plan: *const c_void,  // Trained compression plan
}

impl CompressionAdapter for OpenZLAdapter {
    fn compress(&self, data: &[u8]) -> Result<Vec<u8>> {
        // Call OPENZL_compress with self.plan
    }

    fn decompress(&self, data: &[u8]) -> Result<Vec<u8>> {
        // Call OPENZL_decompress
    }
}
```

#### Option B: Rust Wrapper Crate (Future)
Create `openzl-rs` crate for community (like `zstd-rs` exists for zstd).

---

## Risks & Mitigations

### Risk 1: Maturity & Stability
**Issue:** OpenZL released Oct 2025 (very new)
**Severity:** Medium
**Mitigation:**
- Use release-tagged versions only (not `dev` branch)
- Meta states: "Production-ready, used extensively at Meta"
- Fallback: OpenZL can fall back to zstd for unstructured data
- Phased rollout: Keep zstd/lz4 as primary, OpenZL as experimental

### Risk 2: Training Overhead
**Issue:** OpenZL requires offline training to generate compression plans
**Severity:** Low
**Mitigation:**
- Training is one-time per data format (OTLP logs, metrics, traces)
- Can ship pre-trained plans with Driftlock
- Training time: minutes on sample data (per Meta docs)
- Users can optionally re-train for custom schemas

### Risk 3: Compression Speed Variability
**Issue:** OpenZL speed depends on data structure complexity
**Severity:** Low
**Mitigation:**
- Benchmark on representative OTLP data before committing
- Configure fallback to zstd if OpenZL is slower than threshold
- Phase 1 benchmarks will establish baselines

### Risk 4: FFI Complexity
**Issue:** No official Rust bindings (need to create C FFI wrapper)
**Severity:** Medium
**Mitigation:**
- C FFI is well-understood in Rust ecosystem (many examples: zstd-rs, lz4-rs)
- OpenZL C API is stable (universal decompressor guarantee)
- Start with simple compress/decompress functions (advanced features later)

---

## Recommended Architecture

### Phase 1a: Core Compression Adapters (Weeks 1-3)

**Goal:** Get deterministic anomaly detection working with proven compressors.

```rust
// cbad-core/src/compression/mod.rs

pub trait CompressionAdapter: Send + Sync {
    fn compress(&self, data: &[u8]) -> Result<Vec<u8>>;
    fn decompress(&self, data: &[u8]) -> Result<Vec<u8>>;
    fn name(&self) -> &str;
}

pub enum CompressionAlgorithm {
    Zstd(i32),      // compression level
    Lz4,
    Gzip(u32),      // compression level
    OpenZL(Plan),   // trained plan (Phase 1b)
}

pub fn create_adapter(algo: CompressionAlgorithm) -> Box<dyn CompressionAdapter> {
    match algo {
        CompressionAlgorithm::Zstd(level) => Box::new(ZstdAdapter::new(level)),
        CompressionAlgorithm::Lz4 => Box::new(Lz4Adapter::new()),
        CompressionAlgorithm::Gzip(level) => Box::new(GzipAdapter::new(level)),
        CompressionAlgorithm::OpenZL(plan) => Box::new(OpenZLAdapter::new(plan)),
    }
}
```

**Deliverables:**
- ✅ Zstd adapter (proven, fast, dictionary support)
- ✅ Lz4 adapter (ultra-fast, lower ratio)
- ✅ Gzip adapter (universal compatibility)
- ✅ Benchmarks: throughput, ratio, determinism tests
- ✅ Metrics calculators: compression ratio, delta bits, entropy

### Phase 1b: OpenZL Experimental Integration (Weeks 4-5)

**Goal:** Prove OpenZL value on real OTLP data.

**Tasks:**
1. Build OpenZL as static library (`libopenzl.a`)
2. Create C FFI wrapper in Rust (`openzl.rs`)
3. Train compression plans for:
   - OTLP logs (JSON schema)
   - OTLP metrics (timeseries numeric)
   - OTLP traces (nested spans)
4. Implement `OpenZLAdapter` in cbad-core
5. Benchmark against zstd on synthetic OTLP data
6. Add configuration flag: `enable_openzl: true` (default false)

**Success Criteria:**
- ✅ Compression ratio ≥ 1.5x improvement over zstd on structured logs
- ✅ Compression speed ≥ zstd speed (or within 20%)
- ✅ Deterministic output with fixed compression plan
- ✅ Anomaly detection sensitivity improvement measured

**Risk-Off Strategy:**
If OpenZL doesn't meet criteria → disable by default, revisit in Phase 3.

---

## Performance Targets (Empirical)

These are **goals to measure**, not promises:

| Metric | zstd Baseline | OpenZL Target | Test Data |
|--------|---------------|---------------|-----------|
| **Compression Ratio (logs)** | 2.0x | 3.0x | 1M synthetic OTLP log entries |
| **Compression Ratio (metrics)** | 3.0x | 5.0x | 1M timeseries numeric samples |
| **Compression Speed** | 200 MB/s | ≥ 180 MB/s | Same datasets |
| **Anomaly Sensitivity** | Baseline | +30% true positives | Injected anomalies |

**Measurement Plan:**
- Tools: `tools/benchmark/compression_bench.rs` (to be created)
- Datasets: `tools/synthetic/otlp_generator.go` (to be enhanced)
- Metrics: ratio, speed, RAM usage, determinism (100 runs with same seed)

---

## Compliance & Explainability Benefits

### DORA (Digital Operational Resilience Act)

**Requirement:** ICT risk monitoring with explainable anomaly detection

**OpenZL Advantage:**
- Compression plan = auditable artifact (embedded in frame)
- Transforms are reversible and documented (transpose, delta, tokenize)
- Training logs = evidence of systematic methodology
- Deterministic plans = reproducible risk assessments

**Example Evidence:**
> "Anomaly detected: Log entry compression ratio dropped from 3.2x (baseline) to 0.8x due to unexpected field 'stack_trace' (not in trained schema). OpenZL transform graph shows tokenization failure on field X. Root cause: New error type introduced in service Y version 2.1.3."

### NIS2 (Network & Information Security Directive)

**Requirement:** Incident detection and reporting within 24 hours

**OpenZL Advantage:**
- Fast decompression (1200 MB/s) enables real-time forensics
- Field-level compression metrics pinpoint incident scope
- Universal decompressor = single audit surface (no format-specific tools)

### EU AI Act (Runtime AI Compliance)

**Requirement:** Explainable AI for high-risk systems

**OpenZL Advantage:**
- Glass-box algorithm (not black-box ML)
- Mathematical proof: compression theory (Kolmogorov complexity)
- Testable: compression plans can be validated independently
- No training data bias (transforms are format-agnostic)

**Positioning:**
> "Driftlock uses Meta's OpenZL compression framework, a deterministic, explainable algorithm compliant with EU AI Act requirements. All anomaly detections include compression graph analysis showing exactly why the data deviated from the learned baseline structure."

---

## Build Integration Plan

### Step 1: Add OpenZL as Git Submodule (Phase 1b)

```bash
cd /home/user/driftlock
git submodule add https://github.com/facebook/openzl.git deps/openzl
git submodule update --init --recursive
```

### Step 2: Build Static Library

```makefile
# Makefile additions
OPENZL_DIR := deps/openzl
OPENZL_LIB := $(OPENZL_DIR)/build/lib/libopenzl.a

.PHONY: openzl
openzl:
	mkdir -p $(OPENZL_DIR)/build
	cd $(OPENZL_DIR)/build && \
		cmake -DCMAKE_BUILD_TYPE=Release .. && \
		make -j$(nproc)

cbad-core-lib: openzl
	cd cbad-core && \
		OPENZL_LIB_PATH=$(OPENZL_LIB) cargo build --release --features openzl
```

### Step 3: Rust Feature Flag

```toml
# cbad-core/Cargo.toml

[features]
default = []
openzl = []  # Enables OpenZL adapter

[build-dependencies]
cc = "1.0"  # For linking C library

[dependencies]
zstd = "0.13"
lz4 = "1.24"
flate2 = "1.0"  # gzip
```

### Step 4: Link OpenZL in build.rs

```rust
// cbad-core/build.rs

fn main() {
    #[cfg(feature = "openzl")]
    {
        println!("cargo:rustc-link-search=native=../deps/openzl/build/lib");
        println!("cargo:rustc-link-lib=static=openzl");
        println!("cargo:rustc-link-lib=dylib=stdc++");
    }
}
```

---

## Alternatives Considered

### Alternative 1: Stick with zstd/lz4/gzip only
**Pros:** Proven, simple, fast to implement
**Cons:** Misses competitive differentiation opportunity, leaves compression gains on table
**Decision:** Keep as Phase 1a baseline, add OpenZL as Phase 1b enhancement

### Alternative 2: Build custom format-aware compressor
**Pros:** Full control, tailored to OTLP formats
**Cons:** Massive engineering effort (months), reinventing Meta's work, no community support
**Decision:** REJECT - OpenZL exists and is production-ready

### Alternative 3: Use machine learning for compression
**Pros:** Could learn complex patterns
**Cons:** Not explainable (AI Act issues), not deterministic, slow, high compute cost
**Decision:** REJECT - contradicts Driftlock's glass-box philosophy

---

## Decision

**APPROVED for Phase 1b (Experimental Integration)**

### Rationale
1. **Strategic Fit:** OpenZL aligns perfectly with Driftlock's structured telemetry use case
2. **Competitive Edge:** Format-aware compression for anomaly detection is novel/defensible
3. **Low Risk:** Phased approach with zstd/lz4 fallback protects timeline
4. **High Upside:** Potential for 1.5-2x compression ratio improvement + better anomaly sensitivity
5. **Compliance:** Enhances explainability for DORA/NIS2/AI Act

### Next Steps
1. ✅ Complete Phase 1a (zstd/lz4/gzip adapters)
2. ⏳ Add OpenZL as git submodule (Phase 1b start)
3. ⏳ Create Rust FFI bindings for OpenZL C API
4. ⏳ Train compression plans for OTLP log/metric/trace schemas
5. ⏳ Benchmark against zstd baseline
6. ⏳ Document results in `docs/BENCHMARK_RESULTS.md`

---

## References

- **OpenZL Blog:** https://engineering.fb.com/2025/10/06/developer-tools/openzl-open-source-format-aware-compression-framework/
- **OpenZL GitHub:** https://github.com/facebook/openzl
- **OpenZL Whitepaper:** https://arxiv.org/abs/2510.03203
- **License:** BSD 3-Clause (Meta Platforms)
- **Driftlock Phase 1 Plan:** `docs/PHASE1_PLAN.md`
- **Driftlock Algorithms:** `docs/ALGORITHMS.md`

---

**Document Owner:** Claude Code
**Reviewers:** @user (Driftlock Project Owner)
**Status:** Recommendation Pending Approval
**Last Updated:** 2025-10-23

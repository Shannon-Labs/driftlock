# Driftlock Agents — `cbad-core/`

These instructions apply to all work under `cbad-core/` (Rust CBAD core and FFI).

---

## 1. Purpose of this crate

- `cbad-core` is the **authoritative implementation** of Driftlock's compression-based anomaly detection (CBAD):
  - Compression adapters (generic + OpenZL when enabled).
  - Metrics (NCD, compression ratios, entropy, p‑values, confidence).
  - Sliding window/baseline/hop mechanics.
  - Anomaly detection logic and streaming detector primitives.
- Everything else (Go, HTTP, UIs) depends on this crate remaining **correct, deterministic, and stable**.

---

## 2. Hard invariants

When changing anything under `cbad-core/`, do **not** break:

- **Determinism**
  - Same baseline/window bytes + same config (including permutation count and seed) must always produce the same `AnomalyMetrics` and anomaly flag.
  - All randomness must come from explicitly seeded RNGs; do not use global or time‑based randomness.

- **Safety**
  - Avoid `unsafe`. When it is absolutely required (FFI, OpenZL bindings), keep the scope minimal, add a clear comment explaining why it is safe, and prefer simple, obviously correct pointer arithmetic.
  - Never assume null‑ability, lengths, or alignment from C callers without explicit validation.

- **Public FFI contracts**
  - FFI exports (e.g., `cbad_compute_metrics`, detector FFI) must remain binary compatible with Go bindings in `collector-processor/driftlockcbad/`.
  - If you must change a struct or function signature, introduce a **new** symbol and keep the old one for backward compatibility where feasible.

- **CLI demo compatibility**
  - Do not change behaviour in ways that break `make demo` or `./scripts/verify-launch-readiness.sh` without coordinating updates in the Go demo and HTML report expectations.

---

## 3. Compression & OpenZL

- Generic compressors:
  - `Zstd`, `Lz4`, and `Gzip` adapters must remain available and fully functional in the default build (no special features enabled).
  - Generic adapters should be **deterministic**, stateless, and safe for use in parallel pipelines.

- OpenZL integration:
  - The `openzl` module and feature flag are **optional**, off by default.
  - OpenZL failures (missing symbols, bad plans, compression errors) must surface as `CompressionError` and not crash the process.
  - Never make OpenZL a hard runtime requirement; generic compressors must suffice for correctness and demos.
  - If you change OpenZL bindings, update `docs/OPENZL_ANALYSIS.md` with build and usage notes.

---

## 4. API and config design

- Keep the Rust API small and composable:
  - Prefer a few well‑designed entry points (e.g., `compute_metrics`, `AnomalyDetector`) over many overlapping helpers.
  - Accept explicit configs (`ComputeConfig`, `AnomalyConfig`, `WindowConfig`) rather than magic constants.
- Config fields impacting decisions (thresholds, permutation counts, seeds) must be:
  - Documented in rustdoc.
  - Serializable via serde where practical, so they can be persisted or surfaced to APIs.

---

## 5. Testing and performance

- Add or maintain tests for:
  - Compression round‑trips for all supported adapters.
  - Basic anomaly detection behaviour (anomaly vs non‑anomaly) on synthetic data.
  - Regression tests for previously reported bugs.
- Property/quickcheck‑style tests are encouraged for metrics and windows when feasible, but keep them deterministic.
- When optimizing performance:
  - Never sacrifice determinism.
  - Prefer algorithmic improvements or buffer reuse over unsafe micro‑optimizations.


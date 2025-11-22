# OpenZL Integration Plan: The Next-Gen Core

**Goal:** Transition Driftlock from generic compression (zstd) to **OpenZL** (Meta's format-aware framework) as the default engine.

---

## 1. Why OpenZL?

Driftlock's secret sauce is "Compression as a Sensor."
*   **zstd/lz4:** Good at finding repeated *bytes*.
*   **OpenZL:** Good at finding repeated *structures* (Format-Aware).

OpenZL can "understand" a JSON schema or a Protocol Buffer structure. It separates the "skeleton" from the "meat."
*   **Benefit:** MUCH higher sensitivity to "structural anomalies" (e.g., a JSON field changing type, or a CSV column shifting range) which zstd often misses if the byte entropy is similar.

## 2. Current Status

*   **Module Exists:** `cbad-core/src/compression/openzl.rs` implements the `Compressor` trait.
*   **Feature Flag:** `Cargo.toml` has a `[features] openzl = []` flag.
*   **Submodule:** `deps/openzl` exists in the repo root.
*   **Gap:** It is **disabled by default** and requires complex build-time linking.

## 3. Implementation Roadmap

### Phase 1: The "Hard Link" (Week 1)
*   **Task:** Modify `build.rs` to statically link `libopenzl.a` when the feature is enabled.
*   **Constraint:** Ensure the build falls back gracefully to `zstd` if the OpenZL library isn't found (critical for CI/CD).

### Phase 2: Schema Training (Week 2)
*   **Concept:** OpenZL needs a "dictionary" (or schema).
*   **Feature:** "Adaptive Training."
    *   During the Driftlock "Warmup" phase (first 400 events), feed the data to OpenZL to train a custom dictionary.
    *   Use this custom dictionary for the Detection phase.
*   **Why:** This creates a *hyper-specific* model of "Normal" for that specific stream.

### Phase 3: Default Switch (Week 3)
*   **Task:** Change `Collector-Processor` config to prefer `algo="openzl"` if available.
*   **Update:** `cmd/driftlock-cli` should detect if the server supports OpenZL and request it.

---

## 4. The "Strategic Moat"

By using OpenZL, we align with Meta's open-source stack.
*   **Marketing:** "Driftlock: The first security platform powered by Meta OpenZL."
*   **Performance:** 30% faster decompression means cheaper ingest at scale.
*   **Accuracy:** "Structural Awareness" reduces false positives from random noise (UUIDs) while catching true structural breaks (schema violations).

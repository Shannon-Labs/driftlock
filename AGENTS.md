# Driftlock Agent Guidelines

These instructions are for **AI agents and automation** working in this repository. They apply to the **entire repo** unless a more specific `AGENTS.md` is added in a subdirectory.

---

## 1. Orientation First

- Before making changes, **always read**:
  - `FINAL-STATUS.md`
  - `README.md`
  - `DEMO.md`
  - `docs/ROADMAP_TO_LAUNCH.md`
  - `docs/ARCHITECTURE.md`
  - `docs/ALGORITHMS.md` (if touching core math)
  - `docs/OPENZL_ANALYSIS.md` and `docs/ai-agents/DOCKER-BUILD-STATUS.md` (if touching Docker or compression adapters)
- Understand that this repo currently ships a **demo / engine prototype**, not the full production platform described in the docs. Do not silently expand scope.

---

## 2. Golden Invariants

When editing code, **do not violate** the following without an explicit, documented reason:

- **CLI demo must remain working and deterministic**
  - `make demo` and `./verify-yc-ready.sh` must continue to succeed.
  - The demo pipeline (`cbad-core` + `cmd/demo/main.go`) should keep the same qualitative behaviour: processing synthetic payment data into an HTML report with a small number of anomalies and clear explanations.
- **Determinism**
  - Same input data + same configuration (including seeds and thresholds) must produce the same metrics and anomaly decisions.
  - All randomness (e.g., permutation tests) must be seeded explicitly and configurable.
- **Explainability**
  - Any anomaly path must expose: NCD, compression ratios, entropy, p‑value, confidence, and a short human‑readable explanation string.
  - Do not add opaque “ML-style” black boxes without clear, auditable outputs.

---

## 3. Code‑Level Guidelines

### Rust (`cbad-core/`)

- Treat `cbad-core` as the **source of truth** for CBAD math and metrics:
  - Keep FFI surfaces small, documented, and stable (`src/ffi*.rs`).
  - Minimize `unsafe`; any unsafe block must be simple and obviously correct.
- Do **not** change `crate-type` or features in `Cargo.toml` in ways that break:
  - The Go FFI in `collector-processor/driftlockcbad/`.
  - The CLI demo.
- Compression adapters:
  - Generic compressors (zstd/lz4/gzip) must continue to work in all builds.
  - OpenZL adapter is **optional** and must fail gracefully; never make it a hard runtime dependency.

### Go (`collector-processor/`, `cmd/`)

- Keep `collector-processor/cmd/driftlock-http/main.go` as the **canonical HTTP engine**:
  - `/healthz` should reflect CBAD and compression adapter health.
  - `/v1/detect` is the primary public detection endpoint; changes to its request/response shape must be documented in `docs/API.md`.
- FFI (`collector-processor/driftlockcbad/*.go`):
  - Do not change C signatures without updating the corresponding Rust exports.
  - On error, return clear Go errors; do not panic on normal failure modes.

### Frontend (`playground/`, `landing-page/`)

- Maintain the existing tech choices (Vue 3 + TS for playground; Vue/Tailwind stack for landing page).
- Keep components small and composable; favour derived state and clearly typed props.
- Do not introduce heavy new UI frameworks or state managers without strong justification.

---

## 4. OpenZL and Compression Strategy

- OpenZL is a **format‑aware compressor** that can provide better compression and sharper anomaly signals but:
  - It is **not required** for correctness.
  - It may not be available in all environments (especially Docker).
- Rules:
  - Default builds (especially Docker images) must work with generic compressors only.
  - If you add or modify OpenZL integration:
    - Keep it behind feature flags and/or explicit build args (e.g., `USE_OPENZL=true`).
    - Add clear error paths and fallbacks to zstd when OpenZL libraries, plans, or symbols are unavailable.
    - Update `docs/OPENZL_ANALYSIS.md` with what is supported and how to build it.

---

## 5. Docker and Deployment

- Use the existing Docker files as the primary deployment path:
  - `docker-compose.yml`
  - `collector-processor/cmd/driftlock-http/Dockerfile`
- Goals:
  - `docker compose up` at repo root should:
    - Build and run `driftlock-http` successfully with generic compressors.
    - Pass its health check on `/healthz`.
  - Avoid introducing unnecessary OS‑level or toolchain complexity.
- If you add OpenZL‑enabled images:
  - Do so in **additional paths** (extra Dockerfile or guarded build args), not by breaking the default images.

---

## 6. Testing and Verification

- Before concluding work that touches core logic, FFI, Docker, or the HTTP API:
  - Run `make demo` and `./verify-yc-ready.sh` if available.
  - Run any relevant scripts under `scripts/` (e.g., `test-api.sh`, `test-docker-build.sh`, `test-services.sh`) that cover your changes.
- Only modify or add tests that are clearly related to the behaviour you are changing.
- Keep tests fast and focused; avoid adding slow end‑to‑end suites without a good reason.

---

## 7. Documentation Expectations

- When you change behaviour, configuration, or public interfaces, update:
  - `README.md` and `DEMO.md` if the end‑user flow changes.
  - `docs/API.md` for HTTP/API changes.
  - `docs/OPENZL_ANALYSIS.md` and/or `docs/ai-agents/DOCKER-BUILD-STATUS.md` for compression and Docker changes.
  - `docs/ROADMAP_TO_LAUNCH.md` only when adjusting high‑level roadmap assumptions.
- Keep documentation honest about what is implemented in this repo vs. what is future/roadmap.

---

## 8. Scope and Restraint

- This repo intentionally focuses on:
  - The CBAD core engine.
  - The CLI demo and playground.
  - A thin HTTP detection service and basic Docker story.
- The larger platform (full API server, exporters, multi‑tenant UI, Kafka, ClickHouse, etc.) lives mostly in design docs. Do **not** attempt to fully implement the entire platform here unless explicitly requested.

When in doubt, prefer **small, reversible, well‑documented changes** that move Driftlock closer to a pilot‑ready anomaly detection service while preserving the existing demo and mathematical guarantees.


# Driftlock Agents — `collector-processor/`

These instructions apply to all work under `collector-processor/` (Go FFI bindings, HTTP service, collector, and Docker builds in this tree).

---

## 1. Role of this module

- Exposes `cbad-core` to the outside world via:
  - Go FFI bindings in `driftlockcbad/`.
  - The HTTP detection service `cmd/driftlock-http` (primary engine API for now).
  - The collector/processor scaffolding for future OTel/Kafka integrations.
- This layer must be **stable, observable, and safe**, suitable for running in production containers.

---

## 2. Go code standards

- Follow `docs/CODING_STANDARDS.md`:
  - No `panic`; propagate errors with context, log them with structure.
  - Context‑aware APIs (accept `context.Context` where appropriate).
  - Keep functions small and explicit; prefer clear control flow to cleverness.
- FFI (`driftlockcbad/*.go`):
  - Treat FFI calls as fallible; they must return Go errors on invalid parameters or internal failures.
  - Do not expose raw C pointers beyond the minimal glue code.
  - Ensure every allocation coming from Rust is freed using the correct FFI free function.

---

## 3. HTTP detection service (`cmd/driftlock-http`)

- Treat `driftlock-http` as the **canonical detection API**:
  - `/healthz` must report:
    - Library health (`ValidateLibrary`).
    - Available compression algorithms (zstd/lz4/gzip; optionally openzl if present).
  - `/v1/detect` is the primary anomaly detection endpoint.
    - Do not change its request/response schema without updating `docs/API.md` and any playground callers.
    - Keep behaviour deterministic for a given payload and query parameters.
- Observability:
  - Maintain Prometheus metrics for request counts and durations.
  - Keep structured logging (include request IDs, error categories).
  - Ensure graceful shutdown and no dropped in‑flight requests when possible.

---

## 4. Docker and linking

- Dockerfiles under `collector-processor/cmd/` must:
  - Build and link against `cbad-core` reliably on a standard Debian‑based image.
  - Default to generic compressors; OpenZL must remain optional and guarded (e.g., via `USE_OPENZL` args and Rust feature flags).
- When modifying Dockerfiles:
  - Do not introduce Alpine/musl unless you fully validate the FFI linking story.
  - Keep image size reasonable and avoid unnecessary toolchain bloat in the final runtime image.
  - Preserve a working `docker compose up` path from repo root.

---

## 5. Future collector/processor work

- The OTel / Kafka processors here are **skeletons**:
  - When extending them, keep the configuration model explicit and minimal.
  - Do not expand into a full OTel distro in this repo; keep focus on anomaly detection and evidence emission.
- Any new processors or exporters must:
  - Use `driftlockcbad.Detector` rather than re‑implement the CBAD algorithm.
  - Respect determinism and privacy constraints defined in the core.


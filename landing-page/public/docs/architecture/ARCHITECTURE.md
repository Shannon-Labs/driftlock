Architecture Overview

Core Components
- cbad-core (Rust): Compression-based algorithms with FFI for Go and WASM target.
- collector-processor (Go): `driftlock_cbad` OpenTelemetry Collector processor for logs/metrics.
- llm-receivers (Go): OTel Collector receiver(s) for prompts/responses/tool-calls.
- api-server (Go): Storage and retrieval for anomalies and artifacts (SSE/WebSocket later).
- exporters (Go): Evidence bundle generation (JSON + PDF) aligned to DORA/NIS2.
- ui (Next.js): Minimal dashboard to browse streams, anomalies, and artifacts.
- deploy: Docker Compose and Kubernetes (Helm) with Grafana dashboards.

Data Flow
1. Sources emit OTLP to the Collector.
2. Receivers (including `llm-receivers`) accept streams.
3. `driftlock_cbad` processor computes CBAD metrics and flags anomalies.
4. Exporters persist evidence bundles (and forward to `api-server`).
5. UI queries `api-server` to visualize anomalies and math artifacts.

Determinism
- Use deterministic seeds in permutation testing.
- Configure windows (baseline/window/hop) and thresholds explicitly.
- Avoid non-deterministic concurrency paths in the core algorithm.


Driftlock

A compression-based anomaly detection (CBAD) platform for OpenTelemetry data, powered by Meta's OpenZL format-aware compression framework. Driftlock provides explainable anomaly detection for regulated industries through advanced compression analysis of logs, metrics, traces, and LLM I/O.

This repository emphasizes a Go-based API server, an OpenTelemetry Collector processor, and a Rust CBAD core with OpenZL integration, with OpenTelemetry hooks and a minimal UI path.

What's here (Phase 0/1 skeleton)
- api-server: Go service with health, readiness, version, and basic event ingestion endpoints, OTEL-enabled.
- cbad-core: Rust crate for Compression-Based Anomaly Detection (CBAD) primitives, featuring OpenZL format-aware compression alongside zstd/lz4/gzip adapters.
- collector-processor: OTel Collector processor skeleton `driftlock_cbad` (wires in later with cbad-core).
- llm-receivers: Collector receiver skeletons for LLM prompts/responses/tool calls.
- exporters: Placeholders for evidence bundle exporters (JSON + PDF).
- deploy: Docker Compose with an OTel Collector and API container.
- tools/synthetic: Small Go generator to POST synthetic events to the API.
- ui: Placeholder Next.js app directory (intentionally minimal for now).
- web: Minimal static HTML stub to sanity-check the API.
- deps/openzl: Meta's OpenZL format-aware compression framework (git submodule).

Quick start (API)
- Prereqs: Go 1.22+, Docker (optional), an OTEL collector or backend (optional).
- Local run with OTEL disabled:
  - `make run`
  - Visit http://localhost:8080/healthz
- With OTEL enabled, set for example:
  - `export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318`
  - `export OTEL_SERVICE_NAME=driftlock-api`
  - `make run`

Endpoints
- `GET /healthz` – liveness
- `GET /readyz` – readiness (immediate for now)
- `GET /v1/version` – build/version info
- `POST /v1/events` – ingest JSON event payloads (stub feeding engine)

Configuration (env)
- `PORT` (default: 8080)
- `OTEL_EXPORTER_OTLP_ENDPOINT` (e.g., http://localhost:4318) – enables OTEL when set
- `OTEL_SERVICE_NAME` (default: driftlock-api)
- `OTEL_ENV` (default: dev)
- `DRIFTLOCK_VERSION` (optional build-time/version string)

Project layout
- `api-server/cmd/driftlock-api/` – service entrypoint
- `api-server/internal/api/` – HTTP mux and handlers
- `api-server/internal/engine/` – processing core skeleton
- `api-server/internal/telemetry/` – OpenTelemetry setup
- `pkg/version/` – version string plumbed to the API
- `collector-processor/` – OTel processor skeleton `driftlock_cbad`
- `cbad-core/` – Rust CBAD primitives with OpenZL integration
- `llm-receivers/` – Collector receiver skeleton for LLM I/O
- `exporters/` – Evidence bundle exporters (placeholders)
- `deploy/` – Docker Compose and Collector config
- `tools/synthetic/` – synthetic event generator
- `ui/` – minimal Next.js placeholder
- `web/` – minimal static page stub
- `deps/openzl/` – OpenZL compression framework (git submodule)

Dev notes
- Run `go mod tidy` after first `make run` or `make build` to fetch modules.
- Dockerfile builds the API service binary; adjust base image/labels as needed.
- Keep the UI minimal for MVP; expand once the pipeline is flowing.

## OpenZL Integration: Format-Aware Compression for Anomaly Detection

Driftlock uses Meta's OpenZL as its primary compression engine, providing a significant competitive advantage over traditional anomaly detection systems.

### What is OpenZL?

OpenZL is a format-aware compression framework that understands the structure of your data (JSON logs, timeseries metrics, nested traces) rather than treating it as opaque byte streams. Unlike generic compressors (zstd, gzip), OpenZL:

- Parses data structure and applies intelligent transforms (struct-of-arrays, delta encoding, tokenization)
- Learns optimal compression strategies offline via training on representative data
- Embeds the decode recipe in each compressed frame (no out-of-band coordination needed)
- Provides both better compression ratios AND faster speeds on structured data

**Performance on structured data:**
- 1.5-2x better compression ratios than zstd
- 20-40% faster compression/decompression speeds
- Deterministic output with fixed compression plans

### Why OpenZL for Anomaly Detection?

Format-aware compression provides more sensitive anomaly signals than generic compression:

**Better Baselines:** OpenZL learns the "normal" structure of your telemetry during training. When data deviates from this learned structure, compression ratios drop dramatically, signaling anomalies.

**Granular Insights:** Instead of "this blob compressed poorly," you get "the 'message' field had unusually low compression due to unexpected content length" or "new field 'stack_trace' not in trained schema."

**Structured Explanations:** Field-level compression metrics enable precise root cause analysis. For example, detecting that error logs have a 10x larger message field than baseline logs.

**Example Scenario:**

Normal log: `{"level": "info", "msg": "request completed", "duration_ms": 42}`

Anomalous log: `{"level": "error", "msg": "SEGFAULT...", "stack_trace": "...500 lines..."}`

- Generic compressor (zstd): Modest compression ratio change
- OpenZL: Dramatic compression failure due to structure violation + field-level attribution

### Competitive Differentiation

**Unique Market Position:** Driftlock is the only anomaly detection platform using format-aware compression for OTLP telemetry analysis.

**Key Advantages:**
- Glass-box compression: Explainable WHY compression failed (critical for DORA, NIS2, AI Act compliance)
- Deterministic training: Fixed seed produces reproducible compression plans for audit trails
- No black-box ML: Compression theory (Kolmogorov complexity) provides mathematical foundation
- Novel IP: Format-aware CBAD is defensible differentiation

**Enterprise Value:** "We use Meta's advanced OpenZL compression framework, optimized for your specific OTLP data formats, providing explainable anomaly detection that meets regulatory compliance requirements."

### Technical Architecture

Driftlock implements a multi-algorithm compression adapter pattern:

- **Primary Engine:** OpenZL with pre-trained plans for OTLP logs, metrics, and traces
- **Fallback Options:** zstd, lz4, gzip for unstructured data or baseline comparisons
- **Training:** Compression plans trained offline on representative OTLP schemas
- **Integration:** Rust FFI bindings to OpenZL C/C++ library via static linking

See [docs/OPENZL_ANALYSIS.md](docs/OPENZL_ANALYSIS.md) for detailed technical analysis, benchmarks, and integration roadmap.

Next steps
- **Phase 1 Priority**: Implement OpenZL compression adapter in cbad-core with Rust FFI bindings (see [ROADMAP.md](ROADMAP.md))
- Train OpenZL compression plans for OTLP logs, metrics, and traces
- Benchmark OpenZL vs zstd/lz4 on representative OTLP datasets
- Integrate cbad-core via `driftlock_cbad` in the Collector and route telemetry
- Expand api-server with storage (Postgres) and SSE/WebSockets for live anomalies
- Build minimal Next.js UI to visualize anomalies and compression metrics

## 📋 Enterprise Features

Driftlock includes enterprise-ready compliance and governance features:

### Regulatory Compliance
- **DORA Compliance**: Digital Operational Resilience Act evidence bundles
- **NIS2 Compliance**: EU cybersecurity incident reporting templates  
- **Runtime AI Monitoring**: AI Act compliance for LLM/ML systems
- **Audit Trails**: Cryptographically signed evidence packages

### Documentation Framework
- **[ALGORITHMS.md](docs/ALGORITHMS.md)**: Mathematical foundations and CBAD principles
- **[CODING_STANDARDS.md](docs/CODING_STANDARDS.md)**: Development guidelines and quality standards
- **[BUILD.md](docs/BUILD.md)**: Comprehensive build and deployment instructions
- **[CONTRIBUTING.md](docs/CONTRIBUTING.md)**: Contributor guidelines and workflows

### Advanced Tooling
- **Benchmarking**: Performance validation and regression testing
- **CI/CD Pipeline**: Automated quality gates and validation
- **Decision Logging**: Architectural decision tracking for audit compliance

See the [docs/](docs/) directory for complete enterprise documentation.

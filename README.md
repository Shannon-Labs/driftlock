Driftlock (clean slate)

This repository is a minimal, production-minded scaffold modeled after the GlassBox Monitor blueprint but keeping the Driftlock name. It emphasizes a Go-based API server, a future OpenTelemetry Collector processor, and a Rust CBAD core, with OpenTelemetry hooks and a minimal UI path.

What’s here (Phase 0/1 skeleton)
- api-server: Go service with health, readiness, version, and basic event ingestion endpoints, OTEL-enabled.
- cbad-core: Rust crate placeholder for Compression-Based Anomaly Detection (CBAD) primitives (FFI later).
- collector-processor: OTel Collector processor skeleton `driftlock_cbad` (wires in later with cbad-core).
- llm-receivers: Collector receiver skeletons for LLM prompts/responses/tool calls.
- exporters: Placeholders for evidence bundle exporters (JSON + PDF).
- deploy: Docker Compose with an OTel Collector and API container.
- tools/synthetic: Small Go generator to POST synthetic events to the API.
- ui: Placeholder Next.js app directory (intentionally minimal for now).
- web: Minimal static HTML stub to sanity-check the API.

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
- `cbad-core/` – Rust CBAD primitives (placeholder for now)
- `llm-receivers/` – Collector receiver skeleton for LLM I/O
- `exporters/` – Evidence bundle exporters (placeholders)
- `deploy/` – Docker Compose and Collector config
- `tools/synthetic/` – synthetic event generator
- `ui/` – minimal Next.js placeholder
- `web/` – minimal static page stub

Dev notes
- Run `go mod tidy` after first `make run` or `make build` to fetch modules.
- Dockerfile builds the API service binary; adjust base image/labels as needed.
- Keep the UI minimal for MVP; expand once the pipeline is flowing.

Next steps
- **Phase 1 Priority**: Implement cbad-core metrics and FFI (see [ROADMAP.md](ROADMAP.md))
- Integrate cbad-core via `driftlock_cbad` in the Collector and route logs/metrics
- Expand api-server with storage (Postgres) and SSE/WebSockets for live anomalies
- Build minimal Next.js UI to visualize anomalies and artifacts

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

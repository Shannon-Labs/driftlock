# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Driftlock** is a compression-based anomaly detection (CBAD) platform for OpenTelemetry data, powered by Meta's OpenZL format-aware compression framework. It provides explainable anomaly detection for regulated industries through advanced compression analysis of logs, metrics, traces, and LLM I/O.

## High-Level Architecture

### Core Components

The system follows a streaming architecture with these main components:

1. **cbad-core (Rust)** - Compression-based anomaly detection algorithms with FFI bindings
   - Implements compression ratio analysis, Normalized Compression Distance (NCD), permutation testing
   - Provides `libcbad_core.a` static library for Go integration
   - Supports OpenZL format-aware compression alongside zstd/lz4/gzip

2. **collector-processor (Go)** - OpenTelemetry Collector processor
   - Integrates with cbad-core via CGO
   - Processes OTLP logs, metrics, traces in real-time
   - Computes CBAD metrics and flags anomalies

3. **api-server (Go)** - REST API service
   - Handles event ingestion, anomaly retrieval, evidence export
   - Integrates with PostgreSQL for storage, Kafka for streaming
   - Provides SSE/WebSocket for real-time anomaly streaming
   - Built-in OpenTelemetry instrumentation

4. **llm-receivers (Go)** - OpenTelemetry Collector receivers
   - Specialized receivers for LLM prompts, responses, tool calls
   - Enables AI/ML system monitoring

5. **exporters (Go)** - Evidence bundle generation
   - Creates JSON/PDF reports for regulatory compliance
   - DORA, NIS2, and AI Act compliant evidence bundles

6. **ui (Next.js)** - Minimal dashboard
   - Browse anomalies, streams, and artifacts
   - Real-time monitoring interface

### Data Flow

```
OTLP Sources → Collector Receivers → driftlock_cbad Processor → API Server → UI
                     ↓
                Evidence Bundles (exporters)
                     ↓
                Kafka Stream → Anomaly Storage (PostgreSQL)
```

### Technology Stack

- **Go 1.22+** for API server, collectors, and exporters
- **Rust 1.70+** for CBAD core algorithms
- **PostgreSQL** for primary storage
- **Kafka** for streaming anomalies
- **Redis** for state management
- **OpenTelemetry** for observability
- **Next.js** for UI

## Common Commands

### Development Workflow

```bash
# Start API server in development mode (port 8080)
make run

# Run all tests
make test

# Build all components
make build

# Clean build artifacts
make clean
```

### Building Specific Components

```bash
# Build API server binary
make api

# Build OpenTelemetry collector with CBAD integration
make collector

# Build Rust CBAD static library
make cbad-core-lib

# Build development tools (synthetic data generator)
make tools

# Run full CI validation locally
make ci-check
```

### Code Quality

```bash
# Format Go and Rust code
make fmt

# Run linters
make lint

# Update dependencies
make tidy
```

### Testing

```bash
# Run all Go tests with verbose output
go test ./... -v

# Run Rust tests
cd cbad-core && cargo test

# Run Go benchmarks
go test -bench=. ./...

# Run Rust benchmarks
cd cbad-core && cargo bench

# Run end-to-end tests
go test ./tests/e2e/... -v
```

### Docker & Deployment

```bash
# Build Docker image
make docker

# Start local development stack
docker compose -f deploy/docker-compose.yml up

# Start with observability
docker compose -f deploy/docker-compose.yml \
  -f deploy/docker-compose.observability.yml up

# Build release binaries for multiple platforms
make release
```

## Key Configuration

### Environment Variables

- `PORT` - API server port (default: 8080)
- `OTEL_EXPORTER_OTLP_ENDPOINT` - OTEL endpoint for telemetry export
- `OTEL_SERVICE_NAME` - Service name for tracing (default: driftlock-api)
- `OTEL_ENV` - Environment name (default: dev)
- `DRIFTLOCK_VERSION` - Version string override
- `CBAD_CORE_PROFILE` - Rust build profile: release or dev

### API Endpoints

- `GET /healthz` - Liveness check
- `GET /readyz` - Readiness check
- `GET /v1/version` - Build/version information
- `POST /v1/events` - Ingest JSON event payloads
- `GET /v1/anomalies` - Retrieve anomalies
- `GET /v1/evidence` - Export evidence bundles

## Core Algorithms & Math

### Compression-Based Anomaly Detection (CBAD)

The system uses several mathematical approaches:

1. **Compression Ratio**: `CR = compressed_bytes / raw_bytes`
2. **Delta Bits**: `Δ_bits = ((C_B - C_W) × 8) / |W|`
3. **Shannon Entropy**: `H = -Σ p_i log₂ p_i`
4. **Normalized Compression Distance**: `NCD = (C_{BW} - min(C_B, C_W)) / max(C_B, C_W)`

### Statistical Significance

Permutation testing with deterministic seeds (ChaCha20 RNG) for p-value calculation:
```
p-value = (1 + #{permutations where metric ≥ observed}) / (1 + total_permutations)
```

### Deterministic Rules

- All compression operations use deterministic seeds
- Window sizes (baseline/window/hop) configured explicitly
- Avoid non-deterministic concurrency in core algorithms
- Regex-based PII detection with configurable policies

## Code Organization

### Main Directories

- `api-server/` - Go API service (main entry point: `api-server/cmd/driftlock-api/main.go`)
- `cbad-core/` - Rust crate implementing CBAD algorithms
- `collector-processor/` - OpenTelemetry Collector processor
- `llm-receivers/` - Collector receivers for LLM I/O
- `exporters/` - Evidence bundle exporters
- `pkg/version/` - Version information
- `docs/` - Comprehensive documentation
- `deploy/` - Docker Compose and Kubernetes manifests

### API Server Structure

```
api-server/
├── cmd/driftlock-api/       # Entry point
├── internal/
│   ├── api/                 # HTTP handlers and routing
│   ├── auth/                # Authentication/authorization
│   ├── billing/             # Stripe integration
│   ├── cbad/                # CBAD integration layer
│   ├── compression/         # Compression adapters
│   ├── handlers/            # Business logic handlers
│   ├── middleware/          # HTTP middleware (logging, rate limiting)
│   ├── models/              # Data models
│   ├── services/            # Business services
│   ├── storage/             # Database abstraction (PostgreSQL, Redis)
│   ├── streaming/           # Kafka integration
│   └── telemetry/           # OpenTelemetry setup
└── migrations/              # Database migrations
```

### Key Files

- `Makefile` - Build and development commands
- `go.mod` - Go dependencies and module replacements
- `.env.example` - Environment variable template
- `docs/BUILD.md` - Detailed build instructions
- `docs/ARCHITECTURE.md` - Architecture overview
- `docs/ALGORITHMS.md` - Mathematical foundations

## Development Guidelines

### Coding Standards

**Go Standards:**
- Target Go 1.22+
- Follow OpenTelemetry Collector conventions
- No `panic` - return errors with context
- Enforce gofmt, goimports, staticcheck
- Table-driven tests with deterministic seeds

**Rust Standards:**
- Edition 2021
- No `unsafe` without documented justification
- Use `cargo fmt` and `cargo clippy -D warnings`
- Property tests for compression math
- Deterministic RNG seeds

**Testing Requirements:**
- ≥80% line coverage on core math packages
- Deterministic seeds in benchmarks
- Edge cases and error paths covered

### Build Tags

- `-tags driftlock_cbad_cgo` - Enable CGO integration with CBAD core

### Important Integration Points

1. **Rust-Go FFI Boundary**
   - Located in `collector-processor/`
   - CGO linking with `libcbad_core.a`
   - Use `make collector` to build

2. **OpenTelemetry Integration**
   - All components instrumented with OTEL
   - Trace context propagation across boundaries
   - Configurable exporters via environment variables

3. **Storage Layer**
   - PostgreSQL for persistent storage
   - Redis for ephemeral state
   - Kafka for streaming anomalies
   - Tiered storage pattern in `api-server/internal/storage/`

## Testing & Debugging

### Local Testing

```bash
# Start local stack
docker compose -f deploy/docker-compose.yml up

# Generate synthetic test data
go run ./tools/synthetic

# Check endpoints
curl http://localhost:8080/healthz
curl http://localhost:8080/v1/version
```

### Debugging Tips

- Set `RUST_LOG=debug` for detailed Rust logging
- Use `go run -tags driftlock_cbad_cgo` for testing FFI integration
- Monitor `/metrics` endpoint for Prometheus metrics
- Check `backend.log`, `firebase-debug.log`, `pipeline.log` for diagnostics

### Performance Issues

- Verify CBAD window sizing configuration
- Monitor Go heap and Rust allocations
- Check for memory leaks at FFI boundary
- Profile compression algorithm selection

## Regulatory Compliance

The system includes enterprise-ready compliance features:

- **DORA Compliance** - Digital Operational Resilience Act evidence bundles
- **NIS2 Compliance** - EU cybersecurity incident reporting
- **AI Act Compliance** - Runtime AI monitoring for LLM/ML systems
- **Audit Trails** - Cryptographically signed evidence packages

Evidence bundles generated in `exporters/` directory with JSON/PDF formats.

## OpenZL Integration

**Key Competitive Advantage:** Meta's OpenZL format-aware compression framework provides:
- 1.5-2x better compression ratios than zstd on structured data
- 20-40% faster compression/decompression
- Format-aware transforms (struct-of-arrays, delta encoding, tokenization)
- Glass-box compression with explainable anomaly detection

Located at `deps/openzl/` as a git submodule. Provides deterministic training and embedded decode recipes.

## Database Schema

Initial schema in `api-server/migrations/`:
- `api-server/migrations/001_initial_schema.up.sql` - Creates tables
- `api-server/internal/storage/migrations/` - Migration management

Primary entities:
- Tenants (multi-tenant architecture)
- Anomalies (detected anomalies with metadata)
- Evidence bundles (compliance reports)
- Telemetry streams (raw data references)

## Deployment

### Options

1. **Kubernetes with Helm**
   ```bash
   helm install driftlock ./deploy/helm/driftlock
   ```

2. **Docker Compose**
   ```bash
   docker compose -f deploy/docker-compose.yml up
   ```

3. **Manual Binary Deployment**
   ```bash
   make release
   scp bin/driftlock-api-* target-host:/opt/driftlock/
   ```

### Production Readiness

- Health checks: `/healthz`, `/readyz`, `/metrics`
- Prometheus metrics endpoint
- CBAD performance metrics
- Anomaly detection rate monitoring

## Documentation

Comprehensive documentation in `docs/`:
- `ARCHITECTURE.md` - System design
- `ALGORITHMS.md` - Mathematical foundations
- `BUILD.md` - Build and deployment guide
- `CODING_STANDARDS.md` - Development guidelines
- `CONTRIBUTING.md` - Contribution workflow
- `API.md` - API reference
- `DEPLOYMENT.md` - Production deployment guide

See also phase summaries and roadmaps for implementation progress tracking.

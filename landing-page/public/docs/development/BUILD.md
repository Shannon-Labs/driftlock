# Driftlock Build and Deployment Guide

This document describes how to build the deterministic CBAD core and link it into the Go collector. The steps below assume the required toolchain is installed.

## Prerequisites

### Go Development
- Go 1.24 or newer
- Standard Go toolchain with cgo support

### Rust Development  
- Rust 1.70+ with stable toolchain
- Cargo with staticlib support

### System Dependencies
- Linux: `build-essential`, `pkg-config`
- macOS: Xcode command line tools
- Windows: Visual Studio Build Tools or MinGW-w64

### Node.js Development
- Node.js 18+ with npm or pnpm
- For dashboard development

## Quick Start

### Local Development

```bash
# Clone and setup
git clone https://github.com/Shannon-Labs/driftlock.git
cd driftlock

# Initialize git submodules (OpenZL and nested dependencies)
git submodule update --init --recursive

# Build API server (basic version without CBAD integration)
make run

# Test endpoints
curl http://localhost:8080/healthz
curl http://localhost:8080/v1/version
```

### Full Build with CBAD Integration

```bash
# Ensure git submodules are initialized (OpenZL and nested dependencies)
git submodule update --init --recursive

# Build all components
make build

# Or build individual components:
make cbad-core-lib    # Build Rust CBAD core
make collector        # Build OTel collector processor
make api              # Build API server

# Run tests
make test
```

## Build Targets

### Core Components

- `make api` - Build the API server binary to `bin/driftlock-api`
- `make collector` - Build the OTel collector processor
- `make cbad-core-lib` - Build the Rust CBAD static library
- `make tools` - Build synthetic data generator and benchmarks

### Development Workflow

- `make run` - Start API server locally (port 8080)
- `make test` - Run all Go tests
- `make clean` - Clean build artifacts
- `make ci-check` - Run full CI validation locally
- `make migrate` - Run database migrations
- `make dev` - Start full development environment with Docker Compose

## Rust CBAD Core

### Building the Static Library

The Rust crate produces a `staticlib` artifact (`libcbad_core.a`) for Go integration:

```bash
make cbad-core-lib
```

The library is emitted at `cbad-core/target/release/libcbad_core.a`. 
Set `CBAD_CORE_PROFILE=dev` for debug builds.

### Rust Development

```bash
cd cbad-core

# Format and lint
cargo fmt
cargo clippy --all-targets -- -D warnings

# Test
cargo test

# Benchmark
cargo bench
```

## Go Collector Integration

The Go collector requires the static library and CGO enabled:

```bash
make collector
```

This command:
- Ensures `libcbad_core.a` exists (building it first if necessary)
- Invokes `go build` with CGO enabled so the Rust FFI is linked by default
- Links against the Rust static library via cgo

> Need to disable the Rust dependency temporarily? Use `go build -tags driftlock_no_cbad` to force the stub implementation while keeping the regular build path unchanged.

## Testing

### Unit Tests
```bash
# Go tests
go test ./...

# Rust tests  
cargo test

# Full test suite
make test
```

### Integration Testing
```bash
# Start local stack
docker compose -f deploy/docker-compose.yml up

# Run synthetic data generator
go run ./tools/synthetic

# Check API endpoints
curl http://localhost:8080/v1/events
```

### Benchmarking
```bash
# Rust benchmarks
cargo bench

# Go benchmarks
go test -bench=. ./...

# Performance validation
make benchmark
```

## Docker Deployment

### Building Images

```bash
# Build API server image
docker build -t driftlock:api .

# Build collector image (when ready)
docker build -f deploy/Dockerfile.collector -t driftlock:collector .
```

### Running with Compose

```bash
# Basic stack
docker compose -f deploy/docker-compose.yml up

# With observability
docker compose -f deploy/docker-compose.yml \
  -f deploy/docker-compose.observability.yml up
```

## Environment Configuration

### Required Environment Variables

- `PORT` - API server port (default: 8080)
- `OTEL_EXPORTER_OTLP_ENDPOINT` - OTEL endpoint for telemetry export
- `OTEL_SERVICE_NAME` - Service name for tracing (default: driftlock-api)

### Optional Configuration

- `OTEL_ENV` - Environment name (default: dev)
- `DRIFTLOCK_VERSION` - Version string override
- `CBAD_CORE_PROFILE` - Rust build profile (release/dev)

### Development Setup

```bash
# Copy example environment
cp .env.example .env

# Edit configuration
$EDITOR .env

# Source environment
source .env

# Run with configuration
make run
```

## CI/CD Pipeline

### Local CI Validation

```bash
# Run full CI checks
make ci-check
```

This validates:
- Code formatting (Go and Rust)
- Linting and static analysis
- Unit test coverage
- Build artifact generation
- Integration test readiness

### GitHub Actions

The repository includes workflows for:
- Continuous integration on pull requests
- Security scanning and dependency updates
- Performance regression testing
- Automated releases and tagging

## Troubleshooting

### Common Build Issues

1. **CGO Linking Errors**
   ```bash
   # Ensure Rust library is built
   make cbad-core-lib
   
   # Check library exists
   ls -la cbad-core/target/release/libcbad_core.a
   
   # Rebuild with verbose output
   CGO_LDFLAGS="-v" make collector
   ```

2. **Rust Compilation Errors**
   ```bash
   # Update Rust toolchain
   rustup update stable
   
   # Clean and rebuild
   cargo clean && make cbad-core-lib
   ```

3. **Docker Build Issues**
   ```bash
   # Use multi-stage build caching
   docker build --target builder -t driftlock:builder .
   docker build --from driftlock:builder -t driftlock:api .
   ```

### Performance Issues

1. **Memory Usage**
   - Verify CBAD window sizing configuration
   - Monitor Go heap and Rust allocations
   - Check for memory leaks in FFI boundary

2. **Throughput Issues**
   - Profile compression algorithm selection
   - Validate streaming buffer management
   - Monitor GC pressure and allocation patterns

### Development Tips

- Use `make run` for rapid iteration on API changes
- Rust changes require `make cbad-core-lib` before Go builds
- Set `RUST_LOG=debug` for detailed Rust logging
- Use `go run ./...` (CGO enabled) for FFI integration, or `go run -tags driftlock_no_cbad` to force the stub

## Production Deployment

### Prerequisites

- Kubernetes cluster (1.20+) or Docker Swarm
- Persistent storage for evidence bundles
- Observability stack (Prometheus, Grafana, Jaeger)
- TLS certificates for secure communication

### Deployment Options

1. **Kubernetes with Helm**
   ```bash
   helm install driftlock ./deploy/helm/driftlock
   ```

2. **Docker Compose (Simple)**
   ```bash
   docker compose -f deploy/docker-compose.prod.yml up -d
   ```

3. **Manual Binary Deployment**
   ```bash
   # Build release binaries
   make release
   
   # Deploy to target environment
   scp bin/driftlock-api target-host:/opt/driftlock/
   systemctl start driftlock-api
   ```

### Health Checking

Production deployments should monitor:
- `/healthz` - Basic liveness check
- `/readyz` - Readiness with dependency validation
- `/metrics` - Prometheus metrics endpoint
- CBAD performance metrics and anomaly detection rates

For detailed production deployment guidance, see the deployment documentation in `deploy/README.md`.

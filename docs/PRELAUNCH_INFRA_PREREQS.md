# Driftlock Pre-Launch Infrastructure Prerequisites

This guide captures the runtime setup required to exercise the remaining pre-launch checklist items (sections 2–5 of the master list). It focuses on getting the local environment ready to run Docker Compose, Kubernetes deployments, databases, and the CBAD processor.

## 1. Local Docker & Compose

- **Start Docker Desktop** (macOS) or ensure the Docker daemon is running (`docker info` should succeed).
- Export the environment expected by the API server. Copy `.env.example` and fill in secrets or run
  ```bash
  export $(grep -v '^#' .env.example | xargs)
  ```
- Verify images build with `make docker`. This wraps `docker build -t driftlock:api .` and expects the daemon to be healthy.
- To exercise the production stack locally, use:
  ```bash
  docker compose -f deploy/docker-compose.yml up --build
  ```
  Ensure ports `8080`, `5432`, `9092`, `6379`, `8123`, `9000`, `9090`, `3000`, `4317`, and `4318` are free.

## 2. Kubernetes & Helm

- Install tooling: `kubectl`, `helm`, and a Kubernetes cluster (Kind, Minikube, or managed). Cluster must be **v1.22+** per the checklist.
- Load container images into the cluster’s registry or push them to a registry accessible by the cluster. The `Makefile` release target emits multi-OS binaries under `bin/`.
- Charts live under `helm/`. Typical workflow:
  ```bash
  helm dependency update helm/driftlock
  helm install driftlock helm/driftlock --namespace driftlock --create-namespace
  ```
- Provide secrets (PostgreSQL, Kafka credentials, TLS materials) via `kubectl create secret` before installing the chart.
- For observability, wire the OTEL collector address via environment variables (`OTEL_EXPORTER_OTLP_ENDPOINT`, etc.).

## 3. Databases & Streaming Backends

- **PostgreSQL**: Provision according to `deploy/production/docker-compose.yml`. Run migrations (if present) before smoke tests. Confirm TLS/SSL requirements for production.
- **Redis**: Used for distributed state. Configure address, password, and DB index through `REDIS_*` variables.
- **Kafka**: Ensure brokers listed in `.env.example` are reachable. Topic creation and retention policies are described in `deploy/production/docker-compose.yml`.
- **ClickHouse**: Required for analytics endpoints. Confirm schema migrations.
- **S3 or archival storage**: configure tiered storage if checklist requires hot/warm/cold flows.

## 4. CBAD Processor & Collector

- Build the Rust core once via `make cbad-core-lib` (generates `cbad-core/target/release/libcbad_core.a`).
- The Go collector builds with `make collector`. A placeholder CLI lives in `collector-processor/cmd/driftlock-collector/main.go`; wiring into an OTEL service remains TODO.
- When cross-compiling (`make release`), CGO is disabled; stub implementations in `collector-processor/driftlockcbad` allow builds to succeed, but anomaly detection is inactive without CGO.
- For full functionality, compile with:
  ```bash
  CGO_ENABLED=1 go build -tags driftlock_cbad_cgo ./collector-processor/cmd/driftlock-collector
  ```
  Ensure a system C compiler is installed and the Rust library path exists.

## 5. Observability & Metrics

- Grafana/Prometheus configuration resides in `deploy/production`. Import dashboards as needed.
- Set OTEL exporter endpoints via environment variables before running `make run` or `make collector`.
- Synthetic traffic tool ships under `tools/synthetic`. Build with `make tools` and run to drive `/readyz` and `/v1/version` smoke tests.

## 6. Docker Registry & CI/CD

- Authenticate to the target container registry (`docker login`) before pushing release images.
- Update CI pipelines to use the new `make collector` and `make release` flows that rely on the CGO stubs when cross-compiling.
- Capture release artifacts:
  - `bin/driftlock-api-linux-amd64`
  - `bin/driftlock-api-darwin-amd64`
  - `bin/driftlock-api-windows-amd64.exe`
  - `bin/driftlock-collector` (build with CGO to enable anomaly detection)

## 7. Next Checks

With the prerequisites above in place, proceed with:

- Docker Compose smoke test (`deploy/production/docker-compose.yml`).
- Kubernetes deployment using Helm.
- End-to-end validation using the synthetic traffic generator.
- Performance and load tests (`make benchmark`, `tools/synthetic`).
- Security tooling (vulnerability scanning for containers, TLS verification).

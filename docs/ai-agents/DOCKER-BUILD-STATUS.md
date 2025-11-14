# Docker Build Status - Driftlock

## Current State (2025-03-XX)

- **Rust core:** ✅ `cargo build --release` and the CLI demo both succeed.
- **Docker images:** ✅ `collector-processor/cmd/driftlock-http/Dockerfile` and `collector-processor/cmd/driftlock-collector/Dockerfile` build cleanly on Debian Bookworm.
- **Automation:** ✅ `scripts/test-docker-build.sh` now runs inside CI (GitHub Actions job `docker`) so regressions are caught immediately.
- **Runtime health:** ✅ `/healthz` validates the linked CBAD library and reports whether OpenZL was compiled in.

## How to Build/Run

```bash
# Build the HTTP API image locally
make docker-http  # wraps docker build -f collector-processor/cmd/driftlock-http/Dockerfile .

# Bring up just the HTTP API via Compose
docker compose up --build driftlock-http

# Smoke-test
curl -s http://localhost:8080/healthz | jq
curl -s -X POST http://localhost:8080/v1/detect \
  -H 'Content-Type: application/json' \
  -d @test-data/financial-demo.json | jq '.anomaly_count'
```

`scripts/test-docker-build.sh` builds both HTTP + collector images with generic compressors. Set `ENABLE_OPENZL_BUILD=true` to request the optional OpenZL variants (requires `libopenzl.a` artifacts; see `docs/OPENZL_ANALYSIS.md`).

## Known Issues / Intentional Gaps

1. **OpenZL libraries are proprietary.** The repository does not distribute them. Builds default to zstd/lz4/gzip; OpenZL images are skipped automatically unless the private artifacts are mounted under `openzl/` or `OPENZL_LIB_DIR`.
2. **Collector image is experimental.** It compiles but is only enabled when Compose profile `kafka` is selected.

## Validation Checklist

- `make docker-http` → success ✅
- `scripts/test-docker-build.sh` → success ✅ (OpenZL builds skipped unless toggled)
- `docker compose up driftlock-http` → `/healthz` returns `{ "library_status": "healthy", ... }`
- GitHub Actions `docker` job green ✅

## Next Steps

1. Supply OpenZL artifacts to CI runners (or split into a separate workflow) to exercise the optional feature flag end-to-end.
2. Harden the collector Docker path once Kafka/OTel streaming moves past prototype status.
3. Consider publishing pre-built images (ghcr.io) once driftlock-http’s API is frozen.

With the above workflow, Docker parity is no longer a bottleneck: any commit that breaks the build is rejected automatically, and partners can run `docker compose up --build driftlock-http` to exercise the same stack showcased in the CLI demo.

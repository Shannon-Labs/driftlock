# OpenZL Integration Status

_Last updated: 2025-11-24_

**IMPORTANT: OpenZL is optional and experimental.** Driftlock works with generic compressors (zstd, lz4, gzip) by default, which are always available and sufficient for all demo and production use cases. OpenZL is a format-aware compression adapter (BSD licensed, available at https://github.com/facebook/openzl) that can provide better compression ratios but is **disabled by default** (`USE_OPENZL=false` in all demo scripts). The OpenZL path is strictly opt-in and requires building the submodule.

## Current Behavior

- **OpenZL is disabled by default** - all builds and demos use generic compressors (zstd, lz4, gzip) unless explicitly enabled.
- `cbad-core` builds without OpenZL by default. The FFI call `cbad_has_openzl()` exposes whether the feature flag was compiled in. Go bindings surface this via `driftlockcbad.HasOpenZL()`.
- `cmd/driftlock-http`'s `/healthz` endpoint reports `openzl_available: true/false` and only lists `openzl` in `available_algos` when the Rust core includes those symbols.
- Requests that specify `algo=openzl` automatically fall back to `zstd` when OpenZL isn't compiled in; the HTTP response includes `"fallback_from_algo": "openzl"` so operators can see the downgrade.
- **All demo scripts use `USE_OPENZL=false`** - the default user experience does not require or use OpenZL.

## Quick Start with `just` Commands (Recommended)

The easiest way to work with OpenZL is using the standardized `just` commands:

```bash
# Check OpenZL status
just openzl-status

# Build OpenZL library (first time only)
just build-openzl-lib

# Build cbad-core with OpenZL
just build-core-openzl

# Test OpenZL integration (Rust + Go)
just test-openzl

# Build Docker image with OpenZL
just docker-build-openzl

# Clean OpenZL artifacts
just clean-openzl
```

## Manual Build Steps (Alternative)

If you prefer manual builds or need to customize the process:

1. **Sync the OpenZL submodule** so nested dependencies (like the bundled zstd) are available:
   ```bash
   git submodule update --init --recursive openzl
   ```
2. **Build OpenZL** from source (the repository is included as a git submodule at `openzl/`). You'll need to build `libopenzl.a` and ensure headers are available. Because Driftlock links the static archive into a shared library, force position-independent code during the build:
   ```bash
   (cd openzl && CFLAGS="-fPIC ${CFLAGS:-}" make lib)
   ```
   (The Dockerfiles under `cbad-core/` and `collector-processor/cmd/` patch this automatically; for manual builds you must pass `-fPIC` yourself.)
3. The build system will look for OpenZL under `openzl/` at the repository root _or_ you can set `OPENZL_LIB_DIR=/absolute/path/to/openzl`.
4. Build the Rust core with the feature enabled:
   ```bash
   cd cbad-core
   OPENZL_LIB_DIR=/path/to/openzl cargo build --release --features openzl
   ```
   The build script now fails fast with a clear error if `libopenzl.a` cannot be located, so you are never surprised during linking.
5. Build Go binaries or Docker images with `USE_OPENZL=true` so the Rust stage enables the feature:
   ```bash
   docker build \
     -f collector-processor/cmd/driftlock-http/Dockerfile \
     --build-arg USE_OPENZL=true \
     --build-arg RUST_VERSION=1.82 \
     --build-arg GO_VERSION=1.24 \
     .
   ```

## Docker Compose (OpenZL path)

- **Build + run:** `USE_OPENZL=true docker compose build driftlock-http && USE_OPENZL=true docker compose up driftlock-postgres driftlock-http`
- **Verify locally:** `curl -s http://localhost:8080/healthz | jq '.openzl_available, .available_algos'` should return `true` and include `"openzl"` alongside `zstd`, `lz4`, and `gzip`. If accessing from outside the container fails, run `docker compose exec driftlock-http curl -s http://localhost:8080/healthz`.
- **Fallback:** When `USE_OPENZL` is `false` (default) or OpenZL assets are missing, the image builds with the generic compressors only. Requests that ask for `openzl` will transparently fall back to `zstd` and the HTTP response will include `"fallback_from_algo": "openzl"`.

The Dockerfiles for `driftlock-http` and `driftlock-collector` now copy the bundled `openzl/` tree. Setting `USE_OPENZL=true` triggers the OpenZL build; leaving it `false` avoids the extra work while keeping the default demos fast.

## Prefer OpenZL When Available (opt-in)

- Set `PREFER_OPENZL=true` in the API environment to automatically choose `openzl` as the default compressor when the binary includes OpenZL symbols. Overrides are respected: explicit request payload `compressor` values still win, and we always fall back to `zstd` if OpenZL is missing.
- `/healthz` and Prometheus now export `openzl_available` so operators can see whether the OpenZL path is active.
- A deterministic benchmark (`openzl_bench_test.go`) runs only when OpenZL is compiled; use `go test -run TestDoesNotExist -bench OpenZLPreference ./cmd/driftlock-http` to compare `openzl` vs `zstd` locally.

## Docker & CI Story

- `scripts/test-docker-build.sh` now autodetects whether OpenZL artifacts are present. When they are missing (default), the script builds only the generic images and clearly prints that OpenZL-only variants were skipped. Set `ENABLE_OPENZL_BUILD=true` and expose the libraries to force those builds.
- GitHub Actions runs the script on every push/PR (generic compressors only) so the Dockerfiles never rot. Optional OpenZL jobs can be added to build from the submodule when needed.

## Operator Guidance

- **Default behavior: OpenZL is disabled.** All demos, CI, and default builds use generic compressors only.
- Production deployments should treat OpenZL as an experimental accelerator: enable it only where the library can be built and installed; otherwise rely on zstd (which is the default and recommended path).
- `/healthz` plus Prometheus metrics expose whether OpenZL is active so SREs can alert on misconfigured hosts.
- **Documentation and customer deliverables must clearly state that OpenZL is optional and experimental.** Generic compressors (zstd, lz4, gzip) are the supported baseline and work for all use cases.

For additional details or troubleshooting steps, see the [OpenZL documentation](https://github.com/facebook/openzl) or coordinate with Shannon Labs engineering. The OpenZL source is included as a git submodule and can be built from source.

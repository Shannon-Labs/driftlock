# OpenZL Integration Status

_Last updated: 2025-03-XX_

OpenZL remains an **optional**, proprietary compression adapter. Driftlock ships with deterministic fallback compressors (zstd, lz4, gzip) enabled everywhere. The OpenZL path is now first-class again but strictly opt-in.

## Current Behavior

- `cbad-core` builds without OpenZL by default. The new FFI call `cbad_has_openzl()` exposes whether the feature flag was compiled in. Go bindings surface this via `driftlockcbad.HasOpenZL()`.
- `cmd/driftlock-http`’s `/healthz` endpoint reports `openzl_available: true/false` and only lists `openzl` in `available_algos` when the Rust core includes those symbols.
- Requests that specify `algo=openzl` automatically fall back to `zstd` when OpenZL isn’t compiled in; the HTTP response includes `"fallback_from_algo": "openzl"` so operators can see the downgrade.

## How to Enable OpenZL Locally

1. **Obtain the private artifacts** from Meta/OpenZL (the `libopenzl.a` static library plus headers).
2. Place them under `openzl/` at the repository root _or_ set `OPENZL_LIB_DIR=/absolute/path/to/openzl`.
3. Build the Rust core with the feature enabled:
   ```bash
   cd cbad-core
   OPENZL_LIB_DIR=/path/to/openzl cargo build --release --features openzl
   ```
   The build script now fails fast with a clear error if `libopenzl.a` cannot be located, so you are never surprised during linking.
4. Build Go binaries or Docker images with `USE_OPENZL=true` so the Rust stage enables the feature:
   ```bash
   docker build \
     -f collector-processor/cmd/driftlock-http/Dockerfile \
     --build-arg USE_OPENZL=true \
     --build-arg RUST_VERSION=1.82 \
     --build-arg GO_VERSION=1.24 \
     .
   ```

## Docker & CI Story

- `scripts/test-docker-build.sh` now autodetects whether OpenZL artifacts are present. When they are missing (default), the script builds only the generic images and clearly prints that OpenZL-only variants were skipped. Set `ENABLE_OPENZL_BUILD=true` and expose the libraries to force those builds.
- GitHub Actions runs the script on every push/PR (generic compressors only) so the Dockerfiles never rot. Optional OpenZL jobs can be added later once a redistributable artifact exists.

## Operator Guidance

- Production deployments should treat OpenZL like an accelerator: enable it where the proprietary library can be installed; otherwise rely on zstd.
- `/healthz` plus Prometheus metrics expose whether OpenZL is active so SREs can alert on misconfigured hosts.
- Documentation and customer deliverables must continue to describe OpenZL as optional; deterministic fallbacks are the supported baseline.

For additional details or troubleshooting steps, coordinate with Shannon Labs engineering — the proprietary library cannot be redistributed in this repository.

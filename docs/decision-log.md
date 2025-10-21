# Decision Log

All significant assumptions and architectural decisions for Driftlock.

| Date       | Decision | Rationale | Consequences |
|------------|----------|-----------|--------------|
| 2025-01-09 | Adopt Apache 2.0 license | Aligns with OpenTelemetry ecosystem and encourages regulated adopters to contribute improvements without copyleft concerns. | Requires contributor agreement to Apache 2.0; downstreams can use commercially. |
| 2025-01-09 | Treat Go 1.22 as minimum version | Ensures access to current OTel Collector APIs and generics improvements while staying within LTS support windows. | CI/tooling must install Go 1.22; downstreams on older Go must upgrade. |
| 2025-01-09 | Implement CBAD core in Rust with FFI | Rust offers memory safety, performance, and easy WASM targeting for future UI analytics; FFI keeps Go processor lean. | Requires Cargo toolchain in CI; additional FFI glue for Go/WASM consumers. |
| 2025-01-09 | Default module namespace `github.com/hmbown/driftlock` | Provides stable import path while repo remains private; changeable before broader release. | Need to update if canonical domain differs; documentation must note current location. |
| 2025-01-09 | Use OpenZL in pinned mode by default | Meets deterministic explainability requirements and avoids runtime training variability. | Build pipeline must produce and hash `.zlc` plans; add validation on startup. |
| 2025-01-09 | Establish compressor abstraction with deterministic configs | Trait-based wrapper over zstd/lz4/gzip keeps compression behaviour reproducible and enables future OpenZL integration. | Requires dependency pinning and benchmarking per backend; Go processor can switch compressors via config. |
| 2025-01-09 | Define C FFI API operating on aggregated baseline/window buffers | Simplifies Go integration by avoiding per-chunk allocations across language boundary while keeping math centralised in Rust. | Go-side processor must construct contiguous windows; future streaming API may be added if needed. |
| 2025-01-09 | Provide Go helpers that flatten telemetry windows deterministically | Ensures contiguous buffers and reuse of existing allocations when preparing FFI calls, matching the Rust API contract. | Processor code can assemble baseline/window payloads without repeated allocations; helpers are covered by unit tests. |
| 2025-01-09 | Standardise Makefile-driven build of `libcbad_core.a` for cgo linkage | Keeps build reproducible and enforces the `driftlock_cbad_cgo` tag workflow before the collector can call into Rust. | CI/doc updates require invoking `make cbad-core-lib`/`make collector`; integration tests remain tagged until static library publishing is automated. |
| 2025-01-09 | Add CI target to assert static lib availability and `driftlock_cbad_cgo` tag usage | Prevents regressions where builds skip the Rust artifact or forget the deterministic tag, aligning local dev with CI expectations. | CI runs `make ci-check` (which wraps `tools/ci/verify_cbad_build.sh`) before packaging; developers can mirror the flow locally. |
| 2025-01-09 | Adopt GlassBox Monitor roadmap and compliance framework | Leverage proven architecture and enterprise-ready compliance templates for faster time-to-market in regulated industries. | Inherit comprehensive roadmap with DORA, NIS2, and AI Act compliance; requires adaptation of GlassBox patterns to Driftlock branding. |
| 2025-01-09 | Integrate advanced documentation structure from GlassBox | Adopt comprehensive documentation framework including compliance templates, algorithm documentation, and coding standards. | Improved enterprise readiness and audit trail capabilities; requires maintaining documentation quality standards. |
| 2025-01-09 | Priority focus on enterprise compliance features | Target regulated industries with built-in compliance reporting and evidence bundle generation. | Clear market differentiation but requires deeper regulatory knowledge; may slow initial MVP development. |


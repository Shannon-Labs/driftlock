# Coding Standards

These standards govern all source contributions to Driftlock.

## Go
- Target Go 1.22 or newer.
- Follow OpenTelemetry Collector conventions (context-aware APIs, no global state, observability hooks).
- Prohibit `panic`; return errors with context and use structured logging.
- Enforce gofmt and goimports; staticcheck will run in CI.
- Table-driven tests with deterministic seeds; cover edge cases and error paths.

## Rust
- Edition 2021; forbid `unsafe` without a documented justification and targeted review.
- Use `cargo fmt` + `cargo clippy` with `-D warnings` in CI.
- Provide property tests for compression math where feasible; use deterministic RNG seeds.

## TypeScript / Next.js
- Target strict TypeScript with `tsconfig` `strict=true`.
- Prefer functional components with hooks; keep server/client boundaries explicit.
- Avoid implicit any; use zod or io-ts schemas for external data.

## Documentation
- All public functions require godoc/rustdoc comments including key formulas.
- Update `decision-log.md` for significant assumptions or trade-offs.
- Include ASCII diagrams in architecture docs for portability.

## Security & Privacy
- Treat all inputs as untrusted; validate lengths and schema before processing.
- Redact secrets before storage or compression when `privacy.redact_before_compress=true` is enabled.
- Never write raw payloads to compliance bundles unless explicitly configured.

## Testing & Coverage
- Maintain ≥80% line coverage on core math packages; publish coverage reports in CI.
- For benchmarks, pin dataset seeds and include environment notes in `decision-log.md`.

## Commit Discipline
- Each commit should be self-contained, buildable, and accompanied by updated tests/docs when configuration or behavior changes.
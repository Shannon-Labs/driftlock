# Repository Guidelines

## Project Structure & Module Organization
- `api-server/` hosts the Go service: `cmd/driftlock-api` is the entry point and `internal/{api,engine,telemetry}` split routing, logic, and OTEL wiring.
- Streaming pieces live under `collector-processor/driftlockcbad` and `llm-receivers/llmio`; the Rust anomaly core sits in `cbad-core/src/lib.rs`.
- Supporting assets: `pkg/version` for shared helpers, `deploy/` for compose configs, `docs/` for architecture notes, `tools/synthetic` for traffic generation, and minimal UIs in `ui/` and `web/`.

## Build, Test, and Development Commands
- `make run` starts the API locally (set `OTEL_EXPORTER_*` from `.env.example` when needed); verify at `http://localhost:8080/healthz`.
- `make build` or `make api` emit binaries to `bin/`; use `make collector` for the processor and `make tools` for the synthetic generator.
- `make test` (alias for `go test ./...`) covers all Go packages; run `cargo test` inside `cbad-core/` once Rust logic lands.
- `pnpm install && pnpm dev` (run in `ui/`) starts the Next.js stub; `docker compose -f deploy/docker-compose.yml up` spins up API + collector.

## Coding Style & Naming Conventions
- Format Go code with `gofmt` (`go fmt ./...`); keep package names lowercase and exported identifiers concise.
- Rust modules follow `cargo fmt` and `cargo clippy`; expose deterministic helpers that accept explicit seeds as noted in `docs/ARCHITECTURE.md`.
- Front-end files should stick to the default Next.js/Prettier style (2-space indent, camelCase components); keep environment keys uppercase snake case like `.env.example`.

## Testing Guidelines
- Name Go tests `*_test.go` and prefer table-driven cases for handlers in `api-server/internal/api` and processors in `collector-processor`.
- Seed RNGs in CBAD tests to keep runs deterministic and assert on structured JSON instead of logs.
- For smoke tests, run `make run`, `go run ./tools/synthetic`, then check `/readyz` and `/v1/version`.

## Commit & Pull Request Guidelines
- Write imperative commit subjects with optional scopes (e.g., `feat(api): add anomaly guard`) and keep the first line under 72 characters.
- Bundle related code, tests, and docs per commit; list validation commands in the message body.
- Pull requests need a concise summary, linked issues, and screenshots or sample payloads for UI or collector changes; merge only after green CI and an area reviewer.

## Security & Configuration Tips
- Start from `.env.example`, but inject secrets per environment via shell exports or your secret manager—never commit live credentials.
- When enabling OTEL or Docker, enforce TLS endpoints and rebuild with `docker build -t driftlock:api .` after upgrades, scanning images in your registry.

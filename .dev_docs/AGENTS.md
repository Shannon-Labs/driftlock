# Repository Guidelines

## Project Structure & Module Organization
- `cbad-core/`: Rust NCD/entropy engine; builds shared lib consumed by Go.
- `collector-processor/`: Go HTTP API and billing; Dockerfile in `collector-processor/cmd/driftlock-http`.
- `cmd/`: Go entrypoints; `cmd/demo/main.go` builds the `driftlock-demo` binary.
- `landing-page/`: Vue 3 + Tailwind dashboard (`src` for UI, `src/dataconnect-generated` vendored).
- `functions/` Firebase functions; `extensions/vscode-driftlock/` VS Code extension; `docs/` architecture notes; `scripts/` helper tooling; `test-data/` sample event sets.

## Build, Test, and Development Commands
- `just build` or `make demo`: build Rust core then Go demo; outputs `./driftlock-demo`.
- `DRIFTLOCK_DEV_MODE=true ./scripts/run-api-demo.sh`: start the local API wired to the built core.
- `cd landing-page && npm run dev`: dashboard dev server; `npm run build` for production assets.
- `just test`: cargo + go + Vue type-check + functions build + extension tests.
- `just lint` or `make verify`: go vet, cargo clippy, ESLint; verify adds release-readiness checks.
- Optional Docker: `just docker-build` then `docker-compose up` (or `-f docker-compose.kafka.yml`) for local stack.

## Coding Style & Naming Conventions
- Formatters are authoritative: `just fmt` (`cargo fmt`, `go fmt`), ESLint fixes frontend/extension code.
- Go: package names lower-case; exported symbols CamelCase; tests in `*_test.go`.
- Rust: snake_case modules/functions; keep `#[cfg(test)] mod tests` beside implementations.
- JS/TS/Vue: camelCase functions, PascalCase components/filenames; prefer Composition API with `<script setup>`.
- Place new scripts under `scripts/` and binaries under `bin/` or `cmd/` with descriptive names.

## Testing Guidelines
- `cd cbad-core && cargo test` for detector math; add fixtures under `tests/` or inline modules.
- `go test ./...` per Go module (run after `just build-core`); keep tests deterministic with `test-data/` inputs.
- Frontend: `cd landing-page && npm run type-check`; add component/e2e coverage (Playwright/Vue Test Utils) when altering UI flows.
- Extension: `cd extensions/vscode-driftlock && npm test` for activation/command coverage.
- Add regression cases alongside new logic; prefer writing a failing test before fixes.

## Commit & Pull Request Guidelines
- Conventional commits enforced via `commitlint.config.js`; types like `feat`, `fix`, `docs`, `perf`, `deps`; scopes lower-case; subject no period and ≤100 chars.
- Keep PRs scoped; summarize purpose, list test commands run, and link issues/tasks.
- Include screenshots/GIFs for UI changes (`landing-page`), and call out API/CLI breaking changes.
- Do not commit secrets or generated credentials; use env vars locally and rotate shared tokens via maintainers.
- Before merge, ensure CI parity by running `just test` or targeted suites relevant to your change.

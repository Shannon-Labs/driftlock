# Driftlock Agent Guidelines - SaaS Platform Era

These instructions are for **AI agents and automation** working in this repository. Updated for the **Firebase SaaS Platform** architecture (January 2025).

---

## 1. Current State (Post-SaaS Integration)

**IMPORTANT**: This repo has transformed from a public-facing demo to a **SaaS platform foundation**.

### Agent-First Workflows (NEW)
We have standardized workflows to help you navigate the repo. **ALWAYS** check these first:
-   **Setup**: `.agent/workflows/setup.md`
-   **Test**: `.agent/workflows/test.md`
-   **Build**: `.agent/workflows/build.md`
-   **Lint**: `.agent/workflows/lint.md`
-   **Deploy**: `.agent/workflows/deploy.md`

### Unified Command Runner
We use `just` to standardize commands. Run `just --list` to see all available recipes.
-   `just setup`: Install dependencies
-   `just test`: Run all tests
-   `just build`: Build all artifacts
-   `just lint`: Lint all code

### Architecture Overview:
- **Firebase Hosting** (`landing-page/`) - Professional customer-facing website
- **Firebase Functions** (`functions/`) - API layer with AI integration 
- **Cloud Run Backend** (existing Go/Rust) - Core anomaly detection engine
- **GitHub Repository** - Technical reference for YC/developers, NOT primary user interface

### Before Making Changes, Read:
- `docs/deployment/FIREBASE_SAAS_COMPLETE.md` - SaaS transformation summary
- `docs/deployment/AI_COST_OPTIMIZATION.md` - Cost-efficient AI strategy
- `docs/launch/ROADMAP.md` - Updated SaaS launch roadmap
- `.archive/reports/FINAL-STATUS.md` - Core demo functionality (still must work)
- `README.md` - Project overview

### Developer Tooling Inventory (November 2025 update):
- **`Justfile`** – The single source of truth for running tasks.
- **`cmd/driftlock-cli`** – `driftlock scan` streams NDJSON/STDIN through `pkg/entropywindow`; keep flags (`--format`, `--follow`, `--stdin`) stable.
- **`pkg/entropywindow/`** – Shared Go analyzer for CLI + MCP. Never break its API without updating both callers.
- **`extensions/vscode-driftlock/`** – VS Code Live Radar extension. Use `just build-extension` and `just test-extension`.
- **`cmd/driftlock-mcp/`** – Claude/Cursor MCP server with local entropy fallback. `detect_anomalies` tool must always accept raw strings.
- **Landing page automation** – `scripts/deploy-landing.sh`, `scripts/verify-horizon-datasets.ts`, and Playwright specs under `landing-page/tests/` gate the Horizon Showcase.
- **Live Soak Test** – `scripts/soak_runner.py` and `scripts/crypto_bridge.py` for long-running validation against live Crypto Volatility data (see `docs/launch/SOAK_TEST.md`).
- **Evidence scripts** – `scripts/chaos-report.py` + `docs/launch/THE_CHAOS_REPORT.md` summarize benchmark datasets; regenerate when logs change.

---

## 2. Golden Invariants (NEVER BREAK)

### Core Demo (Must Always Work):
- ✅ `make demo` and `./scripts/verify-launch-readiness.sh` must succeed
- ✅ CLI demo produces deterministic HTML reports
- ✅ Mathematical explanations remain audit-ready
- ✅ Same input = same output (reproducible for compliance)

### SaaS Platform (Must Deploy Clean):
- ✅ `firebase deploy` must work without errors
- ✅ Landing page builds (`cd landing-page && npm run build`)
- ✅ No hard-coded secrets in frontend code
- ✅ API endpoints follow `/api/*` pattern for Firebase routing

### Cost Optimization (Critical):
- ❌ **NO expensive AI calls in default/demo flows**
- ✅ AI analysis is opt-in premium feature only
- ✅ Core anomaly detection must be fast (<2 seconds)
- ✅ Mathematical explanations are primary value, not AI commentary

---

## 3. Development Workflow

### Firebase-First Development:
```bash
# Always test locally first
cd landing-page && npm run dev

# Build before deploy
cd landing-page && npm run build && cd ..

# Deploy complete stack
firebase deploy
```

### Code Changes Priority:
1. **Landing page** (`landing-page/`) - Customer acquisition focus
2. **Firebase Functions** (`functions/`) - API layer and AI integration
3. **Core engine** (`cbad-core/`, `collector-processor/`) - Detection logic

### Authentication Integration:
- Use Firebase Auth for user management
- Integrate with existing API key system in Cloud Run backend
- Maintain backward compatibility with existing `/v1/onboard/signup`

### Landing Page + Showcase Automation:
- Use `make deploy-landing` or `scripts/deploy-landing.sh` (expects `FIREBASE_TOKEN`) for hosting pushes.
- Run `scripts/verify-horizon-datasets.ts` (via `npm run verify:horizon` in repo root) after touching sample data or the Horizon Showcase component.
- UI gating lives in `landing-page/tests/`; Playwright suite (`npx playwright test`) must stay green with the auto-start preview server defined in `playwright.config.ts`.

---

## 4. Business Logic Rules

### Pricing & Features:
- **Free Trial**: Core detection, mathematical explanations, basic compliance
- **Pro ($99/month)**: AI insights, advanced reports, priority support  
- **Enterprise**: Custom pricing, white-glove deployment, unlimited AI

### User Experience:
- **Instant signup** → API key in 30 seconds
- **Fast detection** → Mathematical results in <2 seconds
- **AI as upsell** → "Upgrade for business insights" messaging
- **Compliance first** → DORA/NIS2/AI Act positioning

### Technology Stack:
- **Frontend**: Vue 3 + Firebase Hosting (primary customer interface)
- **API Layer**: Firebase Functions + Gemini Pro (cost-optimized)
- **Backend**: Cloud Run + PostgreSQL (existing anomaly detection)
- **Auth**: Firebase Auth + API key system

---

## 5. Deployment Strategy

### Google-First Architecture:
- ✅ Firebase for frontend hosting and functions
- ✅ Google Cloud Run for backend compute
- ✅ Google Cloud SQL or continue with Supabase
- ✅ Gemini Pro for AI features
- ✅ Google Domains for domain management

### Domain Strategy:
**Recommended**: Move domain to Google Domains + Firebase Hosting
- Simpler management (all Google ecosystem)
- Better integration with Firebase custom domains
- Automatic SSL/CDN via Firebase

**Alternative**: Keep Cloudflare + point CNAME to Firebase
- More complex but keeps Cloudflare features
- DNS: `driftlock.net CNAME -> firebase-project.web.app`

---

## 6. Never Do This

❌ **Break the core demo** (YC reviewers need to see it works)
❌ **Add expensive API calls to default flows** (unit economics)
❌ **Expose technical details in landing page** (customer-facing must be business-focused)
❌ **Hard-code production secrets** (use Firebase config/environment variables)
❌ **Make AI required for core functionality** (mathematical explanations are sufficient)

---

## 7. Integration Points

### Firebase ↔ Cloud Run:
- Firebase Functions proxy requests to Cloud Run backend
- Maintain existing `/v1/*` API structure via `/api/proxy/*`
- Add new `/api/signup`, `/api/analyze` for SaaS features

### Authentication Flow:
- Firebase Auth for user accounts
- Generate API keys via Cloud Run backend
- Store user metadata in Firebase/Firestore
- Validate API keys in Cloud Run for anomaly detection

### Cost Management:
- AI analysis only for premium users
- Cache common insights to reduce Gemini calls
- Use Firebase Functions for rate limiting and user validation

---

## 8. Success Metrics

### Technical:
- Firebase deployment succeeds in <5 minutes
- Landing page loads in <2 seconds
- Signup flow works end-to-end
- Core demo still passes verification

### Business:
- Landing page conversion >10% (visitors → signups)  
- API adoption >50% (signups → first API call)
- Cost per signup <$2 (including AI/infrastructure)
- Customer satisfaction >4.5/5

**Remember**: This is now a real SaaS business, not just a tech demo. Every change should improve customer acquisition, retention, or unit economics.
- **Explainability**
  - Any anomaly path must expose: NCD, compression ratios, entropy, p‑value, confidence, and a short human‑readable explanation string.
  - Do not add opaque “ML-style” black boxes without clear, auditable outputs.

---

## 9. Documentation Structure

We have organized documentation to keep the root clean and context clear:

- **`docs/architecture/`**: Core design, algorithms, API schema (`API.md`, `ARCHITECTURE.md`).
- **`docs/deployment/`**: Cloud Run, Firebase, Cloud SQL setup guides (`DEPLOYMENT.md`, `CLOUDSQL_FIREBASE_SETUP_GUIDE.md`).
- **`docs/compliance/`**: DORA, NIS2, US regulations (`COMPLIANCE_*.md`).
- **`docs/launch/`**: Roadmap, launch checklists, use cases (`LAUNCH_SUMMARY.md`, `ROADMAP.md`).
- **`docs/development/`**: Build guides, coding standards, contributing (`DEVELOPMENT.md`, `CONTRIBUTING.md`).
- **`docs/ai-context/`**: Prompts and AI-specific instructions (`AI_*.md`).
- **`docs/user-guide/`**: End-user documentation and quickstarts.
- **`.archive/`**: Old plans, reports, and superseded docs.

---

## 10. Code‑Level Guidelines

### Rust (`cbad-core/`)

- Treat `cbad-core` as the **source of truth** for CBAD math and metrics:
  - Keep FFI surfaces small, documented, and stable (`src/ffi*.rs`).
  - Minimize `unsafe`; any unsafe block must be simple and obviously correct.
- Do **not** change `crate-type` or features in `Cargo.toml` in ways that break:
  - The Go FFI in `collector-processor/driftlockcbad/`.
  - The CLI demo.
- Compression adapters:
  - Generic compressors (zlab/zstd/lz4/gzip) must continue to work in all builds.
  - OpenZL adapter is **optional** and must fail gracefully; never make it a hard runtime dependency.

### Go (`collector-processor/`, `cmd/`)

- Keep `collector-processor/cmd/driftlock-http/main.go` as the **canonical HTTP engine**:
  - `/healthz` should reflect CBAD and compression adapter health.
  - `/v1/detect` is the primary public detection endpoint; changes to its request/response shape must be documented in `docs/architecture/API.md`.
- FFI (`collector-processor/driftlockcbad/*.go`):
  - Do not change C signatures without updating the corresponding Rust exports.
  - On error, return clear Go errors; do not panic on normal failure modes.

### Developer Tooling (CLI, VS Code, MCP)

- **CLI (`cmd/driftlock-cli`)**: `driftlock scan` is the supported entrypoint. Preserve streaming options (`--format`, `--follow`, `--stdin`, `--baseline-lines`, `--algo`, `--show-all`). Any change to analyzer output must keep `pkg/entropywindow.Result` backward-compatible.
- **VS Code extension (`extensions/vscode-driftlock`)**: Keep analyzer subprocess invocations using `driftlock scan --stdin`. Always run `npm run lint`, `npm run compile`, and `npm run test` before publishing `.vsix` builds.
- **MCP server (`cmd/driftlock-mcp`)**: `detect_anomalies` must accept large raw strings, chunk them with the entropy window, and auto-detect JSON vs raw payloads. Update `docs/integrations/MCP_SETUP.md` when the schema changes.

### Frontend (`playground/`, `landing-page/`)

- Maintain the existing tech choices (Vue 3 + TS for playground; Vue/Tailwind stack for landing page).
- Keep components small and composable; favour derived state and clearly typed props.
- Do not introduce heavy new UI frameworks or state managers without strong justification.

---

## 11. OpenZL and Compression Strategy

- OpenZL is a **format‑aware compressor** that can provide better compression and sharper anomaly signals but:
  - It is **not required** for correctness.
  - It may not be available in all environments (especially Docker).
- Rules:
  - Default builds (especially Docker images) must work with generic compressors only.
  - If you add or modify OpenZL integration:
    - Keep it behind feature flags and/or explicit build args (e.g., `USE_OPENZL=true`).
    - Add clear error paths and fallbacks to zstd when OpenZL libraries, plans, or symbols are unavailable.
    - Update `.archive/reports/OPENZL_ANALYSIS.md` with what is supported and how to build it.

---

## 12. Docker and Deployment

- Use the existing Docker files as the primary deployment path:
  - `docker-compose.yml`
  - `collector-processor/cmd/driftlock-http/Dockerfile`
- Goals:
  - `docker compose up` at repo root should:
    - Build and run `driftlock-http` successfully with generic compressors.
    - Pass its health check on `/healthz`.
  - Avoid introducing unnecessary OS‑level or toolchain complexity.
- If you add OpenZL‑enabled images:
  - Do so in **additional paths** (extra Dockerfile or guarded build args), not by breaking the default images.

---

## 13. Testing and Verification

- Before concluding work that touches core logic, FFI, Docker, or the HTTP API:
  - Run `make demo` and `./scripts/verify-launch-readiness.sh` if available.
  - Run any relevant scripts under `scripts/` (e.g., `test-api.sh`, `test-docker-build.sh`, `test-services.sh`) that cover your changes.
- Frontend + showcase gates:
  - `cd landing-page && npm run lint && npm run type-check && npm run build`.
  - `cd landing-page && npx playwright test` (auto-starts `npm run preview`).
  - `node scripts/verify-horizon-datasets.ts` (or `npm run verify:horizon`) to confirm datasets still load.
- Developer tooling:
  - `cd extensions/vscode-driftlock && npm run lint && npm run compile && npm run test` before shipping.
  - `cd pkg/entropywindow && go test ./...` anytime analyzer math changes; propagate via `go test` inside `cmd/driftlock-cli` and `cmd/driftlock-mcp`.
- Only modify or add tests that are clearly related to the behaviour you are changing.
- Keep tests fast and focused; avoid adding slow end‑to‑end suites without a good reason.

---

## 14. Documentation Expectations

- When you change behaviour, configuration, or public interfaces, update:
  - `README.md` and `docs/launch/DEMO.md` if the end‑user flow changes.
  - `docs/architecture/API.md` for HTTP/API changes.
  - `.archive/reports/OPENZL_ANALYSIS.md` and/or `docs/ai-agents/DOCKER-BUILD-STATUS.md` for compression and Docker changes.
  - `docs/launch/ROADMAP_TO_LAUNCH.md` only when adjusting high‑level roadmap assumptions.
- Keep documentation honest about what is implemented in this repo vs. what is future/roadmap.

---

## 15. Scope and Restraint

- This repo intentionally focuses on:
  - The CBAD core engine.
  - The CLI demo and playground.
  - A thin HTTP detection service and basic Docker story.
- The larger platform (full API server, exporters, multi‑tenant UI, Kafka, ClickHouse, etc.) lives mostly in design docs. Do **not** attempt to fully implement the entire platform here unless explicitly requested.

When in doubt, prefer **small, reversible, well‑documented changes** that move Driftlock closer to a pilot‑ready anomaly detection service while preserving the existing demo and mathematical guarantees.


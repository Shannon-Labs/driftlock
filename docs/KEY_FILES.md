# Key Files Reference

Quick reference for important files in the Driftlock codebase.

## Backend (Rust)

| File | Purpose |
|------|---------|
| `crates/driftlock-api/src/main.rs` | Server setup, Prometheus init, background tasks |
| `crates/driftlock-api/src/routes/mod.rs` | Route registration and router composition |
| `crates/driftlock-api/src/routes/detection.rs` | Detection endpoint handlers |
| `crates/driftlock-api/src/routes/onboarding.rs` | Signup/auth handlers |
| `crates/driftlock-api/src/routes/billing.rs` | Stripe integration |
| `crates/driftlock-api/src/routes/streams.rs` | Stream management |
| `crates/driftlock-api/src/routes/anomalies.rs` | Anomaly endpoints |
| `crates/driftlock-api/src/routes/anchors.rs` | Anchor/drift detection |
| `crates/driftlock-api/src/routes/api_keys.rs` | API key management |
| `crates/driftlock-api/src/middleware/auth.rs` | Authentication middleware |
| `crates/driftlock-api/src/state.rs` | Application state |

## Database Layer (Rust)

| File | Purpose |
|------|---------|
| `crates/driftlock-db/src/lib.rs` | Database module exports |
| `crates/driftlock-db/src/pool.rs` | Connection pool setup |
| `crates/driftlock-db/src/models/tenant.rs` | Tenant model |
| `crates/driftlock-db/src/models/stream.rs` | Stream model |
| `crates/driftlock-db/src/models/anomaly.rs` | Anomaly model |
| `crates/driftlock-db/src/models/api_key.rs` | API key model |
| `crates/driftlock-db/src/models/anchor.rs` | Stream anchor model |
| `crates/driftlock-db/src/repos/` | Repository implementations |

## Authentication (Rust)

| File | Purpose |
|------|---------|
| `crates/driftlock-auth/src/lib.rs` | Auth module exports |
| `crates/driftlock-auth/src/firebase.rs` | Firebase JWT verification |
| `crates/driftlock-auth/src/api_key.rs` | API key authentication |

## Frontend (Vue)

| File | Purpose |
|------|---------|
| `landing-page/src/components/cta/SignupForm.vue` | Main signup form |
| `landing-page/src/views/VerifyEmailView.vue` | Email verification page |
| `landing-page/src/views/DashboardView.vue` | User dashboard with usage chart |
| `landing-page/src/views/HomeView.vue` | Landing page with pricing |
| `landing-page/src/components/playground/PlaygroundShell.vue` | Interactive playground |
| `landing-page/src/components/dashboard/UsageChart.vue` | Daily usage chart |
| `landing-page/src/router/index.ts` | Route definitions |
| `landing-page/src/stores/` | Pinia state stores |

## CBAD Core (Rust)

| File | Purpose |
|------|---------|
| `cbad-core/src/lib.rs` | Main library entry |
| `cbad-core/src/anomaly.rs` | Anomaly detection logic |
| `cbad-core/src/window.rs` | Sliding window implementation |
| `cbad-core/src/metrics/mod.rs` | Statistical metrics |
| `cbad-core/src/ffi.rs` | C FFI bindings |

## Infrastructure

| File | Purpose |
|------|---------|
| `Dockerfile` | Rust API container |
| `Cargo.toml` | Workspace configuration |
| `docker-compose.yml` | Local development stack |
| `.github/workflows/` | CI/CD pipelines |

## Documentation

| File | Purpose |
|------|---------|
| `CLAUDE.md` | AI assistant context |
| `docs/AI_ROUTING.md` | Agent task routing guide |
| `docs/LAUNCH_STATUS.md` | Launch checklist and progress |
| `docs/architecture/ARCHITECTURE.md` | System design |
| `docs/ALGORITHMS.md` | CBAD mathematical foundations |
| `docs/architecture/api/openapi.yaml` | OpenAPI 3.0 spec |

## Configuration

| File | Purpose |
|------|---------|
| `Cargo.toml` | Workspace dependencies |
| `crates/driftlock-api/Cargo.toml` | API dependencies |
| `.env.example` | Environment template |
| `.claudeignore` | Files to exclude from AI indexing |

## Legacy (Archived)

| File | Purpose |
|------|---------|
| `archive/go-backend/` | Original Go implementation (reference only) |

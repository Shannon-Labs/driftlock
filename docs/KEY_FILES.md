# Key Files Reference

Quick reference for important files in the Driftlock codebase.

## Backend (Go)

| File | Purpose |
|------|---------|
| `collector-processor/cmd/driftlock-http/main.go` | Route registration, server setup |
| `collector-processor/cmd/driftlock-http/onboarding.go` | Signup/verify handlers |
| `collector-processor/cmd/driftlock-http/billing.go` | Stripe integration |
| `collector-processor/cmd/driftlock-http/billing_cron.go` | Scheduled billing jobs |
| `collector-processor/cmd/driftlock-http/dashboard.go` | User dashboard endpoints |
| `collector-processor/cmd/driftlock-http/db.go` | Database operations |
| `collector-processor/cmd/driftlock-http/store_auth_ext.go` | API key management |
| `collector-processor/cmd/driftlock-http/auth.go` | Authentication middleware |
| `collector-processor/cmd/driftlock-http/email.go` | SendGrid email service |
| `collector-processor/cmd/driftlock-http/demo.go` | Anonymous demo endpoint |

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

## Database

| File | Purpose |
|------|---------|
| `api/migrations/20250301120000_initial_schema.sql` | Core tables |
| `api/migrations/20250302000000_onboarding.sql` | Onboarding fields |
| `api/migrations/20251119000000_add_stripe_fields.sql` | Stripe columns |
| `api/migrations/20251124000000_launch_enhancements.sql` | Verification & billing |

## Infrastructure

| File | Purpose |
|------|---------|
| `collector-processor/cmd/driftlock-http/Dockerfile` | API container |
| `docker-compose.yml` | Local development stack |
| `infrastructure/` | Terraform configs |
| `.github/workflows/` | CI/CD pipelines |
| `scripts/deploy-production.sh` | Production deploy script |

## Documentation

| File | Purpose |
|------|---------|
| `CLAUDE.md` | AI assistant context (you are here) |
| `docs/AI_ROUTING.md` | Agent task routing guide |
| `docs/LAUNCH_STATUS.md` | Launch checklist and progress |
| `docs/ARCHITECTURE.md` | System design |
| `docs/ALGORITHMS.md` | CBAD mathematical foundations |
| `docs/architecture/api/openapi.yaml` | OpenAPI 3.0 spec |

## CBAD Core (Rust)

| File | Purpose |
|------|---------|
| `cbad-core/src/lib.rs` | Main library entry |
| `cbad-core/src/algorithms/` | Compression algorithms |
| `cbad-core/src/ncd.rs` | Normalized Compression Distance |
| `cbad-core/src/permutation.rs` | Statistical significance testing |

## Configuration

| File | Purpose |
|------|---------|
| `Makefile` | Build commands |
| `go.mod` | Go dependencies |
| `.env.example` | Environment template |
| `.claudeignore` | Files to exclude from AI indexing |

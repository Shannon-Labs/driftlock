# AI Task Routing Guide

Quick reference for which agent handles which tasks in the Driftlock codebase.

## Task Routing Matrix

| Task Type | Primary Agent | When to Use |
|-----------|---------------|-------------|
| Go HTTP handlers | `dl-backend` | API endpoints, request handling |
| Stripe billing | `dl-backend` | Webhooks, checkout, subscriptions |
| API key management | `dl-backend` | Create, rotate, revoke keys |
| Vue components | `dl-frontend` | UI components, Pinia stores |
| TypeScript/Tailwind | `dl-frontend` | Frontend styling, types |
| PostgreSQL schema | `dl-db` | Migrations, indexes, queries |
| Database design | `dl-db` | Schema changes, performance |
| Go unit tests | `dl-testing` | Test files, coverage |
| E2E tests | `dl-testing` | Playwright, integration |
| Docker builds | `dl-devops` | Dockerfile, compose |
| Cloud Run deploy | `dl-devops` | GCP deployment, scaling |
| Project tracking | `dl-devops` | Checklists, status updates |
| Documentation | `dl-docs` | README, API docs, guides |
| Standup reports | `daily-standup` | Sprint status from Linear |
| TODO scanning | `todo-to-linear` | Find TODOs, create issues |

## Decision Tree

```
Is it frontend code (Vue, TypeScript, CSS)?
  └─ Yes → dl-frontend
  └─ No ↓

Is it database-related (SQL, migrations, schema)?
  └─ Yes → dl-db
  └─ No ↓

Is it deployment/infrastructure (Docker, GCP, CI/CD)?
  └─ Yes → dl-devops
  └─ No ↓

Is it testing (unit tests, E2E, coverage)?
  └─ Yes → dl-testing
  └─ No ↓

Is it documentation?
  └─ Yes → dl-docs
  └─ No ↓

Is it Go backend code (handlers, billing, auth)?
  └─ Yes → dl-backend
```

## Agent Capabilities

### dl-backend (Sonnet)
- HTTP handlers in `collector-processor/cmd/driftlock-http/`
- Stripe integration (billing.go, billing_cron.go)
- API key operations (store_auth_ext.go)
- Authentication middleware (auth.go)
- Database operations (db.go)

### dl-frontend (Sonnet)
- Vue 3 components in `landing-page/src/`
- Pinia state management
- Tailwind CSS styling
- Chart.js visualizations
- TypeScript types

### dl-db (Sonnet)
- PostgreSQL schema design
- Goose migrations in `api/migrations/`
- Index optimization
- Query performance
- Data modeling

### dl-devops (Sonnet)
- Docker builds and deployment
- Cloud Run configuration
- GitHub Actions workflows
- Terraform infrastructure
- Project checklist maintenance

### dl-testing (Sonnet)
- Go table-driven tests
- E2E integration tests
- Playwright browser tests
- Test coverage analysis

### dl-docs (Haiku)
- API documentation
- README updates
- User guides
- CLAUDE.md maintenance

### daily-standup (Haiku)
- Linear issue queries
- Sprint progress reports
- Blocker identification

### todo-to-linear (Haiku)
- Codebase TODO scanning
- Linear issue creation
- Technical debt tracking

## Common Workflows

### New Feature
1. `dl-db` - Design schema if needed
2. `dl-backend` - Implement API endpoints
3. `dl-frontend` - Build UI components
4. `dl-testing` - Write tests
5. `dl-devops` - Deploy and update checklist

### Bug Fix
1. `dl-testing` - Reproduce with test
2. `dl-backend` or `dl-frontend` - Fix the bug
3. `dl-testing` - Verify fix
4. `dl-devops` - Deploy

### Documentation Update
1. `dl-docs` - Write/update docs
2. `dl-devops` - Update project checklist if needed

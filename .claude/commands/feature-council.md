---
description: Multi-agent pipeline for end-to-end feature development (db -> backend -> frontend || docs -> testing)
argument-hint: Feature description (e.g., "Add usage dashboard with daily event charts")
---

# Feature Development Council

You are orchestrating a multi-agent pipeline to implement: **$ARGUMENTS**

## Pipeline Architecture

```
Phase 1: @dl-db       → Schema design, migrations
Phase 2: @dl-backend  → API endpoints (depends on Phase 1)
Phase 3: @dl-frontend → UI components (depends on Phase 2)
       + @dl-docs     → Documentation (parallel with frontend)
Phase 4: @dl-testing  → Tests for all components
```

## Execution Process

### Phase 1: Database Schema

Use the `dl-db` agent to:
1. Analyze feature requirements
2. Design necessary schema changes
3. Create migration file(s)
4. Test migration up/down locally

**Output:** Migration file path(s) and schema changes summary

### Phase 2: Backend API

Use the `dl-backend` agent to:
1. Read the schema changes from Phase 1
2. Implement API endpoint(s)
3. Add proper error handling and logging
4. Write handler tests

**Output:** New endpoint(s) and handler file changes

### Phase 3a: Frontend UI (parallel)

Use the `dl-frontend` agent to:
1. Read the API contract from Phase 2
2. Create Vue component(s)
3. Wire up to Pinia store if needed
4. Ensure responsive design

**Output:** New component(s) and view changes

### Phase 3b: Documentation (parallel with 3a)

Use the `dl-docs` agent to:
1. Update API documentation
2. Add usage examples
3. Update CLAUDE.md if needed

**Output:** Documentation updates

### Phase 4: Testing

Use the `dl-testing` agent to:
1. Write unit tests for backend
2. Write E2E tests for full flow
3. Run all tests and verify passing

**Output:** Test files and coverage report

## Progress Tracking

Use TodoWrite to track each phase. Example:

```
1. [completed] Phase 1: Created migration for usage_stats table
2. [in_progress] Phase 2: Implementing GET /v1/me/usage/stats endpoint
3. [pending] Phase 3a: Build UsageStatsChart component
4. [pending] Phase 3b: Document new endpoint
5. [pending] Phase 4: Write tests
```

## When to Pause

Stop and ask the user if:
- Schema design has multiple valid approaches
- API contract needs user input
- UI/UX decisions are ambiguous
- Tests reveal issues that need discussion

## Completion Criteria

Feature is complete when:
- [ ] Migration applies cleanly
- [ ] API endpoint works correctly
- [ ] UI displays data properly
- [ ] Documentation is updated
- [ ] All tests pass

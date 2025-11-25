---
description: Run all tests
---

# Test Workflow

1.  **Run All Tests**
    -   Command: `just test`
    -   This runs unit tests for Go, Rust, and TypeScript components.

2.  **Run Specific Tests**
    -   **Go Core**: `go test ./pkg/... ./cmd/...`
    -   **Rust Core**: `cd cbad-core && cargo test`
    -   **Landing Page**: `cd landing-page && npm run type-check` (and Playwright if configured)
    -   **Functions**: `cd functions && npm run build` (compile check)
    -   **VS Code Extension**: `cd extensions/vscode-driftlock && npm test`

3.  **Integration Tests**
    -   Run `scripts/verify-launch-readiness.sh` to verify the core demo flow.

---
description: Lint all code
---

# Lint Workflow

1.  **Run All Linters**
    -   Command: `just lint`

2.  **Specific Linters**
    -   **Go**: `golangci-lint run` (if installed) or `go vet ./...`
    -   **Rust**: `cd cbad-core && cargo clippy`
    -   **Landing Page**: `cd landing-page && npm run lint`
    -   **Functions**: `cd functions && npm run lint`
    -   **VS Code Extension**: `cd extensions/vscode-driftlock && npm run lint`

3.  **Fix Issues**
    -   Most linters support auto-fix.
    -   JS/TS: `npm run lint -- --fix`
    -   Rust: `cargo clippy --fix`

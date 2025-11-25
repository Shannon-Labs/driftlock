---
description: Build all artifacts
---

# Build Workflow

1.  **Build All**
    -   Command: `just build`

2.  **Build Specific Components**
    -   **Rust Core**: `cd cbad-core && cargo build --release`
    -   **Go Demo**: `go build -o driftlock-demo cmd/demo/main.go`
    -   **Landing Page**: `cd landing-page && npm run build`
    -   **Functions**: `cd functions && npm run build`
    -   **VS Code Extension**: `cd extensions/vscode-driftlock && npm run compile`

3.  **Docker Build**
    -   Command: `just docker-build`
    -   Builds `driftlock-http` image.

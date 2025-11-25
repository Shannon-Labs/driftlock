---
description: Setup the development environment
---

# Setup Workflow

1.  **Install System Dependencies**
    -   Ensure `go`, `rustc`, `cargo`, `node`, `npm`, `just`, and `docker` are installed.
    -   Install `firebase-tools`: `npm install -g firebase-tools`

2.  **Install Project Dependencies**
    -   Root: `npm install`
    -   Landing Page: `cd landing-page && npm install`
    -   Functions: `cd functions && npm install`
    -   VS Code Extension: `cd extensions/vscode-driftlock && npm install`

3.  **Environment Variables**
    -   Copy `.env.example` to `.env` if it doesn't exist.
    -   Ensure `FIREBASE_TOKEN` is set if deploying.

4.  **Verify Setup**
    -   Run `just --list` to see available commands.

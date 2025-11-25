# Institutional Grade Repo & Agent-First Workflow Plan

## Status: ✅ COMPLETE (November 2025)

All phases of the institutional upgrade have been implemented. This document serves as both a historical record and ongoing reference.

---

## Executive Summary
This plan outlines the steps to elevate the `driftlock` repository to an "institutional grade" standard. The goal is to maximize reliability, developer velocity, and "agent-first" ease of use. By standardizing workflows and environments, we ensure that both human developers and AI agents can operate with high confidence and minimal friction.

## Phase 1: Agent-First Foundation ✅ COMPLETE
The primary friction for agents is "figuring out how to run things." We solved this by explicitly codifying workflows.

-   **✅ `.agent/workflows/` Directory**:
    -   Created standardized markdown workflows for common tasks:
        -   `setup.md`: Environment bootstrapping
        -   `test.md`: Running test suites across all languages
        -   `lint.md`: Code quality checks
        -   `build.md`: Building all artifacts (Rust, Go, Docker, Frontend)
        -   `deploy.md`: Safe deployment procedures
    -   **Benefit**: Agents can read these files to know *exactly* how to perform complex tasks without hallucinating commands.

-   **✅ `AGENTS.md`**:
    -   Updated to reference the new workflows as the source of truth
    -   Added "Capability Map" describing which tools/scripts control which parts of the system

## Phase 2: Unified Command Interface & Repo Hygiene ✅ COMPLETE
Replaced ad-hoc shell scripts with a unified, self-documenting command runner.

-   **✅ `Justfile`**:
    -   Cross-platform command runner with `just --list` providing a perfect menu for agents
    -   Recipes implemented: `setup`, `test`, `lint`, `build`, `deploy`, `security`, `fmt`, `clean`
    -   Run `just --list` to see all ~30 available recipes

-   **✅ Pre-commit Hooks (`.pre-commit-config.yaml`)**:
    -   Enforces quality *before* code enters the repo
    -   **Go**: `go fmt`, `go vet`
    -   **Rust**: `cargo fmt`, `cargo clippy`
    -   **JS/TS**: `eslint`, `prettier`
    -   **Secrets**: `detect-secrets` to prevent accidental key leakage

-   **✅ Linters Configuration**:
    -   Standardized configs in appropriate subdirs (`landing-page/.eslintrc.cjs`, `functions/.eslintrc.js`)

## Phase 3: Standardized Development Environment ✅ COMPLETE
Eliminated "it works on my machine" issues.

-   **✅ `.devcontainer/devcontainer.json`**:
    -   Docker-based development environment containing:
        -   Go 1.24 toolchain
        -   Rust toolchain (with clippy, rustfmt)
        -   Node.js 22 & npm
        -   Python 3.11 (for scripts)
        -   Google Cloud SDK & Firebase CLI
        -   Docker-in-Docker support
        -   GitHub CLI
    -   **Port Forwarding**: 3000 (Firebase), 5173 (Vite), 8080 (API)
    -   **VS Code Extensions**: Go, Rust Analyzer, Docker, ESLint, Prettier, Volar, Just syntax

-   **✅ `.devcontainer/post-create.sh`**:
    -   Automatic installation of `just`, Google Cloud SDK, Firebase CLI
    -   Runs `just setup` on container creation

## Phase 4: CI/CD Hardening ✅ COMPLETE
Ensured the pipeline is robust and mirrors local checks.

-   **✅ GitHub Actions** (`.github/workflows/ci.yml`):
    -   **Parity**: CI runs the *exact same* `just` commands as local dev
    -   **Caching**: Aggressive caching for Cargo, Go modules, and npm
    -   **Concurrency**: Cancel obsolete builds on PR updates (`cancel-in-progress: true`)

-   **✅ Demo Verification** (`.github/workflows/yc-ready.yml`):
    -   Automated verification that demo always works
    -   Uses `just build-demo` for consistency

-   **✅ Docker Images** (`.github/workflows/docker.yml`):
    -   Automated image builds with OpenZL variants
    -   Pushes to GHCR on main branch

## Bonus Objectives ✅ COMPLETE

-   **✅ Dependabot** (`.github/dependabot.yml`):
    -   Automated dependency updates for: Go, Rust, npm, GitHub Actions, Docker
    -   Weekly schedule with appropriate PR limits per ecosystem

-   **✅ Security Scanning** (`.github/workflows/security.yml`):
    -   `govulncheck` for Go vulnerabilities
    -   `cargo audit` for Rust vulnerabilities
    -   `npm audit` for JS/TS vulnerabilities
    -   Trivy container scanning
    -   TruffleHog secret scanning
    -   CodeQL static analysis

-   **✅ Local Security** (`just security`):
    -   One command to run all security scans locally
    -   Subcommands: `just security-go`, `just security-rust`, `just security-npm`

---

## Quick Reference

### For Developers
```bash
# First time setup
just setup

# Daily workflow
just lint        # Check code quality
just test        # Run all tests
just build       # Build all artifacts

# Before committing
just fmt         # Format all code
just security    # Run security scans
```

### For AI Agents
1. Read `.agent/workflows/` for detailed procedures
2. Run `just --list` to discover commands
3. Check `AGENTS.md` for project-specific guidelines

---

## "Agent-First" Ease of Use Metrics ✅ ALL MET
-   **Discovery**: Can an agent find the "test" command in <1 step? ✅ Yes, via `just --list` or `.agent/workflows/test.md`
-   **Reliability**: Do commands fail due to missing env vars? ✅ No, `Justfile` validates dependencies
-   **Safety**: Can an agent accidentally deploy to prod? ✅ No, `deploy` requires Firebase auth

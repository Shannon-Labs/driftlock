Driftlock Docs

This directory houses architecture, plans, and compliance templates inspired by the GlassBox Monitor blueprint while keeping the Driftlock name.

## ⚠️ Important: Current vs. Planned Architecture

**What ships today** (as of 2025):
- Rust core (`cbad-core/`) with FFI for Go
- Go CLI demo (`cmd/demo/main.go`) that reads a static JSON file
- Synthetic demo data (`test-data/financial-demo.json`)
- One minimal CI workflow that builds, runs, and verifies the demo
- HTML output (`demo-output.html`) with anomaly cards and explanations

**What's documented here** (planned/future):
- Most docs in this directory describe a **future architecture** including API servers, databases, Docker deployments, Kafka integration, and web dashboards
- These are design documents and planning artifacts, not current implementation

**For the current demo**, refer to:
- `/README.md` - Main project documentation
- `/DEMO.md` - 2-minute walkthrough
- `/FINAL-STATUS.md` - Current status and what ships
- `/VERIFICATION-RESULTS.md` - Current test results

Files
- ARCHITECTURE.md – high-level system map and components (PLANNED)
- PHASE1_PLAN.md – concrete tasks for Phase 1 (CBAD + Collector)
- decision-log.md – running context, tradeoffs, and design decisions

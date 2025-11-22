# Driftlock Agents — `docs/`

These instructions apply to all work under `docs/` (architecture, roadmap, compliance, deployment, and AI‑agent docs).

---

## 1. Single source of truth

- Documentation in this directory is the **authoritative description** of:
  - Architecture and components.
  - Roadmap and phase status.
  - Compliance positioning (DORA, NIS2, AI Act, US regs).
  - Deployment and operational runbooks.
- Keep docs in sync with what the repo actually does:
  - If you implement or remove a feature, update the relevant docs in the same change.

---

## 2. Style and process

- Follow `docs/CODING_STANDARDS.md` for documentation expectations:
  - Use clear, concise prose with enough detail for an experienced engineer.
  - Add or update rustdoc/godoc comments when you change public APIs.
  - For major design decisions or trade‑offs, add an entry to `docs/decision-log.md`.
- Prefer Markdown with simple formatting; avoid heavy tooling or proprietary formats.

---

## 3. Roadmap and status docs

- Files like `ROADMAP_TO_LAUNCH.md`, `PHASE*_STATUS.md`, and `API_DEPLOYMENT_STATUS.md` must be:
  - Honest snapshots of actual progress.
  - Updated only when you have changed the underlying implementation or when plans are deliberately adjusted.
- Do not casually change timelines or completion markers without an accompanying code or infra change that justifies it.

---

## 4. AI‑agent docs

- Under `docs/ai-agents/`, keep:
  - Handoff notes up to date when you change Docker, build, or testing strategies.
  - Status reports accurate (what currently builds, what is disabled, what is future work).
- When you add new automation flows, consider documenting:
  - How agents should run tests and scripts.
  - Environmental assumptions (Docker version, OS, CPU arch).


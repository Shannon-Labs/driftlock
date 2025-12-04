# Driftlock Agents — `playground/`

These instructions apply to all work under `playground/` (Vue 3 + TypeScript detection playground).

---

## 1. Role of the playground

- The playground is a **developer/operator console** for experimenting with CBAD:
  - Paste NDJSON/JSON or load samples.
  - Auto‑derive baselines and window parameters.
  - Call `/v1/detect` and visualize anomalies and metrics.
- It is not a marketing site; keep it focused on clarity and debugging.

---

## 2. Tech stack and constraints

- Use Vue 3 + TypeScript + Vite + Tailwind (existing stack).
- Do not introduce large additional UI frameworks, state managers, or design systems.
- Keep components small and composable (`App.vue` + components in `src/components`).

---

## 3. API usage

- All backend interaction should go through the HTTP detection API:
  - Respect `VITE_API_BASE_URL` for the base URL; default to `http://localhost:8080`.
  - Use the `/v1/detect` endpoint; keep request/response handling in sync with `docs/API.md`.
- Treat the backend as untrusted:
  - Validate input sizes before sending.
  - Handle HTTP errors and JSON parse failures gracefully, surfacing actionable messages to the user.

---

## 4. UX guidelines

- Prioritize:
  - Transparency (show the parameters used: baseline, window, hop, algorithm).
  - Reproducibility (display equivalent curl commands).
  - Responsiveness and accessibility (reasonable defaults, keyboard navigation where sensible).
- Avoid:
  - Hiding or smoothing over raw metrics; advanced users should be able to see exact NCD, p‑values, and compression ratios.
  - Adding features that require a database or complex backend state; persistence should be local (e.g., browser) unless coordinated with backend changes.


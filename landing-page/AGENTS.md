# Driftlock Agents — `landing-page/`

These instructions apply to all work under `landing-page/` (marketing / product landing site).

---

## 1. Purpose

- This site is for **positioning and lead capture**, not for operating the product.
- It must:
  - Accurately describe Driftlock's capabilities and limitations.
  - Emphasize explainable, deterministic anomaly detection and compliance (DORA, NIS2, EU AI Act).
  - Avoid overstating features that are only present in design docs but not yet implemented.

---

## 2. Tech and quality bar

- Keep the existing stack: Vue 3 + TypeScript + Vite + Tailwind.
- Maintain:
  - Strong Lighthouse scores (performance, accessibility, SEO).
  - Responsive design and dark mode support.
  - Clean, modular components in `src/`.
- Do not add heavy client‑side trackers or third‑party widgets without clear value and privacy review.

---

## 3. Content discipline

- When changing copy:
  - Align with the current roadmap and implementation reality as documented in `docs/ROADMAP_TO_LAUNCH.md` and `FINAL-STATUS.md`.
  - Make it clear which capabilities are available **now** vs. **planned**.
- When adding CTAs or forms:
  - Prefer simple, backend‑agnostic integrations (e.g., webhook endpoints provided separately).
  - Do not hardcode secrets or endpoints tied to a specific environment.


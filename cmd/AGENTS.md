# Driftlock Agents — `cmd/`

These instructions apply to all Go programs under `cmd/` (including the CLI demo and any small tools).

---

## 1. CLI demo is sacred

- `cmd/demo` is the **flagship, human‑readable demo**:
  - Reads synthetic payment data from JSON.
  - Streams it through CBAD.
  - Produces `demo-output.html` with anomaly cards and explanations.
- Do not change:
  - The basic usage (`./driftlock-demo <path-to-financial-demo.json>`).
  - The qualitative behaviour (baseline warmup, ~10–30 anomalies, clear anomaly cards) without updating `DEMO.md`, `README.md`, and `FINAL-STATUS.md`.

---

## 2. Behavioural guarantees

- CLI tools should:
  - Fail fast on invalid input with clear error messages.
  - Print concise, operator‑friendly logs (no excessive noise).
  - Exit with non‑zero codes on error conditions.
- Whenever you adjust detection parameters in the demo (baseline count, detection interval, thresholds):
  - Keep them explicitly documented in the code comments and `DEMO.md`.
  - Preserve determinism; avoid runtime randomness.

---

## 3. New commands

- Any new CLI tools under `cmd/` should:
  - Be narrowly focused (e.g., a data converter, small utility, or diagnostic).
  - Reuse shared logic instead of duplicating CBAD plumbing.
  - Avoid user‑visible flags that overlap confusingly with the main detection API unless their semantics are identical.


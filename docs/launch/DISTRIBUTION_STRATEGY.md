# Driftlock Distribution Strategy: 2025-2026

**Objective:** Ubiquity. Make Driftlock the default "health check" for code, data, and infrastructure.

---

## 1. The "Trojan Horse" Integrations

We will not wait for users to download a CLI. We will go where they work.

### A. The AI Agent Ecosystem (MCP)
*   **Target:** Claude Desktop, Cursor, Windsurf, Smithery.
*   **Strategy:** Publish `driftlock-mcp` to the official MCP registries.
*   **The Hook:** "Give your Agent a Volatility Radar."
    *   *Scenario:* A coding agent is writing a loop. It uses Driftlock to check "Am I generating repetitive garbage?" (Self-correction).
*   **Action:** Polish `cmd/driftlock-mcp` and submit to Smithery.

### B. The IDE Extension (VS Code / Cursor)
*   **Target:** VS Code Marketplace.
*   **Strategy:** A lightweight wrapper around `driftlock-cli`.
*   **Feature:** "Live Log Linting." Highlight anomalous log lines in the terminal *as they appear*.
*   **Action:** Build `extensions/vscode-driftlock`.

### C. The Cloud Marketplace (AWS/GCP)
*   **Target:** AWS Marketplace, GCP Marketplace.
*   **Strategy:** "Driftlock Sidecar" container image.
*   **Feature:** "Drop-in anomaly detection for ECS/Cloud Run." 1-click deploy.
*   **Action:** Package `Dockerfile` as a marketplace asset.

---

## 2. Developer Advocacy

### "The Horizon Showcase" (Interactive Marketing)
*   **Asset:** The new "Factory-style" interactive playground on `driftlock.net`.
*   **Strategy:** Allow users to replay historical disasters (Terra Crash, Airline Meltdowns) in the browser.
*   **Conversion:** "See the anomaly with your own eyes -> Sign up to detect it in your data."
*   **Action:** Promote the "Horizon Showcase" on Hacker News and Twitter/X with "Can you spot the crash?" challenges.

### "The Chaos Engineering" Angle
*   **Campaign:** "Don't wait for the crash."
*   **Content:** Blog posts showing Driftlock detecting famous crashes (Terra, FTX, Delta) *before* they made headlines.
*   **Partnerships:** Gremlin (Chaos Engineering), PagerDuty (Incident Response).

### Open Source Community
*   **Strategy:** Keep the Core (`cbad-core`) open and high-performance (Rust).
*   **Monetization:** Charge for the *Auditor* (Gemini 3 Pro context) and *Compliance* (Reports).

---

## 3. Timeline

*   **Q4 2025:** Launch CLI and MCP Server. Publish "Chaos Report."
*   **Q1 2026:** VS Code Extension. First "Enterprise" pilot (FinTech).
*   **Q2 2026:** AWS Marketplace listing. "Driftlock for AI" beta.

# THE CHAOS REPORT — November 2025

## Executive Summary

Driftlock’s compression-based anomaly detector behaves like a universal “chaos radar.” Across finance, infrastructure, cloud traffic, and AI safety datasets, the same deterministic zstd + entropy window catches instability before human operators notice. This report packages those findings for investors, auditors, and partners evaluating the Horizon Showcase vs. the production-grade tooling we are now shipping (Firebase SaaS, CLI streaming, VS Code extension, MCP server).

## Methodology

- All experiments use the open-source `driftlock-cli` pipeline powered by `pkg/entropywindow` (baseline window = 400 lines, hop = 25 lines).
- Compression algorithm: **zstd** only (OpenZL optional path is deliberately disabled for parity across client environments).
- Each dataset is replayed in chronological order. We warm up on 400 records, then compute CBAD metrics on every new line while maintaining deterministic seeds.
- Logs referenced below (`*_results.log`) ship with the repo. Run `./scripts/chaos-report.py` to regenerate the summary table in JSON or Markdown.

## Benchmark Snapshot

| Dataset | Anomalies | Total | % Drift | Baseline median (ms) | p95 (ms) | Runtime (s) |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| airline | 31 | 2000 | 1.55% | 1278 | 29292 | 4.748 |
| nasa | 0 | 192 | 0.00% | 9048 | 9056 | — |
| network | 44 | 2000 | 2.20% | 1184 | 25819 | 3.657 |
| safety | 44 | 1581 | 2.78% | 24 | 38 | 5.101 |
| supply | 64 | 2000 | 3.20% | 748 | 997 | 4.059 |
| terra | 9 | 1076 | 0.84% | 44 | 67 | 1.386 |
| web | 29 | 2000 | 1.45% | 112 | 126 | 3.814 |

## Horizon 1 — FinTech & RegTech

- **Dataset:** `credit_card_fraud` + `terra` crash traces (stablecoin depeg).
- **Signal:** Micro-burst entropy spikes surface long-tail fraud rows at ~1.5% without feature engineering. Terra’s death spiral shows a uniform NCD jump (0.84% drift) hours before price feeds update.
- **Value:** Compliance teams receive deterministic anomaly IDs with compression ratios, which now pipe directly into Firebase Functions for instant audit logs.

## Horizon 2 — Critical Infrastructure

- **Dataset:** `airline`, `supply_chain`, `nasa_turbofan` telemetry.
- **Signal:** The airline meltdown (Delta) shows baseline latency median 1.27s vs. anomaly windows pushing p95 >29s. Supply chain telemetry hits 3.2% drift as soon as warehouse scanners fall out of sync.
- **Value:** Operators can run the new VS Code Live Radar on raw log streams before they ever hit SIEM storage—no JSON normalization required.

## Horizon 3 — Cloud & Cyber

- **Dataset:** `network_intrusion`, `web_traffic`.
- **Signal:** Low-and-slow attacks manifest as compression deltas >2× vs. normal CDN requests while staying under rate-limit thresholds.
- **Value:** MCP integration means any Claude/Cursor agent can `tools/call detect_anomalies` on live traces and receive instant “entropy verdicts” without waiting for SaaS round trips.

## Horizon 4 — AI Safety & Observability

- **Dataset:** `prompt_injection` + `safety_results.log`.
- **Signal:** Jailbreak strings (base64 blobs, homoglyphs) explode entropy despite similar token counts—2.78% drift and entropy deltas >0.4.
- **Value:** The Live Radar extension highlights malicious prompts directly inside the IDE/terminal, upselling AI teams on the Pro plan (Gemini explanations optional, not required).

## Reproducing the Report

1. `./scripts/chaos-report.py > docs/launch/CHAOS_TABLE.md` (or pass `--format json`).
2. `driftlock scan --format raw path/to/log --output pretty` to inspect anomalies interactively.
3. `cmd/driftlock-mcp` for agent integrations; `extensions/vscode-driftlock` for IDE diagnostics; `scripts/deploy-landing.sh` to ship the latest Horizon Showcase.

## Next Steps

- Wire the Firebase Functions onboarding flow so that every new tenant automatically gets a “Chaos Report” PDF generated via exporters.
- Expand the MCP toolset with `summarize_anomaly` once Gemini budgets are approved (kept behind Pro/Enterprise gating).
- Track conversion of VS Code extension users → paid plans by embedding anonymized telemetry (count-only, no payloads) in the extension.

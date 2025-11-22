# Driftlock MCP Server Setup

## Overview

`cmd/driftlock-mcp` exposes the Driftlock anomaly detection pipeline to MCP-compatible agents (Claude Desktop, Cursor, Windsurf, Smithery, etc.). The server now accepts arbitrary strings, performs local entropy analysis when possible, and falls back to the hosted `/v1/detect` API for structured JSON payloads.

## Prerequisites

- Go 1.22+ to build the binary
- Driftlock CLI repository checked out locally
- Optional: `DRIFTLOCK_API_KEY` environment variable for authenticated remote detections

## Build & Run

```bash
cd cmd/driftlock-mcp
go build -o driftlock-mcp
./driftlock-mcp
```

The process reads JSON-RPC 2.0 requests from STDIN/STDOUT per the MCP spec. Register the executable in your agent configuration and point it at the compiled binary.

## Tool Contracts

### `detect_anomalies`

Input object fields:

| Field | Type | Description |
| --- | --- | --- |
| `data` | string (required) | Raw logs, NDJSON, JSON arrays, or any free-form text |
| `mode` | `auto` (default), `raw`, `json` | `auto` inspects the payload; `raw` forces local entropy analysis; `json` proxies to the hosted `/v1/detect` API |
| `format` | `json` (default) or `ndjson` | Content-Type hint when `mode=json` |
| `windowLines` | integer (default 400) | Sliding baseline size used by the local analyzer |
| `threshold` | float (default 0.35) | Score threshold (0–1) for local anomalies |

Behavior:

- **Local mode** (`mode=raw` or inferred): Streams each line through the new `pkg/entropywindow` analyzer (zstd compression + Shannon entropy) and returns an object containing `total_records`, `anomaly_count`, and the serialized anomaly metrics.
- **Remote mode** (`mode=json` or inferred from leading `{`/`[`): Sends the payload to `https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/detect?algo=zstd`. If `DRIFTLOCK_API_KEY` is set, it is forwarded as `X-Api-Key`.

Sample MCP call body:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "detect_anomalies",
    "arguments": {
      "data": "ERROR payment timeout id=8821\nERROR payment timeout id=8822",
      "mode": "raw",
      "windowLines": 200,
      "threshold": 0.4
    }
  }
}
```

### `check_system_health`

Returns the contents of `/healthz` so agents can confirm availability.

## Troubleshooting

- **Binary not found:** Ensure the compiled executable is present on the agent's PATH or register the absolute path.
- **No anomalies returned:** Lower `threshold` or reduce `windowLines` for smaller samples.
- **API errors:** Set `DRIFTLOCK_API_KEY` with a valid key or switch to `mode=raw` for fully local analysis.

## Integration Tips

- Pair `mode=raw` with the VS Code extension’s `driftlock scan` output for zero-latency local checks.
- Use `mode=json` when you need the SaaS explanations or CBAD statistics from `/v1/detect`.
- Cache the MCP response per prompt/tool call to avoid re-analyzing identical payloads.

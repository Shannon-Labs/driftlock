# Driftlock SaaS Ecosystem

Driftlock is more than just an API. It is a complete observability ecosystem designed for the 2025 developer workflow. This document outlines the available tools and how they interact.

## 1. The API (Core)
**Endpoint:** `https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1`
- **Role:** The brain. Handles detection, billing, auth, and storage.
- **Docs:** [API Reference](../architecture/API.md)

## 2. The CLI (`driftlock`)
**Source:** `cmd/driftlock-cli/`
- **Role:** The developer's daily driver.
- **Installation:** `go install github.com/driftlock/driftlock/cmd/driftlock-cli@latest`
- **Key Commands:**
  - `driftlock login`: Authenticates machine with SaaS platform.
  - `driftlock detect logs.json`: Streams local logs for instant analysis.
  - `driftlock whoami`: Checks plan status and usage quotas.

## 3. The MCP Server (`driftlock-mcp`)
**Source:** `cmd/driftlock-mcp/`
- **Role:** The bridge to AI Agents (Claude Desktop, etc.).
- **Protocol:** Model Context Protocol (MCP) over Stdio.
- **Capabilities:**
  - Allows LLMs to "see" anomalies in provided data.
  - Enables agents to query system health.
- **Usage:**
  Add to your `claude_desktop_config.json`:
  ```json
  {
    "mcpServers": {
      "driftlock": {
        "command": "driftlock-mcp",
        "env": {
          "DRIFTLOCK_API_KEY": "dlk_..."
        }
      }
    }
  }
  ```

## 4. The Dashboard (Web)
**URL:** `https://driftlock.web.app`
- **Role:** Management and Visuals.
- **Features:**
  - API Key Management.
  - Billing & Plan Upgrades.
  - Interactive Playground.
  - Historical Anomaly Review.

## 5. SDKs
- **Python:** `sdks/python` - Data Science & Backend integration.
- **TypeScript:** `sdks/typescript` - Frontend & Node.js integration.

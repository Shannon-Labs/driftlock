# Driftlock: The Safety Layer for Autonomous Agents

**Real-time Streaming Telemetry + Explainable Anomaly Detection**

Driftlock provides the mathematical guardrails for autonomous AI agents. When your agents operate at speed, Driftlock monitors their behavior (tool calls, thought traces, outputs) and locks in "normal" patterns. If an agent drifts into hallucinations or rogue behavior, Driftlock flags it instantly with a mathematical explanation.

## The Innovation: Two APIs Working Together

1.  **Streaming API**: Ingest agent telemetry in real-time via WebSocket/SSE.
2.  **Detection API**: Analyze every event against a baseline using compression math (NCD).

**Why Math?**
Black-box ML models can hallucinate too. Driftlock uses **Normalized Compression Distance (NCD)** to provide a verifiable, deterministic proof for every anomaly. "Show your work" is built-in.

## Deployment

Driftlock is designed as a modern SaaS platform that you can self-host or use via our managed service.

### Architecture
- **Frontend**: Vue 3 app deployed to **Firebase Hosting**.
- **Backend**: Go API service deployed to **Google Cloud Run**.
- **Streaming**: Server-Sent Events (SSE) for real-time alerts.

### Quick Start

1.  **Deploy**:
    ```bash
    ./deploy.sh
    ```

2.  **Connect your Agent**:
    ```bash
    # Send telemetry to the API
    curl -X POST https://your-api-url/v1/detect -d @event.json
    ```

3.  **Watch Live**:
    Go to your dashboard to see real-time anomaly streams.

## Documentation

- [Cloud Run Setup](docs/deployment/cloud-run-setup.md)
- [Firebase Hosting Setup](docs/deployment/firebase-hosting-setup.md)
- [Streaming API Guide](docs/STREAMING.md)
- [Architecture Overview](docs/ARCHITECTURE.md)

## License

Driftlock Core (Rust) is Apache 2.0.
API Service & Dashboard are source-available (see LICENSE-COMMERCIAL).

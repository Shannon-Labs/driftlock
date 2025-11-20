# Driftlock TypeScript SDK (beta)

Minimal, dependency-light TypeScript client for Driftlock’s `/v1` API. Designed for Node 18+ where `fetch` and `AbortController` are built in.

## Install

```bash
cd sdks/typescript
npm install
npm run build
# then from the repo root
npm install ./sdks/typescript
```

## Usage

```ts
import { DriftlockClient } from "@driftlock/sdk";

const client = new DriftlockClient({
  apiKey: process.env.DRIFTLOCK_API_KEY!,
  baseUrl: process.env.DRIFTLOCK_BASE_URL || "http://localhost:8080"
});

// Health probe
const health = await client.health();
console.log(health.license?.status, health.database);

// Synchronous detect (payload mirrors docs/API.md)
const detect = await client.detect({
  stream_id: "integration-stream",
  events: [
    { timestamp: new Date().toISOString(), type: "log", body: { message: "demo" } }
  ]
});

console.log("anomalies", detect.anomalies?.length ?? 0);

// Fetch anomaly detail
if (detect.anomalies?.[0]) {
  const anomaly = await client.getAnomaly(detect.anomalies[0].id);
  console.log(anomaly.metrics);
}
```

## API surface

- `health()` → `/healthz`
- `detect(payload)` → `POST /v1/detect`
- `listAnomalies(params?)` → `GET /v1/anomalies`
- `getAnomaly(id)` → `GET /v1/anomalies/{id}`

All methods throw a `DriftlockError` when the HTTP status is non-2xx or when the response is not JSON.

## Notes

- Timeouts default to 10s; override via `timeoutMs`.
- Bring your own fetch implementation via `fetchImpl` if you’re in non-Node environments.
- The SDK is intentionally thin—“math-first” responses stay intact for auditability. AI commentary stays server-side and opt-in. 

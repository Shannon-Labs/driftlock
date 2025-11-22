# API Reference

**Version:** 1.0.0 (Shannon)  
**Base URL:** `https://api.driftlock.net/v1`

The Driftlock API provides direct access to the Universal Radar. It is designed for high-throughput ingestion and synchronous anomaly detection.

## Authentication

Authentication is handled via Bearer tokens (API Keys).

```bash
Authorization: Bearer <your_api_key>
```

To obtain an API key, initialize a pilot on the [dashboard](https://driftlock.net/dashboard).

## Endpoints

### `POST /detect`

Ingest events for anomaly detection. This endpoint accepts a batch of JSON events and returns a forensic verdict.

**Rate Limits:**
*   **Pulse (Free):** 10 req/sec
*   **Radar ($20):** 100 req/sec
*   **Lock ($200):** 1,000 req/sec
*   **Orbit (Custom):** Unlimited

#### Request

```json
{
  "window_size": 50,
  "events": [
    {
      "id": "evt_123",
      "timestamp": "2025-11-21T10:00:00Z",
      "data": { ... }
    }
  ]
}
```

#### Response (200 OK)

```json
{
  "verdict": "ANOMALY_DETECTED",
  "confidence": 0.98,
  "compression_ratio": 2.4,
  "entropy_variance": 3.1,
  "explanation": "Transaction velocity suggests automated skimming attack."
}
```

### `GET /health`

**Status:** Public  
**Response:** `200 OK` - System nominal.

## Error Codes

*   `401 Unauthorized`: Missing or invalid signal key.
*   `402 Payment Required`: Tier limits exceeded. Upgrade to Radar or Lock.
*   `429 Too Many Requests`: Rate limit exceeded.
*   `422 Unprocessable Entity`: Malformed event JSON.

---

*© 2025 Shannon Labs. Entropy does not lie.*

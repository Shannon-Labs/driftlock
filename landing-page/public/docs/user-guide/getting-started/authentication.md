# Authentication

Driftlock uses API keys for programmatic access and Firebase Auth for the dashboard. Use API keys for all `/v1/*` calls; Firebase Auth is only needed in the UI.

## API keys (primary)

1. Sign up at [https://driftlock.net/#signup](https://driftlock.net/#signup).
2. Open **Dashboard → API Keys → Create key**.
3. Choose a role and copy the key (it is shown once). Store it as a secret.

**Roles**

| Role | Best for | Access |
| --- | --- | --- |
| `admin` | Admin automations and internal tools | All detection + management endpoints |
| `stream` | Application ingestion/detection | `/v1/detect`, `/v1/anomalies`, stream-scoped calls |

> Tip: Bind stream-role keys to a specific stream when offered for least-privilege.

### Headers and base URL

```bash
export DRIFTLOCK_API_URL="https://api.driftlock.net/v1"
export DRIFTLOCK_API_KEY="YOUR_API_KEY"
```

Every request should include:

```
Content-Type: application/json
X-Api-Key: YOUR_API_KEY
```

### cURL example

```bash
curl -X POST "$DRIFTLOCK_API_URL/detect" \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -d '{"stream_id":"default","events":[{"body":{"cpu":45,"memory":2048}}]}'
```

### Python example

```python
import os
import requests

api_url = os.getenv("DRIFTLOCK_API_URL", "https://api.driftlock.net/v1")
api_key = os.environ["DRIFTLOCK_API_KEY"]

resp = requests.post(
    f"{api_url}/detect",
    headers={"X-Api-Key": api_key, "Content-Type": "application/json"},
    json={"events": [{"body": {"latency_ms": 950}}]},
    timeout=10,
)
resp.raise_for_status()
print(resp.json())
```

## Rate limits

| Plan | Requests/min | Events per request | Notes |
| --- | --- | --- | --- |
| Pilot (free) | 60 | 256 | Great for testing and CI |
| Radar ($20/mo) | 300 | 256 | Higher throughput, production-ready |
| Lock ($200/mo) | 1000+ | 256 | Enterprise + compliance evidence |

Rate-limit headers: `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`, and `Retry-After` (when 429).

## Firebase Auth (dashboard only)

Use Firebase login to manage your account, billing, and API keys. It is not required for API requests unless an endpoint explicitly requests a bearer token (e.g., user profile pages in the UI).

## Troubleshooting

| Symptom | Likely cause | Fix |
| --- | --- | --- |
| `401 unauthorized` | Missing/invalid `X-Api-Key` | Verify header casing, ensure key is active, no leading/trailing spaces. |
| `403 forbidden` | Wrong role or stream binding | Use an `admin` key for admin calls; verify stream binding. |
| `429 rate_limit_exceeded` | Plan limits hit | Batch events, honor `Retry-After`, or upgrade plan. |
| `invalid_argument` | Payload/schema issue | Ensure `events` array exists (1–256) and payload <10 MB. |

## Next steps

- Run the [Quickstart](./quickstart.md)
- Call [POST /v1/detect](../api/endpoints/detect.md)
- See the full [REST API reference](../api/rest-api.md)
- Try the public [demo endpoint](../api/endpoints/demo.md) (no auth)

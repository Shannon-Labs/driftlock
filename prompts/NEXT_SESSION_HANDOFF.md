# Driftlock End-to-End Testing Handoff

## Session Context
Previous session (Dec 7, 2025) fixed massive storage waste in Docker/Cloud Build:
- **Before:** ~3.8GB uploaded per deploy
- **After:** ~91MB uploaded per deploy (97% reduction)
- Fixed `.dockerignore` and `.gcloudignore` pattern bugs

**Build Status:** Cloud Build `9f100c38-92fe-4673-8695-dad35da64584` should be deployed by now.

---

## Primary Task: Comprehensive API Testing

We need to verify the entire Driftlock service works end-to-end before launch. Test with real data to ensure customers can get an API key and start detecting anomalies.

### 1. Verify Deployment

```bash
# Check service is running
curl https://driftlock-api-o6kjgrsowq-uc.a.run.app/readyz

# Should return: {"ready":true,...}
```

### 2. Test with Kaggle Datasets

The `test-data/` directory contains real datasets for testing:

| Dataset | Path | Description |
|---------|------|-------------|
| Terra Luna Crash | `test-data/terra_luna/` | Crypto crash data with known anomalies |
| NASA Turbofan | `test-data/nasa_turbofan/` | Sensor degradation data |
| Financial Fraud | `test-data/fraud/` | Transaction anomaly data |
| Airline Delays | `test-data/airline/` | Time series with outliers |
| Network Intrusion | `test-data/network/` | Security event data |
| AWS CloudWatch | `test-data/web_traffic/realAWSCloudwatch/` | Real metrics with anomalies |
| Pre-formatted | `test-data/financial-demo.json` | Ready-to-use JSON events |

**Test approach:**
1. Load JSON/CSV datasets
2. Format as `{"events": [{"field": "value"}, ...]}`
3. POST to `/v1/demo/detect` (no auth) or `/v1/events` (with API key)
4. Verify anomalies are detected correctly

### 3. Test Streaming Data (CoinGecko API)

Test real-time anomaly detection with live market data:

```bash
# Get Bitcoin price history
curl -s "https://api.coingecko.com/api/v3/coins/bitcoin/market_chart?vs_currency=usd&days=7&interval=hourly"

# Transform to events format and send to Driftlock
# Look for price spikes, volume anomalies
```

**Why CoinGecko:** Free API, no auth needed, real-time data with natural volatility patterns - perfect for testing anomaly detection.

### 4. API Key Testing

```bash
# Internal API key (from GCP secrets):
gcloud secrets versions access latest --secret=internal-api-key --project=driftlock

# Test authenticated endpoint:
curl -X POST https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/events \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"events": [...]}'
```

### 5. Full User Flow Test

Simulate a new customer:
1. Sign up flow (if implemented)
2. Get API key
3. Send first batch of events
4. Retrieve anomalies
5. Verify explanations work

---

## Key Endpoints

| Endpoint | Auth | Purpose |
|----------|------|---------|
| `GET /readyz` | None | Health check |
| `POST /v1/demo/detect` | None | Demo anomaly detection (rate limited) |
| `POST /v1/events` | API Key | Production event ingestion |
| `GET /v1/anomalies` | API Key | Retrieve detected anomalies |

---

## Expected Demo Request Format

```json
{
  "events": [
    {"message": "Normal transaction", "amount": 50, "timestamp": "2024-01-01T10:00:00Z"},
    {"message": "Normal transaction", "amount": 45, "timestamp": "2024-01-01T10:01:00Z"},
    {"message": "SUSPICIOUS: Large transfer", "amount": 50000, "timestamp": "2024-01-01T10:02:00Z"},
    ...
  ],
  "config_override": {
    "compression_algorithm": "zstd"
  }
}
```

---

## Success Criteria

1. **Demo endpoint works** - Returns anomalies for test data
2. **Kaggle datasets detect anomalies** - Known anomalies in data are flagged
3. **Streaming test works** - CoinGecko volatility detected
4. **API key auth works** - Authenticated requests succeed
5. **Error handling** - Bad requests return helpful errors
6. **AI explanations** - Anomalies include human-readable explanations

---

## Files to Review

- `collector-processor/cmd/driftlock-http/demo.go` - Demo endpoint logic
- `collector-processor/cmd/driftlock-http/billing.go` - API key validation
- `test-data/README.md` - Dataset documentation
- `docs/user-guide/api/` - API documentation

---

## Notes

- Demo endpoint limited to 200 events, 10 req/min
- Production requires valid API key and subscription
- AI explanations use mock provider in current deployment (see `AI_PROVIDER=mock` in cloudbuild.yaml)

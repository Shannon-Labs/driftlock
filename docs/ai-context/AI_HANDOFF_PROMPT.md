# 🚀 HANDOFF PROMPT FOR NEXT CLAUDE SESSION

**Date:** 2025-12-02
**Status:** ~90% Launch Ready - Need to test AI explainability with real data

---

## COMPLETED THIS SESSION

### Verified Working ✅
- **Core Detection API** - CBAD anomaly detection working in production
  - CGO/Rust integration functional (`compression_algo: "zstd"`)
  - 11 anomalies detected in test with config override
  - Rate limiting working (10 req/min demo)
  - Test: `curl -X POST https://driftlock.net/api/v1/demo/detect ...`

### Code Changes (Uncommitted) ✅
1. **Pricing fix:** `landing-page/src/views/HomeView.vue:176` - Tensor $200 → $100
2. **Env example:** `.env.example:33-38` - Updated Stripe price IDs with correct amounts
3. **Docs:** `CLAUDE.md` - Updated pricing table (verified from Stripe)
4. **AI Config:** `collector-processor/internal/ai/config.go` - Updated AI tiers:
   - Trial/Pilot: Z.AI GLM-4 (cheap, no Claude budget tracking)
   - Radar: Claude Haiku 4.5 only ($20/mo budget)
   - Tensor/Lock: Claude Opus 4.5 ($150/mo budget)
   - Orbit: Claude Opus 4.5 (unlimited)
5. **Stale TODOs cleaned:** billing.go, dashboard.go, usage.go

### Linear Issues Created ✅
- SHA-41: Verify AI explainability with paid account
- SHA-42: Deploy pricing fix: Tensor $200 → $100
- SHA-43: Core detection API verified (Done)

---

## PRIORITY 1: Test AI Explainability with Real Market Data

### The Goal
Test that the Z.AI API integration works for AI-powered anomaly explanations. The infrastructure is built, we just need to verify it end-to-end.

### What's Configured
```yaml
Production (Cloud Run):
  AI_PROVIDER: zai
  AI_BASE_URL: https://api.z.ai/api/coding/paas/v4
  AI_MODEL: GLM-4.6
  AI_API_KEY: (secret exists in Secret Manager, updated today)
```

### Test Approach - Use yfinance for Real Data
The US market is open! Pull SPY/QQQ intraday data:

```python
import yfinance as yf
import json

# Get intraday SPY data (last 50 minutes)
spy = yf.download("SPY", period="1d", interval="1m")
events = []
for idx, row in spy.tail(50).iterrows():
    events.append({
        "ts": idx.isoformat(),
        "symbol": "SPY",
        "close": float(row['Close']),
        "volume": int(row['Volume'])
    })

# Format for detection API
payload = {
    "events": events,
    "config_override": {
        "baseline_size": 30,
        "window_size": 10,
        "hop_size": 5
    }
}
print(json.dumps(payload))
```

### Key Files for AI Integration
- `collector-processor/internal/ai/client.go` - AI client factory
- `collector-processor/internal/ai/openai_client.go` - Z.AI compatible client (OpenAI API format)
- `collector-processor/internal/ai/smart_router.go` - Model routing logic
- `collector-processor/internal/ai/config.go` - Per-plan AI settings (just updated!)

### What Success Looks Like
- Detection response includes richer `why` field with AI-generated explanation
- Not just template text like "NCD=0.45, compression ratio changed..."
- Should have contextual analysis from GLM-4

---

## PRIORITY 2: Commit & Deploy

### Uncommitted Changes
```bash
git diff --stat
# Modified:
#   .env.example
#   CLAUDE.md
#   collector-processor/cmd/driftlock-http/billing.go (TODO cleanup)
#   collector-processor/cmd/driftlock-http/dashboard.go (TODO cleanup)
#   collector-processor/cmd/driftlock-http/usage.go (TODO cleanup)
#   collector-processor/internal/ai/config.go (AI tier config)
#   landing-page/src/views/HomeView.vue (pricing fix)
```

### Deploy Steps
1. Commit changes with descriptive message
2. Push to trigger Cloud Build (cloudbuild.yaml configured)
3. Verify Firebase Hosting updates (landing page pricing)
4. Test detection API still works post-deploy

---

## KNOWN STATE

### Stripe Products (Verified Today)
| Tier | Price ID | Price |
|------|----------|-------|
| Radar | `price_1SZkbAL4rhSbUSqA8rWnA0eW` | $15/mo |
| Tensor | `price_1SZjnpL4rhSbUSqALhIjpoR3` | $100/mo |
| Orbit | `price_1SZjnqL4rhSbUSqAr65IPfB1` | $499/mo |

### AI Tier Config (Updated This Session)
| Plan | Model | Budget | Calls/Day |
|------|-------|--------|-----------|
| Trial/Pilot | Z.AI GLM-4 | N/A (cheap) | 100 |
| Radar | Haiku 4.5 | $20/mo | 200 |
| Tensor/Lock | Opus 4.5 | $150/mo | 500 |
| Orbit | Opus 4.5 | Unlimited | Unlimited |

### Key URLs
- Site: https://driftlock.net
- Demo API: https://driftlock.net/api/v1/demo/detect
- Cloud Run: https://driftlock-api-o6kjgrsowq-uc.a.run.app
- Linear: https://linear.app/shannon-labs/project/driftlock-a8c80503816c

---

## QUICK START COMMANDS

```bash
# Check what's changed
git diff --stat

# Test production detection (no auth)
curl -s -X POST "https://driftlock.net/api/v1/demo/detect" \
  -H "Content-Type: application/json" \
  -d '{"events":[{"ts":"2025-01-01T00:00:00Z","v":100},{"ts":"2025-01-01T00:00:01Z","v":9999}]}'

# Get SPY data for realistic test
python3 -c "import yfinance as yf; print(yf.download('SPY', period='1d', interval='1m').tail(10))"

# Run local API (if needed)
cd collector-processor && go run ./cmd/driftlock-http
```

---

## SESSION SUMMARY

| Task | Status |
|------|--------|
| ✅ Core detection API verified | Working in production |
| ✅ Website pricing fixed | Tensor $200 → $100 |
| ✅ .env.example updated | Correct price IDs |
| ✅ CLAUDE.md pricing updated | Verified 2025-12-02 |
| ✅ AI config updated | Trial=Z.AI, Paid=Claude |
| ✅ Linear issues created | SHA-41, 42, 43 |
| ✅ Stale TODOs cleaned | billing, dashboard, usage |
| ⏳ Commit & deploy | Ready to go |
| ⏳ Test AI explainability | Priority 1 for next session |

---

**Previous session by:** Claude Opus 4.5
**Next priority:** Test Z.AI integration with real SPY/QQQ market data

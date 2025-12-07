# CBAD Full Integration - AI Handoff Prompt

> Use this prompt to have the next AI session properly analyze and fix the CBAD detection system.

---

## PROMPT FOR NEXT AI SESSION

```
I need you to do a deep analysis and fix the CBAD (Compression-Based Anomaly Detection) system in this codebase. The previous session made progress but did NOT actually test the real system end-to-end.

## CRITICAL ISSUES TO INVESTIGATE

### 1. Baseline Management - Is it working?

The system claims to use Redis for baseline persistence, but this was never tested.

Files to analyze:
- `collector-processor/cmd/driftlock-http/db.go` - Look for `loadBaseline`, `saveBaseline` functions
- `collector-processor/cmd/driftlock-http/main.go` - Look at `detectHandler` around line 896+

Questions to answer:
- Is the baseline actually being saved to Redis after each detection?
- Is it being loaded on subsequent requests?
- What happens with a fresh stream (no baseline exists)?
- How many events are needed before detection works?

### 2. Demo Endpoint - Does it actually detect anomalies?

The demo endpoint at `POST /v1/demo/detect` was created but:
- Does NOT include AI explanations
- Uses smaller baseline (40 events) but still may not work properly
- Was never tested with real anomalous data

Files to analyze:
- `collector-processor/cmd/driftlock-http/demo.go`
- Check if `demoBaselineSize = 40`, `demoWindowSize = 10` actually allows detection

Test to run:
```bash
# Start the server
cd collector-processor && go run ./cmd/driftlock-http

# In another terminal, send a detection request with mixed normal + anomalous events
curl -X POST http://localhost:8080/api/v1/demo/detect \
  -H "Content-Type: application/json" \
  -d '{
    "events": [
      {"level":"INFO","msg":"normal log 1"},
      {"level":"INFO","msg":"normal log 2"},
      ... (40+ normal events)
      {"level":"CRITICAL","msg":"SYSTEM FAILURE catastrophic error"}
    ]
  }'
```

### 3. CBAD Library - Is it actually computing real metrics?

The Rust FFI library exists at `cbad-core/target/release/libcbad_core.a`

Files to analyze:
- `collector-processor/driftlockcbad/cbad.go` - Real CGO implementation
- `collector-processor/driftlockcbad/cbad_stub.go` - Stub when CGO disabled

Questions:
- When running the server, is CGO enabled? (build tag: `cgo && !driftlock_no_cbad`)
- Are real NCD/p-value metrics being computed or is it using stubs?
- What do the actual metrics look like for normal vs anomalous data?

### 4. AI Explanations - Where are they integrated?

The Ollama client exists at `collector-processor/internal/ai/ollama_client.go`

But:
- The demo endpoint does NOT call the AI client
- Only the authenticated `/v1/detect` endpoint has AI integration (see main.go:1168)

Questions:
- Should demo endpoint also have AI explanations?
- What env vars are needed? (`AI_PROVIDER=ollama`, `OLLAMA_MODEL=ministral-3:3b`)

### 5. The Fundamental Question: What IS the baseline?

In CBAD, the baseline is a corpus of "normal" events that new events are compared against.

Current approach (I think):
- `baseline_size = 400` (production) or `40` (demo)
- Events accumulate until baseline is full
- Then sliding window comparison begins

Problems I see:
- If baseline isn't persisted, every API call starts fresh
- If 400 events needed, most demos will fail
- How does the system know what's "normal" vs "anomalous" in the baseline?

## DELIVERABLES

1. **Analysis Report**: What's actually working vs broken
2. **Fix the issues**: Make baseline persistence work
3. **Create a REAL demo**:
   - Actually run the driftlock API
   - Actually use the CBAD library (not stubs)
   - Actually get real NCD/p-value scores
   - Actually call Ollama for AI explanations
4. **Record it**: Use asciinema to record a real end-to-end demo
5. **Update docs**: Document how the baseline system actually works

## HOW TO RUN THE REAL SYSTEM

```bash
# Terminal 1: Start Postgres
docker compose up -d postgres

# Terminal 2: Start Redis (if baseline persistence is implemented)
docker compose up -d redis

# Terminal 3: Start the API server with CGO enabled
cd collector-processor
CGO_ENABLED=1 go run -tags "cgo" ./cmd/driftlock-http

# Terminal 4: Start Ollama (if not running)
ollama serve

# Terminal 5: Test detection
curl -X POST http://localhost:8080/api/v1/demo/detect ...
```

## ENVIRONMENT VARIABLES NEEDED

```bash
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/driftlock
REDIS_URL=redis://localhost:6379
AI_PROVIDER=ollama
OLLAMA_MODEL=ministral-3:3b
```

## SUCCESS CRITERIA

A truly real demo would show:
1. Server starts with "AI client initialized: ollama"
2. Events sent to demo endpoint
3. Real NCD score computed (not hardcoded)
4. Real p-value from permutation test
5. Real AI explanation from Ministral 3B
6. All of this in a single recorded terminal session

## WHAT THE PREVIOUS SESSION DID (and didn't do)

DID:
- Created Ollama client (`ollama_client.go`)
- Made real Ollama API calls for AI explanations
- Recorded a demo with asciinema

DID NOT:
- Run the actual driftlock API server
- Use the actual CBAD library for metrics
- Test baseline persistence
- Test end-to-end flow
- The "metrics" shown (NCD: 0.847, P-Value: 0.003) were HARDCODED in the demo script

Be thorough. The goal is a REAL, working demo we can put on the website.
```

---

## FILES TO READ FIRST

1. `collector-processor/cmd/driftlock-http/main.go` - Main server, detectHandler
2. `collector-processor/cmd/driftlock-http/demo.go` - Demo endpoint
3. `collector-processor/cmd/driftlock-http/db.go` - Baseline persistence (if exists)
4. `collector-processor/driftlockcbad/cbad.go` - Real CBAD implementation
5. `collector-processor/internal/ai/client.go` - AI provider factory
6. `collector-processor/internal/ai/ollama_client.go` - Ollama integration
7. `docs/ALGORITHMS.md` - How CBAD is supposed to work mathematically

## PLAN FILE LOCATION

There's a plan at `/Users/hunterbown/.claude/plans/buzzing-hugging-wirth.md` that documents what was supposed to be fixed but may not have been completed.

---

*Created: 2025-12-05*
*Previous session: Made AI integration work but didn't test real CBAD flow*

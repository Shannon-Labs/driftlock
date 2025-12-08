# Driftlock Zero-Config UX Verification

**Agent Type**: Claude with full tool access
**Focus**: End-to-end verification of the zero-config UX overhaul
**Goal**: Validate all changes work correctly with real data, update Linear with results

---

## Context

A major UX overhaul was just committed (`48c564e`) implementing "zero-config magic":

1. **Backend defaults changed**: New streams now have `auto_tune_enabled=true` and `adaptive_window_enabled=true`
2. **Rust recommendations**: Rust computes `recommended_ncd_threshold` based on data stability
3. **Auto-apply**: Go automatically applies Rust recommendations after detection
4. **Frontend**: Feedback buttons, sensitivity picker, auto-tune indicator added

The system should now "just work" — users send data, thresholds self-optimize.

---

## Verification Tasks

### Task 1: Verify Backend Defaults

**Test**: Create a new stream and verify defaults are correct.

```bash
# Create a test tenant/stream via the API or check db.go directly
curl -X POST https://driftlock.net/api/v1/demo/detect \
  -H "Content-Type: application/json" \
  -d '{"events": [{"body": {"test": "event"}}]}'
```

**Expected**: Response should show stream was created with adaptive features enabled.

**Verify in code**:
- `collector-processor/cmd/driftlock-http/db.go` lines 293-294 and 392-393
- Both `AutoTuneEnabled: true` and `AdaptiveWindowEnabled: true`

---

### Task 2: Verify Rust Recommendations

**Test**: Run detection and check that recommendations are in the response.

```bash
# Use demo endpoint with enough events to trigger detection
curl -X POST https://driftlock.net/api/v1/demo/detect \
  -H "Content-Type: application/json" \
  -d @test-data/sample_events.json
```

**Expected in response**:
```json
{
  "anomalies": [{
    "metrics": {
      "recommended_ncd_threshold": 0.28,  // Should be non-zero
      "recommended_window_size": 55,       // Should be non-zero
      "data_stability_score": 0.72         // Should be 0.0-1.0
    }
  }]
}
```

**Verify in code**:
- `cbad-core/src/metrics/mod.rs` lines 150-193: `apply_recommendations()` function
- `cbad-core/src/ffi.rs` lines 419-421: FFI export

---

### Task 3: Verify Auto-Apply Logic

**Test**: Confirm `applyRustRecommendation()` is called after detection.

**Check code flow**:
1. `collector-processor/cmd/driftlock-http/main.go` lines 1240-1254: Auto-apply wiring
2. `collector-processor/cmd/driftlock-http/autotune.go` lines 372-437: `applyRustRecommendation()` function

**Verify safeguards**:
- Only adjusts if >2% difference (line 402-404)
- Respects 1hr cooldown (line 390-396)
- Clamps to [0.1, 0.8] bounds (line 408)
- Records in `threshold_tune_history` with `reason: "rust_recommendation"` (line 429)

---

### Task 4: Run Tests

```bash
# Rust tests
cd cbad-core && cargo test

# Go tests
cd collector-processor && go test ./cmd/driftlock-http/...

# Specific auto-tune tests
go test -v -run TestComputeAutoTuneAdjustment ./cmd/driftlock-http/
go test -v -run TestShouldThrottleAutoTune ./cmd/driftlock-http/
```

**Expected**: All tests pass.

---

### Task 5: Verify Frontend Changes

**Check DashboardView.vue**:
- Lines 671, 679: Feedback buttons with `submitFeedback()` calls
- Lines 529-533: "Auto-Tuning Active" indicator
- Lines 536-577: Sensitivity picker (Low/Med/High)
- Lines 1043-1049: `submitFeedback()` function implementation

**Manual verification** (if possible):
1. Load dashboard at `https://app.driftlock.net`
2. Check for "Auto-Tuning Active" indicator with pulsing green dot
3. Check for sensitivity toggle (Low/Med/High)
4. Check anomaly table has thumbs up/down feedback column

---

### Task 6: Documentation Verification

**Check files exist and have collapsible sections**:
- `docs/user-guide/guides/detection-profiles.md`
- `docs/user-guide/guides/auto-tuning.md` (should have `<details>` tag)
- `docs/user-guide/guides/adaptive-windowing.md` (should have `<details>` tag)
- `docs/user-guide/tutorials/profiles-tutorial.md`
- `docs/user-guide/tutorials/feedback-loop.md`
- `docs/user-guide/getting-started/concepts.md` (should have `<details>` tag)

---

### Task 7: End-to-End Test with Real Data

**Use existing test datasets**:
```bash
# Check what test data exists
ls -la test-data/

# Run detection with NASA turbofan data
curl -X POST https://driftlock.net/api/v1/demo/detect \
  -H "Content-Type: application/json" \
  -d @test-data/final_test_results/nasa_payload.json

# Run detection with crypto crash data
curl -X POST https://driftlock.net/api/v1/demo/detect \
  -H "Content-Type: application/json" \
  -d @test-data/final_test_results/terra_payload.json
```

**Verify**:
1. Detection runs without errors
2. Response includes `recommended_ncd_threshold` in metrics
3. Anomalies are detected (if data contains anomalies)

---

## Linear Integration

After verification, update Linear with results:

### If All Passes:
```
Create issue: "QA: Zero-Config UX Verified"
Team: Shipyard
Status: Done
Labels: qa, verified

Body:
Verified commit 48c564e - Zero-config UX overhaul

Checklist:
- [x] Backend defaults (auto_tune=true, adaptive_window=true)
- [x] Rust recommendations computed and exported
- [x] Auto-apply logic with safeguards
- [x] All tests passing (57 Rust, N Go)
- [x] Frontend feedback buttons and sensitivity picker
- [x] Documentation with collapsible sections
- [x] End-to-end test with real data

Ready for production.
```

### If Issues Found:
```
Create issue: "Bug: [Description of issue]"
Team: Shipyard
Status: In Progress
Priority: High
Labels: bug, ux-overhaul

Body:
Found during QA of commit 48c564e

**Issue**: [Describe the problem]
**Expected**: [What should happen]
**Actual**: [What actually happened]
**Steps to reproduce**: [How to reproduce]
**Files involved**: [List relevant files]
```

---

## Key Files Reference

| File | Purpose |
|------|---------|
| `collector-processor/cmd/driftlock-http/db.go` | Stream creation defaults |
| `collector-processor/cmd/driftlock-http/autotune.go` | Auto-tune + Rust apply logic |
| `collector-processor/cmd/driftlock-http/main.go` | Detection handler wiring |
| `cbad-core/src/metrics/mod.rs` | Rust recommendation algorithm |
| `cbad-core/src/ffi.rs` | FFI export of recommendations |
| `landing-page/src/views/DashboardView.vue` | Frontend UI changes |

---

## Success Criteria

1. **All tests pass** (Rust and Go)
2. **Recommendations appear in API response** with non-zero values
3. **Auto-apply respects safeguards** (>2% diff, 1hr cooldown, bounds)
4. **Frontend shows** feedback buttons, sensitivity picker, auto-tune indicator
5. **Docs have** collapsible `<details>` sections for advanced content
6. **Linear updated** with QA results

---

## Notes

- The demo endpoint doesn't require authentication
- For authenticated endpoints, you'll need an API key
- If the server isn't running locally, use the production URL
- Check `CLAUDE.md` for additional context and commands

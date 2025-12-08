# Driftlock Backend: Zero-Config Defaults Overhaul

**Agent Type**: Codex with max reasoning
**Focus**: Backend Go code only
**Goal**: Make the system "just work" out of the box

---

## Problem Statement

Driftlock has sophisticated adaptive detection features that are **OFF by default**. Users must make API calls to enable them. Nobody does. The smart features we built are invisible.

**Current defaults for new streams:**
```go
auto_tune_enabled: false      // User feedback improves detection - OFF
adaptive_window_enabled: false // Smart window sizing - OFF
detection_profile: "balanced"  // This is fine
```

**What we want:**
```go
auto_tune_enabled: true       // Learn from feedback automatically
adaptive_window_enabled: true  // Size windows automatically
detection_profile: "balanced"  // Keep this
```

---

## Primary Task: Change Defaults

### Find and modify stream creation defaults

**Likely locations:**
- `collector-processor/cmd/driftlock-http/db.go` - Database operations
- `collector-processor/cmd/driftlock-http/onboarding.go` - User signup flow
- `collector-processor/cmd/driftlock-http/main.go` - Server initialization

**What to change:**
1. Find where new streams are INSERT'd into the database
2. Change `auto_tune_enabled` default from `false` to `true`
3. Change `adaptive_window_enabled` default from `false` to `true`

**SQL might look like:**
```sql
INSERT INTO streams (tenant_id, stream_id, ..., auto_tune_enabled, adaptive_window_enabled, ...)
VALUES ($1, $2, ..., false, false, ...)  -- CHANGE TO: true, true
```

Or there might be a struct with defaults:
```go
type StreamConfig struct {
    AutoTuneEnabled       bool `default:"false"`  // CHANGE TO: true
    AdaptiveWindowEnabled bool `default:"false"`  // CHANGE TO: true
}
```

---

## Secondary Task: Audit Auto-Tune Readiness

Before enabling auto-tune by default, verify the algorithm is safe:

### Questions to answer:

1. **What happens with zero feedback?**
   - Does auto-tune gracefully handle streams with no feedback?
   - It should do nothing, not crash or produce weird thresholds

2. **Is the cooldown sufficient?**
   - Current: 1 hour between adjustments
   - Is this aggressive enough? Too aggressive?

3. **Are bounds reasonable?**
   - NCD bounds: 0.10 to 0.80
   - P-value bounds: 0.001 to 0.20
   - Can auto-tune ever produce dangerous values?

4. **What's the minimum feedback threshold?**
   - Current: 20 samples before tuning
   - Is this enough to prevent premature tuning?

### Files to audit:
- `collector-processor/cmd/driftlock-http/autotune.go`
- `collector-processor/cmd/driftlock-http/adaptive_windowing.go`

---

## Tertiary Task: Simplification Opportunities

The Codex audit found unused code:

1. **`autotune.go`**: `AvgFPPValue` and `AvgConfirmedPValue` are computed but never used
2. **`adaptive_windowing.go`**: `EventSizeVariance` is computed but never used
3. **`main.go`**: HopSize isn't clamped when using cached adaptive sizes

**Decision needed:** Should we:
- A) Remove the dead code (cleaner, less confusion)
- B) Leave it (might be useful later)
- C) Actually use it (if it improves the algorithm)

Recommend option A unless the unused fields serve a clear future purpose.

---

## Constraints

1. **Don't break existing streams** - Changes should only affect NEW streams
2. **Don't change the API contract** - Endpoints stay the same
3. **Don't change detection behavior** - Only defaults, not algorithms
4. **Backward compatible** - Existing users shouldn't notice anything

---

## Testing Considerations

After changes, verify:

1. **New stream creation** sets `auto_tune_enabled=true` and `adaptive_window_enabled=true`
2. **Existing streams** are NOT modified (their settings persist)
3. **Auto-tune with zero feedback** does nothing (no errors, no changes)
4. **Adaptive windowing** computes reasonable sizes for new streams
5. **Profile switching** still works (can override defaults via API)

---

## File Map

```
collector-processor/cmd/driftlock-http/
├── main.go                 # Server setup, routes
├── db.go                   # Database operations (LIKELY CHANGE HERE)
├── onboarding.go           # Signup flow, stream creation (LIKELY CHANGE HERE)
├── profiles.go             # Profile definitions (reference only)
├── autotune.go             # Auto-tune algorithm (audit for safety)
├── adaptive_windowing.go   # Window sizing (audit for safety)
├── adaptive_handlers.go    # API handlers (reference only)
└── ...
```

---

## Expected Output

1. **Diff showing default changes** - The actual code modification
2. **Audit summary** - Confirmation that auto-tune is safe to enable by default
3. **Dead code recommendation** - Keep or remove unused fields
4. **Test scenarios** - How to verify the changes work

---

## Philosophy

The user should never need to:
- Read documentation to get good defaults
- Make API calls to enable smart features
- Understand NCD thresholds or p-values

The system should be intelligent out of the box. Power users can still override everything via API, but the defaults should be optimal for the 99%.

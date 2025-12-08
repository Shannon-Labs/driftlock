# Driftlock UX Overhaul - Handoff Prompt

## The Vision

**Zero config magic** for 99% of users, **full power-user control** for those who dig deeper.

A new user should be able to:
1. Sign up
2. Send events
3. See anomalies

That's it. No reading docs. No choosing profiles. No understanding NCD thresholds. It just works.

But for the 1% who need it: every knob is available via API, and eventually via an "Advanced Settings" panel.

---

## Current State

### Backend: DONE (and overengineered)
The backend has sophisticated adaptive features:
- **Detection Profiles**: sensitive/balanced/strict/custom with different thresholds
- **Auto-Tuning**: Learns from user feedback, targets 5% false positive rate
- **Adaptive Windowing**: Automatically sizes windows based on stream characteristics
- **Feedback Loop**: API to mark anomalies as false_positive/confirmed/dismissed

All of this works. The API is complete. The algorithms are solid.

### Frontend: BROKEN
The dashboard is missing critical UI:
- **No feedback buttons** - Can't mark anomalies as false positive from UI
- **No profile selector** - Stuck with balanced, no way to change
- **No settings panel** - Can't toggle auto-tune or adaptive windows
- **No tuning history** - Can't see what the algorithm learned

### Defaults: WRONG
Current defaults for new streams:
- `auto_tune_enabled`: **FALSE** (should be TRUE)
- `adaptive_window_enabled`: **FALSE** (should be TRUE)
- `detection_profile`: "balanced" (this is fine)

The system ships with the smart features OFF. Users have to opt-in via API calls they'll never make.

---

## What Needs to Change

### Priority 1: Fix Defaults (Backend)
**File**: `collector-processor/cmd/driftlock-http/db.go` or wherever streams are created

Change new stream defaults:
```go
auto_tune_enabled: true        // was: false
adaptive_window_enabled: true  // was: false
```

This makes the system "learn and adapt" by default. Zero config.

### Priority 2: Add Feedback UI (Frontend)
**File**: `landing-page/src/views/DashboardView.vue`

Add to each anomaly row:
- Thumbs up button → POST `/v1/anomalies/{id}/feedback` with `{"feedback_type": "confirmed"}`
- Thumbs down button → POST `/v1/anomalies/{id}/feedback` with `{"feedback_type": "false_positive"}`

Simple. One click. The backend handles the rest.

### Priority 3: Surface Auto-Tune Status (Frontend)
**File**: `landing-page/src/views/DashboardView.vue`

Add a small status indicator:
- "Learning: 8/20 feedback samples" (not enough data yet)
- "Auto-tuning active" (system is adapting)
- "Last adjusted: 2 hours ago" (show it's working)

Users should feel the system is alive and learning, not static.

### Priority 4: Simple Profile Picker (Frontend)
**File**: New component or add to DashboardView

Add a simple dropdown or toggle:
- "Detection Sensitivity: Low / Medium / High"
- Maps to: strict / balanced / sensitive profiles
- One click to change, immediate effect

Don't expose the technical names. Use human language.

### Priority 5: Hide the Complex Docs
The detailed documentation I created (detection-profiles.md, auto-tuning.md, adaptive-windowing.md) should be:
- Moved to an "Advanced" or "API Reference" section
- Not linked from main navigation
- Available for power users who search for it

The main docs should be:
1. Quick Start (send events, see anomalies)
2. Dashboard Guide (what each chart means)
3. Troubleshooting (common issues)

---

## File Reference

### Backend (Go)
- `collector-processor/cmd/driftlock-http/main.go` - Main server
- `collector-processor/cmd/driftlock-http/db.go` - Database operations, stream creation
- `collector-processor/cmd/driftlock-http/profiles.go` - Profile definitions
- `collector-processor/cmd/driftlock-http/autotune.go` - Auto-tune algorithm
- `collector-processor/cmd/driftlock-http/adaptive_windowing.go` - Window sizing
- `collector-processor/cmd/driftlock-http/adaptive_handlers.go` - API handlers

### Frontend (Vue)
- `landing-page/src/views/DashboardView.vue` - Main dashboard (1000+ lines)
- `landing-page/src/views/HomeView.vue` - Landing page
- `landing-page/src/views/PlaygroundView.vue` - Demo playground
- `landing-page/src/components/dashboard/` - Dashboard widgets
- `landing-page/src/stores/auth.ts` - Auth state

### Docs
- `docs/user-guide/guides/` - Detailed guides (hide these)
- `docs/user-guide/getting-started/` - Quickstart (keep visible)

---

## Success Criteria

1. **New user signs up** → Streams have auto_tune and adaptive_window ON by default
2. **User sees anomaly** → Can click thumbs up/down to provide feedback
3. **User wants less noise** → Can click "Low sensitivity" without reading docs
4. **System learns** → Visible indicator shows "Learning from your feedback"
5. **Power user** → Can still access all API endpoints and detailed docs

---

## Anti-Patterns to Avoid

- Don't make users read documentation to use basic features
- Don't expose technical terms (NCD, p-value, baseline) in main UI
- Don't require API calls for common operations
- Don't add wizards or modals - keep it simple
- Don't over-explain - the product should be self-evident

---

## Agent Routing

Use these specialized agents for implementation:
- **dl-frontend**: Vue dashboard components, feedback buttons, profile picker
- **dl-backend**: Default changes, any API modifications
- **dl-docs**: Reorganizing documentation hierarchy

See `docs/AI_ROUTING.md` for full agent matrix.

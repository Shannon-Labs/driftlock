# The "Streaming Variability" Prompt (Adaptive CBAD)

**Context:** Driftlock uses Compression-Based Anomaly Detection (CBAD). Currently, it uses fixed-width sliding windows (e.g., 400 events). This is predictable. A sophisticated attacker could "pulse" their attack to hide between checks.

**The Goal:** Implement **Adaptive Windowing** in the Rust core (`cbad-core`).
**The Hook:** "Ungameable Telemetry."

---

## The Task for the Next AI Agent

**1. Implement Randomized Window Sizes in `cbad-core`**
   - **Where:** `cbad-core/src/window.rs` (or similar).
   - **Logic:** Instead of `window_size = 50`, allow `window_range = [40, 60]`.
   - **Mechanism:** For each step, pick a random `k` within the range using a ChaCha20 RNG seeded by the `tenant_id` + `stream_id` (Deterministic Randomness).
   - **Why:** The attacker doesn't know the seed. They can't predict if the window is 42 or 58 events wide.

**2. Implement Overlapping Checks ("Jitter")**
   - **Logic:** Instead of checking every `N` events (stride), introduce `jitter`.
   - **Example:** Check at index `100`, then `100 + (stride + jitter)`.
   - **Result:** The sampling frequency becomes aperiodic.

**3. Update the Go FFI (`collector-processor`)**
   - Expose the `window_min`, `window_max`, and `jitter_pct` config options to the Go layer.

**4. Add "Entropy Fingerprinting"**
   - Calculate a rolling hash of the *compression ratios* themselves.
   - If the compression ratio curve becomes "too flat" (someone padding data to look normal), flag an "Anomaly of the Anomaly Detector."

---

## Technical Constraints
- **Determinism:** The "randomness" MUST be deterministic based on the stream ID. If we replay the same log file, we must get the exact same "random" windows. **Use `rand_chacha` with a seed derived from the stream ID.**
- **Performance:** Re-allocating windows is expensive. Use a ring buffer that can handle variable logical sizes without re-allocating memory.

---

## The Marketing Payoff
*"Driftlock locks onto the unique entropy signature of your stream. By randomizing the observation windows deterministically, we create a mathematical trap that no bot can evade."*

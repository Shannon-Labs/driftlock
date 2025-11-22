Phase 1 Plan – Core + Collector

Goals
- Implement `cbad-core` primitives (sliding window, entropy, compression metrics, NCD).
- Expose FFI surface for Go integration.
- Implement `driftlock_cbad` OTel Collector processor and call cbad.ComputeMetrics.

Tasks
- cbad-core
  - Define window buffer, adapters (zstd/lz4/gzip/OpenZL placeholders), and calculators.
  - Add deterministic permutation testing harness with seeded RNG.
- collector-processor
  - Implement config (window sizes, thresholds, stream kinds).
  - Wire metrics/logs processing; produce glass-box explanations.
- api-server
  - Seed API endpoints for write/read anomalies and artifacts.
  - Temporary in-memory store; plan Postgres schema.
- ui
  - Minimal list/detail for anomalies; fetch from api-server.

Exit Criteria
- Synthetic streams through Collector produce anomalies with explanations.
- Deterministic outputs with fixed seeds.


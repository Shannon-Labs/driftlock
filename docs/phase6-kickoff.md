# Phase 6 Kickoff Notes

Date: 2025-02-14

## Kafka Streaming Foundation

- Added `Streaming.Kafka` configuration block (`api-server/internal/config/config.go`) with env overrides for brokers, topics, and TLS.
- Scaffolding for publisher/subscriber abstractions lives in `api-server/internal/streaming/kafka` with an in-memory broker used for tests and future integration prototyping.
- API anomaly handler now emits `anomaly.created` events through the pluggable publisher interface; this currently targets the Kafka topic configured in `Streaming.Kafka.EventsTopic`.

## Immediate Next Steps

1. Implement real Kafka clients (producer and consumer) that satisfy the new interfaces, using configurable injection for unit tests.
2. Extend the collector processor to emit OTLP events onto the Kafka `EventsTopic` when streaming is enabled.
3. Build a dedicated ingestion worker (or extend the API engine) that consumes from Kafka, reusing CBAD detection logic.
4. Add load-test scenarios to validate publisher throughput before introducing external dependencies.

## Open Questions

- Do we require exactly-once semantics immediately, or can we start with at-least-once plus idempotent storage writes?
- Should anomaly exports be routed through a separate topic or share the events topic with filtering downstream?
- What is the secure default for TLS bootstrap (e.g. mTLS, SASL)? Need alignment with security team before production rollout.

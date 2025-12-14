# Kafka Integration with DriftLock

This document explains how to use Kafka ingestion with the DriftLock Rust API for streaming anomaly detection.

## Overview

DriftLock's Rust API includes an optional Kafka consumer that enables high-throughput streaming ingestion. When enabled, the API consumes messages from a Kafka topic and processes them through the CBAD anomaly detection engine.

## Architecture

```
┌─────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   Kafka     │────▶│  Driftlock API  │────▶│   PostgreSQL    │
│   Topic     │     │  (Rust/Axum)    │     │   Database      │
└─────────────┘     └────────┬────────┘     └─────────────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │   CBAD Core     │
                    │   (Detection)   │
                    └─────────────────┘
```

## Building with Kafka Support

The Kafka consumer is feature-gated. Build with the `kafka` feature enabled:

```bash
# Build with Kafka support
cargo build -p driftlock-api --features kafka --release

# Or run directly
cargo run -p driftlock-api --features kafka --release
```

## Configuration

Configure Kafka via environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `KAFKA_ENABLED` | Enable Kafka consumer | `false` |
| `KAFKA_BROKERS` | Comma-separated broker list | `localhost:9092` |
| `KAFKA_TOPIC` | Topic to consume from | `driftlock-events` |
| `KAFKA_GROUP_ID` | Consumer group ID | `driftlock-api` |
| `KAFKA_STREAM_ID_HEADER` | Header key for stream ID | (optional) |
| `KAFKA_STREAM_ID_FIELD` | JSON field path for stream ID | (optional) |
| `KAFKA_MAX_PAYLOAD_BYTES` | Max message size | `1048576` (1MB) |
| `KAFKA_MAX_IN_FLIGHT` | Concurrent message processing | `100` |

### Example Configuration

```bash
export KAFKA_ENABLED=true
export KAFKA_BROKERS=kafka:9092
export KAFKA_TOPIC=otlp-events
export KAFKA_GROUP_ID=driftlock-consumer
export KAFKA_STREAM_ID_HEADER=X-Stream-Id
export KAFKA_MAX_IN_FLIGHT=50
```

## Stream ID Resolution

The consumer needs to map each Kafka message to a DriftLock stream. Two methods are supported:

### 1. Header-Based (Recommended)

Set `KAFKA_STREAM_ID_HEADER` to the header key containing the stream UUID:

```
Header: X-Stream-Id: 550e8400-e29b-41d4-a716-446655440000
```

### 2. JSON Field-Based

Set `KAFKA_STREAM_ID_FIELD` to the JSON path:

```bash
export KAFKA_STREAM_ID_FIELD=metadata.stream_id
```

For a message like:
```json
{
  "metadata": {
    "stream_id": "550e8400-e29b-41d4-a716-446655440000"
  },
  "data": "..."
}
```

## Message Format

Messages should be raw bytes (typically JSON). The entire payload is processed by CBAD:

```json
{
  "timestamp": "2025-12-11T10:00:00Z",
  "level": "ERROR",
  "message": "Connection timeout to database",
  "service": "api-gateway",
  "trace_id": "abc123"
}
```

## Backpressure Handling

The consumer uses a semaphore-based backpressure mechanism:

- `KAFKA_MAX_IN_FLIGHT` controls concurrent message processing
- Messages are processed asynchronously with offset commits
- If processing falls behind, new messages wait for permits

## Metrics

When Kafka is enabled, these Prometheus metrics are emitted:

| Metric | Description |
|--------|-------------|
| `driftlock_events_processed_total` | Total events processed |
| `driftlock_anomalies_detected_total` | Total anomalies detected |
| `driftlock_stream_events_total{stream_id}` | Events per stream |
| `driftlock_stream_anomalies_total{stream_id}` | Anomalies per stream |

## Local Development

### Using Docker Compose

Start Kafka locally:

```bash
# Start Kafka and Zookeeper
docker compose -f docker-compose.kafka.yml up -d

# Verify Kafka is running
docker compose -f docker-compose.kafka.yml ps
```

### Testing the Consumer

```bash
# Run API with Kafka enabled
KAFKA_ENABLED=true \
KAFKA_BROKERS=localhost:9092 \
KAFKA_TOPIC=test-events \
KAFKA_STREAM_ID_HEADER=X-Stream-Id \
cargo run -p driftlock-api --features kafka

# In another terminal, produce test messages
docker exec -it kafka kafka-console-producer \
  --broker-list localhost:9092 \
  --topic test-events \
  --property "parse.headers=true"
```

## Production Considerations

### Security

- Enable TLS for broker connections (configure in `rdkafka` settings)
- Use SASL authentication for secured clusters
- Restrict topic access via ACLs

### Scaling

- Increase `KAFKA_MAX_IN_FLIGHT` for higher throughput
- Use multiple partitions for parallel consumption
- Deploy multiple API instances with the same `KAFKA_GROUP_ID` for horizontal scaling

### Reliability

- Offsets are committed after successful processing
- Failed messages are logged but don't block consumption
- Use dead-letter topics for failed message handling (configure separately)

## Troubleshooting

### Consumer Not Starting

Check logs for connection errors:
```bash
cargo run -p driftlock-api --features kafka 2>&1 | grep -i kafka
```

Common issues:
- Brokers unreachable (check `KAFKA_BROKERS`)
- Topic doesn't exist (auto-create may be disabled)
- Network/firewall issues

### Messages Not Processing

1. Verify stream exists in database
2. Check stream ID resolution (header vs field)
3. Ensure payload doesn't exceed `KAFKA_MAX_PAYLOAD_BYTES`

### High Latency

- Reduce `KAFKA_MAX_IN_FLIGHT` if CPU-bound
- Increase if I/O-bound and CPU available
- Check detector manager cache hit rate

## Implementation Details

Source: `crates/driftlock-api/src/kafka.rs`

Key components:
- `spawn_kafka_consumer()`: Spawns async consumer task
- `run_consumer()`: Main consume loop with backpressure
- `handle_message()`: Process single message through CBAD
- `resolve_stream_id()`: Extract stream ID from message

# Phase 6: Kafka Configuration Guide

This document outlines how to configure Apache Kafka for Driftlock's streaming capabilities in Phase 6.

## Environment Variables

### API Server
The API server can publish anomaly events to Kafka when enabled:

```bash
# Enable Kafka publishing
KAFKA_ENABLED=true

# Kafka broker addresses (comma-separated)
KAFKA_BROKERS=localhost:9092,localhost:9093

# Client and group IDs
KAFKA_CLIENT_ID=driftlock-api
KAFKA_GROUP_ID=driftlock-api

# Topic names
KAFKA_EVENTS_TOPIC=otlp-events
KAFKA_ANOMALIES_TOPIC=anomaly-events

# TLS configuration
KAFKA_TLS_ENABLED=false
```

### Collector-Processor
The collector-processor can publish raw OTLP events to Kafka:

```yaml
driftlock_cbad:
  window_size: 1024
  hop_size: 256
  threshold: 0.9
  determinism: true
  kafka:
    enabled: true
    brokers: ["localhost:9092"]
    client_id: "driftlock-collector"
    events_topic: "otlp-events"
    tls_enabled: false
    batch_size: 100
    batch_timeout_ms: 5
```

## Feature Flags

- When `KAFKA_ENABLED=false` (default), the API server uses an in-memory broker for testing
- When `kafka.enabled=false` in collector config, no Kafka publishing occurs

## Testing

To run tests without Kafka:
- Set `KAFKA_ENABLED=false` in your environment
- The in-memory broker will be used instead

To run tests with Kafka:
- Set up a Kafka instance (local or remote)
- Configure the environment variables appropriately
- Set `KAFKA_ENABLED=true`

## Production Considerations

- Enable TLS in production: `KAFKA_TLS_ENABLED=true`
- Adjust batch settings based on throughput requirements
- Use appropriate client and group IDs
- Monitor consumer lag for anomaly detection workflows
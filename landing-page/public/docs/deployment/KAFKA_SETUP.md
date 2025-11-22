# Kafka Setup for DriftLock Project

This document explains how to set up and manage Kafka with Zookeeper for the driftlock project.

## Prerequisites

- Colima installed and running (recommended for macOS)
- Bash shell

## Quick Start

### Option 1: Using Colima (Recommended for macOS)

1. Start Colima if not already running:
   ```bash
   colima start --cpu 4 --memory 4
   ```

2. Start the complete DriftLock stack:
   ```bash
   ./scripts/start-stack.sh start
   ```

3. Check the status:
   ```bash
   ./scripts/start-stack.sh status
   ```

### Option 2: Using Docker Desktop (Alternative)

1. Start Kafka and Zookeeper with the required topics:
   ```bash
   ./scripts/kafka-setup.sh start
   ```

2. Check the status:
   ```bash
   ./scripts/kafka-setup.sh status
   ```

## Available Commands

The `start-stack.sh` script provides the following commands:

- `start` - Start Kafka and Zookeeper, create topics, and start the complete stack
- `stop` - Stop all services
- `restart` - Restart all services
- `status` - Show status of all services
- `logs` - Show logs from all services

## Topics

The following topics are automatically created:

- `otlp-events` - For OpenTelemetry events
- `anomaly-events` - For anomaly detection events

## Configuration

The Kafka broker is configured with:
- Broker ID: 1
- Advertised listener: `PLAINTEXT://localhost:9092`
- Auto topic creation: enabled
- Offset topic replication factor: 1

## Integration with DriftLock

Once Kafka is running, you can configure your driftlock components to connect to:

- Kafka broker: `localhost:9092`
- Zookeeper: `localhost:2181`

The topics `otlp-events` and `anomaly-events` are ready for use by your driftlock application.

## Colima-Specific Notes

When using Colima instead of Docker Desktop:

1. Container networking works differently - use container IPs instead of localhost
2. The `start-stack.sh` script automatically detects and uses the appropriate Kafka broker address
3. For direct container access, use `colima exec` instead of `docker exec`
4. Port forwarding is handled automatically by Colima

## Troubleshooting

### Docker Commands Hanging

If Docker commands are hanging, try these steps:

1. Restart Colima:
   ```bash
   colima stop && colima start
   ```

2. Check if Colima is running:
   ```bash
   colima list
   ```

3. Verify Docker is working with Colima:
   ```bash
   docker ps
   ```

### Kafka Connection Issues

1. Check if Kafka is running:
   ```bash
   ./scripts/kafka-native-setup.sh status
   ```

2. Check Kafka logs:
   ```bash
   ./scripts/kafka-native-setup.sh logs
   ```

3. Test connectivity:
   ```bash
   ./scripts/kafka-native-setup.sh test
   ```

### Topic Issues

1. List all topics:
   ```bash
   ./scripts/kafka-native-setup.sh topics
   ```

2. Manually create a topic:
   ```bash
   docker exec driftlock-kafka kafka-topics --create --topic your-topic --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1
   ```

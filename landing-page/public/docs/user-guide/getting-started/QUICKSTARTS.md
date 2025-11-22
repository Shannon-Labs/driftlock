# Driftlock Quickstarts

These quickstarts help you try Driftlock across all surfaces: CLI, REST API, and Kafka (compose).

## Prerequisites
- Docker 24+
- Make, Go 1.24+ (optional for local builds)
- Rust 1.70+ (for building CBAD core)

**Important:** If cloning the repository, initialize git submodules first:
```bash
git clone https://github.com/Shannon-Labs/driftlock.git
cd driftlock
git submodule update --init --recursive  # Required for OpenZL support
```

## 1) CLI

Run on NDJSON or JSON array files using the streaming detector.

```bash
# Ensure submodules are initialized (if not already done)
git submodule update --init --recursive

# Build binary inside the module
cd collector-processor
go build -o ../bin/driftlock-cli ./cmd/driftlock-cli
cd ..

# NDJSON
./bin/driftlock-cli --input test-data/normal-transactions.jsonl --format ndjson --output -

# JSON array (auto-detected)
./bin/driftlock-cli --input test-data/test-demo.json --output results.json
```

Flags:
- --baseline: number of events to build the baseline (default 400)
- --window: events per detection window (default 1)
- --hop: hop size between windows (default 1)
- --algo: zstd|lz4|gzip|openzl (default zstd)

## 2) REST API (Docker Compose - Recommended)

The unified `docker-compose.yml` includes both the HTTP API and optional Kafka collector:

```bash
# Start HTTP API server (default)
docker compose up -d driftlock-http

# Health check
curl -s http://localhost:8080/healthz

# NDJSON detection
curl -s -X POST "http://localhost:8080/v1/detect?format=ndjson" \
  -H "Content-Type: application/json" \
  --data-binary @test-data/mixed-transactions.jsonl | jq .
```

### Standalone Docker Build

```bash
# Build image (without OpenZL)
docker build -t driftlock-http:latest -f collector-processor/cmd/driftlock-http/Dockerfile .

# Run
docker run --rm -p 8080:8080 -e CORS_ALLOW_ORIGINS=http://localhost:5174 driftlock-http:latest
```

Optional OpenZL:
```bash
docker build -t driftlock-http:openzl \
  --build-arg USE_OPENZL=true \
  -f collector-processor/cmd/driftlock-http/Dockerfile .
```

## 3) Kafka Collector (docker compose)

The unified `docker-compose.yml` supports Kafka with profiles:

```bash
# Start HTTP API + Kafka + Collector together
docker compose --profile kafka up -d

# Or use the dedicated Kafka compose file
docker compose -f docker-compose.kafka.yml up -d

# Produce sample events
docker run --rm -it --network host edenhill/kafkacat:1.7.1 \
  -b localhost:9092 -t driftlock-events -P test-data/mixed-transactions.jsonl
```

Notes:
- Use `--profile kafka` to enable Kafka services in the unified compose file
- The collector processes OpenTelemetry streams via Kafka
- The build includes the CBAD Rust core so you can extend the collector quickly

## 4) Hosted API

Follow the Render or Cloud Run deployment guide:
- docs/deployment/hosted-api.md

## 5) Web Playground

```bash
cd playground
npm install
cp .env.example .env   # set VITE_API_BASE_URL to your API, e.g. http://localhost:8080
npm run dev            # opens http://localhost:5174
```




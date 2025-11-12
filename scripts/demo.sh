#!/bin/bash
# One-command demo setup for Driftlock

set -e

echo "🎬 Driftlock Demo Setup"
echo "======================"
echo ""

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check Docker
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Docker is not installed. Please install Docker first.${NC}"
    exit 1
fi

# Start services
echo -e "${BLUE}Starting Driftlock services...${NC}"
docker compose up -d driftlock-http

echo ""
echo "Waiting for services to be ready..."
sleep 5

# Check health
echo ""
echo -e "${BLUE}Checking service health...${NC}"
if curl -sf http://localhost:8080/healthz > /dev/null 2>&1; then
    echo -e "${GREEN}✓ API server is healthy${NC}"
else
    echo -e "${YELLOW}⚠ API server is not responding yet. Waiting...${NC}"
    sleep 10
    if curl -sf http://localhost:8080/healthz > /dev/null 2>&1; then
        echo -e "${GREEN}✓ API server is now healthy${NC}"
    else
        echo -e "${RED}✗ API server failed to start${NC}"
        echo "Check logs: docker compose logs driftlock-http"
        exit 1
    fi
fi

# Display service information
echo ""
echo -e "${BLUE}Service Information:${NC}"
echo "  API Server: http://localhost:8080"
echo "  Health Check: http://localhost:8080/healthz"
echo "  Metrics: http://localhost:8080/metrics"
echo "  API Endpoint: http://localhost:8080/v1/detect"

# Test with sample data
echo ""
echo -e "${BLUE}Running sample detection...${NC}"
SAMPLE_DATA='{"timestamp":"2025-01-01T00:00:00Z","value":100}
{"timestamp":"2025-01-01T00:00:01Z","value":101}
{"timestamp":"2025-01-01T00:00:02Z","value":99999}'

RESPONSE=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    --data-binary "$SAMPLE_DATA" \
    "http://localhost:8080/v1/detect?format=ndjson")

ANOMALY_COUNT=$(echo "$RESPONSE" | jq -r '.anomaly_count // 0')
TOTAL_EVENTS=$(echo "$RESPONSE" | jq -r '.total_events // 0')

echo ""
echo -e "${GREEN}Sample Detection Results:${NC}"
echo "  Total Events: $TOTAL_EVENTS"
echo "  Anomalies Detected: $ANOMALY_COUNT"

# Playground instructions
echo ""
echo -e "${BLUE}Playground Setup:${NC}"
echo "  To start the playground UI:"
echo "    1. cd playground"
echo "    2. npm install"
echo "    3. cp .env.example .env"
echo "    4. npm run dev"
echo ""
echo "  The playground will be available at http://localhost:5174"

# Optional Kafka setup
echo ""
read -p "Start Kafka services? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${BLUE}Starting Kafka services...${NC}"
    docker compose --profile kafka up -d
    echo -e "${GREEN}✓ Kafka services started${NC}"
    echo "  Kafka: localhost:9092"
    echo "  Zookeeper: localhost:2181"
fi

echo ""
echo -e "${GREEN}Demo setup complete!${NC}"
echo ""
echo "Next steps:"
echo "  1. Test the API: curl http://localhost:8080/healthz"
echo "  2. View metrics: curl http://localhost:8080/metrics"
echo "  3. Run detection: curl -X POST http://localhost:8080/v1/detect -H 'Content-Type: application/json' --data-binary @test-data/mixed-transactions.jsonl"
echo "  4. Start playground: cd playground && npm run dev"


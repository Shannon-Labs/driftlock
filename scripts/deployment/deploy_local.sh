#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}Starting Driftlock Local Deployment...${NC}"

# Check for .env file
if [ ! -f .env ]; then
    echo -e "${RED}Error: .env file not found.${NC}"
    echo "Please copy .env.example to .env and fill in the required values."
    exit 1
fi

# Check for ANTHROPIC_API_KEY
if ! grep -q "ANTHROPIC_API_KEY" .env; then
    echo -e "${RED}Error: ANTHROPIC_API_KEY not found in .env.${NC}"
    echo "Please add your Anthropic API key to .env."
    exit 1
fi

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}Error: Docker is not running.${NC}"
    echo "Please start Docker Desktop and try again."
    exit 1
fi

echo "Building and starting services..."
docker compose up --build -d

echo -e "${GREEN}Deployment successful!${NC}"
echo "Driftlock HTTP API is running at http://localhost:8080"
echo "To view logs: docker compose logs -f"

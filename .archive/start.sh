#!/bin/bash

set -euo pipefail

echo "🚀 Starting Driftlock DORA Compliance Demo..."
echo ""
echo "📋 This demo uses pre-configured settings for quick evaluation"
echo "🔐 API Key: demo-key-123 (hardcoded for demo)"
echo "💾 Database: PostgreSQL with demo data auto-loaded"
echo ""

# Start services with docker-compose
echo "🐳 Starting services with Docker Compose..."
if ! docker compose up -d; then
    echo "❌ Failed to start services"
    echo "💡 Check Docker is running and ports 3000, 5432, 8080 are available"
    exit 1
fi

# Wait for services to initialize
echo "⏳ Waiting for services to initialize..."
sleep 15

echo ""
echo "🎉 Driftlock started successfully!"
echo ""
echo "📊 Web Dashboard: http://localhost:3000"
echo "🔧 API Server: http://localhost:8080"
echo "📈 API Health: http://localhost:8080/healthz"
echo "📋 API Documentation: http://localhost:8080/swagger/"
echo ""
echo "🔑 Dashboard Login: Use API key 'demo-key-123'"
echo "💡 Demo data is loading automatically (1,600 transactions)"
echo "💡 To stop services: docker compose down"
echo "💡 To view logs: docker compose logs -f [service-name]"
echo ""
echo "🚀 Ready for DORA compliance demo!"
echo ""
echo "⏱️  Dashboard should show anomalies within 60 seconds..."
echo "   Look for flagged payment latency spikes with explanations"
#!/bin/bash

set -euo pipefail

# Driftlock Open Source System Startup Script

echo "🚀 Starting Driftlock Open Source System..."

# Check if .env file exists
if [ ! -f .env ]; then
	echo "❌ .env file not found. Please copy .env.example to .env and configure your values."
    echo ""
    echo "📋 Quick setup:"
    echo "   cp .env.example .env"
    echo "   # Edit .env and set your API key and database password"
    exit 1
fi

# Load environment variables
source .env

# Check if API key is configured
DEFAULT_API_KEY_VALUE="${DEFAULT_API_KEY:-}"
if [ -z "$DEFAULT_API_KEY_VALUE" ] || [ "$DEFAULT_API_KEY_VALUE" = "your_api_key_here_for_dashboard_access" ]; then
	echo "⚠️  API key not configured in .env"
	echo "   Please set DEFAULT_API_KEY to secure your dashboard access"
	echo "   Example: DEFAULT_API_KEY=your-secret-api-key-here"
	echo ""
	echo "🔐 This key will be used to log into the web dashboard"
	exit 1
fi

# Check if database password is configured
DB_PASSWORD_VALUE="${DB_PASSWORD:-}"
if [ -z "$DB_PASSWORD_VALUE" ] || [ "$DB_PASSWORD_VALUE" = "your_secure_password_here" ]; then
	echo "⚠️  Database password not configured in .env"
	echo "   Please set DB_PASSWORD to secure your database"
	echo "   Example: DB_PASSWORD=your-secure-db-password"
	exit 1
fi

export DEFAULT_ORG_ID="${DEFAULT_ORG_ID:-default}"
if [ -z "${DRIFTLOCK_DEV_API_KEY:-}" ]; then
	export DRIFTLOCK_DEV_API_KEY="$DEFAULT_API_KEY_VALUE"
	echo "ℹ️  DRIFTLOCK_DEV_API_KEY not set; mirroring DEFAULT_API_KEY for local access"
fi

echo "✅ Environment loaded successfully"
echo ""

# Start services with docker-compose
echo "🐳 Starting services with Docker Compose..."
if ! docker compose up -d; then
	echo "❌ Failed to start services"
	echo "💡 Check Docker is running and ports 3000, 5432, 8080 are available"
	exit 1
fi

# Wait a moment for services to initialize
echo "⏳ Waiting for services to initialize..."
sleep 10

echo ""
echo "🎉 Driftlock started successfully!"
echo ""
echo "📊 Web Dashboard: http://localhost:3000"
echo "🔧 API Server: http://localhost:8080"
echo "📈 API Health: http://localhost:8080/healthz"
echo "📋 API Documentation: http://localhost:8080/swagger/"
echo ""
echo "🔑 Dashboard Login: Use your API key from .env"
echo "💡 To stop services: docker compose down"
echo "💡 To view logs: docker compose logs -f [service-name]"
echo ""
echo "🚀 Ready for anomaly detection!"

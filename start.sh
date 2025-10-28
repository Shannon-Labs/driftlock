#!/bin/bash

# DriftLock Integrated System Startup Script

echo "🚀 Starting DriftLock Integrated System..."

# Check if .env file exists
if [ ! -f .env ]; then
    echo "❌ .env file not found. Please copy .env.example to .env and configure your values."
    exit 1
fi

# Load environment variables
source .env

# Check if Supabase is configured
if [ -z "$SUPABASE_PROJECT_ID" ] || [ -z "$SUPABASE_ANON_KEY" ]; then
    echo "⚠️  Supabase configuration not found in .env"
    echo "   Please set SUPABASE_PROJECT_ID, SUPABASE_ANON_KEY, and SUPABASE_SERVICE_ROLE_KEY"
    echo "   You can get these values from your Supabase project settings."
    echo ""
    echo "📖 For help, see INTEGRATION_README.md"
    exit 1
fi

echo "✅ Environment loaded successfully"
echo ""

# Start services with docker-compose
echo "🐳 Starting services with Docker Compose..."
docker-compose up -d

# Check if services started successfully
if [ $? -eq 0 ]; then
    echo ""
    echo "🎉 Services started successfully!"
    echo ""
    echo "📊 Web Frontend: http://localhost:3000"
    echo "🔧 Go API Server: http://localhost:8080"
    echo "📈 Prometheus Metrics: http://localhost:9090/metrics"
    echo ""
    echo "💡 To stop services: docker-compose down"
    echo "💡 To view logs: docker-compose logs -f [service-name]"
else
    echo "❌ Failed to start services"
    exit 1
fi

#!/bin/bash

# Driftlock Local Development Setup Script
# This script sets up a complete local development environment

set -e

echo "🚀 Setting up Driftlock Local Development Environment..."
echo ""

# Function to install dependencies
install_dependencies() {
    echo "📦 Installing dependencies..."

    # Check and install Docker
    if ! command -v docker &> /dev/null; then
        echo "❌ Docker not found. Please install Docker Desktop first:"
        echo "   macOS: https://docs.docker.com/docker-for-mac/install/"
        echo "   Ubuntu: sudo apt-get install docker.io"
        exit 1
    fi

    # Check and install Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        echo "❌ Docker Compose not found. Please install Docker Compose:"
        echo "   macOS: Included with Docker Desktop"
        echo "   Ubuntu: sudo apt-get install docker-compose"
        exit 1
    fi

    # Check and install Go
    if ! command -v go &> /dev/null; then
        echo "❌ Go not found. Please install Go 1.24+:"
        echo "   macOS: brew install go"
        echo "   Ubuntu: sudo apt-get install golang-go"
        exit 1
    fi

    # Check and install Rust
    if ! command -v cargo &> /dev/null; then
        echo "❌ Rust not found. Installing Rust..."
        curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
        source "$HOME/.cargo/env"
    fi

    # Install Node.js for frontend
    if ! command -v node &> /dev/null; then
        echo "❌ Node.js not found. Installing Node.js..."
        curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
        sudo apt-get install -y nodejs
    fi

    echo "✅ Dependencies installed"
}

# Function to set up local database
setup_local_database() {
    echo "🗄️  Setting up local database..."

    # Start PostgreSQL and Redis via Docker Compose
    echo "Starting local services with Docker Compose..."
    docker-compose up -d driftlock-postgres

    # Wait for PostgreSQL to be ready
    echo "Waiting for PostgreSQL to be ready..."
    for i in {1..30}; do
        if docker-compose exec -T driftlock-postgres pg_isready -U driftlock -d driftlock &>/dev/null; then
            echo "✅ PostgreSQL is ready"
            break
        fi
        sleep 1
    done

    # Run database migrations
    echo "🔄 Running database migrations..."
    if [ -d "api/migrations" ]; then
        # Find the latest migration file
        latest_migration=$(ls -1 api/migrations/*.sql | tail -1)
        if [ -n "$latest_migration" ]; then
            echo "Applying migration: $latest_migration"
            docker-compose exec -T driftlock-postgres psql -U driftlock -d driftlock < "$latest_migration"
            echo "✅ Migrations applied"
        fi
    fi

    echo "✅ Local database setup complete"
}

# Function to build the backend
build_backend() {
    echo "🔨 Building Driftlock backend..."

    cd collector-processor/cmd/driftlock-http

    # Build the Go binary
    echo "Building Go binary..."
    go mod tidy
    go build -o ../../driftlock-http .

    cd ../../..

    echo "✅ Backend built successfully"
}

# Function to setup frontend
setup_frontend() {
    echo "🎨 Setting up frontend..."

    cd landing-page

    # Install dependencies
    if [ ! -d "node_modules" ]; then
        echo "Installing frontend dependencies..."
        npm install
    fi

    # Build for development
    echo "Building frontend for development..."
    npm run build

    cd ..
    echo "✅ Frontend setup complete"
}

# Function to create environment file
create_env_file() {
    echo "📝 Creating environment configuration..."

    cat > .env.local << 'EOF'
# Driftlock Local Development Environment

# Database Configuration
DATABASE_URL=postgres://driftlock:driftlock@localhost:7543/driftlock?sslmode=disable

# Service Configuration
PORT=8080
LOG_LEVEL=debug
DRIFTLOCK_DEV_MODE=true
DRIFTLOCK_LICENSE_KEY=dev-mode
CORS_ALLOW_ORIGINS=http://localhost:3000,http://localhost:5173

# External Services (Use test keys)
STRIPE_SECRET_KEY=sk_test_placeholder
STRIPE_PRICE_ID_PRO=price_test_placeholder
SENDGRID_API_KEY=SG_test_placeholder

# Admin Configuration
ADMIN_KEY=admin-dev-key-12345

# Frontend Development
VITE_API_BASE_URL=http://localhost:8080
VITE_STRIPE_PUBLISHABLE_KEY=pk_test_placeholder
EOF

    echo "✅ Environment file created: .env.local"
}

# Function to create development scripts
create_dev_scripts() {
    echo "📜 Creating development scripts..."

    # Script to start API server
    cat > scripts/start-api.sh << 'EOF'
#!/bin/bash
# Start Driftlock API server locally

set -e

cd "$(dirname "$0")/.."

# Load environment variables
if [ -f .env.local ]; then
    export $(cat .env.local | grep -v '^#' | xargs)
fi

# Build if binary doesn't exist
if [ ! -f "collector-processor/driftlock-http" ]; then
    echo "Building API server..."
    cd collector-processor/cmd/driftlock-http
    go build -o ../../driftlock-http .
    cd ../..
fi

# Start the API server
echo "🚀 Starting Driftlock API server on http://localhost:8080"
cd collector-processor
./driftlock-http
EOF

    # Script to start frontend
    cat > scripts/start-frontend.sh << 'EOF'
#!/bin/bash
# Start Driftlock frontend locally

set -e

cd "$(dirname "$0")/../landing-page"

# Load environment variables
if [ -f ../.env.local ]; then
    export $(cat ../.env.local | grep '^VITE_' | xargs)
fi

# Start development server
echo "🎨 Starting frontend development server on http://localhost:5173"
npm run dev
EOF

    # Script to run tests
    cat > scripts/run-tests.sh << 'EOF'
#!/bin/bash
# Run Driftlock tests

set -e

cd "$(dirname "$0")/.."

echo "🧪 Running backend tests..."
cd collector-processor
go test ./...

echo "🧪 Running frontend tests..."
cd ../landing-page
npm test

echo "✅ All tests completed"
EOF

    # Make scripts executable
    chmod +x scripts/start-api.sh scripts/start-frontend.sh scripts/run-tests.sh

    echo "✅ Development scripts created"
}

# Function to create local Supabase setup
setup_local_supabase() {
    echo "🔥 Setting up local Supabase (alternative to PostgreSQL)..."

    if ! command -v supabase &> /dev/null; then
        echo "❌ Supabase CLI not found. Install with: brew install supabase/tap/supabase"
        return
    fi

    # Initialize Supabase if not already done
    if [ ! -d "supabase" ]; then
        supabase init
    fi

    # Start local Supabase
    echo "Starting local Supabase..."
    supabase start

    # Get the local database URL
    db_url=$(supabase status | grep "DB URL" | awk '{print $3}')

    # Update environment file with Supabase URL
    sed -i.bak "s|DATABASE_URL=.*|DATABASE_URL=$db_url|" .env.local

    echo "✅ Local Supabase started with URL: $db_url"
}

# Function to run local services
run_local_services() {
    echo "🚀 Starting all local services..."

    # Start database services
    echo "Starting database services..."
    docker-compose up -d driftlock-postgres

    # Start API server in background
    echo "Starting API server..."
    scripts/start-api.sh &
    API_PID=$!

    # Start frontend in background
    echo "Starting frontend..."
    scripts/start-frontend.sh &
    FRONTEND_PID=$!

    echo ""
    echo "🎉 Local development environment started!"
    echo "======================================="
    echo ""
    echo "📍 API Server: http://localhost:8080"
    echo "📍 Frontend:   http://localhost:5173"
    echo "📍 Database:   localhost:7543 (PostgreSQL)"
    echo ""
    echo "🛑 To stop services:"
    echo "  kill $API_PID $FRONTEND_PID"
    echo "  docker-compose down"
    echo ""
    echo "🧪 To test API:"
    echo "  curl http://localhost:8080/healthz"
    echo ""

    # Wait for user to stop
    read -p "Press Enter to stop all services..."

    # Clean up
    echo "🛑 Stopping services..."
    kill $API_PID $FRONTEND_PID 2>/dev/null || true
    docker-compose down

    echo "✅ All services stopped"
}

# Function to show development tips
show_dev_tips() {
    echo "💡 Local Development Tips:"
    echo "=========================="
    echo ""
    echo "🔧 Environment Variables:"
    echo "  All settings are in .env.local"
    echo "  Modify as needed for your setup"
    echo ""
    echo "🗄️  Database Access:"
    echo "  psql -h localhost -p 7543 -U driftlock -d driftlock"
    echo "  Password: driftlock"
    echo ""
    echo "🚀 Quick Start Commands:"
    echo "  ./scripts/start-api.sh      # Start API server"
    echo "  ./scripts/start-frontend.sh  # Start frontend dev server"
    echo "  ./scripts/run-tests.sh       # Run all tests"
    echo ""
    echo "🐛 Debugging:"
    echo "  API logs: Check console output"
    echo "  Frontend logs: Browser dev tools"
    echo "  Database: docker-compose logs driftlock-postgres"
    echo ""
    echo "🔄 Hot Reload:"
    echo "  Backend: Rebuild and restart manually"
    echo "  Frontend: Automatic with Vite"
    echo ""
}

# Main execution
echo "Choose setup option:"
echo "1) Complete local development setup (recommended)"
echo "2) Quick start (assumes dependencies are installed)"
echo "3) Set up local Supabase instead of PostgreSQL"
echo "4) Start all services"
echo "5) Show development tips"

read -p "Enter your choice (1-5): " choice

case $choice in
    1)
        install_dependencies
        setup_local_database
        build_backend
        setup_frontend
        create_env_file
        create_dev_scripts
        show_dev_tips
        ;;
    2)
        setup_local_database
        build_backend
        setup_frontend
        create_env_file
        create_dev_scripts
        ;;
    3)
        setup_local_supabase
        build_backend
        setup_frontend
        create_env_file
        create_dev_scripts
        ;;
    4)
        run_local_services
        ;;
    5)
        show_dev_tips
        ;;
    *)
        echo "❌ Invalid choice"
        exit 1
        ;;
esac

echo ""
echo "🎯 Local development setup complete!"
echo ""
echo "Next steps:"
echo "1. Review .env.local for configuration"
echo "2. Start services with: ./scripts/start-api.sh"
echo "3. Start frontend with: ./scripts/start-frontend.sh"
echo "4. Test API: curl http://localhost:8080/healthz"
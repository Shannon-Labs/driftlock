#!/bin/bash

# Driftlock Supabase Setup Script
# This script helps set up Supabase for Driftlock and provides the connection string

set -e

echo "🚀 Setting up Supabase for Driftlock..."
echo ""

# Function to check if supabase CLI is installed
check_supabase_cli() {
    if ! command -v supabase &> /dev/null; then
        echo "❌ Supabase CLI not found. Please install it first:"
        echo "   macOS: brew install supabase/tap/supabase"
        echo "   Linux: curl -L https://github.com/supabase/cli/releases/latest/download/supabase_linux_amd64.tar.gz | tar xz"
        echo "   Windows: scoop install supabase"
        exit 1
    fi
    echo "✅ Supabase CLI found"
}

# Function to set up local Supabase (for development)
setup_local_supabase() {
    echo "🔧 Setting up local Supabase instance..."

    supabase init
    supabase start

    echo "✅ Local Supabase started"
    echo ""
    echo "📋 Local Connection Details:"
    echo "Host: localhost"
    echo "Port: 54322"
    echo "Database: postgres"
    echo "User: postgres"
    echo "Password: postgres"
    echo ""

    # Run migrations
    echo "🔄 Running database migrations..."
    supabase db push
    echo "✅ Migrations completed"
}

# Function to guide through Supabase Cloud setup
setup_cloud_supabase() {
    echo "🌐 Setting up Supabase Cloud..."
    echo ""
    echo "Please follow these steps:"
    echo ""
    echo "1. Go to https://supabase.com/dashboard"
    echo "2. Click 'New Project'"
    echo "3. Choose your organization"
    echo "4. Enter project details:"
    echo "   - Name: driftlock"
    echo "   - Database Password: [generate a strong password]"
    echo "   - Region: Choose nearest to your users"
    echo ""
    echo "5. Wait for project creation (2-3 minutes)"
    echo ""

    read -p "Press Enter once you've created your Supabase project..."

    echo ""
    echo "📋 To get your connection string:"
    echo "1. In Supabase Dashboard, go to Settings > Database"
    echo "2. Scroll down to 'Connection string'"
    echo "3. Copy the 'Transaction pooler' connection string"
    echo "4. Replace [YOUR-PASSWORD] with your actual database password"
    echo ""

    echo "The connection string should look like:"
    echo "postgresql://postgres.[ref]:[password]@aws-0-[region].pooler.supabase.com:6543/postgres?sslmode=require"
    echo ""

    read -p "Enter your Supabase connection string: " connection_string

    if [ -z "$connection_string" ]; then
        echo "❌ No connection string provided"
        exit 1
    fi

    echo "✅ Connection string received"
    echo ""
    echo "🔐 To add this to GCP Secret Manager, run:"
    echo "echo -n '$connection_string' | gcloud secrets create driftlock-db-url --data-file=- --project=driftlock"
    echo ""

    # Save to temporary file for later use
    echo "$connection_string" > /tmp/supabase-connection-string.txt
    echo "💾 Connection string saved to /tmp/supabase-connection-string.txt"
}

# Function to run migrations
run_migrations() {
    echo "🔄 Running database migrations..."

    # Find migration files
    if [ -d "api/migrations" ]; then
        echo "Found migrations in api/migrations/"
        ls -la api/migrations/

        # For local development
        if supabase status &> /dev/null; then
            echo "Running migrations on local Supabase..."
            supabase db reset
            echo "✅ Local migrations completed"
        else
            echo "⚠️  Local Supabase not running. Please run migrations manually in Supabase dashboard:"
            echo "1. Go to Supabase Dashboard > SQL Editor"
            echo "2. Copy and run the contents of api/migrations/*.sql files"
        fi
    else
        echo "⚠️  No migrations directory found at api/migrations/"
    fi
}

# Main execution
check_supabase_cli

echo "Choose setup option:"
echo "1) Local development (Docker-based Supabase)"
echo "2) Cloud Supabase (production-ready)"
echo "3) Just run migrations (assuming Supabase is already set up)"

read -p "Enter your choice (1-3): " choice

case $choice in
    1)
        setup_local_supabase
        ;;
    2)
        setup_cloud_supabase
        ;;
    3)
        run_migrations
        ;;
    *)
        echo "❌ Invalid choice"
        exit 1
        ;;
esac

echo ""
echo "🎉 Supabase setup complete!"
echo ""
echo "Next steps:"
echo "1. Update your GCP secrets with the connection string"
echo "2. Run the GCP secrets setup script: ./scripts/setup-gcp-secrets.sh"
echo "3. Deploy to GCP: gcloud builds submit --config=cloudbuild.yaml"
#!/bin/bash

# Driftlock GCP Secrets Setup Script
# This script creates all required secrets in Google Secret Manager

set -e

PROJECT_ID="driftlock"
REGION="us-central1"

echo "🚀 Setting up GCP secrets for Driftlock..."
echo "Project: $PROJECT_ID"
echo "Region: $REGION"
echo ""

# Function to create a secret from user input
create_secret_from_input() {
    local secret_name="$1"
    local description="$2"
    local example_value="$3"

    echo "📝 Setting up secret: $secret_name"
    echo "Description: $description"
    echo "Example: $example_value"
    echo ""

    read -s -p "Enter value for $secret_name (or press Enter to use placeholder): " secret_value
    echo ""

    if [ -z "$secret_value" ]; then
        echo "⚠️  Using placeholder value - you'll need to update this later"
        secret_value="placeholder-update-me"
    fi

    echo -n "$secret_value" | gcloud secrets create "$secret_name" \
        --data-file=- \
        --project="$PROJECT_ID" \
        --description="$description" \
        --replication-policy="automatic"

    echo "✅ Secret '$secret_name' created successfully"
    echo ""
}

# Function to create a secret with a known value
create_known_secret() {
    local secret_name="$1"
    local value="$2"
    local description="$3"

    echo "📝 Creating secret: $secret_name"
    echo "Description: $description"

    echo -n "$value" | gcloud secrets create "$secret_name" \
        --data-file=- \
        --project="$PROJECT_ID" \
        --description="$description" \
        --replication-policy="automatic"

    echo "✅ Secret '$secret_name' created successfully"
    echo ""
}

# Check if gcloud is authenticated
echo "🔐 Checking GCP authentication..."
if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q .; then
    echo "❌ Please run: gcloud auth login"
    exit 1
fi

echo "✅ GCP authentication confirmed"
echo ""

# Enable required APIs
echo "🔧 Enabling required APIs..."
gcloud services enable secretmanager.googleapis.com --project="$PROJECT_ID"
gcloud services enable run.googleapis.com --project="$PROJECT_ID"
gcloud services enable cloudbuild.googleapis.com --project="$PROJECT_ID"
gcloud services enable firebase.googleapis.com --project="$PROJECT_ID"
echo "✅ APIs enabled"
echo ""

# Create known secrets first
echo "Creating known secrets..."
create_known_secret "driftlock-license-key" "dev-mode" "License key for Driftlock (using dev-mode for testing)"

# Create secrets that need user input
create_secret_from_input "driftlock-db-url" \
    "PostgreSQL Connection String (Supabase Transaction Pooler)" \
    "postgresql://postgres.[ref]:[pass]@aws-0-[region].pooler.supabase.com:6543/postgres?sslmode=require"

create_secret_from_input "sendgrid-api-key" \
    "SendGrid API Key for email delivery" \
    "SG.xxxxxxxx..."

create_secret_from_input "stripe-secret-key" \
    "Stripe Secret Key (use sk_test_ for testing)" \
    "sk_test_xxxxxxxx..."

create_secret_from_input "stripe-price-id-pro" \
    "Stripe Price ID for the Pro plan" \
    "price_1xxxxxxxx..."

# Optional admin key
echo "🔐 Setting up optional admin key..."
read -s -p "Enter admin key for dashboard access (leave empty to generate random): " admin_key

if [ -z "$admin_key" ]; then
    admin_key=$(openssl rand -hex 32)
    echo "Generated random admin key: $admin_key"
fi

echo -n "$admin_key" | gcloud secrets create "admin-key" \
    --data-file=- \
    --project="$PROJECT_ID" \
    --description="Static key for admin dashboard access" \
    --replication-policy="automatic"

echo "✅ Admin key secret created"
echo ""

# List all created secrets
echo "📋 All secrets created:"
gcloud secrets list --project="$PROJECT_ID" --filter="name~driftlock OR name=admin-key" --format="table(name,createTime)"

echo ""
echo "🎉 GCP secrets setup complete!"
echo ""
echo "Next steps:"
echo "1. Update any placeholder secrets with real values"
echo "2. Set up Supabase database (if not done already)"
echo "3. Configure Stripe products and prices"
echo "4. Run: gcloud builds submit --config=cloudbuild.yaml"
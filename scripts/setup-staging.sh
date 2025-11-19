#!/bin/bash

# Driftlock Staging Environment Setup Script
# This script sets up a separate staging environment for testing

set -e

STAGING_PROJECT="driftlock-staging"
REGION="us-central1"

echo "🚀 Setting up Driftlock Staging Environment..."
echo "Staging Project: $STAGING_PROJECT"
echo "Region: $REGION"
echo ""

# Function to create staging GCP project
create_staging_project() {
    echo "🏗️  Creating staging GCP project..."

    # Check if project already exists
    if gcloud projects describe "$STAGING_PROJECT" &> /dev/null; then
        echo "✅ Staging project already exists"
    else
        echo "Creating new project: $STAGING_PROJECT"
        gcloud projects create "$STAGING_PROJECT" --name="Driftlock Staging"
        echo "✅ Staging project created"
    fi

    # Set as current project
    gcloud config set project "$STAGING_PROJECT"

    # Enable required APIs
    echo "🔧 Enabling APIs..."
    gcloud services enable secretmanager.googleapis.com
    gcloud services enable run.googleapis.com
    gcloud services enable cloudbuild.googleapis.com
    gcloud services enable firebase.googleapis.com

    echo "✅ APIs enabled"
}

# Function to create staging secrets
create_staging_secrets() {
    echo "🔐 Creating staging secrets..."

    # Use test/staging values
    echo -n "dev-mode" | gcloud secrets create driftlock-license-key \
        --data-file=- \
        --project="$STAGING_PROJECT" \
        --description="License key for Driftlock staging (dev-mode)" \
        --replication-policy="automatic"

    # Use placeholder for database (will be updated later)
    echo -n "staging-db-connection-string" | gcloud secrets create driftlock-db-url \
        --data-file=- \
        --project="$STAGING_PROJECT" \
        --description="Staging database connection string" \
        --replication-policy="automatic"

    # Test Stripe key (if available)
    read -p "Enter Stripe test key (sk_test_...) or press Enter to use placeholder: " stripe_key
    if [ -z "$stripe_key" ]; then
        stripe_key="sk_test_placeholder"
    fi

    echo -n "$stripe_key" | gcloud secrets create stripe-secret-key \
        --data-file=- \
        --project="$STAGING_PROJECT" \
        --description="Stripe test key for staging" \
        --replication-policy="automatic"

    # Use test price ID
    echo -n "price_test_placeholder" | gcloud secrets create stripe-price-id-pro \
        --data-file=- \
        --project="$STAGING_PROJECT" \
        --description="Test price ID for staging" \
        --replication-policy="automatic"

    # Test SendGrid key
    echo -n "SG_test_placeholder" | gcloud secrets create sendgrid-api-key \
        --data-file=- \
        --project="$STAGING_PROJECT" \
        --description="SendGrid test key for staging" \
        --replication-policy="automatic"

    echo "✅ Staging secrets created"
}

# Function to create staging Cloud Build config
create_staging_cloudbuild() {
    echo "📝 Creating staging Cloud Build configuration..."

    cat > cloudbuild-staging.yaml << 'EOF'
steps:
  # Build Docker image
  - name: 'gcr.io/cloud-builders/docker'
    id: 'build-image'
    args:
      - 'build'
      - '-f'
      - 'collector-processor/cmd/driftlock-http/Dockerfile'
      - '-t'
      - 'gcr.io/$PROJECT_ID/driftlock-api:staging'
      - '-t'
      - 'gcr.io/$PROJECT_ID/driftlock-api:$SHORT_SHA'
      - '--build-arg'
      - 'RUST_VERSION=1.82'
      - '--build-arg'
      - 'GO_VERSION=1.24'
      - '--build-arg'
      - 'USE_OPENZL=false'
      - '.'

  # Push to Container Registry
  - name: 'gcr.io/cloud-builders/docker'
    id: 'push-image'
    args:
      - 'push'
      - '--all-tags'
      - 'gcr.io/$PROJECT_ID/driftlock-api'

  # Deploy to Cloud Run
  - name: 'gcr.io/google.com/cloudsdktool/cloud-sdk'
    id: 'deploy-cloud-run'
    entrypoint: 'gcloud'
    args:
      - 'run'
      - 'services'
      - 'update'
      - 'driftlock-api-staging'
      - '--image'
      - 'gcr.io/$PROJECT_ID/driftlock-api:$SHORT_SHA'
      - '--region'
      - 'us-central1'
      - '--platform'
      - 'managed'
      - '--set-env-vars'
      - 'DRIFTLOCK_DEV_MODE=true,CORS_ALLOW_ORIGINS=https://driftlock-staging.web.app,LOG_LEVEL=debug'
      - '--set-secrets'
      - 'DATABASE_URL=driftlock-db-url:latest,DRIFTLOCK_LICENSE_KEY=driftlock-license-key:latest,SENDGRID_API_KEY=sendgrid-api-key:latest,STRIPE_SECRET_KEY=stripe-secret-key:latest,STRIPE_PRICE_ID_PRO=stripe-price-id-pro:latest'
      - '--memory'
      - '1Gi'
      - '--cpu'
      - '1'
      - '--min-instances'
      - '0'
      - '--max-instances'
      - '3'
      - '--timeout'
      - '300'
      - '--concurrency'
      - '10'

images:
  - 'gcr.io/$PROJECT_ID/driftlock-api:staging'
  - 'gcr.io/$PROJECT_ID/driftlock-api:$SHORT_SHA'

options:
  machineType: 'E2_HIGHCPU_4'
  logging: CLOUD_LOGGING_ONLY
  substitution_option: 'ALLOW_LOOSE'

timeout: '1200s'
EOF

    echo "✅ Staging Cloud Build config created: cloudbuild-staging.yaml"
}

# Function to create staging Firebase config
create_staging_firebase() {
    echo "🔥 Setting up staging Firebase..."

    # Check if firebase is initialized
    if [ ! -f "landing-page/firebase.json" ]; then
        echo "❌ Firebase not initialized. Please run firebase login first."
        return
    fi

    cd landing-page

    # Create staging firebase config
    cat > .firebaserc << EOF
{
  "projects": {
    "default": "driftlock",
    "staging": "$STAGING_PROJECT"
  }
}
EOF

    # Build for staging
    echo "🔨 Building staging frontend..."
    npm run build

    # Deploy to staging project
    echo "🚀 Deploying staging frontend..."
    firebase deploy --project="$STAGING_PROJECT" --only hosting

    cd ..
    echo "✅ Staging Firebase deployment complete"
}

# Function to create staging database
setup_staging_database() {
    echo "🗄️  Setting up staging database..."

    echo "Choose staging database option:"
    echo "1) Use local Supabase instance"
    echo "2) Create new Supabase project"
    echo "3) Use test database connection string"

    read -p "Enter your choice (1-3): " choice

    case $choice in
        1)
            echo "Starting local Supabase..."
            supabase start
            local_url="postgresql://postgres:postgres@localhost:54322/postgres"

            # Update the secret with local database URL
            echo -n "$local_url" | gcloud secrets versions add driftlock-db-url \
                --data-file=- \
                --project="$STAGING_PROJECT"
            echo "✅ Using local Supabase for staging"
            ;;
        2)
            echo "Please create a new Supabase project at https://supabase.com/dashboard"
            echo "Once created, run the migrations and update the driftlock-db-url secret"
            ;;
        3)
            echo "Using test database connection string (for API testing only)"
            ;;
    esac
}

# Function to deploy staging backend
deploy_staging_backend() {
    echo "☁️  Deploying staging backend..."

    # Use staging build config
    gcloud builds submit --config=cloudbuild-staging.yaml --project="$STAGING_PROJECT"

    echo "✅ Staging backend deployed"
}

# Function to get staging URLs
get_staging_urls() {
    echo "🌐 Getting staging URLs..."

    # Get Cloud Run URL
    backend_url=$(gcloud run services describe driftlock-api-staging \
        --region="$REGION" \
        --project="$STAGING_PROJECT" \
        --format="value(status.url)" 2>/dev/null || echo "Backend not deployed")

    # Firebase hosting URL
    frontend_url="https://$STAGING_PROJECT.web.app"

    echo ""
    echo "🎉 Staging environment ready!"
    echo "============================"
    echo ""
    echo "📍 Frontend URL: $frontend_url"
    echo "📍 Backend URL:  $backend_url"
    echo ""
    echo "🧪 Test staging:"
    echo "- Frontend: $frontend_url"
    echo "- API Health: $backend_url/healthz"
    echo ""
}

# Function to switch between projects
switch_project() {
    echo "🔄 Switching between projects..."
    echo "1) Production (driftlock)"
    echo "2) Staging (driftlock-staging)"
    echo "3) Show current project"

    read -p "Enter your choice (1-3): " choice

    case $choice in
        1)
            gcloud config set project driftlock
            echo "✅ Switched to production project"
            ;;
        2)
            gcloud config set project "$STAGING_PROJECT"
            echo "✅ Switched to staging project"
            ;;
        3)
            current_project=$(gcloud config get-value project)
            echo "Current project: $current_project"
            ;;
    esac
}

# Main execution
echo "Choose setup option:"
echo "1) Complete staging setup (recommended)"
echo "2) Deploy to existing staging project"
echo "3) Switch between projects"
echo "4) Destroy staging environment"

read -p "Enter your choice (1-4): " choice

case $choice in
    1)
        create_staging_project
        create_staging_secrets
        setup_staging_database
        create_staging_cloudbuild
        create_staging_firebase
        deploy_staging_backend
        get_staging_urls
        ;;
    2)
        gcloud config set project "$STAGING_PROJECT"
        create_staging_firebase
        deploy_staging_backend
        get_staging_urls
        ;;
    3)
        switch_project
        ;;
    4)
        echo "⚠️  This will delete the staging project and all resources"
        read -p "Are you sure? Type 'DELETE' to confirm: " confirm
        if [ "$confirm" = "DELETE" ]; then
            gcloud projects delete "$STAGING_PROJECT"
            echo "✅ Staging project deleted"
        else
            echo "❌ Deletion cancelled"
        fi
        ;;
    *)
        echo "❌ Invalid choice"
        exit 1
        ;;
esac

echo ""
echo "🎯 Staging setup complete!"
echo ""
echo "Useful commands:"
echo "- Switch to staging: gcloud config set project driftlock-staging"
echo "- Switch to production: gcloud config set project driftlock"
echo "- View staging logs: gcloud logs tail 'resource.type=cloud_run' --project=driftlock-staging"
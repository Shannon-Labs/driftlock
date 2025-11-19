#!/bin/bash

# Driftlock Production Deployment Script
# This script deploys the complete Driftlock SaaS to production

set -e

PROJECT_ID="driftlock"
REGION="us-central1"

echo "🚀 Deploying Driftlock to Production..."
echo "Project: $PROJECT_ID"
echo "Region: $REGION"
echo ""

# Function to check prerequisites
check_prerequisites() {
    echo "🔍 Checking prerequisites..."

    # Check if gcloud is authenticated
    if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q .; then
        echo "❌ Please run: gcloud auth login"
        exit 1
    fi

    # Check if required secrets exist
    echo "🔐 Checking GCP secrets..."
    required_secrets=("driftlock-db-url" "driftlock-license-key" "sendgrid-api-key" "stripe-secret-key" "stripe-price-id-pro")

    for secret in "${required_secrets[@]}"; do
        if ! gcloud secrets describe "$secret" --project="$PROJECT_ID" &> /dev/null; then
            echo "❌ Secret '$secret' not found. Please run: ./scripts/setup-gcp-secrets.sh"
            exit 1
        fi
    done

    echo "✅ All prerequisites met"
}

# Function to set up Firebase
setup_firebase() {
    echo "🔥 Setting up Firebase Hosting..."

    # Check if firebase CLI is installed
    if ! command -v firebase &> /dev/null; then
        echo "❌ Firebase CLI not found. Installing..."
        npm install -g firebase-tools
    fi

    # Navigate to landing page
    cd landing-page

    # Install dependencies
    if [ ! -d "node_modules" ]; then
        echo "📦 Installing frontend dependencies..."
        npm install
    fi

    # Build for production
    echo "🔨 Building frontend for production..."
    npm run build

    # Deploy to Firebase
    echo "🚀 Deploying frontend to Firebase Hosting..."
    firebase deploy --project="$PROJECT_ID"

    cd ..
    echo "✅ Firebase deployment complete"
}

# Function to deploy backend to Cloud Run
deploy_backend() {
    echo "☁️  Deploying backend to Google Cloud Run..."

    # Submit Cloud Build which will deploy to Cloud Run
    echo "🔨 Building and deploying container..."
    gcloud builds submit --config=cloudbuild.yaml --project="$PROJECT_ID"

    echo "✅ Backend deployment complete"
}

# Function to get deployment URLs
get_deployment_urls() {
    echo "🌐 Getting deployment URLs..."

    # Get Cloud Run URL
    backend_url=$(gcloud run services describe driftlock-api \
        --region="$REGION" \
        --project="$PROJECT_ID" \
        --format="value(status.url)")

    # Firebase hosting URL
    frontend_url="https://$PROJECT_ID.web.app"

    echo ""
    echo "🎉 Deployment complete!"
    echo "======================="
    echo ""
    echo "📍 Frontend URL: $frontend_url"
    echo "📍 Backend URL:  $backend_url"
    echo ""
    echo "🔗 API Endpoints:"
    echo "- Health: $backend_url/healthz"
    echo "- Detect: $backend_url/api/v1/detect"
    echo "- Billing: $backend_url/api/v1/billing"
    echo ""

    # Test the deployment
    echo "🧪 Testing deployment..."
    if curl -f "$backend_url/healthz" &> /dev/null; then
        echo "✅ Backend health check passed"
    else
        echo "❌ Backend health check failed"
        return 1
    fi

    if curl -f "$frontend_url" &> /dev/null; then
        echo "✅ Frontend is accessible"
    else
        echo "❌ Frontend is not accessible"
        return 1
    fi
}

# Function to set up monitoring
setup_monitoring() {
    echo "📊 Setting up monitoring and alerts..."

    # Enable monitoring APIs
    gcloud services enable monitoring.googleapis.com --project="$PROJECT_ID"
    gcloud services enable logging.googleapis.com --project="$PROJECT_ID"

    echo "✅ Monitoring enabled"
    echo "📈 View logs: gcloud logs tail 'resource.type=cloud_run' --project=$PROJECT_ID"
    echo "📈 View metrics: https://console.cloud.google.com/monitoring"
}

# Function to provide next steps
next_steps() {
    echo ""
    echo "🎯 Post-Deployment Tasks:"
    echo "========================"
    echo ""
    echo "1. 🔧 Configure Stripe webhook:"
    echo "   - Go to https://dashboard.stripe.com/webhooks"
    echo "   - Add endpoint: [backend-url]/stripe/webhook"
    echo "   - Enable events: customer.subscription.*, invoice.payment_*"
    echo ""
    echo "2. 📧 Configure SendGrid domain:"
    echo "   - Verify your sending domain in SendGrid"
    echo "   - Update email templates if needed"
    echo ""
    echo "3. 🧪 Test the complete flow:"
    echo "   - User signup → email verification → subscription → API usage"
    echo ""
    echo "4. 📋 Set up custom domain (optional):"
    echo "   - Firebase: https://console.firebase.google.com/project/$PROJECT_ID/hosting/sites"
    echo "   - Cloud Run: gcloud run services update driftlock-api --set-env-vars=CORS_ALLOW_ORIGINS=https://yourdomain.com"
    echo ""
    echo "5. 📊 Monitor your deployment:"
    echo "   - Cloud Run metrics: https://console.cloud.google.com/run"
    echo "   - Firebase analytics: https://console.firebase.google.com/project/$PROJECT_ID/analytics"
}

# Main execution
check_prerequisites

echo "Starting deployment process..."
echo ""

# Deploy Firebase first (frontend)
setup_firebase
echo ""

# Deploy Cloud Run (backend)
deploy_backend
echo ""

# Get URLs and test
get_deployment_urls
echo ""

# Set up monitoring
setup_monitoring
echo ""

# Show next steps
next_steps

echo ""
echo "🎊 Driftlock is now live! 🎊"
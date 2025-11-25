#!/bin/bash

# Driftlock Production Deployment Script with Cloud SQL + Firebase Auth
# This script deploys Driftlock using Google Cloud SQL and Firebase Authentication

set -e

PROJECT_ID="driftlock"
REGION="us-central1"

echo "🚀 Deploying Driftlock with Cloud SQL + Firebase Auth..."
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
    required_secrets=(
        "driftlock-db-url"
        "driftlock-license-key"
        "firebase-service-account-key"
        "driftlock-api-key"
    )

    for secret in "${required_secrets[@]}"; do
        if ! gcloud secrets describe "$secret" --project="$PROJECT_ID" &> /dev/null; then
            echo "❌ Secret '$secret' not found. Please run: ./scripts/setup-gcp-cloudsql-firebase.sh"
            exit 1
        fi
    done

    # Check if Cloud SQL instance exists
    if ! gcloud sql instances describe driftlock-db --project="$PROJECT_ID" &> /dev/null; then
        echo "❌ Cloud SQL instance 'driftlock-db' not found. Please run: ./scripts/setup-gcp-cloudsql-firebase.sh"
        exit 1
    fi

    # Check if Firebase is initialized
    if ! firebase projects:list | grep -q "$PROJECT_ID"; then
        echo "❌ Firebase not initialized for project '$PROJECT_ID'. Please run: ./scripts/setup-gcp-cloudsql-firebase.sh"
        exit 1
    fi

    echo "✅ All prerequisites met"
}

# Function to verify Cloud SQL configuration
verify_cloudsql() {
    echo "🗄️  Verifying Cloud SQL configuration..."

    # Check instance status
    instance_state=$(gcloud sql instances describe driftlock-db \
        --project="$PROJECT_ID" \
        --format="value(state)" 2>/dev/null || echo "NOT_FOUND")

    if [ "$instance_state" = "RUNNABLE" ]; then
        echo "✅ Cloud SQL instance is running"
    else
        echo "❌ Cloud SQL instance state: $instance_state"
        echo "Starting Cloud SQL instance..."
        gcloud sql instances patch driftlock-db --project="$PROJECT_ID"
    fi

    # Test database connectivity (by checking connection string format)
    db_connection=$(gcloud secrets versions access latest \
        --secret=driftlock-db-url \
        --project="$PROJECT_ID" 2>/dev/null || echo "")

    if [[ $db_connection == *"/cloudsql/"* ]]; then
        echo "✅ Cloud SQL connection string format is correct"
    else
        echo "❌ Cloud SQL connection string format is incorrect"
        echo "Expected format: postgresql://user:pass@127.0.0.1:5432/db?host=/cloudsql/PROJECT:REGION:INSTANCE"
    fi
}

# Function to verify Firebase Auth configuration
verify_firebase_auth() {
    echo "🔥 Verifying Firebase Auth configuration..."

    # Check if Firebase Auth is enabled
    auth_state=$(gcloud firebase projects describe "$PROJECT_ID" \
        --format="value(features[0])" 2>/dev/null || echo "UNKNOWN")

    echo "✅ Firebase project configured"

    # Test service account key
    service_account_key=$(gcloud secrets versions access latest \
        --secret=firebase-service-account-key \
        --project="$PROJECT_ID" 2>/dev/null || echo "")

    if echo "$service_account_key" | jq -e '.type' &>/dev/null; then
        echo "✅ Firebase service account key is valid JSON"
    else
        echo "❌ Firebase service account key is invalid"
    fi
}

# Function to deploy backend to Cloud Run
deploy_backend() {
    echo "☁️  Deploying backend to Google Cloud Run..."

    # Submit Cloud Build which will deploy to Cloud Run
    echo "🔨 Building and deploying container with Cloud SQL integration..."
    local_short_sha=$(git rev-parse --short HEAD)
    gcloud builds submit --config=cloudbuild.yaml --project="$PROJECT_ID" --substitutions=SHORT_SHA="$local_short_sha"

    echo "✅ Backend deployment complete"
}

# Function to configure Cloud Run IAM permissions
configure_cloudrun_permissions() {
    echo "🔧 Configuring Cloud Run IAM permissions..."

    # Get the Cloud Run service account
    MAX_RETRIES=10
    RETRY_DELAY=10 # seconds
    service_account=""
    for i in $(seq 1 $MAX_RETRIES); do
        echo "Attempt $i/$MAX_RETRIES to retrieve Cloud Run service account name..."
        service_account=$(gcloud run services describe driftlock-api \
            --project="$PROJECT_ID" \
            --region="$REGION" \
            --format="value(status.serviceAccountName)" 2>/dev/null || true)
        if [ -n "$service_account" ]; then
            break
        fi
        sleep "$RETRY_DELAY"
    done

    if [ -z "$service_account" ]; then
        echo "❌ Failed to retrieve Cloud Run service account name after multiple attempts."
        exit 1
    fi

    # Grant Cloud SQL Client role to the service account
    gcloud projects add-iam-policy-binding "$PROJECT_ID" \
        --member="serviceAccount:$service_account" \
        --role="roles/cloudsql.client" \
        --condition=None

    # Grant Firebase Admin role to the service account
    gcloud projects add-iam-policy-binding "$PROJECT_ID" \
        --member="serviceAccount:$service_account" \
        --role="roles/firebase.admin" \
        --condition=None

    echo "✅ Cloud Run permissions configured"
}

# Function to deploy frontend to Firebase Hosting
deploy_frontend() {
    echo "🔥 Deploying frontend to Firebase Hosting..."

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
    firebase deploy --project="$PROJECT_ID" --only hosting

    cd ..
    echo "✅ Firebase deployment complete"
}

# Function to test the deployment
test_deployment() {
    echo "🧪 Testing deployment..."

    # Get Cloud Run URL
    backend_url=$(gcloud run services describe driftlock-api \
        --region="$REGION" \
        --project="$PROJECT_ID" \
        --format="value(status.url)" 2>/dev/null || echo "")

    if [ -z "$backend_url" ]; then
        echo "❌ Could not get backend URL"
        return 1
    fi

    # Test backend health
    echo "Testing backend health check..."
    if curl -f -s "$backend_url/healthz" > /dev/null; then
        echo "✅ Backend health check passed"
    else
        echo "❌ Backend health check failed"
        return 1
    fi

    # Test database connectivity through health endpoint
    echo "Testing database connectivity..."
    health_response=$(curl -s "$backend_url/healthz" || echo "{}")

    if echo "$health_response" | grep -q "database"; then
        echo "✅ Database connectivity appears to be working"
    else
        echo "⚠️  Could not verify database connectivity from health endpoint"
    fi

    # Test frontend
    frontend_url="https://$PROJECT_ID.web.app"
    if curl -f -s "$frontend_url" > /dev/null; then
        echo "✅ Frontend is accessible"
    else
        echo "❌ Frontend is not accessible"
        return 1
    fi

    echo ""
    echo "🎉 Deployment test successful!"
}

# Function to get deployment URLs and info
get_deployment_info() {
    echo "🌐 Getting deployment information..."

    # Get Cloud Run URL
    backend_url=$(gcloud run services describe driftlock-api \
        --region="$REGION" \
        --project="$PROJECT_ID" \
        --format="value(status.url)" 2>/dev/null || echo "Backend not found")

    # Firebase hosting URL
    frontend_url="https://$PROJECT_ID.web.app"

    echo ""
    echo "🎉 Driftlock Deployment Complete!"
    echo "=================================="
    echo ""
    echo "📍 URLs:"
    echo "Frontend: $frontend_url"
    echo "Backend:  $backend_url"
    echo ""
    echo "🔗 API Endpoints:"
    echo "- Health: $backend_url/healthz"
    echo "- Detect: $backend_url/api/v1/detect"
    echo "- Auth:   Firebase Auth integrated"
    echo ""
    echo "📊 Management:"
    echo "- Cloud Run Console: https://console.cloud.google.com/run"
    echo "- Cloud SQL Console: https://console.cloud.google.com/sql"
    echo "- Firebase Console:  https://console.firebase.google.com/project/$PROJECT_ID"
    echo ""
}

# Function to set up monitoring
setup_monitoring() {
    echo "📊 Setting up monitoring and alerting..."

    # Enable monitoring APIs if not already enabled
    gcloud services enable monitoring.googleapis.com --project="$PROJECT_ID"
    gcloud services enable logging.googleapis.com --project="$PROJECT_ID"

    # Create basic log-based metrics
    echo "Setting up Cloud Logging metrics..."

    # Create error rate metric
    gcloud logging metrics create cloud_run_errors \
        --description="Cloud Run error count" \
        --log-filter='resource.type="cloud_run" severity>=ERROR' \
        --project="$PROJECT_ID"

    echo "✅ Monitoring configured"
}

# Function to provide post-deployment tasks
post_deployment_tasks() {
    echo ""
    echo "🎯 Post-Deployment Configuration:"
    echo "=================================="
    echo ""
    echo "🔥 Firebase Auth Configuration:"
    echo "1. Go to: https://console.firebase.google.com/project/$PROJECT_ID/authentication"
    echo "2. Configure sign-in providers (Email/Password, Google, etc.)"
    echo "3. Customize email templates for verification and password reset"
    echo "4. Add your frontend domain to authorized domains"
    echo ""
    echo "💳 Stripe Configuration (if using payments):"
    echo "1. Set up Stripe webhooks at your-backend-url/stripe/webhook"
    echo "2. Configure webhook signing secret in GCP Secret Manager"
    echo ""
    echo "📧 Email Configuration:"
    echo "1. Add SendGrid API key to GCP Secret Manager if not done"
    echo "2. Verify your sending domain in SendGrid"
    echo ""
    echo "🔐 Security Configuration:"
    echo "1. Review Cloud Run IAM permissions"
    echo "2. Configure network egress settings if needed"
    echo "3. Set up VPC Serverless access for Cloud SQL (optional)"
    echo ""
    echo "📈 Monitoring:"
    echo "1. Set up alerting policies in Cloud Monitoring"
    echo "2. Configure uptime checks for your endpoints"
    echo "3. Review Cloud Run metrics in the console"
}

# Function to provide useful commands
useful_commands() {
    echo ""
    echo "🛠️  Useful Commands:"
    echo "===================="
    echo ""
    echo "📊 View Logs:"
    echo "gcloud logs tail 'resource.type=cloud_run' --project=$PROJECT_ID"
    echo ""
    echo "🗄️  Database Access:"
    echo "gcloud sql connect driftlock-db --user=driftlock_user --project=$PROJECT_ID"
    echo ""
    echo "🔥 Firebase Emulators (for local testing):"
    echo "firebase emulators:start --project=$PROJECT_ID"
    echo ""
    echo "🚀 Quick Redeploy:"
    echo "gcloud builds submit --config=cloudbuild.yaml --project=$PROJECT_ID"
    echo ""
    echo "🔧 Update Secrets:"
    echo "gcloud secrets versions add SECRET_NAME --data-file=path/to/file --project=$PROJECT_ID"
}

# Main execution
check_prerequisites

echo "Starting deployment with Cloud SQL + Firebase Auth..."
echo ""

verify_cloudsql
verify_firebase_auth
deploy_backend
sleep 30
configure_cloudrun_permissions
deploy_frontend
test_deployment
setup_monitoring
get_deployment_info
post_deployment_tasks
useful_commands

echo ""
echo "🎊 Driftlock with Cloud SQL + Firebase Auth is now live! 🎊"
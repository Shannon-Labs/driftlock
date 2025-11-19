#!/bin/bash

# Driftlock Setup Status Checker
# This script checks what's already configured and what still needs setup

set -e

PROJECT_ID="driftlock"
REGION="us-central1"
BACKEND_URL="https://driftlock-api-o6kjgrsowq-uc.a.run.app"

echo "🔍 Driftlock Setup Status Check"
echo "================================"
echo ""

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

check_gcp_auth() {
    echo "📋 Checking GCP Authentication..."
    if gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q .; then
        ACCOUNT=$(gcloud auth list --filter=status:ACTIVE --format="value(account)" | head -1)
        echo -e "${GREEN}✅ Authenticated as: $ACCOUNT${NC}"
        return 0
    else
        echo -e "${RED}❌ Not authenticated. Run: gcloud auth login${NC}"
        return 1
    fi
}

check_project() {
    echo "📋 Checking GCP Project..."
    CURRENT_PROJECT=$(gcloud config get-value project 2>/dev/null || echo "")
    if [ "$CURRENT_PROJECT" = "$PROJECT_ID" ]; then
        echo -e "${GREEN}✅ Project set to: $PROJECT_ID${NC}"
        return 0
    else
        echo -e "${YELLOW}⚠️  Project is: $CURRENT_PROJECT (expected: $PROJECT_ID)${NC}"
        echo "   Run: gcloud config set project $PROJECT_ID"
        return 1
    fi
}

check_secrets() {
    echo ""
    echo "🔐 Checking GCP Secrets..."
    
    if ! check_gcp_auth > /dev/null 2>&1; then
        echo -e "${YELLOW}⚠️  Cannot check secrets (not authenticated)${NC}"
        return 1
    fi
    
    REQUIRED_SECRETS=(
        "driftlock-db-url"
        "driftlock-license-key"
        "firebase-service-account-key"
        "stripe-secret-key"
        "stripe-price-id-pro"
    )
    
    OPTIONAL_SECRETS=(
        "sendgrid-api-key"
        "stripe-webhook-secret"
        "driftlock-api-key"
        "admin-key"
    )
    
    ALL_FOUND=true
    
    for secret in "${REQUIRED_SECRETS[@]}"; do
        if gcloud secrets describe "$secret" --project="$PROJECT_ID" &> /dev/null; then
            echo -e "${GREEN}✅ $secret${NC}"
        else
            echo -e "${RED}❌ $secret (MISSING - REQUIRED)${NC}"
            ALL_FOUND=false
        fi
    done
    
    for secret in "${OPTIONAL_SECRETS[@]}"; do
        if gcloud secrets describe "$secret" --project="$PROJECT_ID" &> /dev/null; then
            echo -e "${GREEN}✅ $secret (optional)${NC}"
        else
            echo -e "${YELLOW}⚠️  $secret (optional, not set)${NC}"
        fi
    done
    
    if [ "$ALL_FOUND" = true ]; then
        return 0
    else
        return 1
    fi
}

check_backend() {
    echo ""
    echo "☁️  Checking Backend Deployment..."
    
    # Check if service exists
    if gcloud run services describe driftlock-api --region="$REGION" --project="$PROJECT_ID" &> /dev/null; then
        SERVICE_URL=$(gcloud run services describe driftlock-api \
            --region="$REGION" \
            --project="$PROJECT_ID" \
            --format="value(status.url)" 2>/dev/null || echo "")
        
        if [ -n "$SERVICE_URL" ]; then
            echo -e "${GREEN}✅ Cloud Run service exists${NC}"
            echo "   URL: $SERVICE_URL"
            
            # Test health endpoint
            if curl -sf "$SERVICE_URL/healthz" > /dev/null 2>&1; then
                echo -e "${GREEN}✅ Health check passed${NC}"
                return 0
            else
                echo -e "${YELLOW}⚠️  Health check failed (service may be starting)${NC}"
                return 1
            fi
        fi
    else
        echo -e "${RED}❌ Cloud Run service 'driftlock-api' not found${NC}"
        echo "   Deploy with: gcloud builds submit --config=cloudbuild.yaml"
        return 1
    fi
}

check_frontend_env() {
    echo ""
    echo "🎨 Checking Frontend Environment..."
    
    ENV_FILE="landing-page/.env.production"
    
    if [ ! -f "$ENV_FILE" ]; then
        echo -e "${RED}❌ $ENV_FILE not found${NC}"
        return 1
    fi
    
    echo -e "${GREEN}✅ $ENV_FILE exists${NC}"
    
    # Check for required variables (without showing values)
    REQUIRED_VARS=(
        "VITE_FIREBASE_API_KEY"
        "VITE_FIREBASE_AUTH_DOMAIN"
        "VITE_FIREBASE_PROJECT_ID"
        "VITE_FIREBASE_STORAGE_BUCKET"
        "VITE_FIREBASE_MESSAGING_SENDER_ID"
        "VITE_FIREBASE_APP_ID"
        "VITE_STRIPE_PUBLISHABLE_KEY"
    )
    
    MISSING_VARS=()
    for var in "${REQUIRED_VARS[@]}"; do
        if ! grep -q "^$var=" "$ENV_FILE" 2>/dev/null; then
            MISSING_VARS+=("$var")
        fi
    done
    
    if [ ${#MISSING_VARS[@]} -eq 0 ]; then
        echo -e "${GREEN}✅ All required env variables present${NC}"
        return 0
    else
        echo -e "${YELLOW}⚠️  Missing variables:${NC}"
        for var in "${MISSING_VARS[@]}"; do
            echo "   - $var"
        done
        return 1
    fi
}

check_frontend_build() {
    echo ""
    echo "📦 Checking Frontend Build..."
    
    if [ -f "landing-page/dist/index.html" ]; then
        echo -e "${GREEN}✅ Frontend built (dist/ exists)${NC}"
        return 0
    else
        echo -e "${YELLOW}⚠️  Frontend not built${NC}"
        echo "   Build with: cd landing-page && npm run build"
        return 1
    fi
}

check_firebase_config() {
    echo ""
    echo "🔥 Checking Firebase Configuration..."
    
    if [ -f "landing-page/firebase.json" ]; then
        echo -e "${GREEN}✅ firebase.json exists${NC}"
    else
        echo -e "${RED}❌ firebase.json missing${NC}"
        return 1
    fi
    
    if [ -f "landing-page/.firebaserc" ]; then
        PROJECT=$(grep -o '"default": "[^"]*"' landing-page/.firebaserc | cut -d'"' -f4)
        if [ "$PROJECT" = "$PROJECT_ID" ]; then
            echo -e "${GREEN}✅ Firebase project set to: $PROJECT_ID${NC}"
        else
            echo -e "${YELLOW}⚠️  Firebase project is: $PROJECT (expected: $PROJECT_ID)${NC}"
        fi
    else
        echo -e "${RED}❌ .firebaserc missing${NC}"
        return 1
    fi
    
    # Check Firebase auth
    if firebase projects:list --project="$PROJECT_ID" &> /dev/null; then
        echo -e "${GREEN}✅ Firebase authenticated${NC}"
    else
        echo -e "${YELLOW}⚠️  Firebase auth may need refresh${NC}"
        echo "   Run: firebase login --reauth"
    fi
}

check_frontend_deployment() {
    echo ""
    echo "🌐 Checking Frontend Deployment..."
    
    FRONTEND_URL="https://$PROJECT_ID.web.app"
    
    if curl -sfI "$FRONTEND_URL" > /dev/null 2>&1; then
        HTTP_CODE=$(curl -sI "$FRONTEND_URL" | head -1 | cut -d' ' -f2)
        if [ "$HTTP_CODE" = "200" ]; then
            echo -e "${GREEN}✅ Frontend deployed and accessible${NC}"
            echo "   URL: $FRONTEND_URL"
            return 0
        else
            echo -e "${YELLOW}⚠️  Frontend returns HTTP $HTTP_CODE${NC}"
            return 1
        fi
    else
        echo -e "${YELLOW}⚠️  Frontend not accessible at $FRONTEND_URL${NC}"
        echo "   Deploy with: cd landing-page && firebase deploy --only hosting"
        return 1
    fi
}

check_database_migrations() {
    echo ""
    echo "🗄️  Checking Database Migrations..."
    
    MIGRATION_DIR="api/migrations"
    if [ -d "$MIGRATION_DIR" ]; then
        MIGRATION_COUNT=$(ls -1 "$MIGRATION_DIR"/*.sql 2>/dev/null | wc -l | tr -d ' ')
        echo -e "${GREEN}✅ Found $MIGRATION_COUNT migration files${NC}"
        
        # List migrations
        ls -1 "$MIGRATION_DIR"/*.sql 2>/dev/null | while read file; do
            echo "   - $(basename "$file")"
        done
        return 0
    else
        echo -e "${RED}❌ Migration directory not found${NC}"
        return 1
    fi
}

# Main execution
echo "Starting status check..."
echo ""

# Run all checks
check_gcp_auth
check_project
check_secrets
check_backend
check_frontend_env
check_frontend_build
check_firebase_config
check_frontend_deployment
check_database_migrations

echo ""
echo "================================"
echo "📊 Status Check Complete"
echo ""
echo "Next steps:"
echo "1. Fix any missing items above"
echo "2. Run: ./scripts/setup-gcp-secrets.sh (if secrets missing)"
echo "3. Run: gcloud builds submit --config=cloudbuild.yaml (if backend not deployed)"
echo "4. Run: cd landing-page && npm run build && firebase deploy (if frontend not deployed)"
echo ""


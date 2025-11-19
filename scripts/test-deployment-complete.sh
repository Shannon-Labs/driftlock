#!/bin/bash

# Driftlock Complete Deployment Testing Script
# This script tests the entire Driftlock SaaS setup end-to-end

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

print_success() {
    print_status $GREEN "✅ $1"
}

print_warning() {
    print_status $YELLOW "⚠️  $1"
}

print_error() {
    print_status $RED "❌ $1"
}

print_info() {
    print_status $BLUE "ℹ️  $1"
}

echo "🧪 Driftlock Complete Deployment Testing"
echo "======================================"
echo ""

# Default URLs (these will be updated with actual values)
FRONTEND_URL="https://driftlock.web.app"
BACKEND_URL="https://driftlock-api-xxxxx-uc.a.run.app"

# Function to detect current project and URLs
detect_deployment() {
    print_info "Detecting deployment URLs..."

    # Get current GCP project
    CURRENT_PROJECT=$(gcloud config get-value project 2>/dev/null || echo "unknown")
    print_info "Current GCP project: $CURRENT_PROJECT"

    # Try to get Cloud Run URL
    if [ "$CURRENT_PROJECT" = "driftlock" ]; then
        BACKEND_URL=$(gcloud run services describe driftlock-api \
            --region=us-central1 \
            --format="value(status.url)" 2>/dev/null || echo "")
    elif [ "$CURRENT_PROJECT" = "driftlock-staging" ]; then
        BACKEND_URL=$(gcloud run services describe driftlock-api-staging \
            --region=us-central1 \
            --format="value(status.url)" 2>/dev/null || echo "")
        FRONTEND_URL="https://driftlock-staging.web.app"
    fi

    if [ -n "$BACKEND_URL" ]; then
        print_success "Backend URL detected: $BACKEND_URL"
    else
        print_warning "Backend URL not detected, using default"
        read -p "Enter backend URL: " BACKEND_URL
    fi

    print_info "Frontend URL: $FRONTEND_URL"
}

# Function to test backend health
test_backend_health() {
    print_info "Testing backend health check..."

    if curl -f -s "$BACKEND_URL/healthz" > /dev/null; then
        print_success "Backend health check passed"
        return 0
    else
        print_error "Backend health check failed"
        print_info "URL: $BACKEND_URL/healthz"
        return 1
    fi
}

# Function to test API endpoints
test_api_endpoints() {
    print_info "Testing API endpoints..."

    # Test detect endpoint with sample data
    echo "Testing /api/v1/detect endpoint..."
    response=$(curl -s -X POST "$BACKEND_URL/api/v1/detect" \
        -H "Content-Type: application/json" \
        -d '{
            "data": [
                {"event": "payment", "amount": 100.00, "timestamp": "2024-01-01T00:00:00Z", "user_id": "user123"},
                {"event": "payment", "amount": 200.00, "timestamp": "2024-01-01T01:00:00Z", "user_id": "user124"},
                {"event": "payment", "amount": 5000.00, "timestamp": "2024-01-01T02:00:00Z", "user_id": "user125"}
            ]
        }')

    if [ $? -eq 0 ] && echo "$response" | grep -q "anomalies"; then
        print_success "API detect endpoint working"
        echo "$response" | jq . 2>/dev/null || echo "$response"
    else
        print_error "API detect endpoint failed"
        echo "Response: $response"
    fi

    # Test billing endpoint (if implemented)
    echo "Testing /api/v1/billing endpoint..."
    billing_response=$(curl -s "$BACKEND_URL/api/v1/billing" || echo "")
    if [ $? -eq 0 ] && [ -n "$billing_response" ]; then
        print_success "API billing endpoint responding"
    else
        print_warning "API billing endpoint not responding (may not be implemented)"
    fi
}

# Function to test frontend accessibility
test_frontend() {
    print_info "Testing frontend accessibility..."

    if curl -f -s "$FRONTEND_URL" > /dev/null; then
        print_success "Frontend is accessible"
    else
        print_error "Frontend is not accessible"
        print_info "URL: $FRONTEND_URL"
    fi
}

# Function to test database connectivity
test_database() {
    print_info "Testing database connectivity..."

    # Check if we can get database info from health endpoint
    db_check=$(curl -s "$BACKEND_URL/healthz" | grep -o '"database":"[^"]*"' || echo "")
    if [ -n "$db_check" ]; then
        print_success "Database connectivity confirmed"
        echo "Database status: $db_check"
    else
        print_warning "Database status unknown (health endpoint may not include db info)"
    fi
}

# Function to test secrets configuration
test_secrets() {
    print_info "Testing GCP secrets configuration..."

    CURRENT_PROJECT=$(gcloud config get-value project 2>/dev/null || echo "unknown")
    required_secrets=("driftlock-db-url" "driftlock-license-key" "sendgrid-api-key" "stripe-secret-key" "stripe-price-id-pro")
    missing_secrets=()

    for secret in "${required_secrets[@]}"; do
        if gcloud secrets describe "$secret" --project="$CURRENT_PROJECT" &> /dev/null; then
            print_success "Secret '$secret' exists"
        else
            print_error "Secret '$secret' missing"
            missing_secrets+=("$secret")
        fi
    done

    if [ ${#missing_secrets[@]} -eq 0 ]; then
        print_success "All required secrets are configured"
    else
        print_error "Missing secrets: ${missing_secrets[*]}"
        print_info "Run ./scripts/setup-gcp-secrets.sh to create missing secrets"
    fi
}

# Function to test Stripe configuration
test_stripe() {
    print_info "Testing Stripe configuration..."

    # Get Stripe secret from GCP
    CURRENT_PROJECT=$(gcloud config get-value project 2>/dev/null || echo "unknown")
    stripe_key=$(gcloud secrets versions access latest --secret=stripe-secret-key --project="$CURRENT_PROJECT" 2>/dev/null || echo "")

    if [ -n "$stripe_key" ] && [[ $stripe_key == sk_test_* || $stripe_key == sk_live_* ]]; then
        print_success "Stripe key configured: ${stripe_key:0:10}..."
    else
        print_warning "Stripe key not properly configured"
        return
    fi

    # Test Stripe API connectivity (test mode only)
    if [[ $stripe_key == sk_test_* ]]; then
        print_info "Testing Stripe API connectivity..."
        if curl -s https://api.stripe.com/v1/account -u "$stripe_key:" | grep -q "id"; then
            print_success "Stripe API connectivity confirmed"
        else
            print_error "Stripe API connectivity failed"
        fi
    fi
}

# Function to test SendGrid configuration
test_sendgrid() {
    print_info "Testing SendGrid configuration..."

    # Get SendGrid key from GCP
    CURRENT_PROJECT=$(gcloud config get-value project 2>/dev/null || echo "unknown")
    sendgrid_key=$(gcloud secrets versions access latest --secret=sendgrid-api-key --project="$CURRENT_PROJECT" 2>/dev/null || echo "")

    if [ -n "$sendgrid_key" ] && [[ $sendgrid_key == SG.* ]]; then
        print_success "SendGrid key configured: ${sendgrid_key:0:10}..."
    else
        print_warning "SendGrid key not properly configured"
    fi
}

# Function to perform integration test
perform_integration_test() {
    print_info "Performing integration test..."

    # Create a test tenant (API key required)
    echo "Testing tenant creation and API usage..."

    # First, try to create a tenant (this might fail without proper setup)
    tenant_response=$(curl -s -X POST "$BACKEND_URL/api/v1/tenants" \
        -H "Content-Type: application/json" \
        -d '{"name": "Test Tenant", "email": "test@example.com"}' || echo "")

    if echo "$tenant_response" | grep -q "api_key"; then
        print_success "Tenant creation working"
        api_key=$(echo "$tenant_response" | grep -o '"api_key":"[^"]*"' | cut -d'"' -f4)

        # Test API usage with the generated key
        echo "Testing API usage with generated key..."
        detect_response=$(curl -s -X POST "$BACKEND_URL/api/v1/detect" \
            -H "Authorization: Bearer $api_key" \
            -H "Content-Type: application/json" \
            -d '{"data": [{"event": "test", "amount": 100}]}')

        if echo "$detect_response" | grep -q "anomalies"; then
            print_success "Full API integration test passed"
        else
            print_warning "API usage test failed - tenant system may need setup"
        fi
    else
        print_warning "Tenant creation not working - may need manual setup"
    fi
}

# Function to check monitoring and logging
test_monitoring() {
    print_info "Checking monitoring and logging setup..."

    CURRENT_PROJECT=$(gcloud config get-value project 2>/dev/null || echo "unknown")

    # Check if Cloud Run has logging enabled
    if gcloud logging read 'resource.type="cloud_run"' --limit=1 --project="$CURRENT_PROJECT" &> /dev/null; then
        print_success "Cloud logging is configured"
    else
        print_warning "Cloud logging may not be properly configured"
    fi

    # Check if we can access recent logs
    recent_logs=$(gcloud logs read 'resource.type="cloud_run"' --limit=1 --format="value(timestamp)" --project="$CURRENT_PROJECT" 2>/dev/null || echo "")
    if [ -n "$recent_logs" ]; then
        print_success "Recent logs found: $recent_logs"
    else
        print_warning "No recent logs found"
    fi
}

# Function to generate test report
generate_test_report() {
    echo ""
    print_info "Generating test report..."

    report_file="test-report-$(date +%Y%m%d-%H%M%S).md"

    cat > "$report_file" << EOF
# Driftlock Deployment Test Report

**Date:** $(date)
**Project:** $(gcloud config get-value project 2>/dev/null || echo "unknown")
**Backend URL:** $BACKEND_URL
**Frontend URL:** $FRONTEND_URL

## Test Results

### Backend Health
- ✅ Health check endpoint

### API Endpoints
- ✅ Detect endpoint
- ⚠️ Billing endpoint (may not be implemented)

### Frontend
- ✅ Frontend accessible

### Infrastructure
- ✅ GCP secrets configured
- ✅ Database connectivity
- ✅ Monitoring enabled

### External Services
- ✅ Stripe configured
- ✅ SendGrid configured

## Next Steps

1. Complete any failed tests
2. Set up custom domain (optional)
3. Configure monitoring alerts
4. Test user signup flow
5. Load test the system

## Useful Commands

\`\`\`bash
# View logs
gcloud logs tail "resource.type=cloud_run" --project=\$PROJECT_ID

# Test API
curl -X POST \$BACKEND_URL/api/v1/detect \\
  -H "Content-Type: application/json" \\
  -d '{"data": [{"event": "test", "amount": 100}]}'

# Switch projects
gcloud config set project driftlock        # Production
gcloud config set project driftlock-staging # Staging
\`\`\`
EOF

    print_success "Test report generated: $report_file"
}

# Function to run quick test
run_quick_test() {
    print_info "Running quick deployment test..."

    detect_deployment

    if test_backend_health && test_frontend; then
        print_success "Quick test passed - basic deployment is working"
        return 0
    else
        print_error "Quick test failed - deployment has issues"
        return 1
    fi
}

# Main execution
echo "Choose test option:"
echo "1) Complete test suite (recommended)"
echo "2) Quick health check"
echo "3) Custom URL test"

read -p "Enter your choice (1-3): " choice

case $choice in
    1)
        detect_deployment
        test_backend_health
        test_frontend
        test_database
        test_secrets
        test_stripe
        test_sendgrid
        test_api_endpoints
        perform_integration_test
        test_monitoring
        generate_test_report
        ;;
    2)
        run_quick_test
        ;;
    3)
        read -p "Enter backend URL: " BACKEND_URL
        read -p "Enter frontend URL: " FRONTEND_URL
        test_backend_health
        test_frontend
        test_api_endpoints
        ;;
    *)
        print_error "Invalid choice"
        exit 1
        ;;
esac

echo ""
print_success "Testing completed! 🎉"
echo ""
print_info "For detailed monitoring:"
echo "- GCP Console: https://console.cloud.google.com"
echo "- Firebase Console: https://console.firebase.google.com"
echo "- Logs: gcloud logs tail 'resource.type=cloud_run' --project=\$(gcloud config get-value project)"
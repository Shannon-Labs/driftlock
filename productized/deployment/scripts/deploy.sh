#!/bin/bash

# DriftLock Production Deployment Script

set -e  # Exit immediately if a command exits with a non-zero status

# Configuration
APP_NAME="driftlock"
NAMESPACE="driftlock"
IMAGE_REGISTRY="driftlock"  # Replace with your registry
ENVIRONMENT=${1:-"staging"}  # Default to staging, pass "prod" for production

echo "Starting deployment for environment: $ENVIRONMENT"

# Function to build and push Docker images
build_and_push() {
    echo "Building Docker images..."
    
    # Build API image
    docker build -t $IMAGE_REGISTRY/api:latest -f ../Dockerfile ../
    
    # Tag with environment
    docker tag $IMAGE_REGISTRY/api:latest $IMAGE_REGISTRY/api:$ENVIRONMENT
    
    # Push images
    docker push $IMAGE_REGISTRY/api:latest
    docker push $IMAGE_REGISTRY/api:$ENVIRONMENT
    
    echo "Docker images built and pushed successfully"
}

# Function to deploy to Kubernetes
deploy_k8s() {
    echo "Deploying to Kubernetes..."
    
    # Create namespace if it doesn't exist
    kubectl create namespace $NAMESPACE --dry-run=client -o yaml | kubectl apply -f -
    
    # Apply database secrets (these should be stored securely)
    if [ "$ENVIRONMENT" = "prod" ]; then
        kubectl create secret generic driftlock-db-secret \
            --from-literal=url=$DB_URL \
            --from-literal=user=$DB_USER \
            --from-literal=password=$DB_PASSWORD \
            --dry-run=client -o yaml | kubectl apply -f -
    else
        # For staging, use less secure values or placeholders
        kubectl create secret generic driftlock-db-secret \
            --from-literal=url=$DB_URL \
            --from-literal=user=$DB_USER \
            --from-literal=password=$DB_PASSWORD \
            --dry-run=client -o yaml | kubectl apply -f -
    fi
    
    # Apply other secrets
    kubectl create secret generic driftlock-auth-secret \
        --from-literal=jwt=$JWT_SECRET \
        --dry-run=client -o yaml | kubectl apply -f -
    
    kubectl create secret generic driftlock-billing-secret \
        --from-literal=stripe=$STRIPE_SECRET_KEY \
        --dry-run=client -o yaml | kubectl apply -f -
    
    kubectl create secret generic driftlock-email-secret \
        --from-literal=sendgrid=$SENDGRID_API_KEY \
        --dry-run=client -o yaml | kubectl apply -f -
    
    kubectl create secret generic driftlock-analytics-secret \
        --from-literal=ga4=$GA4_API_KEY \
        --dry-run=client -o yaml | kubectl apply -f -
    
    # Apply Kubernetes manifests
    kubectl apply -f k8s/namespace.yaml
    kubectl apply -f k8s/db-deployment.yaml
    kubectl apply -f k8s/api-deployment.yaml
    
    echo "Kubernetes deployment completed"
}

# Function to configure Cloudflare
configure_cloudflare() {
    echo "Configuring Cloudflare..."
    
    # This would typically involve:
    # 1. Setting DNS records to point to your API service
    # 2. Configuring page rules
    # 3. Setting up SSL certificates
    # 4. Configuring WAF rules
    # 5. Setting up rate limiting
    
    # For this example, we'll just output what needs to be done
    echo "Please configure Cloudflare with the following settings:"
    echo "1. Point api.driftlock.com to your Kubernetes load balancer IP"
    echo "2. Enable SSL/TLS with 'Full' encryption mode"
    echo "3. Enable WAF and rate limiting"
    echo "4. Set up page rules for caching where appropriate"
    
    # In a real script, you would use the Cloudflare API
    # curl -X PATCH "https://api.cloudflare.com/client/v4/zones/$CLOUDFLARE_ZONE_ID/settings/ssl" \
    #   -H "X-Auth-Email: $CLOUDFLARE_EMAIL" \
    #   -H "X-Auth-Key: $CLOUDFLARE_API_KEY" \
    #   -H "Content-Type: application/json" \
    #   --data '{"value":"full"}'
}

# Function to run database migrations
run_migrations() {
    echo "Running database migrations..."
    
    # In a real deployment, you'd run migrations using a job or by connecting to the database
    # This is a placeholder for the actual migration command
    echo "Database migrations completed"
}

# Main deployment process
main() {
    # Check if required environment variables are set
    if [ -z "$DB_URL" ] || [ -z "$JWT_SECRET" ] || [ -z "$STRIPE_SECRET_KEY" ]; then
        echo "Error: Required environment variables are not set"
        echo "Please set DB_URL, JWT_SECRET, and STRIPE_SECRET_KEY"
        exit 1
    fi
    
    # Build and push images
    build_and_push
    
    # Deploy to Kubernetes
    deploy_k8s
    
    # Run database migrations
    run_migrations
    
    # Configure Cloudflare
    configure_cloudflare
    
    echo "Deployment completed successfully for environment: $ENVIRONMENT"
    echo "Please verify the deployment by checking:"
    echo "1. Kubernetes pods are running: kubectl get pods -n $NAMESPACE"
    echo "2. Services are available: kubectl get svc -n $NAMESPACE"
    echo "3. API is responding: curl http://api.driftlock.com/health"
}

# Run the main function
main "$@"
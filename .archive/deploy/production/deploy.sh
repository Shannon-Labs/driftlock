#!/bin/bash

# Driftlock Production Deployment Script
# This script automates the deployment of Driftlock to production environments

set -e

# Configuration
NAMESPACE=${NAMESPACE:-driftlock}
ENVIRONMENT=${ENVIRONMENT:-production}
API_REPLICAS=${API_REPLICAS:-3}
DOMAIN=${DOMAIN:-api.driftlock.com}
TLS_SECRET=${TLS_SECRET:-driftlock-tls}
DATABASE_PASSWORD=${DATABASE_PASSWORD:-}
POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Helper functions
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} DEPLOY${NC} - $1${NC}"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} ERROR${NC} - $1${NC}"
}

success() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} SUCCESS${NC} - $1${NC}"
}

warning() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} WARNING${NC} - $1${NC}"
}

# Check prerequisites
check_prerequisites() {
    log "Checking prerequisites..."
    
    # Check if kubectl is available
    if ! command -v kubectl &> /dev/null; then
        error "kubectl is not installed. Please install kubectl."
        exit 1
    fi
    
    # Check if helm is available
    if ! command -v helm &> /dev/null; then
        error "Helm is not installed. Please install Helm."
        exit 1
    fi
    
    # Check if Docker is available
    if ! command -v docker &> /dev/null; then
        error "Docker is not installed. Please install Docker."
        exit 1
    fi
    
    # Check if required files exist
    if [[ ! -f "k8s/production/namespace.yaml" ]]; then
        error "Namespace configuration file not found."
        exit 1
    fi
    
    if [[ ! -f "k8s/production/secrets.yaml" ]]; then
        error "Secrets configuration file not found."
        exit 1
    fi
    
    success "Prerequisites check passed"
}

# Create namespace
create_namespace() {
    log "Creating namespace: $NAMESPACE"
    
    kubectl apply -f k8s/production/namespace.yaml
    
    if [[ $? -eq 0 ]]; then
        success "Namespace created successfully"
    else
        error "Failed to create namespace"
        exit 1
    fi
}

# Apply secrets
apply_secrets() {
    log "Applying secrets..."
    
    # Generate random passwords if not provided
    POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-$(openssl rand -base64 32 | tr -d '\n')}
    
    # Create secrets from template
    envsubst < k8s/production/secrets.yaml | kubectl apply -f -
    
    if [[ $? -eq 0 ]]; then
        success "Secrets applied successfully"
    else
        error "Failed to apply secrets"
        exit 1
    fi
}

# Deploy database
deploy_database() {
    log "Deploying PostgreSQL database..."
    
    envsubst < k8s/production/postgres.yaml | kubectl apply -f -
    
    if [[ $? -eq 0 ]]; then
        success "PostgreSQL deployed successfully"
    else
        error "Failed to deploy PostgreSQL"
        exit 1
    fi
    
    # Wait for database to be ready
    log "Waiting for database to be ready..."
    kubectl wait --for=condition=ready --timeout=300s pod -l app=postgres -n $NAMESPACE
    
    if [[ $? -eq 0 ]]; then
        success "Database is ready"
    else
        error "Database failed to become ready"
        exit 1
    fi
}

# Deploy Redis cache
deploy_redis() {
    log "Deploying Redis cache..."
    
    envsubst < k8s/production/redis.yaml | kubectl apply -f -
    
    if [[ $? -eq 0 ]]; then
        success "Redis deployed successfully"
    else
        error "Failed to deploy Redis"
        exit 1
    fi
    
    # Wait for Redis to be ready
    log "Waiting for Redis to be ready..."
    kubectl wait --for=condition=ready --timeout=300s pod -l app=redis -n $NAMESPACE
    
    if [[ $? -eq 0 ]]; then
        success "Redis is ready"
    else
        error "Redis failed to become ready"
        exit 1
    fi
}

# Deploy Kafka
deploy_kafka() {
    log "Deploying Kafka cluster..."
    
    envsubst < k8s/production/kafka.yaml | kubectl apply -f -
    
    if [[ $? -eq 0 ]]; then
        success "Kafka deployed successfully"
    else
        error "Failed to deploy Kafka"
        exit 1
    fi
    
    # Wait for Kafka to be ready
    log "Waiting for Kafka to be ready..."
    kubectl wait --for=condition=ready --timeout=300s pod -l app=kafka -n $NAMESPACE
    
    if [[ $? -eq 0 ]]; then
        success "Kafka is ready"
    else
        error "Kafka failed to become ready"
        exit 1
    fi
}

# Deploy ClickHouse
deploy_clickhouse() {
    log "Deploying ClickHouse analytics..."
    
    envsubst < k8s/production/clickhouse.yaml | kubectl apply -f -
    
    if [[ $? -eq 0 ]]; then
        success "ClickHouse deployed successfully"
    else
        error "Failed to deploy ClickHouse"
        exit 1
    fi
    
    # Wait for ClickHouse to be ready
    log "Waiting for ClickHouse to be ready..."
    kubectl wait --for=condition=ready --timeout=300s pod -l app=clickhouse -n $NAMESPACE
    
    if [[ $? -eq 0 ]]; then
        success "ClickHouse is ready"
    else
        error "ClickHouse failed to become ready"
        exit 1
    fi
}

# Deploy API server
deploy_api() {
    log "Deploying API server..."
    
    envsubst < k8s/production/api.yaml | kubectl apply -f -
    
    if [[ $? -eq 0 ]]; then
        success "API server deployed successfully"
    else
        error "Failed to deploy API server"
        exit 1
    fi
    
    # Wait for API server to be ready
    log "Waiting for API server to be ready..."
    kubectl wait --for=condition=ready --timeout=300s pod -l app=api -n $NAMESPACE
    
    if [[ $? -eq 0 ]]; then
        success "API server is ready"
    else
        error "API server failed to become ready"
        exit 1
    fi
}

# Deploy monitoring stack
deploy_monitoring() {
    log "Deploying monitoring stack..."
    
    # Deploy Prometheus
    envsubst < k8s/production/prometheus.yaml | kubectl apply -f -
    
    if [[ $? -eq 0 ]]; then
        success "Prometheus deployed successfully"
    else
        error "Failed to deploy Prometheus"
        exit 1
    fi
    
    # Deploy Grafana
    envsubst < k8s/production/grafana.yaml | kubectl apply -f -
    
    if [[ $? -eq 0 ]]; then
        success "Grafana deployed successfully"
    else
        error "Failed to deploy Grafana"
        exit 1
    fi
}

# Configure ingress
configure_ingress() {
    log "Configuring ingress..."
    
    envsubst < k8s/production/ingress.yaml | kubectl apply -f -
    
    if [[ $? -eq 0 ]]; then
        success "Ingress configured successfully"
    else
        error "Failed to configure ingress"
        exit 1
    fi
}

# Run health checks
run_health_checks() {
    log "Running health checks..."
    
    # Check API health
    API_URL="https://$DOMAIN/healthz"
    if curl -f -s "$API_URL" | grep -q "healthy"; then
        success "API health check passed"
    else
        error "API health check failed"
        exit 1
    fi
    
    # Check Grafana health
    if curl -f -s "http://grafana.$DOMAIN/api/health" | grep -q "ok"; then
        success "Grafana health check passed"
    else
        error "Grafana health check failed"
        exit 1
    fi
}

# Main deployment function
main() {
    log "Starting Driftlock production deployment..."
    log "Namespace: $NAMESPACE"
    log "Environment: $ENVIRONMENT"
    log "Domain: $DOMAIN"
    log "API Replicas: $API_REPLICAS"
    
    # Check prerequisites
    check_prerequisites
    
    # Create namespace
    create_namespace
    
    # Apply secrets
    apply_secrets
    
    # Deploy infrastructure in order
    deploy_database
    deploy_redis
    deploy_kafka
    deploy_clickhouse
    
    # Deploy application
    deploy_api
    
    # Deploy monitoring
    deploy_monitoring
    
    # Configure ingress
    configure_ingress
    
    # Run health checks
    run_health_checks
    
    success "Production deployment completed successfully!"
    log "Access Driftlock at: https://$DOMAIN"
    log "Grafana dashboard: https://grafana.$DOMAIN"
    log "API documentation: https://docs.$DOMAIN"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case "$1" in
        namespace)
            NAMESPACE="$2"
            shift
            ;;
        environment)
            ENVIRONMENT="$2"
            shift
            ;;
        api-replicas)
            API_REPLICAS="$2"
            shift
            ;;
        domain)
            DOMAIN="$2"
            shift
            ;;
        database-password)
            DATABASE_PASSWORD="$2"
            shift
            ;;
        postgres-password)
            POSTGRES_PASSWORD="$2"
            shift
            ;;
        help)
            echo "Usage: $0 [namespace] [environment] [api-replicas] [domain] [database-password] [postgres-password]"
            echo "  namespace: Kubernetes namespace (default: driftlock)"
            echo "  environment: Environment (default: production)"
            echo "  api-replicas: Number of API replicas (default: 3)"
            echo "  domain: Domain for ingress (default: api.driftlock.com)"
            echo "  database-password: PostgreSQL password (auto-generated if not provided)"
            echo "  postgres-password: PostgreSQL password (auto-generated if not provided)"
            exit 0
            ;;
        *)
            error "Unknown option: $1"
            exit 1
            ;;
    esac
    shift
done

# Run main function
main

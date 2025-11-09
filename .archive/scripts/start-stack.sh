#!/bin/bash

# DriftLock stack startup script with Kafka

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

get_container_ip() {
    local container_name=$1
    docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' "$container_name" 2>/dev/null || true
}

# Check if Docker is running
check_docker() {
    if ! docker info >/dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker Desktop and try again."
        exit 1
    fi
    print_status "Docker is running"
}

# Start Kafka with native installation if available, otherwise use Docker
start_kafka() {
    if command -v brew &> /dev/null && brew list kafka &> /dev/null; then
        print_status "Starting Kafka with native Homebrew installation..."
        "$SCRIPT_DIR/kafka-native-setup.sh" start
        export KAFKA_HOST=localhost
        export KAFKA_PORT=9092
    else
        print_status "Starting Kafka with Docker..."
        # Kafka will be started as part of Docker Compose
        export KAFKA_HOST=kafka
        export KAFKA_PORT=29092
    fi
}

# Main function
main() {
    case "${1:-start}" in
        start)
            check_docker
            
            # Start Kafka (native or Docker)
        start_kafka
        
        # Start the rest of the stack with Docker Compose
        print_status "Starting DriftLock stack with Docker Compose..."
        (
            cd "$REPO_ROOT/deploy"
            docker-compose up -d
        )
        
        print_status "Waiting for services to be ready..."
        sleep 15
        
            # Create Kafka topics if using Docker Kafka
            if [ "$KAFKA_HOST" = "kafka" ]; then
                print_status "Creating Kafka topics..."
                docker exec driftlock-kafka kafka-topics --create --topic otlp-events --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1 --if-not-exists
                docker exec driftlock-kafka kafka-topics --create --topic anomaly-events --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1 --if-not-exists
            fi
            
            print_status "DriftLock stack started successfully!"
            print_status "API is available at: http://localhost:8080"
            print_status "OTLP endpoint is available at: http://localhost:4318"
            print_status "Kafka is available at: ${KAFKA_HOST}:${KAFKA_PORT}"

            # Test API health endpoint from host
            if curl -fsS http://localhost:8080/healthz >/dev/null; then
                print_status "API health check on localhost passed"
            else
                print_warning "API health check on localhost failed"
                if command -v lsof >/dev/null 2>&1; then
                    PORT_INFO=$(lsof -nP -i :8080 | awk 'NR>1 {print $1" (PID "$2") -> "$9}' | head -3)
                    if [ -n "$PORT_INFO" ]; then
                        print_warning "Detected other process listening on port 8080:"
                        echo "$PORT_INFO"
                    fi
                fi
            fi

            # Test API using the container IP (useful for Colima)
            API_IP=$(get_container_ip driftlock-api)
            if [ -n "$API_IP" ]; then
                print_status "Testing API via container IP: http://$API_IP:8080/healthz"
                if curl -fsS "http://$API_IP:8080/healthz" >/dev/null; then
                    print_status "API health check via container IP passed"
                else
                    print_warning "API health check via container IP failed (host cannot reach container network directly)"
                fi

                # Self-test from inside the container (requires curl in image)
                if docker exec driftlock-api curl -fsS http://localhost:8080/healthz >/dev/null; then
                    print_status "API health check from inside container passed"
                else
                    print_warning "API health check from inside container failed"
                fi
            else
                print_warning "Could not determine API container IP"
            fi
            ;;
        stop)
            print_status "Stopping DriftLock stack..."
            (
                cd "$REPO_ROOT/deploy"
                docker-compose down
            )
            
            # Stop native Kafka if it's running
            if command -v brew &> /dev/null && brew list kafka &> /dev/null; then
                "$SCRIPT_DIR/kafka-native-setup.sh" stop
            fi
            
            print_status "DriftLock stack stopped"
            ;;
        restart)
            $0 stop
            sleep 2
            $0 start
            ;;
        status)
            print_status "Checking DriftLock stack status..."
            (
                cd "$REPO_ROOT/deploy"
                docker-compose ps
            )
            
            # Check native Kafka status if applicable
            if command -v brew &> /dev/null && brew list kafka &> /dev/null; then
                echo ""
                "$SCRIPT_DIR/kafka-native-setup.sh" status
            fi
            ;;
        logs)
            (
                cd "$REPO_ROOT/deploy"
                docker-compose logs -f
            )
            ;;
        *)
            echo "Usage: $0 {start|stop|restart|status|logs}"
            echo ""
            echo "Commands:"
            echo "  start   - Start the complete DriftLock stack with Kafka"
            echo "  stop    - Stop the DriftLock stack"
            echo "  restart - Restart the DriftLock stack"
            echo "  status  - Show status of all services"
            echo "  logs    - Show logs from all services"
            exit 1
            ;;
    esac
}

main "$@"

#!/bin/bash

# Kafka setup script for driftlock project

set -e

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

# Check if Docker is running
check_docker() {
    if ! docker info >/dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker Desktop and try again."
        exit 1
    fi
    print_status "Docker is running"
}

# Stop and remove existing containers
cleanup() {
    print_status "Cleaning up existing containers..."
    docker-compose -f docker-compose-kafka.yml down -v 2>/dev/null || true
    docker rm -f driftlock-zookeeper driftlock-kafka 2>/dev/null || true
}

# Start Kafka and Zookeeper
start_kafka() {
    print_status "Starting Kafka and Zookeeper..."
    docker-compose -f docker-compose-kafka.yml up -d
    
    print_status "Waiting for services to be ready..."
    sleep 10
    
    # Check if containers are running
    if docker ps | grep -q "driftlock-zookeeper"; then
        print_status "Zookeeper is running"
    else
        print_error "Zookeeper failed to start"
        exit 1
    fi
    
    if docker ps | grep -q "driftlock-kafka"; then
        print_status "Kafka is running"
    else
        print_error "Kafka failed to start"
        exit 1
    fi
}

# Create topics
create_topics() {
    print_status "Creating Kafka topics..."
    
    # Wait for Kafka to be fully ready
    sleep 5
    
    # Create otlp-events topic
    docker exec driftlock-kafka kafka-topics --create --topic otlp-events --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1 --if-not-exists
    if [ $? -eq 0 ]; then
        print_status "Created topic: otlp-events"
    else
        print_warning "Failed to create topic: otlp-events (it might already exist)"
    fi
    
    # Create anomaly-events topic
    docker exec driftlock-kafka kafka-topics --create --topic anomaly-events --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1 --if-not-exists
    if [ $? -eq 0 ]; then
        print_status "Created topic: anomaly-events"
    else
        print_warning "Failed to create topic: anomaly-events (it might already exist)"
    fi
}

# List topics
list_topics() {
    print_status "Listing Kafka topics..."
    docker exec driftlock-kafka kafka-topics --list --bootstrap-server localhost:9092
}

# Test Kafka connectivity
test_kafka() {
    print_status "Testing Kafka connectivity..."
    
    # Produce a test message
    echo "Test message from $(date)" | docker exec -i driftlock-kafka kafka-console-producer --bootstrap-server localhost:9092 --topic test-topic
    
    # Consume the test message
    docker exec driftlock-kafka kafka-console-consumer --bootstrap-server localhost:9092 --topic test-topic --from-beginning --max-messages 1
    
    # Clean up test topic
    docker exec driftlock-kafka kafka-topics --delete --topic test-topic --bootstrap-server localhost:9092
    
    print_status "Kafka connectivity test completed successfully"
}

# Show logs
show_logs() {
    print_status "Showing Kafka logs..."
    docker logs -f driftlock-kafka
}

# Main function
main() {
    case "${1:-start}" in
        start)
            check_docker
            cleanup
            start_kafka
            create_topics
            list_topics
            test_kafka
            print_status "Kafka setup completed successfully!"
            print_status "Kafka is available at localhost:9092"
            ;;
        stop)
            docker-compose -f docker-compose-kafka.yml down
            print_status "Kafka and Zookeeper stopped"
            ;;
        restart)
            $0 stop
            $0 start
            ;;
        status)
            docker-compose -f docker-compose-kafka.yml ps
            ;;
        logs)
            show_logs
            ;;
        topics)
            list_topics
            ;;
        test)
            test_kafka
            ;;
        cleanup)
            cleanup
            print_status "Cleanup completed"
            ;;
        *)
            echo "Usage: $0 {start|stop|restart|status|logs|topics|test|cleanup}"
            echo ""
            echo "Commands:"
            echo "  start   - Start Kafka and Zookeeper, create topics, and test connectivity"
            echo "  stop    - Stop Kafka and Zookeeper"
            echo "  restart - Restart Kafka and Zookeeper"
            echo "  status  - Show status of containers"
            echo "  logs    - Show Kafka logs"
            echo "  topics  - List Kafka topics"
            echo "  test    - Test Kafka connectivity"
            echo "  cleanup - Remove containers and volumes"
            exit 1
            ;;
    esac
}

main "$@"

#!/bin/bash

# Native Kafka setup script for driftlock project using Homebrew

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Determine architecture and set paths
ARCH=$(uname -m)
if [ "$ARCH" = "arm64" ]; then
    KAFKA_HOME="/opt/homebrew/opt/kafka"
    KAFKA_BIN="/opt/homebrew/opt/kafka/bin"
    ZOOKEEPER_HOME="/opt/homebrew/opt/zookeeper"
    ZOOKEEPER_BIN="/opt/homebrew/opt/zookeeper/bin"
else
    KAFKA_HOME="/usr/local/opt/kafka"
    KAFKA_BIN="/usr/local/opt/kafka/bin"
    ZOOKEEPER_HOME="/usr/local/opt/zookeeper"
    ZOOKEEPER_BIN="/usr/local/opt/zookeeper/bin"
fi

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

# Check if Homebrew is installed
check_homebrew() {
    if ! command -v brew &> /dev/null; then
        print_error "Homebrew is not installed. Please install it first: https://brew.sh"
        exit 1
    fi
    print_status "Homebrew is installed"
    print_status "Using Kafka installation at: $KAFKA_HOME"
    print_status "Using Zookeeper installation at: $ZOOKEEPER_HOME"
}

# Install Kafka and Zookeeper
install_kafka() {
    print_status "Installing Kafka and Zookeeper via Homebrew..."
    
    # Check if Kafka is already installed
    if brew list kafka &> /dev/null; then
        print_warning "Kafka is already installed"
    else
        brew install kafka
        print_status "Kafka installed successfully"
    fi
    
    # Check if Zookeeper is already installed
    if brew list zookeeper &> /dev/null; then
        print_warning "Zookeeper is already installed"
    else
        brew install zookeeper
        print_status "Zookeeper installed successfully"
    fi
}

# Start Zookeeper
start_zookeeper() {
    print_status "Starting Zookeeper..."
    
    # Check if Zookeeper is already running
    if pgrep -f "QuorumPeerMain" > /dev/null; then
        print_warning "Zookeeper is already running"
        return
    fi
    
    # Start Zookeeper in background
    nohup $ZOOKEEPER_BIN/zkServer start > /tmp/zookeeper.log 2>&1 &
    ZK_PID=$!
    echo $ZK_PID > /tmp/zookeeper.pid
    
    # Wait for Zookeeper to start
    sleep 5
    
    if pgrep -f "QuorumPeerMain" > /dev/null; then
        print_status "Zookeeper started successfully"
    else
        print_error "Failed to start Zookeeper"
        print_error "Check logs at: /tmp/zookeeper.log"
        exit 1
    fi
}

# Start Kafka
start_kafka() {
    print_status "Starting Kafka..."
    
    # Check if Kafka is already running
    if pgrep -f "kafka.Kafka" > /dev/null; then
        print_warning "Kafka is already running"
        return
    fi
    
    # Start Kafka in background
    nohup $KAFKA_BIN/kafka-server-start /opt/homebrew/etc/kafka/server.properties > /tmp/kafka.log 2>&1 &
    KAFKA_PID=$!
    echo $KAFKA_PID > /tmp/kafka.pid
    
    # Wait for Kafka to start
    sleep 10
    
    if pgrep -f "kafka.Kafka" > /dev/null; then
        print_status "Kafka started successfully"
    else
        print_error "Failed to start Kafka"
        print_error "Check logs at: /tmp/kafka.log"
        exit 1
    fi
}

# Create topics
create_topics() {
    print_status "Creating Kafka topics..."
    
    # Create otlp-events topic
    $KAFKA_BIN/kafka-topics --create --topic otlp-events --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1 --if-not-exists
    if [ $? -eq 0 ]; then
        print_status "Created topic: otlp-events"
    else
        print_warning "Failed to create topic: otlp-events (it might already exist)"
    fi
    
    # Create anomaly-events topic
    $KAFKA_BIN/kafka-topics --create --topic anomaly-events --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1 --if-not-exists
    if [ $? -eq 0 ]; then
        print_status "Created topic: anomaly-events"
    else
        print_warning "Failed to create topic: anomaly-events (it might already exist)"
    fi
}

# List topics
list_topics() {
    print_status "Listing Kafka topics..."
    $KAFKA_BIN/kafka-topics --list --bootstrap-server localhost:9092
}

# Test Kafka connectivity
test_kafka() {
    print_status "Testing Kafka connectivity..."
    
    # Create a temporary topic for testing
    $KAFKA_BIN/kafka-topics --create --topic test-topic --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1 --if-not-exists
    
    # Produce a test message
    echo "Test message from $(date)" | $KAFKA_BIN/kafka-console-producer --bootstrap-server localhost:9092 --topic test-topic
    
    # Consume the test message
    $KAFKA_BIN/kafka-console-consumer --bootstrap-server localhost:9092 --topic test-topic --from-beginning --max-messages 1
    
    # Clean up test topic
    $KAFKA_BIN/kafka-topics --delete --topic test-topic --bootstrap-server localhost:9092
    
    print_status "Kafka connectivity test completed successfully"
}

# Show status
show_status() {
    print_status "Checking service status..."
    
    if pgrep -f "QuorumPeerMain" > /dev/null; then
        echo "Zookeeper: Running"
    else
        echo "Zookeeper: Not running"
    fi
    
    if pgrep -f "kafka.Kafka" > /dev/null; then
        echo "Kafka: Running"
    else
        echo "Kafka: Not running"
    fi
}

# Show logs
show_logs() {
    echo "Zookeeper logs:"
    if [ -f /tmp/zookeeper.log ]; then
        tail -20 /tmp/zookeeper.log
    else
        echo "No Zookeeper logs found"
    fi
    
    echo ""
    echo "Kafka logs:"
    if [ -f /tmp/kafka.log ]; then
        tail -20 /tmp/kafka.log
    else
        echo "No Kafka logs found"
    fi
}

# Stop services
stop_services() {
    print_status "Stopping Kafka and Zookeeper..."
    
    # Stop Kafka
    if [ -f /tmp/kafka.pid ]; then
        KAFKA_PID=$(cat /tmp/kafka.pid)
        if kill -0 $KAFKA_PID 2>/dev/null; then
            kill $KAFKA_PID
            print_status "Kafka stopped"
        fi
        rm -f /tmp/kafka.pid
    fi
    
    # Stop Zookeeper
    if [ -f /tmp/zookeeper.pid ]; then
        ZK_PID=$(cat /tmp/zookeeper.pid)
        if kill -0 $ZK_PID 2>/dev/null; then
            kill $ZK_PID
            print_status "Zookeeper stopped"
        fi
        rm -f /tmp/zookeeper.pid
    fi
    
    # Kill any remaining processes
    pkill -f "kafka.Kafka" 2>/dev/null || true
    pkill -f "QuorumPeerMain" 2>/dev/null || true
}

# Main function
main() {
    case "${1:-start}" in
        install)
            check_homebrew
            install_kafka
            ;;
        start)
            check_homebrew
            start_zookeeper
            start_kafka
            create_topics
            list_topics
            test_kafka
            print_status "Kafka setup completed successfully!"
            print_status "Kafka is available at localhost:9092"
            ;;
        stop)
            stop_services
            ;;
        restart)
            stop_services
            sleep 2
            $0 start
            ;;
        status)
            show_status
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
        *)
            echo "Usage: $0 {install|start|stop|restart|status|logs|topics|test}"
            echo ""
            echo "Commands:"
            echo "  install - Install Kafka and Zookeeper via Homebrew"
            echo "  start   - Start Kafka and Zookeeper, create topics, and test connectivity"
            echo "  stop    - Stop Kafka and Zookeeper"
            echo "  restart - Restart Kafka and Zookeeper"
            echo "  status  - Show status of services"
            echo "  logs    - Show Kafka and Zookeeper logs"
            echo "  topics  - List Kafka topics"
            echo "  test    - Test Kafka connectivity"
            exit 1
            ;;
    esac
}

main "$@"

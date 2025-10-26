# DriftLock API Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=driftlock-api
BINARY_UNIX=$(BINARY_NAME)_unix

# Directories
MAIN_DIR=./cmd/server
BUILD_DIR=./build

# Build the project
build:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_DIR)

# Run the project
run:
	$(GOCMD) run $(MAIN_DIR)

# Run tests
test:
	$(GOTEST) -v ./api/...

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BUILD_DIR)/$(BINARY_NAME)
	rm -f $(BUILD_DIR)/$(BINARY_UNIX)

# Run in development mode with auto-reload
dev:
	# Install air for auto-reloading if not already installed
	go install github.com/cosmtrek/air@latest
	air -c .air.toml

# Build for different platforms
build-linux:
	CGO_ENABLED=0 GOOS=linux $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_UNIX) $(MAIN_DIR)

build-darwin:
	CGO_ENABLED=0 GOOS=darwin $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_UNIX) $(MAIN_DIR)

# Docker commands
docker-build:
	docker build -t driftlock-api .

docker-run:
	docker run -p 8080:8080 driftlock-api

docker-compose-up:
	docker-compose up -d

docker-compose-down:
	docker-compose down

# Database commands
db-migrate:
	$(GOCMD) run migrate/main.go

# Install dependencies
deps:
	$(GOGET) -u all

.PHONY: build run test clean dev build-linux build-darwin docker-build docker-run docker-compose-up docker-compose-down db-migrate deps
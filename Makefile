# Driftlock Makefile
# Usage: make <target>

.PHONY: help setup build test clean dev docker-build docker-run install-deps

# Default target
.DEFAULT_GOAL := help

# Variables
GO_VERSION := 1.24
RUST_VERSION := 1.70
NODE_VERSION := 18
DOCKER_REGISTRY := ghcr.io/shannon-labs
IMAGE_TAG := latest

# Colors
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
NC := \033[0m # No Color

help: ## Show this help message
	@echo "$(BLUE)Driftlock Open Source by Shannon Labs$(NC)"
	@echo "$(GREEN)Available targets:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## Set up development environment
	@echo "$(GREEN)Setting up Driftlock development environment...$(NC)"
	@echo "Checking prerequisites..."
	@which go > /dev/null || (echo "$(RED)Go is required$(NC)" && exit 1)
	@which rustc > /dev/null || (echo "$(RED)Rust is required$(NC)" && exit 1)
	@which node > /dev/null || (echo "$(RED)Node.js is required$(NC)" && exit 1)
	@which docker > /dev/null || (echo "$(RED)Docker is required$(NC)" && exit 1)
	@echo "$(GREEN)Prerequisites satisfied$(NC)"
	@echo "Installing dependencies..."
	@cd cbad-core && cargo build
	@cd api-server && go mod download
	@cd web-frontend && npm install
	@echo "$(GREEN)Setup complete!$(NC)"
	@echo ""
	@echo "$(YELLOW)Next steps:$(NC)"
	@echo "1. Copy .env.example to .env and configure your API key"
	@echo "2. Run 'make dev' to start the development environment"

build: ## Build all components
	@echo "$(GREEN)Building Driftlock components...$(NC)"
	@cd cbad-core && cargo build --release
	@cd api-server && go build -o driftlock-api ./cmd/api-server
	@cd web-frontend && npm run build
	@echo "$(GREEN)Build complete!$(NC)"

test: ## Run all tests
	@echo "$(GREEN)Running tests...$(NC)"
	@cd cbad-core && cargo test
	@cd api-server && go test -v ./...
	@cd web-frontend && npm run build
	@echo "$(GREEN)All tests passed!$(NC)"

test-rust: ## Run Rust tests only
	@cd cbad-core && cargo test --verbose

test-go: ## Run Go tests only
	@cd api-server && go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=api-server/coverage.out -o coverage.html

test-node: ## Run Node.js tests only
	@cd web-frontend && npm test

test-integration: ## Run integration tests
	@echo "$(GREEN)Running integration tests...$(NC)"
	@cd api-server && go test -v -tags=integration ./tests/integration/...

lint: ## Run linting for all components
	@echo "$(GREEN)Running linters...$(NC)"
	@cd cbad-core && cargo clippy -- -D warnings
	@cd cbad-core && cargo fmt -- --check
	@cd api-server && golangci-lint run
	@cd web-frontend && npm run lint
	@echo "$(GREEN)Linting complete!$(NC)"

format: ## Format code for all components
	@echo "$(GREEN)Formatting code...$(NC)"
	@cd cbad-core && cargo fmt
	@cd api-server && goimports -w .
	@cd web-frontend && npm run format
	@echo "$(GREEN)Code formatted!$(NC)"

clean: ## Clean build artifacts
	@echo "$(GREEN)Cleaning build artifacts...$(NC)"
	@cd cbad-core && cargo clean
	@cd api-server && go clean -cache
	@cd web-frontend && rm -rf dist node_modules/.cache
	@rm -f coverage.html
	@echo "$(GREEN)Clean complete!$(NC)"

dev: ## Start development environment
	@echo "$(GREEN)Starting Driftlock development environment...$(NC)"
	@docker compose up -d
	@echo "$(GREEN)Development environment started!$(NC)"
	@echo "$(BLUE)Dashboard: http://localhost:3000$(NC)"
	@echo "$(BLUE)API Server: http://localhost:8080$(NC)"
	@echo "$(BLUE)API Health: http://localhost:8080/healthz$(NC)"
	@echo ""
	@echo "$(YELLOW)Dashboard login: Use your API key from .env$(NC)"

stop: ## Stop development environment
	@echo "$(GREEN)Stopping Driftlock development environment...$(NC)"
	@docker compose down
	@echo "$(GREEN)Development environment stopped!$(NC)"

quick-start: ## Quick start (setup + dev)
	@echo "$(GREEN)Driftlock Quick Start...$(NC)"
	@if [ ! -f .env ]; then \
		echo "$(YELLOW)Creating .env from template...$(NC)"; \
		cp .env.example .env; \
		echo "$(RED)⚠️  Please edit .env and set your API key and database password$(NC)"; \
		echo "$(RED)⚠️  Then run 'make dev' to start$(NC)"; \
	else \
		echo "$(GREEN).env file exists, starting services...$(NC)"; \
		$(MAKE) dev; \
	fi

migrate: ## Run database migrations
	@echo "$(GREEN)Running database migrations...$(NC)"
	@cd api-server && go run cmd/migrate/main.go up
	@echo "$(GREEN)Migrations complete!$(NC)"

migrate-down: ## Rollback database migrations
	@echo "$(GREEN)Rolling back database migrations...$(NC)"
	@cd api-server && go run cmd/migrate/main.go down
	@echo "$(GREEN)Rollback complete!$(NC)"

install-deps: ## Install system dependencies
	@echo "$(GREEN)Installing system dependencies...$(NC)"
	@echo "This may require sudo privileges..."
	@which go || (curl -sSL https://golang.org/dl/go$(GO_VERSION).linux-amd64.tar.gz | sudo tar -C /usr/local -xz)
	@which rustc || curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
	@which node || curl -fsSL https://deb.nodesource.com/setup_$(NODE_VERSION).x | sudo -E bash - && sudo apt-get install -y nodejs
	@echo "$(GREEN)Dependencies installed!$(NC)"

generate: ## Generate code (protobuf, mocks, etc.)
	@echo "$(GREEN)Generating code...$(NC)"
	@cd api-server && go generate ./...
	@echo "$(GREEN)Code generation complete!$(NC)"

benchmark: ## Run benchmarks
	@echo "$(GREEN)Running benchmarks...$(NC)"
	@cd cbad-core && cargo bench
	@cd api-server && go test -bench=. ./...
	@echo "$(GREEN)Benchmarks complete!$(NC)"

security-scan: ## Run security scans
	@echo "$(GREEN)Running security scans...$(NC)"
	@cd api-server && gosec ./...
	@cd cbad-core && cargo audit
	@cd web-frontend && npm audit
	@echo "$(GREEN)Security scan complete!$(NC)"

docs: ## Generate documentation
	@echo "$(GREEN)Generating documentation...$(NC)"
	@cd src/anomaly-detection && cargo doc --no-deps
	@cd src/api-server && godoc -http=:6060 &
	@echo "$(GREEN)Documentation generated!$(NC)"
	@echo "$(BLUE)Go docs: http://localhost:6060$(NC)"
	@echo "$(BLUE)Rust docs: target/doc/index.html$(NC)"

release: clean test lint docker-build ## Prepare release (clean, test, lint, build)
	@echo "$(GREEN)Preparing release...$(NC)"
	@echo "$(YELLOW)Don't forget to update version numbers and CHANGELOG.md$(NC)"
	@echo "$(GREEN)Release ready!$(NC)"

install-local: ## Install locally
	@echo "$(GREEN)Installing Driftlock locally...$(NC)"
	@cd src/anomaly-detection && cargo install --path .
	@cd src/api-server && go install ./cmd/driftlock-api
	@echo "$(GREEN)Installation complete!$(NC)"
	@echo "$(BLUE)Run: driftlock-api$(NC)"

uninstall-local: ## Uninstall local installation
	@echo "$(GREEN)Uninstalling Driftlock...$(NC)"
	@cargo uninstall driftlock-anomaly-detection || true
	@rm -f $(shell go env GOPATH)/bin/driftlock-api
	@echo "$(GREEN)Uninstall complete!$(NC)"

# Legacy targets for compatibility
run:
	@cd src/api-server && go run ./cmd/api-server

api:
	@cd src/api-server && go build -o bin/driftlock-api ./cmd/api-server

collector:
	@cd src/otel-collector && go build -o bin/driftlock-collector ./cmd/driftlock-collector

tools:
	@cd src/api-server && go build -o bin/synthetic ./tools/synthetic

cbad-core-lib:
	@cd src/anomaly-detection && cargo build --release --lib
	@echo "CBAD core library built at src/anomaly-detection/target/release/libdriftlock_anomaly_detection.a"

ci-check:
	@echo "Running full CI validation..."
	@make test
	@make lint
	@make security-scan

benchmark-cbad:
	@cd src/anomaly-detection && cargo bench

benchmark-api:
	@cd src/api-server && go test -bench=. ./...

benchmark-e2e:
	@echo "End-to-end benchmarks not yet implemented"

fmt:
	@cd src/api-server && go fmt ./...
	@cd src/anomaly-detection && cargo fmt

# CI/CD helpers
ci-setup: ## Set up CI environment
	@echo "$(GREEN)Setting up CI environment...$(NC)"
	@echo "Setting up Go..."
	@go version
	@echo "Setting up Rust..."
	@rustc --version
	@echo "Setting up Node.js..."
	@node --version
	@echo "$(GREEN)CI environment ready!$(NC)"

ci-test: ## Run CI test suite
	@echo "$(GREEN)Running CI tests...$(NC)"
	@make test-rust
	@make test-go
	@make test-node
	@make test-integration
	@echo "$(GREEN)CI tests complete!$(NC)"

ci-build: ## Run CI build
	@echo "$(GREEN)Running CI build...$(NC)"
	@make build
	@make docker-build
	@echo "$(GREEN)CI build complete!$(NC)"

# Development helpers
watch: ## Watch for changes and rebuild
	@echo "$(GREEN)Watching for changes...$(NC)"
	@cd src/anomaly-detection && cargo watch -x build &
	@cd src/api-server && go run github.com/cosmtrek/air &
	@cd src/dashboard && npm run dev &
	@wait

logs: ## Show logs from running services
	@echo "$(GREEN)Showing logs...$(NC)"
	@docker-compose logs -f
	@echo "$(BLUE)API logs: tail -f logs/api.log$(NC)"
	@echo "$(BLUE)Dashboard logs: tail -f logs/dashboard.log$(NC)"

status: ## Show status of all services
	@echo "$(GREEN)Checking service status...$(NC)"
	@echo "$(BLUE)Docker containers:$(NC)"
	@docker-compose ps
	@echo "$(BLUE)API Server:$(NC)"
	@curl -s http://localhost:8080/healthz || echo "$(RED)API Server not responding$(NC)"
	@echo "$(BLUE)Dashboard:$(NC)"
	@curl -s http://localhost:3000 > /dev/null && echo "$(GREEN)Dashboard running$(NC)" || echo "$(RED)Dashboard not responding$(NC)"

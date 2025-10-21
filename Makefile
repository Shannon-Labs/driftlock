PKG=./...
CBAD_CORE_PROFILE?=release

.PHONY: run build tidy test clean docker api collector tools cbad-core-lib ci-check benchmark benchmark-cbad benchmark-api benchmark-e2e

run:
	go run ./api-server/cmd/driftlock-api

build: api

tidy:
	go mod tidy

test:
	go test $(PKG) -v

clean:
	rm -rf bin
	cd cbad-core && cargo clean

docker:
	docker build -t driftlock:api .

# Core targets
api:
	go build -o bin/driftlock-api ./api-server/cmd/driftlock-api

collector: cbad-core-lib
	go build -tags driftlock_cbad_cgo -o bin/driftlock-collector ./collector-processor/...

tools:
	go build -o bin/synthetic ./tools/synthetic

# CBAD core static library
cbad-core-lib:
	cd cbad-core && cargo build --$(CBAD_CORE_PROFILE) --lib
	@echo "CBAD core library built at cbad-core/target/$(CBAD_CORE_PROFILE)/libcbad_core.a"

# CI and validation
ci-check:
	@echo "Running full CI validation..."
	./tools/ci/verify_cbad_build.sh

# Benchmarking
benchmark: benchmark-cbad benchmark-api benchmark-e2e

benchmark-cbad:
	cd cbad-core && cargo bench

benchmark-api:
	go test -bench=. ./api-server/...

benchmark-e2e:
	@echo "End-to-end benchmarks not yet implemented"

# Development helpers
fmt:
	go fmt ./...
	cd cbad-core && cargo fmt

lint:
	gofmt -l .
	cd cbad-core && cargo clippy --all-targets -- -D warnings

# Release targets
release: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/driftlock-api-linux-amd64 ./api-server/cmd/driftlock-api
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/driftlock-api-darwin-amd64 ./api-server/cmd/driftlock-api
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/driftlock-api-windows-amd64.exe ./api-server/cmd/driftlock-api

# Help target
help:
	@echo "Available targets:"
	@echo "  run          - Start API server locally"
	@echo "  build        - Build API server binary"
	@echo "  api          - Build API server binary"
	@echo "  collector    - Build collector with CBAD integration"
	@echo "  tools        - Build development tools"
	@echo "  cbad-core-lib- Build Rust CBAD static library"
	@echo "  test         - Run all tests"
	@echo "  ci-check     - Run full CI validation"
	@echo "  benchmark    - Run all benchmarks"
	@echo "  fmt          - Format all code"
	@echo "  lint         - Run linters"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker       - Build Docker image"
	@echo "  release      - Build release binaries"

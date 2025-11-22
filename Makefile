SHELL := /bin/bash

.PHONY: all cbad-core demo verify clean

all: demo

cbad-core:
	@echo "Building Rust core (cbad_core) ..."
	@cd cbad-core && cargo build --release

demo: cbad-core
	@echo "Building Go demo ..."
	@go build -o driftlock-demo cmd/demo/main.go

.PHONY: docker-http docker-test

docker-http:
	@echo "Building driftlock-http Docker image ..."
	@docker build -t driftlock-http:dev -f collector-processor/cmd/driftlock-http/Dockerfile .

docker-test:
	@./scripts/test-docker-build.sh

verify: demo
	@echo "Running verification script ..."
	@./scripts/verify-launch-readiness.sh

clean:
	@rm -f driftlock-demo demo-output.html verify.log
	@cd cbad-core && cargo clean

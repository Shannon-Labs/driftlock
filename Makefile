SHELL := /bin/bash

.PHONY: all cbad-core demo verify clean

all: demo

cbad-core:
	@echo "Building Rust core (cbad_core) ..."
	@cd cbad-core && cargo build --release

demo: cbad-core
	@echo "Building Go demo ..."
	@go build -o driftlock-demo cmd/demo/main.go

verify: demo
	@echo "Running verification script ..."
	@./verify-yc-ready.sh

clean:
	@rm -f driftlock-demo demo-output.html verify.log
	@cd cbad-core && cargo clean

PKG=./...

.PHONY: run build tidy test clean docker api collector tools

run:
	go run ./api-server/cmd/driftlock-api

build:
	CGO_ENABLED=0 go build -o bin/driftlock-api ./api-server/cmd/driftlock-api

tidy:
	go mod tidy

test:
	go test $(PKG) -v

clean:
	rm -rf bin

docker:
	docker build -t driftlock:api .

# Focused targets
api:
	go build -o bin/driftlock-api ./api-server/cmd/driftlock-api

collector:
	go build ./collector-processor/...

tools:
	go build ./tools/...

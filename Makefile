.PHONY: build run test lint fmt clean coverage help

help:
	@echo "Available targets:"
	@echo "  make build    - Build binary"
	@echo "  make run      - Run server"
	@echo "  make test     - Run tests"
	@echo "  make lint     - Lint code"
	@echo "  make fmt      - Format code"
	@echo "  make clean    - Clean build"

build:
	mkdir -p bin
	go build -o bin/mcp-proxy ./cmd/proxy

run:
	go run ./cmd/proxy

test:
	go test -v -race ./...

lint:
	golangci-lint run

fmt:
	gofumpt -w . && goimports -w .

clean:
	rm -rf bin/ && go clean

.DEFAULT_GOAL := help

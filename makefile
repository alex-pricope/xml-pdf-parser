.PHONY: all lint build test build-lint test-build-lint

all: lint test build

lint:
	@command -v golangci-lint >/dev/null 2>&1 || { echo >&2 "golangci-lint is not installed. Please install it: https://golangci-lint.run/usage/install/"; exit 1; }
	golangci-lint run ./...

build:
	mkdir -p bin
	go build -o bin/parser .

test:
	go test ./...

.PHONY: all build test lint
ENVFLAGS = GO111MODULE=on CGO_ENABLED=0 GOOS=$(shell go env GOOS) GOARCH=$(shell go env GOARCH)
SRV = "wsc"

help:
	@echo "usage:"
	@echo "make build            -- build $(SRV) one executable"
	@echo "make test             -- test all packages"
	@echo "make lint             -- check the code for common errors"

all: lint test build

build:
	@echo "Building the binary..."
	$(ENVFLAGS) go build -a -tags netgo -ldflags="-w -s" -o build/bin/$(SRV) ./cmd/cli
	@echo "$(SRV) is built in build/$(SRV)"
	@echo "DONE!"

test:
	@echo "Running tests"
	go test ./... -v -race
	@echo "DONE!"

lint:
	@echo "Running Linter..."
	go vet ./...
	@echo "DONE!"

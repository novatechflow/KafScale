.PHONY: proto build test tidy lint generate

proto: ## Generate protobuf + gRPC stubs
	buf generate

generate: proto

build: ## Build all binaries
	go build ./...

test: ## Run unit tests
	go test ./...

tidy:
	go mod tidy

lint:
	golangci-lint run

help: ## Show targets
	@grep -E '^[a-zA-Z_-]+:.*?##' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "%-20s %s\n", $$1, $$2}'

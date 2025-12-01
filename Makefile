.PHONY: help build test lint clean run-monitor run-client docker-build docker-run docker-stop install

# Variables
BINARY_MONITOR=bin/monitor
BINARY_CLIENT=bin/client
DOCKER_IMAGE=system-monitor
DOCKER_TAG=latest
GO_VERSION=1.23

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build monitor and client binaries
	@echo "Building binaries..."
	@mkdir -p bin
	go build -o $(BINARY_MONITOR) ./cmd/monitor
	go build -o $(BINARY_CLIENT) ./cmd/client
	@echo "Build complete: $(BINARY_MONITOR), $(BINARY_CLIENT)"

test: ## Run tests
	@echo "Running tests..."
	go test -v -race -count=1 ./...

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	go test -v -tags=integration ./internal/server/...

test-all: test test-integration ## Run all tests

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	rm -f *.log
	@echo "Clean complete"

run-monitor: build ## Run monitor daemon
	@echo "Starting monitor..."
	./$(BINARY_MONITOR)

run-monitor-config: build ## Run monitor with config file
	@echo "Starting monitor with config..."
	./$(BINARY_MONITOR) --config configs/config.yaml

run-client: build ## Run client
	@echo "Starting client..."
	./$(BINARY_CLIENT)

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)"

docker-run: docker-build ## Run monitor in Docker container
	@echo "Running monitor in Docker..."
	docker run -d \
		--name system-monitor \
		--network host \
		-v /proc:/host/proc:ro \
		-v /sys:/host/sys:ro \
		$(DOCKER_IMAGE):$(DOCKER_TAG)
	@echo "Container started. Use 'make docker-logs' to view logs."

docker-run-config: docker-build ## Run monitor in Docker with config file
	@echo "Running monitor in Docker with config..."
	docker run -d \
		--name system-monitor \
		--network host \
		-v /proc:/host/proc:ro \
		-v /sys:/host/sys:ro \
		-v $(PWD)/configs/config.yaml:/app/configs/config.yaml:ro \
		$(DOCKER_IMAGE):$(DOCKER_TAG) --config /app/configs/config.yaml
	@echo "Container started. Use 'make docker-logs' to view logs."

docker-logs: ## View Docker container logs
	docker logs -f system-monitor

docker-stop: ## Stop Docker container
	@echo "Stopping container..."
	docker stop system-monitor || true
	docker rm system-monitor || true
	@echo "Container stopped"

docker-rm: ## Remove Docker image
	@echo "Removing Docker image..."
	docker rmi $(DOCKER_IMAGE):$(DOCKER_TAG) || true
	@echo "Image removed"

install: build ## Install binaries to /usr/local/bin
	@echo "Installing binaries..."
	sudo cp $(BINARY_MONITOR) /usr/local/bin/monitor
	sudo cp $(BINARY_CLIENT) /usr/local/bin/system-monitor-client
	@echo "Installation complete"

proto: ## Generate protobuf code
	@echo "Generating protobuf code..."
	protoc --go_out=. --go-grpc_out=. api/proto/metrics.proto
	@echo "Protobuf code generated"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

vendor: ## Create vendor directory
	@echo "Creating vendor directory..."
	go mod vendor

.DEFAULT_GOAL := help


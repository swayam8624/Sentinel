# Sentinel Makefile

# Variables
APP_NAME := sentinel
VERSION := 0.1.0
DOCKER_IMAGE := sentinel/gateway:$(VERSION)
DOCKER_IMAGE_LATEST := sentinel/gateway:latest

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	$(GOBUILD) -o $(APP_NAME) .

# Run the application
.PHONY: run
run:
	$(GOBUILD) -o $(APP_NAME) .
	./$(APP_NAME)

# Run tests
.PHONY: test
test:
	$(GOTEST) -v ./...

# Run unit tests
.PHONY: test-unit
test-unit:
	$(GOTEST) -v ./tests/unit/...

# Run integration tests
.PHONY: test-integration
test-integration:
	$(GOTEST) -v ./tests/integration/...

# Run security tests
.PHONY: test-security
test-security:
	$(GOTEST) -v ./tests/security/...

# Run performance tests
.PHONY: test-performance
test-performance:
	$(GOTEST) -v ./tests/performance/...

# Run comprehensive tests
.PHONY: test-comprehensive
test-comprehensive:
	$(GOTEST) -v ./tests/comprehensive/...

# Run benchmarks
.PHONY: bench
bench:
	$(GOTEST) -bench=. ./tests/performance/...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

# Clean build artifacts
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(APP_NAME)

# Install dependencies
.PHONY: deps
deps:
	$(GOMOD) download

# Update dependencies
.PHONY: deps-update
deps-update:
	$(GOMOD) tidy

# Build Docker image
.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_IMAGE) .
	docker tag $(DOCKER_IMAGE) $(DOCKER_IMAGE_LATEST)

# Run Docker container
.PHONY: docker-run
docker-run:
	docker-compose up -d

# Stop Docker containers
.PHONY: docker-stop
docker-stop:
	docker-compose down

# Run integration tests
.PHONY: integration-test
integration-test:
	# Start services
	docker-compose up -d
	# Wait for services to be ready
	sleep 10
	# Run integration tests
	$(GOTEST) -v ./integration/...
	# Stop services
	docker-compose down

# Lint the code
.PHONY: lint
lint:
	golangci-lint run

# Generate documentation
.PHONY: docs
docs:
	# Generate API documentation
	# Generate code documentation
	$(GOCMD) doc ./...

# Install development tools
.PHONY: dev-tools
dev-tools:
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOGET) -u github.com/swaggo/swag/cmd/swag@latest

# Create a new release
.PHONY: release
release:
	@if [ -z "$(VERSION)" ]; then \
		echo "VERSION is not set. Usage: make release VERSION=1.0.0"; \
		exit 1; \
	fi
	./scripts/release.sh $(VERSION)

# Publish to Docker Hub
.PHONY: publish-docker
publish-docker:
	@if [ -z "$(DOCKERHUB_USERNAME)" ]; then \
		echo "DOCKERHUB_USERNAME is not set. Usage: make publish-docker DOCKERHUB_USERNAME=yourusername"; \
		exit 1; \
	fi
	./scripts/publish-to-dockerhub.sh $(DOCKERHUB_USERNAME)

# Publish Helm charts
.PHONY: publish-helm
publish-helm:
	./scripts/publish-helm-charts.sh

# Publish Node.js SDK
.PHONY: publish-nodejs
publish-nodejs:
	./scripts/publish-nodejs-sdk.sh

# Publish Python SDK
.PHONY: publish-python
publish-python:
	./scripts/publish-python-sdk.sh

# Run all tests
.PHONY: test-all
test-all: test-unit test-integration test-security test-performance test-comprehensive

# Help
.PHONY: help
help:
	@echo "Sentinel Makefile"
	@echo "================"
	@echo "build          - Build the application"
	@echo "run            - Run the application"
	@echo "test           - Run all tests"
	@echo "test-unit      - Run unit tests"
	@echo "test-integration - Run integration tests"
	@echo "test-security  - Run security tests"
	@echo "test-performance - Run performance tests"
	@echo "test-comprehensive - Run comprehensive tests"
	@echo "test-all       - Run all test suites"
	@echo "bench          - Run benchmarks"
	@echo "test-coverage  - Run tests with coverage report"
	@echo "clean          - Clean build artifacts"
	@echo "deps           - Install dependencies"
	@echo "deps-update    - Update dependencies"
	@echo "docker-build   - Build Docker image"
	@echo "docker-run     - Run Docker containers"
	@echo "docker-stop    - Stop Docker containers"
	@echo "integration-test - Run integration tests"
	@echo "lint           - Lint the code"
	@echo "docs           - Generate documentation"
	@echo "dev-tools      - Install development tools"
	@echo "release        - Create a new release (requires version)"
	@echo "publish-docker - Publish to Docker Hub (requires username)"
	@echo "publish-helm   - Publish Helm charts"
	@echo "publish-nodejs - Publish Node.js SDK"
	@echo "publish-python - Publish Python SDK"
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

# Help
.PHONY: help
help:
	@echo "Sentinel Makefile"
	@echo "================"
	@echo "build          - Build the application"
	@echo "run            - Run the application"
	@echo "test           - Run unit tests"
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
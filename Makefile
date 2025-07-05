# Makefile for spemu - Spanner Emulator DML Inserter

.PHONY: help test test-unit test-integration build clean lint fmt vet install dev-setup emulator-start emulator-stop

# Variables
BINARY_NAME=spemu
MAIN_PATH=.
PKG_LIST := $(shell go list ./...)
TEST_PKG_LIST := $(shell go list ./pkg/...)
VERSION := $(shell git describe --tags --always --dirty)

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*?##/ { printf "  %-15s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

# Development setup
dev-setup: ## Install development dependencies
	@echo "Installing development dependencies..."
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go mod download

# Build targets
build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	go build -ldflags "-X main.Version=$(VERSION)" -o $(BINARY_NAME) $(MAIN_PATH)


install: ## Install the binary to GOPATH/bin
	go install -ldflags "-X main.Version=$(VERSION)" .

# Code quality targets
fmt: ## Format Go code
	@echo "Formatting code..."
	gofmt -s -w .
	goimports -w .

lint: ## Run linters
	@echo "Running linters..."
	go vet ./...
	staticcheck ./...

vet: ## Run go vet
	go vet ./...

# Test targets
test: test-unit ## Run all tests (default: unit tests only)

test-unit: ## Run unit tests
	@echo "Running unit tests..."
	go test -v -race -coverprofile=coverage.out $(TEST_PKG_LIST)

test-integration: ## Run integration tests (automatically starts emulator)
	@echo "Running integration tests..."
	@echo "Note: This will automatically setup Spanner emulator if needed"
	SPANNER_EMULATOR_HOST=localhost:9010 go test -v ./test/...

test-all: test-unit test-integration ## Run all tests including integration tests

test-coverage: ## Run tests and show coverage
	go test -v -race -coverprofile=coverage.out $(TEST_PKG_LIST)
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report saved to coverage.html"


# Docker and Spanner Emulator targets
emulator-start: ## Start Spanner emulator with Docker
	@echo "Starting Spanner emulator with Docker..."
	docker compose up -d spanner-emulator
	@echo "Waiting for emulator to be ready..."
	@until docker compose exec -T spanner-emulator curl -f http://localhost:9010 >/dev/null 2>&1; do \
		echo "Waiting for emulator..."; \
		sleep 2; \
	done
	@echo "Spanner emulator is ready on localhost:9010"

emulator-init: ## Initialize database schema
	@echo "Initializing database schema..."
	docker compose --profile init up spanner-init
	docker compose --profile init down spanner-init

emulator-setup: emulator-start emulator-init ## Start emulator and initialize database

emulator-stop: ## Stop Spanner emulator
	@echo "Stopping Spanner emulator..."
	docker compose down

emulator-reset: ## Reset and reinitialize emulator
	@echo "Resetting emulator..."
	docker compose down
	$(MAKE) emulator-setup

emulator-logs: ## Show emulator logs
	docker compose logs -f spanner-emulator

# Example and demo targets
demo: build ## Run demo with example data
	@echo "Running demo with example data..."
	@if [ ! -f $(BINARY_NAME) ]; then $(MAKE) build; fi
	@echo "Testing dry run..."
	./$(BINARY_NAME) -p test-project -i test-instance -d test-database --dry-run examples/seed.sql
	@echo "\nExecuting statements (requires emulator)..."
	SPANNER_EMULATOR_HOST=localhost:9010 ./$(BINARY_NAME) -p test-project -i test-instance -d test-database examples/seed.sql

# Development workflow targets
dev-test: fmt lint test-unit ## Format, lint, and run unit tests
	@echo "Development tests completed successfully!"

dev-full: fmt lint test-all ## Format, lint, and run all tests
	@echo "Full development cycle completed successfully!"

ci: lint test-unit build ## Run CI pipeline (lint, test, build)
	@echo "CI pipeline completed successfully!"

# Clean targets
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*
	rm -f coverage.out coverage.html
	go clean -cache
	go clean -testcache

clean-all: clean ## Clean everything including module cache
	go clean -modcache



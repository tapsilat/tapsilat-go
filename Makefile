# Tapsilat Go SDK Makefile

.PHONY: test test-unit test-integration test-coverage clean build fmt vet lint help

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

test: test-unit test-integration ## Run all tests

test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	go test -v ./tests/unit/...

test-integration: ## Run integration tests only (requires real token in test files)
	@echo "Running integration tests..."
	go test -v ./tests/integration/...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./tests/unit/... ./tests/integration/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-validators: ## Run only validator tests
	@echo "Running validator tests..."
	go test -v ./tests/unit/ -run TestValidate

test-orders: ## Run only order tests
	@echo "Running order tests..."
	go test -v ./tests/unit/ -run TestOrder

test-api: ## Run only API tests
	@echo "Running API tests..."
	go test -v ./tests/unit/ -run TestAPI

build: ## Build the package
	@echo "Building package..."
	go build -v ./...

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

lint: fmt vet ## Run formatting and vetting

clean: ## Clean build artifacts and coverage files
	@echo "Cleaning..."
	rm -f coverage.out coverage.html
	go clean -testcache
	go mod tidy

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Quick test run without verbose output
test-quick: ## Run tests quickly without verbose output
	go test ./tests/unit/... ./tests/integration/...

# Test specific functions
test-gsm: ## Test GSM validation specifically
	go test -v ./tests/unit/ -run TestValidateGSMNumber

test-installments: ## Test installment validation specifically
	go test -v ./tests/unit/ -run TestValidateInstallments

# Development helpers
dev-setup: deps ## Setup development environment
	@echo "Setting up development environment..."
	@echo "Development environment ready. Update test files with real tokens for integration tests."

# Example test with real API (use with caution)
test-create-order: ## Test order creation with real API (requires real token in test files)
	go test -v ./tests/integration/ -run TestCreateOrder

# Run usage examples
run-examples: ## Run usage examples
	@echo "Running usage examples..."
	go run examples/usage.go

run-usage: run-examples ## Alias for run-examples

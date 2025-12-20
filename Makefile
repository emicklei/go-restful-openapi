# go-restful-openapi Makefile

# =============================================================================
# Test Targets
# =============================================================================

## test: Run tests with coverage
.PHONY: test
test:
	@echo "Running tests..."
	go test -cover ./...

## test-race: Run tests with race detector
.PHONY: test-race
test-race:
	@echo "Running tests with race detector..."
	go test -race -short ./...

## coverage: Run tests with HTML coverage report
.PHONY: coverage
coverage:
	@echo "Generating coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# =============================================================================
# Code Quality Targets
# =============================================================================

## lint: Run linter
.PHONY: lint
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install it from https://golangci-lint.run/usage/install/"; \
		exit 1; \
	fi

## fmt: Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

## vet: Run go vet
.PHONY: vet
vet:
	@echo "Running go vet..."
	go vet ./...

## check: Run tidy, fmt, lint, test, and git status
.PHONY: check
check: tidy fmt lint test
	@echo "Running git status..."
	@git status

# =============================================================================
# Dependency Targets
# =============================================================================

## tidy: Tidy go modules
.PHONY: tidy
tidy:
	@echo "Tidying go modules..."
	go mod tidy

## outdated: Check for outdated dependencies
.PHONY: outdated
outdated:
	@echo "Checking for outdated dependencies..."
	go list -u -m -json all | go-mod-outdated -update -direct

# =============================================================================
# Cleanup Targets
# =============================================================================

## clean: Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -f coverage.out

# =============================================================================
# Help Target
# =============================================================================

## help: Show this help message
.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

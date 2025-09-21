# Carthooks Go SDK Makefile

.PHONY: test build clean fmt vet lint example

# Default target
all: fmt vet test

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Build example
build:
	@echo "Building example..."
	go build -o bin/example ./examples/basic_usage.go

# Run example (requires valid credentials)
example: build
	@echo "Running example..."
	./bin/example

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Run golint (if available)
lint:
	@echo "Running golint..."
	@if command -v golint >/dev/null 2>&1; then \
		golint ./...; \
	else \
		echo "golint not installed. Install with: go install golang.org/x/lint/golint@latest"; \
	fi

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Install dependencies for development
dev-deps:
	@echo "Installing development dependencies..."
	go install golang.org/x/lint/golint@latest

# Run all checks
check: fmt vet lint test

# Show help
help:
	@echo "Available targets:"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  build        - Build example"
	@echo "  example      - Build and run example"
	@echo "  fmt          - Format code"
	@echo "  vet          - Run go vet"
	@echo "  lint         - Run golint"
	@echo "  clean        - Clean build artifacts"
	@echo "  dev-deps     - Install development dependencies"
	@echo "  check        - Run all checks (fmt, vet, lint, test)"
	@echo "  help         - Show this help"

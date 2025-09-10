# ADK Go SDK Makefile

.PHONY: build test run clean fmt vet lint deps example

# Go settings
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOVET=$(GOCMD) vet

# Default target
all: deps fmt vet test build

# Build the example
build:
	mkdir -p bin
	$(GOBUILD) -o bin/adk-example ./examples/main.go

# Run tests
test:
	$(GOTEST) -v ./google/adk/...

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -cover ./google/adk/...

# Run the example
run: build
	./bin/adk-example

# Run the example directly
example:
	$(GOCMD) run ./examples/main.go

# Format code
fmt:
	$(GOFMT) -s -w .

# Vet code
vet:
	$(GOVET) ./...

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Clean build artifacts
clean:
	rm -rf bin/
	$(GOCMD) clean

# Development helpers
dev-setup:
	@echo "Setting up development environment..."
	$(GOMOD) tidy
	@echo "Development environment ready!"

# Check for common issues
check: fmt vet test

# Build for multiple platforms
build-all:
	mkdir -p bin
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/adk-example-linux-amd64 ./examples/main.go
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o bin/adk-example-darwin-amd64 ./examples/main.go
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o bin/adk-example-windows-amd64.exe ./examples/main.go

# Run benchmark tests
bench:
	$(GOTEST) -bench=. ./google/adk/...

# Help
help:
	@echo "Available targets:"
	@echo "  all          - Run deps, fmt, vet, test, and build"
	@echo "  build        - Build the example application"
	@echo "  test         - Run all tests"
	@echo "  test-coverage- Run tests with coverage"
	@echo "  run          - Build and run the example"
	@echo "  example      - Run the example directly"
	@echo "  fmt          - Format code"
	@echo "  vet          - Vet code"
	@echo "  deps         - Download dependencies"
	@echo "  clean        - Clean build artifacts"
	@echo "  dev-setup    - Set up development environment"
	@echo "  check        - Run fmt, vet, and test"
	@echo "  build-all    - Build for multiple platforms"
	@echo "  bench        - Run benchmark tests"
	@echo "  help         - Show this help message"
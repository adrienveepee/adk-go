# Makefile for ADK Go project

.PHONY: build test clean run fmt vet lint

# Default target
all: fmt vet test build

# Build the application
build:
	go build -o bin/adk-server ./cmd/adk-server

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -cover ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Run the server
run: build
	./bin/adk-server

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Install dependencies
deps:
	go mod download
	go mod tidy

# Development server with auto-reload (requires air)
dev:
	air

# Install air for development
install-air:
	go install github.com/air-verse/air@latest

# Print help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean         - Clean build artifacts"
	@echo "  run           - Build and run the server"
	@echo "  fmt           - Format code"
	@echo "  vet           - Vet code"
	@echo "  deps          - Install dependencies"
	@echo "  dev           - Run development server (requires air)"
	@echo "  install-air   - Install air for development"
	@echo "  help          - Show this help"
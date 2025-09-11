# Makefile for agent-ollama-gin

# Variables
BINARY_NAME=agent-ollama-gin
MAIN_PATH=./main.go
BUILD_DIR=./bin

# Go commands
.PHONY: build run clean test deps install-tools install-genkit dev watch

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	go run $(MAIN_PATH)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	go clean

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Install development tools
install-tools: install-genkit
	@echo "Installing development tools..."
	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install Firebase Genkit CLI
install-genkit:
	@echo "Installing Firebase Genkit CLI..."
	@if command -v genkit >/dev/null 2>&1; then \
		echo "Genkit already installed"; \
	else \
		curl -sL cli.genkit.dev | bash; \
	fi

# Development mode with hot reload
watch: install-tools
	@echo "Starting development server with hot reload..."
	air

# Development setup (install everything needed)
dev-setup: deps install-tools
	@echo "Development environment setup complete!"
	@echo "Run 'make watch' to start development server"
	@echo "Run 'genkit --help' to see Genkit commands"

# Start genkit development server
genkit-dev:
	@echo "Starting Genkit development server..."
	genkit start

# Initialize genkit in the project
genkit-init:
	@echo "Initializing Genkit..."
	genkit init

# Genkit flow deployment
genkit-deploy:
	@echo "Deploying with Genkit..."
	genkit deploy

# Help
help:
	@echo "Available commands:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"
	@echo "  deps         - Install Go dependencies"
	@echo "  install-tools- Install development tools (Air, Genkit, etc.)"
	@echo "  install-genkit- Install Firebase Genkit CLI"
	@echo "  watch        - Start development server with hot reload"
	@echo "  dev-setup    - Complete development environment setup"
	@echo "  genkit-dev   - Start Genkit development server"
	@echo "  genkit-init  - Initialize Genkit in the project"
	@echo "  genkit-deploy- Deploy with Genkit"
	@echo "  help         - Show this help message"
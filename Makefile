# Makefile for agent-ollama-gin

# Default target
help:
	@echo "Available commands:"
	@echo "  build         - Build the Go application"
	@echo "  build-cli     - Build the encyclopedia CLI"
	@echo "  run           - Run the application locally"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run with Docker Compose"
	@echo "  docker-stop   - Stop Docker services"
	@echo "  install-deps  - Install Go dependencies"
	@echo "  lint          - Run linter"
	@echo "  format        - Format Go code"

# Go commands
.PHONY: build run clean test deps install-tools install-genkit dev watch

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Build the encyclopedia CLI
build-cli: install-deps
	go build -o bin/encyclopedia cmd/encyclopedia/main.go

# Run locally
run: install-deps
	go run main.go

# Run tests
test: install-deps
	go test -v ./...

# Run tests with coverage
test-coverage: install-deps
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out

# Build Docker image
docker-build:
	docker build -t llama-api .

# Run with Docker Compose
docker-run:
	docker-compose up -d

# Stop Docker services
docker-stop:
	docker-compose down

# Run with Docker Compose and rebuild
docker-rebuild:
	docker-compose up -d --build

# View logs
docker-logs:
	docker-compose logs -f

# Lint code
lint: install-deps
	golangci-lint run

# Format code
format: install-deps
	go fmt ./...
	go vet ./...

# Install Ollama (macOS)
install-ollama-mac:
	curl -fsSL https://ollama.ai/install.sh | sh

# Install Ollama (Linux)
install-ollama-linux:
	curl -fsSL https://ollama.ai/install.sh | sh

# Pull default model
pull-model:
	ollama pull llama2

# Start Ollama service
start-ollama:
	ollama serve

# Check Ollama status
check-ollama:
	ollama list

# Development setup
dev-setup: install-deps install-ollama-mac pull-model
	@echo "Development environment setup complete!"
	@echo "Run 'make start-ollama' in one terminal"
	@echo "Run 'make run' in another terminal"

# Production build
prod-build: install-deps
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/llama-api main.go

# Show project info
info:
	@echo "Llama API Project"
	@echo "================="
	@echo "Go version: $(shell go version)"
	@echo "Project path: $(shell pwd)"
	@echo "Go modules: $(shell go list -m all | head -5)"

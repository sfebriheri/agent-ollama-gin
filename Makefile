.PHONY: help build run test clean docker-build docker-run docker-stop install-deps lint format

# Default target
help:
	@echo "Available commands:"
	@echo "  build         - Build the Go application"
	@echo "  run           - Run the application locally"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  deep-clean    - Deep clean (all generated files)"
	@echo "  clean-docker  - Clean Docker resources"
	@echo "  clean-all     - Complete project cleanup"
	@echo "  clean-interactive - Interactive cleaning menu"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run with Docker Compose"
	@echo "  docker-stop   - Stop Docker services"
	@echo "  install-deps  - Install Go dependencies"
	@echo "  lint          - Run linter"
	@echo "  format        - Format Go code"
	@echo "  setup         - Run complete setup script"
	@echo "  quick-start   - Quick development setup"
	@echo "  dev           - Start development environment"
	@echo "  stop-dev      - Stop development environment"
	@echo "  test-api      - Test API endpoints"
	@echo "  watch         - Run with file watching (requires air)"
	@echo "  install-tools - Install development tools"
	@echo "  status        - Check development environment status"
	@echo "  workflow      - Interactive development workflow"

# Install dependencies
install-deps:
	go mod download
	go mod tidy

# Build the application
build: install-deps
	go build -o bin/llama-api main.go

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
	@echo "🧹 Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out
	@echo "✅ Build artifacts cleaned"

# Deep clean (removes all generated files)
deep-clean: clean
	@echo "🧹 Deep cleaning project..."
	rm -rf tmp/
	rm -rf logs/
	rm -f .env
	rm -f build-errors.log
	rm -f *.log
	go clean -modcache
	go clean -cache
	go clean -testcache
	@echo "✅ Deep clean completed"

# Clean Docker resources
clean-docker:
	@echo "🐳 Cleaning Docker resources..."
	docker-compose down --volumes --remove-orphans 2>/dev/null || true
	docker system prune -f
	docker image prune -f
	@echo "✅ Docker resources cleaned"

# Clean all (build + Docker + generated files)
clean-all: clean-docker deep-clean
	@echo "🧹 Complete project cleanup..."
	rm -rf .air.toml.bak 2>/dev/null || true
	rm -rf .DS_Store 2>/dev/null || true
	rm -rf *.tmp 2>/dev/null || true
	@echo "✅ Complete cleanup finished"

# Interactive cleaning script
clean-interactive:
	@echo "Starting interactive cleaning script..."
	@chmod +x llama-api.sh
	@./llama-api.sh

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

# Check development environment status
status:
	@echo "Checking development environment status..."
	@chmod +x llama-api.sh
	@./llama-api.sh

# Interactive development workflow
workflow:
	@echo "Starting interactive development workflow..."
	@chmod +x llama-api.sh
	@./llama-api.sh

# Setup and development commands
setup:
	@echo "Running complete setup..."
	@chmod +x llama-api.sh
	@./llama-api.sh

quick-start:
	@echo "Running quick start..."
	@chmod +x llama-api.sh
	@./llama-api.sh

dev:
	@echo "Starting development environment..."
	@chmod +x llama-api.sh
	@./llama-api.sh

stop-dev:
	@echo "Stopping development environment..."
	@chmod +x llama-api.sh
	@./llama-api.sh

test-api:
	@echo "Testing API endpoints..."
	@chmod +x llama-api.sh
	@./llama-api.sh

# Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Development tools installed"

# Run with file watching (requires air)
watch:
	@echo "Starting with file watching..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not found. Installing..."; \
		make install-tools; \
		air; \
	fi

# Development shortcuts
dev-build: install-deps build
dev-test: install-deps test
dev-lint: install-deps lint
dev-format: install-deps format
dev-clean: clean

# Docker development
docker-dev:
	@echo "Starting development environment with Docker..."
	@docker-compose -f docker-compose.yml up -d ollama
	@echo "Ollama started. Run 'make run' to start the API"

# Health checks
health-check:
	@echo "Checking service health..."
	@curl -s http://localhost:8080/api/v1/health | jq . || echo "API not running"
	@curl -s http://localhost:11434/api/tags > /dev/null && echo "Ollama: OK" || echo "Ollama: Not running"

# Model management
list-models:
	@echo "Available models:"
	@ollama list

pull-model:
	@echo "Pulling default model..."
	@ollama pull llama2

# Logs
logs-api:
	@echo "API logs (if running with Docker):"
	@docker-compose logs api || echo "API not running in Docker"

logs-ollama:
	@echo "Ollama logs (if running with Docker):"
	@docker-compose logs ollama || echo "Ollama not running in Docker"

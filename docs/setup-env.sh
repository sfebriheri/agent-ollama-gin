#!/bin/bash

# Setup script for agent-ollama-gin development environment
# This script helps configure the environment after cloning the repository

set -e

echo "🚀 Setting up agent-ollama-gin development environment..."
echo "=================================================="

# Check if .env file exists
if [ ! -f ".env" ]; then
    echo "📝 Creating .env file from template..."
    cp env.example .env
    echo "✅ .env file created from env.example"
    echo "⚠️  Please edit .env file with your actual configuration values"
else
    echo "✅ .env file already exists"
fi

# Create necessary directories
echo "📁 Creating necessary directories..."
mkdir -p bin
mkdir -p tmp
mkdir -p logs
mkdir -p static/css
mkdir -p static/js
mkdir -p static/images
echo "✅ Directories created"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or later"
    exit 1
fi

echo "✅ Go is installed: $(go version)"

# Install dependencies
echo "📦 Installing Go dependencies..."
go mod download
go mod tidy
echo "✅ Go dependencies installed"

# Check if Ollama is installed
if ! command -v ollama &> /dev/null; then
    echo "⚠️  Ollama is not installed. Please install Ollama from https://ollama.ai"
    echo "   After installation, run: ollama serve"
else
    echo "✅ Ollama is installed: $(ollama --version)"
fi

# Install development tools
echo "🔧 Installing development tools..."

# Install Air for hot reloading
if ! command -v air &> /dev/null; then
    echo "Installing Air for hot reloading..."
    go install github.com/air-verse/air@latest
    echo "✅ Air installed"
else
    echo "✅ Air already installed"
fi

# Install golangci-lint for linting
if ! command -v golangci-lint &> /dev/null; then
    echo "Installing golangci-lint for code linting..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    echo "✅ golangci-lint installed"
else
    echo "✅ golangci-lint already installed"
fi

# Check if Firebase Genkit is installed
if ! command -v genkit &> /dev/null; then
    echo "⚠️  Firebase Genkit CLI is not installed."
    echo "   To install: curl -sL cli.genkit.dev | bash"
else
    echo "✅ Firebase Genkit CLI is installed"
fi

echo ""
echo "🎉 Setup complete!"
echo ""
echo "Next steps:"
echo "1. Edit .env file with your configuration"
echo "2. Start Ollama: ollama serve"
echo "3. Pull a model: ollama pull llama2"
echo "4. Start development server: make watch"
echo "   Or: go run main.go"
echo ""
echo "Available commands:"
echo "  make build       - Build the application"
echo "  make run         - Run the application"
echo "  make watch       - Start with hot reload"
echo "  make test        - Run tests"
echo "  make clean       - Clean build artifacts"
echo ""
echo "For more information, see README.md"
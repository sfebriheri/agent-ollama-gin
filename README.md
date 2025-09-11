# Agent Ollama Gin

A Go-based AI agent application built with Gin framework and integrated with Ollama for local LLM capabilities.

## Features

- 🚀 Fast HTTP server with Gin framework
- 🤖 Local LLM integration with Ollama
- 🔥 Hot reload development with Air
- 🛠️ AI development toolkit with Firebase Genkit
- 📝 RESTful API design
- 🔄 Real-time capabilities

## Prerequisites

- [Go 1.21+](https://golang.org/doc/install)
- [Ollama](https://ollama.ai/) installed and running
- [Node.js](https://nodejs.org/) (optional, for Genkit features)

## Quick Start

### Automated Setup (Recommended)

```bash
# Clone the repository
git clone https://github.com/sfebriheri/agent-ollama-gin.git
cd agent-ollama-gin

# Run the setup script
chmod +x setup.sh
./setup.sh

# Copy environment configuration
cp .env.example .env

# Start development server
make watch
```

### Manual Setup

1. **Install Dependencies**
   ```bash
   go mod download
   go mod tidy
   ```

2. **Install Development Tools**
   ```bash
   # Install Air for hot reload
   go install github.com/air-verse/air@latest
   
   # Install Firebase Genkit CLI
   curl -sL cli.genkit.dev | bash
   ```

3. **Start Ollama**
   ```bash
   ollama serve
   ```

4. **Run the Application**
   ```bash
   # Development mode with hot reload
   make watch
   
   # Or run directly
   go run main.go
   ```

## Development Commands

### Using Makefile

```bash
# Build the application
make build

# Run the application
make run

# Development with hot reload
make watch

# Install all development tools
make dev-setup

# Clean build artifacts
make clean

# Run tests
make test
```

### Using Genkit

```bash
# Start Genkit development server
genkit start

# Initialize Genkit (if not done automatically)
genkit init

# Deploy with Genkit
genkit deploy

# Show Genkit help
genkit --help
```

## Project Structure

```
agent-ollama-gin/
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
├── go.sum                  # Go dependencies
├── Makefile                # Build and development commands
├── setup.sh                # Automated setup script
├── .air.toml              # Air configuration for hot reload
├── .env.example           # Environment variables template
├── .env                   # Environment variables (create from example)
├── genkit.config.js       # Genkit configuration (auto-generated)
├── bin/                   # Built binaries
├── cmd/                   # Application commands
├── internal/              # Private application code
├── pkg/                   # Public packages
├── api/                   # API handlers and routes
├── models/                # Data models
├── services/              # Business logic
├── config/                # Configuration files
├── templates/             # HTML templates (if using)
├── static/                # Static assets
└── docs/                  # Documentation
```

## Configuration

Copy `.env.example` to `.env` and configure:

```bash
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=localhost

# Ollama Configuration
OLLAMA_URL=http://localhost:11434
OLLAMA_MODEL=llama2

# Logging
LOG_LEVEL=info
```

## API Endpoints

```bash
# Health check
GET /health

# Agent endpoints (example)
POST /api/v1/chat
GET /api/v1/models
POST /api/v1/generate
```

## Genkit Integration

This project uses [Firebase Genkit](https://firebase.google.com/docs/genkit) for AI application development:

- **Flow Development**: Create and manage AI flows
- **Model Integration**: Easy LLM model switching
- **Development UI**: Visual development environment
- **Deployment**: Streamlined deployment process

### Genkit Features

- Visual flow builder and debugger
- Built-in evaluation and testing tools  
- Multiple model provider support
- Production deployment capabilities

## Development Workflow

1. **Start Ollama**: `ollama serve`
2. **Start Development**: `make watch` or `genkit start`
3. **Make Changes**: Edit code with hot reload
4. **Test**: Run tests with `make test`
5. **Build**: Create production build with `make build`

## Docker Support (Optional)

```dockerfile
# Dockerfile example
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o bin/agent-ollama-gin main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bin/agent-ollama-gin .
EXPOSE 8080
CMD ["./agent-ollama-gin"]
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Submit a pull request

## Troubleshooting

### Common Issues

1. **Genkit not found after installation**
   ```bash
   # Add to your PATH
   export PATH=$PATH:$HOME/.local/bin
   source ~/.bashrc  # or ~/.zshrc
   ```

2. **Ollama connection issues**
   ```bash
   # Check Ollama status
   ollama list
   
   # Start Ollama service
   ollama serve
   ```

3. **Module path conflicts**
   ```bash
   # Clean module cache
   go clean -modcache
   go mod tidy
   ```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Gin Web Framework](https://gin-gonic.com/)
- [Ollama](https://ollama.ai/)
- [Firebase Genkit](https://firebase.google.com/docs/genkit)
- [Air - Live reload](https://github.com/air-verse/air)
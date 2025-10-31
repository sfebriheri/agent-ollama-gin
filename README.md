# Agent Ollama Gin - Ollama Cloud Integration

A powerful REST API service built with Go and Gin framework that provides seamless integration with Ollama cloud models. This service has been migrated from Genkit to support Ollama's cloud infrastructure for enhanced AI model access.

## 🚀 Features

- 🚀 **Fast & Lightweight**: Built with Go and Gin for high performance
- 🤖 **LLM Integration**: Seamless integration with Llama models via Ollama
- 💬 **Chat Completion**: Full chat conversation support
- ✍️ **Text Completion**: Generate text completions
- 🔍 **Embeddings**: Generate text embeddings for semantic search
- 📡 **Streaming Support**: Real-time streaming responses
- 🔍 **Encyclopedia Agent**: AI-powered encyclopedia search and content generation
- 💻 **CLI Interface**: Command-line tool for encyclopedia access
- 🌐 **Web Interface**: Modern, responsive web interface
- 🐳 **Docker Ready**: Complete containerized setup
- 🔧 **Configurable**: Environment-based configuration
- 📊 **Health Monitoring**: Built-in health checks
- 🔒 **CORS Support**: Cross-origin request handling
- 🧹 **Auto-Cleaning**: Comprehensive development environment management
- 🔄 **Hot Reloading**: Development with auto-restart on file changes

## 📋 Prerequisites

- Go 1.25.1 or higher
- Ollama installed locally (for local models)
- Ollama cloud account (for cloud models)

## 🛠️ Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd agent-ollama-gin
```

2. Install dependencies:
```bash
go mod tidy
```

3. Copy environment configuration:
```bash
cp env.example .env
```

4. Configure your environment variables in `.env`:
```env
# Server Configuration
PORT=8080

# Ollama Configuration
OLLAMA_HOST=http://localhost:11434

# Ollama Cloud Configuration
LLAMA_CLOUD_ENABLED=true
LLAMA_CLOUD_API_URL=https://api.ollama.com
LLAMA_CLOUD_API_KEY=your_api_key_here
LLAMA_SIGNED_IN=false
```

## 🚀 Running the Service

1. Start the server:
```bash
go run main.go
```

2. The API will be available at `http://localhost:8080`

3. Check the health endpoint:
```bash
curl http://localhost:8080/
```

## 📚 API Endpoints

### Core Endpoints

#### Chat Completion
```bash
POST /api/v1/llama/chat
Content-Type: application/json

{
  "model": "llama3.2:1b",
  "messages": [
    {
      "role": "user",
      "content": "Hello! How are you?"
    }
  ],
  "stream": false
}
```

#### Text Completion
```bash
POST /api/v1/llama/completion
Content-Type: application/json

{
  "model": "llama3.2:1b",
  "prompt": "The future of artificial intelligence is",
  "stream": false
}
```

#### Generate Embeddings
```bash
POST /api/v1/llama/embedding
Content-Type: application/json

{
  "model": "nomic-embed-text",
  "input": "Text to generate embeddings for"
}
```

#### List Models
```bash
GET /api/v1/llama/models
```

### Streaming Endpoints

#### Streaming Chat
```bash
POST /api/v1/llama/stream/chat
Content-Type: application/json

{
  "model": "llama3.2:1b",
  "messages": [
    {
      "role": "user",
      "content": "Tell me a story"
    }
  ],
  "stream": true
}
```

### Model Management

#### Pull Model
```bash
POST /api/v1/llama/pull
Content-Type: application/json

{
  "name": "llama3.2:1b"
}
```

### Cloud Authentication

#### Sign In to Ollama Cloud
```bash
POST /api/v1/llama/cloud/signin
Content-Type: application/json

{
  "email": "your-email@example.com",
  "password": "your-password"
}
```

#### Sign Out from Ollama Cloud
```bash
POST /api/v1/llama/cloud/signout
```

#### List Cloud Models
```bash
GET /api/v1/llama/cloud/models
```

## 🧪 Testing

### Run the Test Suite

1. Start the server:
```bash
go run main.go
```

2. In another terminal, run the test script:
```bash
go run tests/ollama_cloud_test.go
```

The test script will verify:
- Server health
- Model listing
- Chat completions
- Text completions
- Embeddings generation
- Cloud authentication
- Streaming responses

### Manual Testing with cURL

Test chat completion:
```bash
curl -X POST http://localhost:8080/api/v1/llama/chat \
  -H "Content-Type: application/json" \
  -d '{
    "model": "llama3.2:1b",
    "messages": [
      {
        "role": "user",
        "content": "Hello! Tell me about Ollama."
      }
    ],
    "stream": false
  }'
```

Test streaming chat:
```bash
curl -X POST http://localhost:8080/api/v1/llama/stream/chat \
  -H "Content-Type: application/json" \
  -d '{
    "model": "llama3.2:1b",
    "messages": [
      {
        "role": "user",
        "content": "Tell me a short story"
      }
    ],
    "stream": true
  }'
```

## 🏗️ Architecture

### Project Structure
```
agent-ollama-gin/
├── config/          # Configuration management
├── handlers/        # HTTP request handlers
├── models/          # Data models and structures
├── services/        # Business logic and Ollama integration
├── tests/           # Test files
├── main.go          # Application entry point
├── go.mod           # Go module dependencies
└── env.example      # Environment configuration template
```

### Key Components

- **Config**: Manages environment variables and application configuration
- **Services**: Handles Ollama API communication and cloud integration
- **Handlers**: Processes HTTP requests and responses
- **Models**: Defines data structures for requests and responses

## 🔧 Configuration Options

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `OLLAMA_HOST` | Local Ollama host URL | `http://localhost:11434` |
| `LLAMA_CLOUD_ENABLED` | Enable cloud models | `false` |
| `LLAMA_CLOUD_API_URL` | Ollama cloud API URL | `https://api.ollama.com` |
| `LLAMA_CLOUD_API_KEY` | Your Ollama cloud API key | - |
| `LLAMA_SIGNED_IN` | Cloud authentication status | `false` |

## 🌟 Migration from Genkit

This service has been completely migrated from Google's Genkit framework to native Ollama cloud integration:

### What Changed:
- ✅ Removed all Genkit dependencies
- ✅ Implemented direct Ollama cloud API integration
- ✅ Added cloud authentication support
- ✅ Enhanced model management capabilities
- ✅ Improved streaming responses
- ✅ Added comprehensive testing

# Health check
make health-check
```

### **Hot Reloading (Recommended)**
```bash
# Install development tools
make install-tools

# Run with file watching
make watch
```

### **Docker Development**
```bash
# Start Ollama with Docker
make docker-dev

# Run full stack
make docker-run

# Stop services
make docker-stop
```

### **Cleaning & Maintenance**
```bash
# Basic cleaning
make clean

# Deep cleaning
make deep-clean

# Docker cleaning
make clean-docker

# Complete cleanup
make clean-all

# Interactive cleaning
make clean-interactive
```

### **Monitoring & Debugging**
```bash
# Check environment status
make status

# Health checks
make health-check

# View logs
make logs-api
make logs-ollama

# Model management
make list-models
make pull-model
```



## 🚨 Troubleshooting

### **Common Issues**

#### **1. Port Already in Use**
```bash
# Check what's using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>
```

#### **2. Ollama Not Starting**
```bash
# Check if Ollama is installed
which ollama

# Reinstall Ollama
curl -fsSL https://ollama.ai/install.sh | sh

# Start Ollama manually
ollama serve
```

#### **3. Go Dependencies Issues**
```bash
# Clean and reinstall
make clean
make install-deps

# Or manually
go clean -modcache
go mod download
go mod tidy
```

#### **4. Docker Issues**
```bash
# Reset Docker environment
make docker-stop
docker system prune -f
make docker-run
```

### **Debug Mode**
```bash
# Enable debug logging
export GIN_MODE=debug
make run

# Or set in .env file
LOG_LEVEL=debug
```

## 📚 Useful Aliases

After running `./llama-api.sh`, you can add these aliases to your shell:

```bash
# Add to ~/.bashrc or ~/.zshrc
source /path/to/your/project/.bashrc-additions
```

Available aliases:
- `llama-dev` - Start development environment
- `llama-stop` - Stop development environment
- `llama-test` - Test API endpoints
- `llama-build` - Build application
- `llama-run` - Run application
- `llama-test-all` - Run all tests
- `llama-clean` - Clean build artifacts
- `llama-deps` - Install dependencies
- `llama-docker` - Start Docker services
- `llama-docker-stop` - Stop Docker services
- `llama-logs` - View Docker logs

## 🔄 Development Cycle

### **Typical Development Session**
1. **Start**: `make quick-start` or `make dev`
2. **Develop**: Edit code with hot reloading (`make watch`)
3. **Test**: `make test-api` or `make test`
4. **Build**: `make build` to verify compilation
5. **Stop**: `make stop-dev` when done

### **Before Committing**
```bash
# Ensure code quality
make format
make lint
make test

# Build to verify
make build

# Clean up
make clean
```

## 🌟 Pro Tips

1. **Use hot reloading**: `make watch` for faster development
2. **Keep Ollama running**: Start it once and leave it running
3. **Use health checks**: `make health-check` to verify services
4. **Monitor logs**: Use `make logs-api` to debug issues
5. **Test frequently**: Run `make test-api` after changes
6. **Clean regularly**: Use `make clean` after builds
7. **Use interactive tools**: `./llama-api.sh` for guided development



## 🔒 Security

- CORS configuration for cross-origin requests
- Input validation and sanitization
- Rate limiting support (can be added)
- API key authentication (optional)

## 🚀 Quick Commands Reference

| Command | Description |
|---------|-------------|
| `./llama-api.sh` | Complete environment setup |
| `./llama-api.sh` | Quick development start |
| `./llama-api.sh` | Start development environment |
| `./llama-api.sh` | Stop development environment |
| `./llama-api.sh` | Test API endpoints |
| `./llama-api.sh` | Check environment status |
| `./llama-api.sh` | Interactive development menu |
| `./llama-api.sh` | Interactive cleaning menu |
| `make help` | Show all commands |
| `make build-cli` | Build encyclopedia CLI |
| `make watch` | Hot reloading development |
| `make status` | Environment status |
| `make health-check` | Service health check |
| `make clean` | Basic cleaning |
| `make deep-clean` | Deep cleaning |
| `make clean-all` | Complete cleanup |

## 🔍 Encyclopedia Agent CLI

The project now includes a powerful command-line interface for encyclopedia access:

### **Quick CLI Start**
```bash
# Build the CLI
make build-cli

# Interactive mode
./bin/encyclopedia

# Command line usage
./bin/encyclopedia search "artificial intelligence"
./bin/encyclopedia article "Machine Learning"
./bin/encyclopedia prompt "neural networks"
```

### **CLI Features**
- 🔍 **Search**: Multi-source encyclopedia search
- 📖 **Articles**: Retrieve full articles with custom length
- ✍️ **Prompts**: AI-generated encyclopedia-style prompts
- 🌍 **Multi-language**: Support for 10+ languages
- 💻 **Interactive**: Command-line interface with help system
- 🔗 **Integration**: Easy to use in scripts and automation

### **CLI Documentation**
- **Usage Guide**: `CLI_USAGE.md`
- **Demo Script**: `./demo_cli.sh`
- **Help**: `./bin/encyclopedia help`

## 🆘 Need Help?

### **Quick Start**
1. **Start here**: `./llama-api.sh` (interactive menu)
2. **Quick help**: `make help`
3. **Status check**: `make status`
4. **Quick start**: `./llama-api.sh`

### **Development Issues**
1. **Environment problems**: `make status`
2. **Service issues**: `make health-check`
3. **Cleaning needed**: `make clean-interactive`
4. **Docker issues**: `make clean-docker`

### **Documentation**
- **This README** - Complete project overview
- **`make help`** - All available commands
- **Interactive menus** - Guided development experience

## 🌟 Roadmap

- [ ] Authentication middleware
- [ ] Rate limiting
- [ ] Request/response caching
- [ ] Metrics and monitoring
- [ ] WebSocket support
- [ ] Model fine-tuning endpoints
- [ ] Batch processing
- [ ] Multi-model load balancing

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License.

## 🆘 Support

For issues and questions:
1. Check the [Ollama documentation](https://ollama.com/docs)
2. Review the [Ollama cloud models guide](https://ollama.com/blog/cloud-models)
3. Open an issue in this repository

## 🔗 Related Links

- [Ollama Official Website](https://ollama.com)
- [Ollama Cloud Models](https://ollama.com/blog/cloud-models)
- [Gin Web Framework](https://gin-gonic.com)
- [Go Documentation](https://golang.org/doc)
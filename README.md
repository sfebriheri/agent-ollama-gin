# Agent Ollama Gin - Ollama Cloud Integration

A powerful REST API service built with Go and Gin framework that provides seamless integration with Ollama cloud models. This service has been migrated from Genkit to support Ollama's cloud infrastructure for enhanced AI model access.

## 🚀 Features

- **Ollama Cloud Integration**: Full support for Ollama cloud models and authentication
- **Local & Cloud Models**: Seamlessly switch between local Ollama and cloud-hosted models
- **Chat Completions**: Interactive chat with AI models
- **Text Completions**: Generate text completions from prompts
- **Embeddings**: Generate vector embeddings for text
- **Streaming Responses**: Real-time streaming for chat and completions
- **Model Management**: List, pull, and manage both local and cloud models
- **Authentication**: Secure sign-in/sign-out for Ollama cloud services

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

### Benefits:
- **Better Performance**: Direct API calls without framework overhead
- **Cloud Support**: Native integration with Ollama cloud models
- **Flexibility**: Support for both local and cloud models
- **Scalability**: Enhanced for production deployments

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
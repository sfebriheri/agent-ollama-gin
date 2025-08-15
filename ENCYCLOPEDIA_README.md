# 🔍 Encyclopedia Agent

An AI-powered encyclopedia agent that combines Llama LLM with Wikipedia and Britannica sources for intelligent content search, retrieval, and prompt generation.

## ✨ Features

- **🔍 Multi-Source Search**: Wikipedia, Britannica, or both simultaneously
- **🌍 Multi-Language Support**: 10+ languages including EN, ES, FR, DE, IT, PT, RU, JA, ZH, AR
- **📖 Article Retrieval**: Full articles with customizable length limits
- **✍️ AI Prompt Generation**: Encyclopedia-style prompts using Llama LLM
- **🎯 Intelligent Suggestions**: Related topics and keyword extraction
- **🚀 RESTful API**: Clean, well-documented endpoints
- **💻 Modern Web Interface**: Responsive interface for easy interaction

## 🏗️ Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Web Client    │    │  Encyclopedia    │    │   Llama LLM     │
│                 │◄──►│     Service      │◄──►│   (via Ollama)  │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌──────────────────┐
                       │  Encyclopedia    │
                       │     Sources      │
                       │  (Wikipedia,     │
                       │   Britannica)    │
                       └──────────────────┘
```

## 🚀 Quick Start

### Prerequisites

- Go 1.19+
- Ollama with Llama models
- Internet connection for encyclopedia APIs

### Installation

1. **Clone and setup**
   ```bash
   git clone <repository-url>
   cd agent-ollama-gin
   cp env.example .env
   # Edit .env with your configuration
   ```

2. **Start Ollama**
   ```bash
   ollama run llama2
   ```

3. **Build and run**
   ```bash
   go build -o bin/llama-api .
   ./bin/llama-api
   ```

4. **Access web interface**
   ```
   http://localhost:8080/examples/encyclopedia_interface.html
   ```

## 📚 API Endpoints

### Encyclopedia Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/encyclopedia/health` | Service health check |
| `GET` | `/api/v1/encyclopedia/sources` | Available sources |
| `GET` | `/api/v1/encyclopedia/languages` | Supported languages |
| `POST` | `/api/v1/encyclopedia/search` | Search articles |
| `POST` | `/api/v1/encyclopedia/article` | Retrieve article |
| `POST` | `/api/v1/encyclopedia/prompt` | Generate prompts |

### Example Requests

#### Search Encyclopedia
```json
POST /api/v1/encyclopedia/search
{
  "query": "artificial intelligence",
  "source": "all",
  "max_results": 5,
  "language": "en"
}
```

#### Get Article
```json
POST /api/v1/encyclopedia/article
{
  "title": "Artificial Intelligence",
  "source": "wikipedia",
  "language": "en",
  "max_length": 2000
}
```

#### Generate Prompt
```json
POST /api/v1/encyclopedia/prompt
{
  "topic": "quantum computing",
  "style": "educational",
  "length": "medium",
  "language": "en"
}
```

## 🎯 Usage Examples

### Command Line Testing

```bash
# Run the comprehensive test suite
./test_encyclopedia.sh

# Test individual endpoints
curl -X POST http://localhost:8080/api/v1/encyclopedia/search \
  -H "Content-Type: application/json" \
  -d '{"query": "machine learning", "source": "wikipedia"}'
```

### Web Interface

The modern web interface provides:
- **Search Tab**: Multi-source encyclopedia search
- **Article Tab**: Retrieve specific articles
- **Prompt Tab**: AI-powered prompt generation
- **Info Tab**: View sources and languages

## ⚙️ Configuration

### Environment Variables

```bash
# Encyclopedia Configuration
WIKIPEDIA_API_URL=https://en.wikipedia.org/api/rest_v1
BRITANNICA_API_URL=https://api.britannica.com
BRITANNICA_API_KEY=your_api_key_here
ENCYCLOPEDIA_TIMEOUT=30

# Llama Configuration
LLAMA_BASE_URL=http://localhost:11434
LLAMA_DEFAULT_MODEL=llama2
LLAMA_TIMEOUT=60

# Server Configuration
PORT=8080
HOST=0.0.0.0
```

### Supported Sources & Languages

- **Sources**: Wikipedia (free), Britannica (API key required)
- **Languages**: EN, ES, FR, DE, IT, PT, RU, JA, ZH, AR

## 🔧 Development

### Project Structure

```
agent-ollama-gin/
├── handlers/
│   ├── llama_handler.go
│   └── encyclopedia_handler.go
├── services/
│   ├── llama_service.go
│   ├── encyclopedia_service.go
│   └── britannica_service.go
├── models/
│   └── llama_models.go
├── examples/
│   └── encyclopedia_interface.html
└── main.go
```

### Adding New Sources

1. Update models in `models/llama_models.go`
2. Add source methods in `services/encyclopedia_service.go`
3. Update handler in `handlers/encyclopedia_handler.go`
4. Add to web interface

## 🧪 Testing

### Automated Testing

```bash
# Run comprehensive test suite
./test_encyclopedia.sh

# Test web interface
open http://localhost:8080/examples/encyclopedia_interface.html
```

### Manual Testing

```bash
# Test health endpoint
curl http://localhost:8080/api/v1/encyclopedia/health

# Test search functionality
curl -X POST http://localhost:8080/api/v1/encyclopedia/search \
  -H "Content-Type: application/json" \
  -d '{"query": "test", "source": "wikipedia"}'
```

## 🚀 Deployment

### Docker

```dockerfile
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

### Production Environment

```bash
# Production settings
PORT=8080
HOST=0.0.0.0
LLAMA_BASE_URL=http://your-ollama-instance:11434
WIKIPEDIA_API_URL=https://en.wikipedia.org/api/rest_v1
BRITANNICA_API_URL=https://api.britannica.com
BRITANNICA_API_KEY=your_production_key
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

MIT License - see LICENSE file for details.

## 🔮 Future Enhancements

- [ ] Additional encyclopedia sources
- [ ] Article caching system
- [ ] Advanced search filters
- [ ] Media support (images, videos)
- [ ] Mobile application
- [ ] User authentication
- [ ] Article recommendations
- [ ] Translation support
- [ ] Browser extension

## 📞 Support

- Check [Issues](https://github.com/your-repo/issues) page
- Create new issue with environment details
- Include error messages and logs

---

**Happy exploring with your Encyclopedia Agent! 🚀**

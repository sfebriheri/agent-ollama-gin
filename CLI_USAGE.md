# 🔍 Encyclopedia Agent CLI

A powerful command-line interface for the Encyclopedia Agent that provides direct access to encyclopedia search, article retrieval, and AI prompt generation.

## 🚀 Quick Start

### Build the CLI

```bash
# Build the encyclopedia CLI
make build-cli

# Or build manually
go build -o bin/encyclopedia cmd/encyclopedia/main.go
```

### Start the Server

```bash
# Start the encyclopedia server
make run
# or
go run main.go
```

### Use the CLI

```bash
# Interactive mode
./bin/encyclopedia

# Command line mode
./bin/encyclopedia search "artificial intelligence"
```

## 📚 Available Commands

### 🔍 Search Encyclopedia

Search for articles across multiple sources:

```bash
# Basic search
./bin/encyclopedia search "machine learning"

# Search with specific source
./bin/encyclopedia search "quantum computing" wikipedia

# Search with language and results limit
./bin/encyclopedia search "artificial intelligence" all en 10
```

**Parameters:**
- `query` (required): Search term
- `source` (optional): `wikipedia`, `britannica`, or `all` (default: `all`)
- `language` (optional): Language code (default: `en`)
- `max_results` (optional): Maximum number of results (default: `5`)

### 📖 Get Article

Retrieve a specific article by title:

```bash
# Basic article retrieval
./bin/encyclopedia article "Machine Learning"

# Article with specific source and language
./bin/encyclopedia article "Artificial Intelligence" wikipedia en

# Article with custom length limit
./bin/encyclopedia article "Neural Networks" wikipedia en 1500
```

**Parameters:**
- `title` (required): Article title
- `source` (optional): `wikipedia` or `britannica` (default: `wikipedia`)
- `language` (optional): Language code (default: `en`)
- `max_length` (optional): Maximum content length (default: `2000`)

### ✍️ Generate Prompt

Generate AI-powered encyclopedia-style prompts:

```bash
# Basic prompt generation
./bin/encyclopedia prompt "deep learning"

# Prompt with specific style and length
./bin/encyclopedia prompt "computer vision" academic long

# Prompt with custom language
./bin/encyclopedia prompt "natural language processing" educational medium es
```

**Parameters:**
- `topic` (required): Topic for prompt generation
- `style` (optional): `academic`, `casual`, or `educational` (default: `educational`)
- `length` (optional): `short`, `medium`, or `long` (default: `medium`)
- `language` (optional): Language code (default: `en`)

### 📚 Information Commands

```bash
# Show available sources
./bin/encyclopedia sources

# Show supported languages
./bin/encyclopedia languages

# Check service health
./bin/encyclopedia health

# Open web interface
./bin/encyclopedia web

# Show help
./bin/encyclopedia help
```

## 🎯 Interactive Mode

Run the CLI without arguments to enter interactive mode:

```bash
./bin/encyclopedia
```

**Interactive Commands:**
- `search <query> [source] [language] [max_results]`
- `article <title> [source] [language] [max_length]`
- `prompt <topic> [style] [length] [language]`
- `sources` - Show available sources
- `languages` - Show supported languages
- `health` - Check service health
- `web` - Open web interface
- `clear` - Clear screen
- `help` - Show help
- `quit` or `exit` - Exit CLI

## 🌍 Supported Languages

- **EN** - English (default)
- **ES** - Spanish
- **FR** - French
- **DE** - German
- **IT** - Italian
- **PT** - Portuguese
- **RU** - Russian
- **JA** - Japanese
- **ZH** - Chinese
- **AR** - Arabic

## 📊 Output Examples

### Search Results
```
🔍 Search Results for: artificial intelligence
📊 Found 5 results from all

1. 📖 Artificial Intelligence
   🌐 https://en.wikipedia.org/wiki/Artificial_intelligence
   📝 Artificial intelligence (AI) is intelligence demonstrated by machines...
   🏷️  wikipedia (en)
   ⭐ Relevance: 0.95

2. 📖 Machine Learning
   🌐 https://en.wikipedia.org/wiki/Machine_learning
   📝 Machine learning is a subset of artificial intelligence...
   🏷️  wikipedia (en)
   ⭐ Relevance: 0.92
```

### Article Content
```
📖 Article: Machine Learning
🌐 Source: wikipedia (en)
📅 Updated: 2024-01-15
📊 Word Count: 1500
🏷️  Categories: artificial intelligence, computer science
🔗 URL: https://en.wikipedia.org/wiki/Machine_learning

📝 Summary:
Machine learning is a subset of artificial intelligence...

📄 Content:
Machine learning (ML) is a subset of artificial intelligence...
```

### Generated Prompt
```
✍️  Generated Prompt for: neural networks
🎨 Style: educational
📏 Length: medium
🌍 Language: en

📝 Prompt:
Neural networks are computational models inspired by biological...

💡 Suggestions:
   • History of neural networks
   • Modern developments in neural networks
   • Key figures in neural networks

🔑 Keywords:
   • neural
   • networks
   • learning
   • artificial
```

## ⚙️ Configuration

The CLI connects to the Encyclopedia Agent server running on `localhost:8080`. Make sure:

1. **Server is running**: `go run main.go`
2. **Ollama is available**: `ollama serve`
3. **Models are pulled**: `ollama pull llama2`

## 🔧 Advanced Usage

### Batch Processing

```bash
# Search multiple topics
for topic in "AI" "ML" "DL"; do
    ./bin/encyclopedia search "$topic" wikipedia en 3
done
```

### Integration with Other Tools

```bash
# Save search results to file
./bin/encyclopedia search "quantum computing" > quantum_results.txt

# Generate prompt and pipe to another tool
./bin/encyclopedia prompt "blockchain" academic long | grep -i "technology"
```

### Custom Scripts

```bash
#!/bin/bash
# encyclopedia_research.sh

TOPIC="$1"
SOURCE="${2:-wikipedia}"
LANGUAGE="${3:-en}"

echo "🔬 Researching: $TOPIC"
echo "======================"

# Search for articles
./bin/encyclopedia search "$TOPIC" "$SOURCE" "$LANGUAGE" 5

# Get main article
./bin/encyclopedia article "$TOPIC" "$SOURCE" "$LANGUAGE" 2000

# Generate research prompt
./bin/encyclopedia prompt "$TOPIC" academic long "$LANGUAGE"
```

## 🐛 Troubleshooting

### Common Issues

1. **"Encyclopedia service is not running"**
   - Start the server: `go run main.go`

2. **"Connection refused"**
   - Check if server is on port 8080
   - Verify firewall settings

3. **"API error 500"**
   - Check server logs
   - Verify Ollama is running
   - Check API key configuration

### Debug Mode

```bash
# Run with verbose output
go run cmd/encyclopedia/main.go search "test" 2>&1 | tee debug.log
```

## 🚀 Performance Tips

1. **Use specific sources** when you know which encyclopedia you want
2. **Limit results** for faster searches
3. **Use appropriate content length** for articles
4. **Batch operations** for multiple queries

## 📝 Examples

### Research Workflow

```bash
# 1. Explore a topic
./bin/encyclopedia search "quantum computing" all en 10

# 2. Get detailed article
./bin/encyclopedia article "Quantum Computing" wikipedia en 3000

# 3. Generate research prompt
./bin/encyclopedia prompt "quantum algorithms" academic long en

# 4. Check related topics
./bin/encyclopedia search "quantum cryptography" wikipedia en 5
```

### Educational Content Creation

```bash
# Generate educational prompts for different styles
for style in academic casual educational; do
    ./bin/encyclopedia prompt "climate change" "$style" medium en
done
```

### Multi-language Research

```bash
# Compare articles in different languages
for lang in en es fr; do
    ./bin/encyclopedia search "artificial intelligence" wikipedia "$lang" 3
done
```

---

**Happy exploring with the Encyclopedia Agent CLI! 🚀**

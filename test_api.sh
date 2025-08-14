#!/bin/bash

# Test script for Llama API
# Make sure the API is running on localhost:8080

BASE_URL="http://localhost:8080/api/v1"

echo "🧪 Testing Llama API endpoints..."
echo "=================================="

# Test health check
echo -e "\n1. Testing health check..."
curl -s "${BASE_URL}/health" | jq '.'

# Test list models
echo -e "\n2. Testing list models..."
curl -s "${BASE_URL}/llama/models" | jq '.'

# Test chat completion
echo -e "\n3. Testing chat completion..."
curl -s -X POST "${BASE_URL}/llama/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [
      {"role": "user", "content": "Hello! Can you tell me a short joke?"}
    ],
    "model": "llama2",
    "temperature": 0.7,
    "max_tokens": 100
  }' | jq '.'

# Test text completion
echo -e "\n4. Testing text completion..."
curl -s -X POST "${BASE_URL}/llama/completion" \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "The future of artificial intelligence is",
    "model": "llama2",
    "temperature": 0.8,
    "max_tokens": 50
  }' | jq '.'

# Test embedding
echo -e "\n5. Testing text embedding..."
curl -s -X POST "${BASE_URL}/llama/embedding" \
  -H "Content-Type: application/json" \
  -d '{
    "input": "This is a sample text for embedding generation",
    "model": "llama2"
  }' | jq '.'

echo -e "\n✅ API testing completed!"
echo -e "\nNote: Make sure Ollama is running and you have the llama2 model installed:"
echo "   ollama pull llama2"

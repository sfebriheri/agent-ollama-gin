#!/usr/bin/env python3
"""
Python client example for Llama API
Demonstrates how to interact with the API endpoints
"""

import requests
import json
import time
from typing import List, Dict, Any

class LlamaAPIClient:
    def __init__(self, base_url: str = "http://localhost:8080"):
        self.base_url = base_url.rstrip('/')
        self.api_base = f"{self.base_url}/api/v1"
        
    def health_check(self) -> Dict[str, Any]:
        """Check API health status"""
        response = requests.get(f"{self.api_base}/health")
        response.raise_for_status()
        return response.json()
    
    def list_models(self) -> Dict[str, Any]:
        """List available models"""
        response = requests.get(f"{self.api_base}/llama/models")
        response.raise_for_status()
        return response.json()
    
    def chat_completion(self, messages: List[Dict[str, str]], 
                       model: str = "llama2", 
                       temperature: float = 0.7,
                       max_tokens: int = 100) -> Dict[str, Any]:
        """Send chat completion request"""
        payload = {
            "messages": messages,
            "model": model,
            "temperature": temperature,
            "max_tokens": max_tokens
        }
        
        response = requests.post(
            f"{self.api_base}/llama/chat",
            json=payload,
            headers={"Content-Type": "application/json"}
        )
        response.raise_for_status()
        return response.json()
    
    def text_completion(self, prompt: str, 
                       model: str = "llama2",
                       temperature: float = 0.8,
                       max_tokens: int = 50) -> Dict[str, Any]:
        """Send text completion request"""
        payload = {
            "prompt": prompt,
            "model": model,
            "temperature": temperature,
            "max_tokens": max_tokens
        }
        
        response = requests.post(
            f"{self.api_base}/llama/completion",
            json=payload,
            headers={"Content-Type": "application/json"}
        )
        response.raise_for_status()
        return response.json()
    
    def generate_embedding(self, text: str, model: str = "llama2") -> Dict[str, Any]:
        """Generate text embedding"""
        payload = {
            "input": text,
            "model": model
        }
        
        response = requests.post(
            f"{self.api_base}/llama/embedding",
            json=payload,
            headers={"Content-Type": "application/json"}
        )
        response.raise_for_status()
        return response.json()

def main():
    """Example usage of the Llama API client"""
    print("🚀 Llama API Python Client Example")
    print("=" * 40)
    
    # Initialize client
    client = LlamaAPIClient()
    
    try:
        # Health check
        print("\n1. Health Check:")
        health = client.health_check()
        print(f"   Status: {health['status']}")
        print(f"   Message: {health['message']}")
        
        # List models
        print("\n2. Available Models:")
        models = client.list_models()
        for model in models.get('models', []):
            print(f"   - {model['id']}")
        
        # Chat completion
        print("\n3. Chat Completion:")
        messages = [
            {"role": "user", "content": "Hello! Can you tell me a short joke?"}
        ]
        chat_response = client.chat_completion(messages, temperature=0.7)
        print(f"   Response: {chat_response['choices'][0]['message']['content']}")
        
        # Text completion
        print("\n4. Text Completion:")
        prompt = "The future of artificial intelligence is"
        completion_response = client.text_completion(prompt, temperature=0.8)
        print(f"   Prompt: {prompt}")
        print(f"   Completion: {completion_response['choices'][0]['message']['content']}")
        
        # Embedding
        print("\n5. Text Embedding:")
        text = "This is a sample text for embedding generation"
        embedding_response = client.generate_embedding(text)
        embedding_vector = embedding_response['data'][0]['embedding']
        print(f"   Text: {text}")
        print(f"   Embedding dimensions: {len(embedding_vector)}")
        print(f"   First 5 values: {embedding_vector[:5]}")
        
        print("\n✅ All API calls completed successfully!")
        
    except requests.exceptions.ConnectionError:
        print("❌ Error: Could not connect to the API. Make sure it's running on localhost:8080")
    except requests.exceptions.HTTPError as e:
        print(f"❌ HTTP Error: {e}")
    except Exception as e:
        print(f"❌ Unexpected error: {e}")

if __name__ == "__main__":
    main()

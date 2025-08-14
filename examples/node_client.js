#!/usr/bin/env node
/**
 * Node.js client example for Llama API
 * Demonstrates how to interact with the API endpoints
 */

const axios = require('axios');

class LlamaAPIClient {
    constructor(baseURL = 'http://localhost:8080') {
        this.baseURL = baseURL.replace(/\/$/, '');
        this.apiBase = `${this.baseURL}/api/v1`;
    }

    /**
     * Check API health status
     */
    async healthCheck() {
        try {
            const response = await axios.get(`${this.apiBase}/health`);
            return response.data;
        } catch (error) {
            throw new Error(`Health check failed: ${error.message}`);
        }
    }

    /**
     * List available models
     */
    async listModels() {
        try {
            const response = await axios.get(`${this.apiBase}/llama/models`);
            return response.data;
        } catch (error) {
            throw new Error(`Failed to list models: ${error.message}`);
        }
    }

    /**
     * Send chat completion request
     */
    async chatCompletion(messages, model = 'llama2', temperature = 0.7, maxTokens = 100) {
        try {
            const payload = {
                messages,
                model,
                temperature,
                max_tokens: maxTokens
            };

            const response = await axios.post(`${this.apiBase}/llama/chat`, payload, {
                headers: { 'Content-Type': 'application/json' }
            });
            return response.data;
        } catch (error) {
            throw new Error(`Chat completion failed: ${error.message}`);
        }
    }

    /**
     * Send text completion request
     */
    async textCompletion(prompt, model = 'llama2', temperature = 0.8, maxTokens = 50) {
        try {
            const payload = {
                prompt,
                model,
                temperature,
                max_tokens: maxTokens
            };

            const response = await axios.post(`${this.apiBase}/llama/completion`, payload, {
                headers: { 'Content-Type': 'application/json' }
            });
            return response.data;
        } catch (error) {
            throw new Error(`Text completion failed: ${error.message}`);
        }
    }

    /**
     * Generate text embedding
     */
    async generateEmbedding(text, model = 'llama2') {
        try {
            const payload = {
                input: text,
                model
            };

            const response = await axios.post(`${this.apiBase}/llama/embedding`, payload, {
                headers: { 'Content-Type': 'application/json' }
            });
            return response.data;
        } catch (error) {
            throw new Error(`Embedding generation failed: ${error.message}`);
        }
    }

    /**
     * Stream chat responses
     */
    async streamChat(messages, model = 'llama2', temperature = 0.7) {
        try {
            const payload = {
                messages,
                model,
                temperature,
                stream: true
            };

            const response = await axios.post(`${this.apiBase}/llama/stream-chat`, payload, {
                headers: { 'Content-Type': 'application/json' },
                responseType: 'stream'
            });

            return response.data;
        } catch (error) {
            throw new Error(`Stream chat failed: ${error.message}`);
        }
    }
}

/**
 * Example usage of the Llama API client
 */
async function main() {
    console.log('🚀 Llama API Node.js Client Example');
    console.log('=' .repeat(40));

    // Initialize client
    const client = new LlamaAPIClient();

    try {
        // Health check
        console.log('\n1. Health Check:');
        const health = await client.healthCheck();
        console.log(`   Status: ${health.status}`);
        console.log(`   Message: ${health.message}`);

        // List models
        console.log('\n2. Available Models:');
        const models = await client.listModels();
        if (models.models && models.models.length > 0) {
            models.models.forEach(model => {
                console.log(`   - ${model.id}`);
            });
        } else {
            console.log('   No models found');
        }

        // Chat completion
        console.log('\n3. Chat Completion:');
        const messages = [
            { role: 'user', content: 'Hello! Can you tell me a short joke?' }
        ];
        const chatResponse = await client.chatCompletion(messages, 'llama2', 0.7);
        console.log(`   Response: ${chatResponse.choices[0].message.content}`);

        // Text completion
        console.log('\n4. Text Completion:');
        const prompt = 'The future of artificial intelligence is';
        const completionResponse = await client.textCompletion(prompt, 'llama2', 0.8);
        console.log(`   Prompt: ${prompt}`);
        console.log(`   Completion: ${completionResponse.choices[0].message.content}`);

        // Embedding
        console.log('\n5. Text Embedding:');
        const text = 'This is a sample text for embedding generation';
        const embeddingResponse = await client.generateEmbedding(text, 'llama2');
        const embeddingVector = embeddingResponse.data[0].embedding;
        console.log(`   Text: ${text}`);
        console.log(`   Embedding dimensions: ${embeddingVector.length}`);
        console.log(`   First 5 values: [${embeddingVector.slice(0, 5).join(', ')}]`);

        console.log('\n✅ All API calls completed successfully!');

    } catch (error) {
        console.error(`❌ Error: ${error.message}`);
        if (error.code === 'ECONNREFUSED') {
            console.log('\n💡 Make sure the API is running on localhost:8080');
        }
    }
}

// Run the example if this file is executed directly
if (require.main === module) {
    main().catch(console.error);
}

module.exports = LlamaAPIClient;

package models

import "time"

// Message represents a chat message
type Message struct {
	Role    string `json:"role" binding:"required"` // "system", "user", "assistant"
	Content string `json:"content" binding:"required"`
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Messages    []Message `json:"messages" binding:"required"`
	Model       string    `json:"model,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

// ChatResponse represents a chat completion response
type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice represents a completion choice
type Choice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
	Delta   Message `json:"delta,omitempty"` // For streaming
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// CompletionRequest represents a text completion request
type CompletionRequest struct {
	Prompt      string  `json:"prompt" binding:"required"`
	Model       string  `json:"model,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
	MaxTokens   int     `json:"max_tokens,omitempty"`
	Stop        string  `json:"stop,omitempty"`
}

// CompletionResponse represents a text completion response
type CompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// EmbeddingRequest represents an embedding request
type EmbeddingRequest struct {
	Input string `json:"input" binding:"required"`
	Model string `json:"model,omitempty"`
}

// EmbeddingResponse represents an embedding response
type EmbeddingResponse struct {
	Object string      `json:"object"`
	Data   []Embedding `json:"data"`
	Model  string      `json:"model"`
	Usage  Usage       `json:"usage"`
}

// Embedding represents a text embedding
type Embedding struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

// Model represents a Llama model
type Model struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
	IsCloud bool   `json:"is_cloud,omitempty"`
	Size    string `json:"size,omitempty"`
}

// CloudModel represents available Ollama cloud models
type CloudModel struct {
	Name        string `json:"name"`
	ID          string `json:"id"`
	Size        string `json:"size"`
	Description string `json:"description"`
	Available   bool   `json:"available"`
}

// AuthRequest represents authentication request for Ollama cloud
type AuthRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token,omitempty"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token,omitempty"`
	Message string `json:"message"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
	Code    int    `json:"code"`
}

// StreamResponse represents a streaming response chunk
type StreamResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
}

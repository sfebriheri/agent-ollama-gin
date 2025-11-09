package services

import "agent-ollama-gin/models"

// LlamaServiceInterface defines the interface for Llama service operations
type LlamaServiceInterface interface {
	Chat(request models.ChatRequest) (*models.ChatResponse, error)
	Completion(request models.CompletionRequest) (*models.CompletionResponse, error)
	Embedding(request models.EmbeddingRequest) (*models.EmbeddingResponse, error)
	ListModels() ([]models.Model, error)
	SignIn(username, password string) (*models.AuthResponse, error)
	SignOut() error
	PullModel(modelName string) error
	StreamChat(request models.ChatRequest, responseChan chan<- string)
}

// Ensure LlamaService implements the interface
var _ LlamaServiceInterface = (*LlamaService)(nil)

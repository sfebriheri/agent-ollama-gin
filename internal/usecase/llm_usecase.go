package usecase

import (
	"context"

	"llama-api/internal/domain"
	"llama-api/internal/infrastructure"
	"llama-api/pkg/logger"
)

// LLMUsecase implements domain.LLMUsecase interface
type LLMUsecase struct {
	client *infrastructure.OllamaClient
	cache  domain.Cache
	logger *logger.Logger
}

// NewLLMUsecase creates a new LLM use case
func NewLLMUsecase(client *infrastructure.OllamaClient, cache domain.Cache, logger *logger.Logger) domain.LLMUsecase {
	return &LLMUsecase{
		client: client,
		cache:  cache,
		logger: logger,
	}
}

// Chat performs a chat completion
func (u *LLMUsecase) Chat(ctx context.Context, request domain.LLMRequest) (*domain.LLMResponse, error) {
	// Try to get from cache first
	cacheKey := "chat:" + hashRequest(request)
	if cached, err := u.cache.Get(ctx, cacheKey); err == nil && cached != nil {
		if response, ok := cached.(*domain.LLMResponse); ok {
			u.logger.Debug("cache hit for chat request", cacheKey)
			return response, nil
		}
	}

	// Call the service
	response, err := u.client.Chat(ctx, request)
	if err != nil {
		return nil, err
	}

	// Cache the response for 1 hour
	_ = u.cache.Set(ctx, cacheKey, response, 3600)

	return response, nil
}

// Completion performs a text completion
func (u *LLMUsecase) Completion(ctx context.Context, request domain.CompletionRequest) (*domain.CompletionResponse, error) {
	// Try to get from cache first
	cacheKey := "completion:" + hashCompletionRequest(request)
	if cached, err := u.cache.Get(ctx, cacheKey); err == nil && cached != nil {
		if response, ok := cached.(*domain.CompletionResponse); ok {
			u.logger.Debug("cache hit for completion request", cacheKey)
			return response, nil
		}
	}

	// Call the service
	response, err := u.client.Completion(ctx, request)
	if err != nil {
		return nil, err
	}

	// Cache the response for 1 hour
	_ = u.cache.Set(ctx, cacheKey, response, 3600)

	return response, nil
}

// Embedding performs an embedding request
func (u *LLMUsecase) Embedding(ctx context.Context, request domain.EmbeddingRequest) (*domain.EmbeddingResponse, error) {
	// Try to get from cache first
	cacheKey := "embedding:" + request.Input
	if cached, err := u.cache.Get(ctx, cacheKey); err == nil && cached != nil {
		if response, ok := cached.(*domain.EmbeddingResponse); ok {
			u.logger.Debug("cache hit for embedding request", cacheKey)
			return response, nil
		}
	}

	// Call the service
	response, err := u.client.Embedding(ctx, request)
	if err != nil {
		return nil, err
	}

	// Cache the response for 24 hours (embeddings don't change often)
	_ = u.cache.Set(ctx, cacheKey, response, 86400)

	return response, nil
}

// ListModels returns available models
func (u *LLMUsecase) ListModels(ctx context.Context) ([]domain.LLMModel, error) {
	// Try to get from cache first
	cacheKey := "models:all"
	if cached, err := u.cache.Get(ctx, cacheKey); err == nil && cached != nil {
		if models, ok := cached.([]domain.LLMModel); ok {
			u.logger.Debug("cache hit for models list", cacheKey)
			return models, nil
		}
	}

	// Call the service
	models, err := u.client.ListModels(ctx)
	if err != nil {
		return nil, err
	}

	// Cache the response for 1 hour
	_ = u.cache.Set(ctx, cacheKey, models, 3600)

	return models, nil
}

// StreamChat streams chat responses
func (u *LLMUsecase) StreamChat(ctx context.Context, request domain.LLMRequest, responseChan chan<- string) error {
	return u.client.StreamChat(ctx, request, responseChan)
}

// Helper functions

func hashRequest(req domain.LLMRequest) string {
	// Simple hash - in production use crypto/sha256
	hash := ""
	for _, msg := range req.Messages {
		hash += msg.Role + ":" + msg.Content + "|"
	}
	hash += req.Model
	return hash
}

func hashCompletionRequest(req domain.CompletionRequest) string {
	return req.Prompt + ":" + req.Model
}

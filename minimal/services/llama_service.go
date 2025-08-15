package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"llama-api/models"
)

// OptimizedLlamaService provides better performance and memory management
type OptimizedLlamaService struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	modelCache map[string]bool
	cacheMutex sync.RWMutex
}

// NewOptimizedLlamaService creates an optimized service instance
func NewOptimizedLlamaService() *OptimizedLlamaService {
	baseURL := os.Getenv("LLAMA_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	apiKey := os.Getenv("LLAMA_API_KEY")

	// Get timeout from environment or use default
	timeout := 120 * time.Second
	if timeoutStr := os.Getenv("LLAMA_TIMEOUT"); timeoutStr != "" {
		if timeoutSec, err := strconv.Atoi(timeoutStr); err == nil {
			timeout = time.Duration(timeoutSec) * time.Second
		}
	}

	return &OptimizedLlamaService{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		modelCache: make(map[string]bool),
	}
}

// Chat handles chat completion with optimized memory usage
func (s *OptimizedLlamaService) Chat(request models.ChatRequest) (*models.ChatResponse, error) {
	// Validate model availability
	if !s.isModelAvailable(request.Model) {
		return nil, fmt.Errorf("model %s not available", request.Model)
	}

	// Convert to Ollama format
	ollamaRequest := map[string]interface{}{
		"model":    s.getModel(request.Model),
		"messages": request.Messages,
		"stream":   false,
	}

	if request.Temperature > 0 {
		ollamaRequest["temperature"] = request.Temperature
	}

	// Make request with context for better timeout handling
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	resp, err := s.makeRequestWithContext(ctx, "POST", "/api/chat", ollamaRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to make chat request: %w", err)
	}
	defer resp.Body.Close()

	// Parse response with memory-efficient streaming
	var ollamaResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to our format
	response := &models.ChatResponse{
		ID:      generateID(),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   s.getModel(request.Model),
		Choices: []models.Choice{
			{
				Index: 0,
				Message: models.Message{
					Role:    "assistant",
					Content: ollamaResp["message"].(map[string]interface{})["content"].(string),
				},
			},
		},
		Usage: models.Usage{
			PromptTokens:     int(ollamaResp["prompt_eval_count"].(float64)),
			CompletionTokens: int(ollamaResp["eval_count"].(float64)),
			TotalTokens:      int(ollamaResp["prompt_eval_count"].(float64) + ollamaResp["eval_count"].(float64)),
		},
	}

	return response, nil
}

// Completion handles text completion with optimization
func (s *OptimizedLlamaService) Completion(request models.CompletionRequest) (*models.CompletionResponse, error) {
	// Validate model availability
	if !s.isModelAvailable(request.Model) {
		return nil, fmt.Errorf("model %s not available", request.Model)
	}

	ollamaRequest := map[string]interface{}{
		"model":  s.getModel(request.Model),
		"prompt": request.Prompt,
		"stream": false,
	}

	if request.Temperature > 0 {
		ollamaRequest["temperature"] = request.Temperature
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	resp, err := s.makeRequestWithContext(ctx, "POST", "/api/generate", ollamaRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to make completion request: %w", err)
	}
	defer resp.Body.Close()

	var ollamaResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	response := &models.CompletionResponse{
		ID:      generateID(),
		Object:  "text_completion",
		Created: time.Now().Unix(),
		Model:   s.getModel(request.Model),
		Choices: []models.Choice{
			{
				Index: 0,
				Message: models.Message{
					Role:    "assistant",
					Content: ollamaResp["response"].(string),
				},
			},
		},
		Usage: models.Usage{
			PromptTokens:     int(ollamaResp["prompt_eval_count"].(float64)),
			CompletionTokens: int(ollamaResp["eval_count"].(float64)),
			TotalTokens:      int(ollamaResp["prompt_eval_count"].(float64) + ollamaResp["eval_count"].(float64)),
		},
	}

	return response, nil
}

// Embedding handles text embedding with optimization
func (s *OptimizedLlamaService) Embedding(request models.EmbeddingRequest) (*models.EmbeddingResponse, error) {
	ollamaRequest := map[string]interface{}{
		"model":  s.getModel(request.Model),
		"prompt": request.Input,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := s.makeRequestWithContext(ctx, "POST", "/api/embeddings", ollamaRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to make embedding request: %w", err)
	}
	defer resp.Body.Close()

	var ollamaResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	embedding := models.Embedding{
		Object:    "embedding",
		Index:     0,
		Embedding: convertToFloat64Slice(ollamaResp["embedding"].([]interface{})),
	}

	response := &models.EmbeddingResponse{
		Object: "list",
		Data:   []models.Embedding{embedding},
		Model:  s.getModel(request.Model),
		Usage: models.Usage{
			PromptTokens: 1,
			TotalTokens:  1,
		},
	}

	return response, nil
}

// ListModels returns available models with caching
func (s *OptimizedLlamaService) ListModels() ([]models.Model, error) {
	resp, err := s.makeRequest("GET", "/api/tags", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}
	defer resp.Body.Close()

	var ollamaResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	modelList := []models.Model{}
	if modelsData, ok := ollamaResp["models"].([]interface{}); ok {
		for _, modelData := range modelsData {
			model := modelData.(map[string]interface{})
			modelName := model["name"].(string)
			modelList = append(modelList, models.Model{
				ID:      modelName,
				Object:  "model",
				Created: time.Now().Unix(),
				OwnedBy: "ollama",
			})

			// Cache model availability
			s.cacheModel(modelName)
		}
	}

	return modelList, nil
}

// Helper methods
func (s *OptimizedLlamaService) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	return s.makeRequestWithContext(context.Background(), method, endpoint, body)
}

func (s *OptimizedLlamaService) makeRequestWithContext(ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, s.baseURL+endpoint, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if s.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiKey)
	}

	return s.httpClient.Do(req)
}

func (s *OptimizedLlamaService) getModel(requestedModel string) string {
	if requestedModel != "" {
		return requestedModel
	}
	return os.Getenv("LLAMA_DEFAULT_MODEL")
}

func (s *OptimizedLlamaService) isModelAvailable(modelName string) bool {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()
	return s.modelCache[modelName]
}

func (s *OptimizedLlamaService) cacheModel(modelName string) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()
	s.modelCache[modelName] = true
}

// Helper functions
func generateID() string {
	return "chatcmpl-" + strconv.FormatInt(time.Now().UnixNano(), 10)
}

func convertToFloat64Slice(interfaceSlice []interface{}) []float64 {
	result := make([]float64, len(interfaceSlice))
	for i, v := range interfaceSlice {
		result[i] = v.(float64)
	}
	return result
}

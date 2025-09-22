package services

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"agent-ollama-gin/config"
	"agent-ollama-gin/models"
)

type LlamaService struct {
	config     *config.LlamaConfig
	httpClient *http.Client
	isSignedIn bool
}

// Available cloud models based on Ollama cloud documentation
var CloudModels = []models.CloudModel{
	{
		Name:        "qwen3-coder:480b-cloud",
		ID:          "qwen3-coder:480b-cloud",
		Size:        "480B",
		Description: "Qwen3 Coder model optimized for code generation",
		Available:   true,
	},
	{
		Name:        "gpt-oss:120b-cloud",
		ID:          "gpt-oss:120b-cloud",
		Size:        "120B",
		Description: "GPT OSS large model for general purpose tasks",
		Available:   true,
	},
	{
		Name:        "gpt-oss:20b-cloud",
		ID:          "gpt-oss:20b-cloud",
		Size:        "20B",
		Description: "GPT OSS medium model for efficient processing",
		Available:   true,
	},
	{
		Name:        "deepseek-v3.1:671b-cloud",
		ID:          "deepseek-v3.1:671b-cloud",
		Size:        "671B",
		Description: "DeepSeek v3.1 ultra-large model for complex reasoning",
		Available:   true,
	},
}

func NewLlamaService() *LlamaService {
	cfg := config.Load()

	// Get timeout from environment or use default
	timeout := time.Duration(cfg.Llama.Timeout) * time.Second

	service := &LlamaService{
		config: &cfg.Llama,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		isSignedIn: cfg.Llama.SignedIn,
	}

	// Auto-signin if cloud is enabled and credentials are available
	if cfg.Llama.CloudEnabled && cfg.Llama.CloudAPIKey != "" {
		service.isSignedIn = true
	}

	return service
}

// SignIn authenticates with Ollama cloud
func (s *LlamaService) SignIn(username, password string) (*models.AuthResponse, error) {
	if !s.config.CloudEnabled {
		return &models.AuthResponse{
			Success: true,
			Message: "Cloud mode is not enabled",
		}, nil
	}

	// For now, we'll simulate sign-in since the actual Ollama cloud auth API isn't fully documented
	// In a real implementation, this would make an actual API call to ollama.com
	if username != "" && password != "" {
		s.isSignedIn = true
		return &models.AuthResponse{
			Success: true,
			Token:   "simulated-token",
			Message: "Successfully signed in to Ollama cloud",
		}, nil
	}

	return &models.AuthResponse{
		Success: false,
		Message: "Invalid credentials",
	}, nil
}

// SignOut signs out from Ollama cloud
func (s *LlamaService) SignOut() error {
	s.isSignedIn = false
	return nil
}

// IsCloudModel checks if a model is a cloud model
func (s *LlamaService) IsCloudModel(modelName string) bool {
	return strings.HasSuffix(modelName, "-cloud")
}

// PullModel pulls a model (cloud or local)
func (s *LlamaService) PullModel(modelName string) error {
	if s.IsCloudModel(modelName) && !s.isSignedIn {
		return fmt.Errorf("must be signed in to use cloud models")
	}

	pullRequest := map[string]interface{}{
		"name": modelName,
	}

	baseURL := s.config.BaseURL
	if s.IsCloudModel(modelName) && s.config.CloudEnabled {
		baseURL = s.config.CloudAPIURL
	}

	resp, err := s.makeRequest("POST", "/api/pull", pullRequest, baseURL)
	if err != nil {
		return fmt.Errorf("failed to pull model: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

// Chat handles chat completion using Ollama (local or cloud)
func (s *LlamaService) Chat(request models.ChatRequest) (*models.ChatResponse, error) {
	model := s.getModel(request.Model)

	// Check if cloud model and authentication
	if s.IsCloudModel(model) && !s.isSignedIn {
		return nil, fmt.Errorf("must be signed in to use cloud model: %s", model)
	}

	// Convert to Ollama format
	ollamaRequest := map[string]interface{}{
		"model":    model,
		"messages": request.Messages,
		"stream":   false,
	}

	if request.Temperature > 0 {
		ollamaRequest["temperature"] = request.Temperature
	}
	if request.MaxTokens > 0 {
		ollamaRequest["max_tokens"] = request.MaxTokens
	}

	// Determine which API to use
	baseURL := s.config.BaseURL
	if s.IsCloudModel(model) && s.config.CloudEnabled {
		baseURL = s.config.CloudAPIURL
	}

	// Make request to Ollama
	resp, err := s.makeRequest("POST", "/api/chat", ollamaRequest, baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make chat request: %w", err)
	}
	defer resp.Body.Close()

	// Parse Ollama response
	var ollamaResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to our format
	response := &models.ChatResponse{
		ID:      generateID(),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []models.Choice{
			{
				Index: 0,
				Message: models.Message{
					Role:    "assistant",
					Content: s.extractContent(ollamaResp),
				},
			},
		},
		Usage: s.extractUsage(ollamaResp),
	}

	return response, nil
}

// Completion handles text completion using Ollama
func (s *LlamaService) Completion(request models.CompletionRequest) (*models.CompletionResponse, error) {
	model := s.getModel(request.Model)

	// Check if cloud model and authentication
	if s.IsCloudModel(model) && !s.isSignedIn {
		return nil, fmt.Errorf("must be signed in to use cloud model: %s", model)
	}

	// Convert to Ollama format
	ollamaRequest := map[string]interface{}{
		"model":  model,
		"prompt": request.Prompt,
		"stream": false,
	}

	if request.Temperature > 0 {
		ollamaRequest["temperature"] = request.Temperature
	}
	if request.MaxTokens > 0 {
		ollamaRequest["max_tokens"] = request.MaxTokens
	}
	if request.Stop != "" {
		ollamaRequest["stop"] = request.Stop
	}

	// Determine which API to use
	baseURL := s.config.BaseURL
	if s.IsCloudModel(model) && s.config.CloudEnabled {
		baseURL = s.config.CloudAPIURL
	}

	// Make request to Ollama
	resp, err := s.makeRequest("POST", "/api/generate", ollamaRequest, baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make completion request: %w", err)
	}
	defer resp.Body.Close()

	// Parse Ollama response
	var ollamaResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to our format
	response := &models.CompletionResponse{
		ID:      generateID(),
		Object:  "text_completion",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []models.Choice{
			{
				Index: 0,
				Message: models.Message{
					Role:    "assistant",
					Content: s.extractResponse(ollamaResp),
				},
			},
		},
		Usage: s.extractUsage(ollamaResp),
	}

	return response, nil
}

// Embedding handles embedding generation using Ollama
func (s *LlamaService) Embedding(request models.EmbeddingRequest) (*models.EmbeddingResponse, error) {
	model := s.getModel(request.Model)

	// Check if cloud model and authentication
	if s.IsCloudModel(model) && !s.isSignedIn {
		return nil, fmt.Errorf("must be signed in to use cloud model: %s", model)
	}

	// Convert to Ollama format
	ollamaRequest := map[string]interface{}{
		"model":  model,
		"prompt": request.Input,
	}

	// Determine which API to use
	baseURL := s.config.BaseURL
	if s.IsCloudModel(model) && s.config.CloudEnabled {
		baseURL = s.config.CloudAPIURL
	}

	// Make request to Ollama
	resp, err := s.makeRequest("POST", "/api/embeddings", ollamaRequest, baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make embedding request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body for debugging
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check if the response indicates an error (like model not found)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse Ollama response
	var ollamaResp map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Extract embedding - handle different possible response formats
	var embeddingData []interface{}
	var ok bool

	// Try different possible field names for embedding data
	if embeddingData, ok = ollamaResp["embedding"].([]interface{}); !ok {
		if embeddingData, ok = ollamaResp["embeddings"].([]interface{}); !ok {
			if data, exists := ollamaResp["data"]; exists {
				if dataArray, isArray := data.([]interface{}); isArray && len(dataArray) > 0 {
					if firstItem, isMap := dataArray[0].(map[string]interface{}); isMap {
						if embedding, hasEmbedding := firstItem["embedding"].([]interface{}); hasEmbedding {
							embeddingData = embedding
							ok = true
						}
					}
				}
			}
		}
	}

	if !ok || embeddingData == nil {
		// If we can't find embedding data, return a more informative error
		return nil, fmt.Errorf("invalid embedding response format - no embedding data found in response: %v", ollamaResp)
	}

	// Convert to our format
	response := &models.EmbeddingResponse{
		Object: "list",
		Data: []models.Embedding{
			{
				Object:    "embedding",
				Embedding: convertToFloat64Slice(embeddingData),
				Index:     0,
			},
		},
		Model: model,
		Usage: s.extractUsage(ollamaResp),
	}

	return response, nil
}

// ListModels returns available models (local and cloud)
func (s *LlamaService) ListModels() ([]models.Model, error) {
	var allModels []models.Model

	// Get local models
	resp, err := s.makeRequest("GET", "/api/tags", nil, s.config.BaseURL)
	if err == nil {
		defer resp.Body.Close()
		var localResp map[string]interface{}
		if json.NewDecoder(resp.Body).Decode(&localResp) == nil {
			if modelsData, ok := localResp["models"].([]interface{}); ok {
				for _, modelData := range modelsData {
					if modelMap, ok := modelData.(map[string]interface{}); ok {
						model := models.Model{
							ID:      modelMap["name"].(string),
							Object:  "model",
							Created: time.Now().Unix(),
							OwnedBy: "ollama",
							IsCloud: false,
						}
						if size, ok := modelMap["size"].(string); ok {
							model.Size = size
						}
						allModels = append(allModels, model)
					}
				}
			}
		}
	}

	// Add cloud models if enabled and signed in
	if s.config.CloudEnabled && s.isSignedIn {
		for _, cloudModel := range CloudModels {
			if cloudModel.Available {
				model := models.Model{
					ID:      cloudModel.ID,
					Object:  "model",
					Created: time.Now().Unix(),
					OwnedBy: "ollama-cloud",
					IsCloud: true,
					Size:    cloudModel.Size,
				}
				allModels = append(allModels, model)
			}
		}
	}

	return allModels, nil
}

// StreamChat handles streaming chat completion
func (s *LlamaService) StreamChat(request models.ChatRequest, responseChan chan<- string) {
	defer close(responseChan)

	model := s.getModel(request.Model)

	// Check if cloud model and authentication
	if s.IsCloudModel(model) && !s.isSignedIn {
		responseChan <- fmt.Sprintf("Error: must be signed in to use cloud model: %s", model)
		return
	}

	// Convert to Ollama format
	ollamaRequest := map[string]interface{}{
		"model":    model,
		"messages": request.Messages,
		"stream":   true,
	}

	if request.Temperature > 0 {
		ollamaRequest["temperature"] = request.Temperature
	}

	// Determine which API to use
	baseURL := s.config.BaseURL
	if s.IsCloudModel(model) && s.config.CloudEnabled {
		baseURL = s.config.CloudAPIURL
	}

	// Make request to Ollama
	resp, err := s.makeRequest("POST", "/api/chat", ollamaRequest, baseURL)
	if err != nil {
		responseChan <- fmt.Sprintf("Error: %v", err)
		return
	}
	defer resp.Body.Close()

	// Read streaming response
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var streamResp map[string]interface{}
		if err := json.Unmarshal([]byte(line), &streamResp); err != nil {
			continue
		}

		if message, ok := streamResp["message"].(map[string]interface{}); ok {
			if content, ok := message["content"].(string); ok {
				responseChan <- content
			}
		}
	}
}

// makeRequest makes HTTP request to Ollama API
func (s *LlamaService) makeRequest(method, endpoint string, body interface{}, baseURL string) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, baseURL+endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Add authentication for cloud requests
	if strings.Contains(baseURL, "api.ollama.com") && s.config.CloudAPIKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.config.CloudAPIKey)
	}

	return s.httpClient.Do(req)
}

// Helper functions
func (s *LlamaService) getModel(requestedModel string) string {
	if requestedModel == "" {
		return s.config.DefaultModel
	}
	return requestedModel
}

func (s *LlamaService) extractContent(response map[string]interface{}) string {
	if message, ok := response["message"].(map[string]interface{}); ok {
		if content, ok := message["content"].(string); ok {
			return content
		}
	}
	return ""
}

func (s *LlamaService) extractResponse(response map[string]interface{}) string {
	if content, ok := response["response"].(string); ok {
		return content
	}
	return ""
}

func (s *LlamaService) extractUsage(response map[string]interface{}) models.Usage {
	usage := models.Usage{}

	if promptTokens, ok := response["prompt_eval_count"].(float64); ok {
		usage.PromptTokens = int(promptTokens)
	}
	if completionTokens, ok := response["eval_count"].(float64); ok {
		usage.CompletionTokens = int(completionTokens)
	}
	usage.TotalTokens = usage.PromptTokens + usage.CompletionTokens

	return usage
}

func generateID() string {
	return fmt.Sprintf("chatcmpl-%d", time.Now().UnixNano())
}

func convertToFloat64Slice(interfaceSlice []interface{}) []float64 {
	float64Slice := make([]float64, len(interfaceSlice))
	for i, v := range interfaceSlice {
		if f, ok := v.(float64); ok {
			float64Slice[i] = f
		}
	}
	return float64Slice
}

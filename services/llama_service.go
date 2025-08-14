package services

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"llama-api/models"
)

type LlamaService struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewLlamaService() *LlamaService {
	baseURL := os.Getenv("LLAMA_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:11434" // Default Ollama endpoint
	}

	apiKey := os.Getenv("LLAMA_API_KEY")

	return &LlamaService{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Chat handles chat completion using Llama
func (s *LlamaService) Chat(request models.ChatRequest) (*models.ChatResponse, error) {
	// Convert to Ollama format
	ollamaRequest := map[string]interface{}{
		"model":    s.getModel(request.Model),
		"messages": request.Messages,
		"stream":   false,
	}

	if request.Temperature > 0 {
		ollamaRequest["temperature"] = request.Temperature
	}

	// Make request to Ollama
	resp, err := s.makeRequest("POST", "/api/chat", ollamaRequest)
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

// Completion handles text completion using Llama
func (s *LlamaService) Completion(request models.CompletionRequest) (*models.CompletionResponse, error) {
	// Convert to Ollama format
	ollamaRequest := map[string]interface{}{
		"model":  s.getModel(request.Model),
		"prompt": request.Prompt,
		"stream": false,
	}

	if request.Temperature > 0 {
		ollamaRequest["temperature"] = request.Temperature
	}

	// Make request to Ollama
	resp, err := s.makeRequest("POST", "/api/generate", ollamaRequest)
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

// Embedding handles text embedding using Llama
func (s *LlamaService) Embedding(request models.EmbeddingRequest) (*models.EmbeddingResponse, error) {
	// Convert to Ollama format
	ollamaRequest := map[string]interface{}{
		"model": s.getModel(request.Model),
		"prompt": request.Input,
	}

	// Make request to Ollama
	resp, err := s.makeRequest("POST", "/api/embeddings", ollamaRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to make embedding request: %w", err)
	}
	defer resp.Body.Close()

	// Parse Ollama response
	var ollamaResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to our format
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

// ListModels returns available models
func (s *LlamaService) ListModels() ([]models.Model, error) {
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
			modelList = append(modelList, models.Model{
				ID:      model["name"].(string),
				Object:  "model",
				Created: time.Now().Unix(),
				OwnedBy: "ollama",
			})
		}
	}

	return modelList, nil
}

// StreamChat handles streaming chat responses
func (s *LlamaService) StreamChat(request models.ChatRequest, responseChan chan<- string) {
	defer close(responseChan)

	// Convert to Ollama format
	ollamaRequest := map[string]interface{}{
		"model":    s.getModel(request.Model),
		"messages": request.Messages,
		"stream":   true,
	}

	if request.Temperature > 0 {
		ollamaRequest["temperature"] = request.Temperature
	}

	// Make streaming request to Ollama
	resp, err := s.makeRequest("POST", "/api/chat", ollamaRequest)
	if err != nil {
		responseChan <- fmt.Sprintf("Error: %v", err)
		return
	}
	defer resp.Body.Close()

	// Stream the response
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			responseChan <- fmt.Sprintf("Error reading stream: %v", err)
			break
		}

		// Parse the line and send to channel
		if line = strings.TrimSpace(line); line != "" {
			responseChan <- line
		}
	}
}

// Helper methods
func (s *LlamaService) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, s.baseURL+endpoint, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if s.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiKey)
	}

	return s.httpClient.Do(req)
}

func (s *LlamaService) getModel(requestedModel string) string {
	if requestedModel != "" {
		return requestedModel
	}
	return os.Getenv("LLAMA_DEFAULT_MODEL")
}

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

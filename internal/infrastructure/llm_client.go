package infrastructure

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"llama-api/internal/domain"
	"llama-api/pkg/errors"
	"llama-api/pkg/logger"
)

// OllamaClient handles communication with Ollama API
type OllamaClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	logger     *logger.Logger
	timeout    time.Duration
}

// NewOllamaClient creates a new Ollama client
func NewOllamaClient(logger *logger.Logger) *OllamaClient {
	baseURL := os.Getenv("LLAMA_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	apiKey := os.Getenv("LLAMA_API_KEY")
	timeout := 60 * time.Second

	return &OllamaClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		timeout: timeout,
		logger:  logger,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// Chat performs a chat completion request
func (c *OllamaClient) Chat(ctx context.Context, request domain.LLMRequest) (*domain.LLMResponse, error) {
	ollamaRequest := c.buildChatRequest(request)

	resp, err := c.makeRequest(ctx, "POST", "/api/chat", ollamaRequest)
	if err != nil {
		return nil, errors.Wrap(errors.CodeLLMService, "failed to make chat request", err)
	}
	defer resp.Body.Close()

	return c.parseChatResponse(resp, request.Model)
}

// Completion performs a text completion request
func (c *OllamaClient) Completion(ctx context.Context, request domain.CompletionRequest) (*domain.CompletionResponse, error) {
	ollamaRequest := c.buildCompletionRequest(request)

	resp, err := c.makeRequest(ctx, "POST", "/api/generate", ollamaRequest)
	if err != nil {
		return nil, errors.Wrap(errors.CodeLLMService, "failed to make completion request", err)
	}
	defer resp.Body.Close()

	return c.parseCompletionResponse(resp, request.Model)
}

// Embedding performs an embedding request
func (c *OllamaClient) Embedding(ctx context.Context, request domain.EmbeddingRequest) (*domain.EmbeddingResponse, error) {
	ollamaRequest := c.buildEmbeddingRequest(request)

	resp, err := c.makeRequest(ctx, "POST", "/api/embeddings", ollamaRequest)
	if err != nil {
		return nil, errors.Wrap(errors.CodeLLMService, "failed to make embedding request", err)
	}
	defer resp.Body.Close()

	return c.parseEmbeddingResponse(resp, request.Model)
}

// ListModels returns available models from Ollama
func (c *OllamaClient) ListModels(ctx context.Context) ([]domain.LLMModel, error) {
	resp, err := c.makeRequest(ctx, "GET", "/api/tags", nil)
	if err != nil {
		return nil, errors.Wrap(errors.CodeLLMService, "failed to list models", err)
	}
	defer resp.Body.Close()

	var ollamaResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, errors.Wrap(errors.CodeLLMService, "failed to decode models response", err)
	}

	modelsSlice, err := errors.SafeSliceAssert(ollamaResp["models"], "models")
	if err != nil {
		return nil, errors.Wrap(errors.CodeLLMService, "models field missing or invalid", err)
	}

	var modelList []domain.LLMModel
	for _, modelData := range modelsSlice {
		modelMap, err := errors.SafeMapAssert(modelData, "model")
		if err != nil {
			c.logger.Warn("skipping invalid model entry", err)
			continue
		}

		name, err := errors.SafeStringAssert(modelMap["name"], "name")
		if err != nil {
			c.logger.Warn("model name missing", err)
			continue
		}

		modelList = append(modelList, domain.LLMModel{
			ID:      name,
			Object:  "model",
			Created: time.Now().Unix(),
			OwnedBy: "ollama",
		})
	}

	return modelList, nil
}

// StreamChat streams chat responses
func (c *OllamaClient) StreamChat(ctx context.Context, request domain.LLMRequest, responseChan chan<- string) error {
	defer close(responseChan)

	ollamaRequest := c.buildChatRequest(request)
	ollamaRequest["stream"] = true

	resp, err := c.makeRequest(ctx, "POST", "/api/chat", ollamaRequest)
	if err != nil {
		responseChan <- fmt.Sprintf("Error: %v", err)
		return errors.Wrap(errors.CodeLLMService, "failed to stream chat", err)
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return errors.Wrap(errors.CodeLLMService, "failed to read stream", err)
		}

		if line = strings.TrimSpace(line); line != "" {
			responseChan <- line
		}
	}

	return nil
}

// Helper methods

func (c *OllamaClient) buildChatRequest(request domain.LLMRequest) map[string]interface{} {
	ollamaRequest := map[string]interface{}{
		"model":    c.getModel(request.Model),
		"messages": c.convertMessages(request.Messages),
		"stream":   request.Stream,
	}

	if request.Temperature > 0 {
		ollamaRequest["temperature"] = request.Temperature
	}

	if request.MaxTokens > 0 {
		ollamaRequest["num_predict"] = request.MaxTokens
	}

	return ollamaRequest
}

func (c *OllamaClient) buildCompletionRequest(request domain.CompletionRequest) map[string]interface{} {
	ollamaRequest := map[string]interface{}{
		"model":  c.getModel(request.Model),
		"prompt": request.Prompt,
		"stream": false,
	}

	if request.Temperature > 0 {
		ollamaRequest["temperature"] = request.Temperature
	}

	if request.MaxTokens > 0 {
		ollamaRequest["num_predict"] = request.MaxTokens
	}

	return ollamaRequest
}

func (c *OllamaClient) buildEmbeddingRequest(request domain.EmbeddingRequest) map[string]interface{} {
	return map[string]interface{}{
		"model":  c.getModel(request.Model),
		"prompt": request.Input,
	}
}

func (c *OllamaClient) convertMessages(messages []domain.Message) []map[string]string {
	converted := make([]map[string]string, len(messages))
	for i, msg := range messages {
		converted[i] = map[string]string{
			"role":    msg.Role,
			"content": msg.Content,
		}
	}
	return converted
}

func (c *OllamaClient) getModel(requestedModel string) string {
	if requestedModel != "" {
		return requestedModel
	}
	model := os.Getenv("LLAMA_DEFAULT_MODEL")
	if model == "" {
		model = "llama2"
	}
	return model
}

func (c *OllamaClient) makeRequest(ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+endpoint, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	return c.httpClient.Do(req)
}

func (c *OllamaClient) parseChatResponse(resp *http.Response, model string) (*domain.LLMResponse, error) {
	var ollamaResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, errors.Wrap(errors.CodeLLMService, "failed to decode chat response", err)
	}

	// Safely extract message with proper error handling
	messageData, err := errors.SafeMapAssert(ollamaResp["message"], "message")
	if err != nil {
		return nil, err
	}

	content, err := errors.SafeStringAssert(messageData["content"], "content")
	if err != nil {
		return nil, err
	}

	promptTokens, _ := errors.SafeFloat64Assert(ollamaResp["prompt_eval_count"], "prompt_eval_count")
	completionTokens, _ := errors.SafeFloat64Assert(ollamaResp["eval_count"], "eval_count")

	response := &domain.LLMResponse{
		ID:      c.generateID(),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []domain.Choice{
			{
				Index: 0,
				Message: domain.Message{
					Role:    "assistant",
					Content: content,
				},
			},
		},
		Usage: domain.Usage{
			PromptTokens:     int(promptTokens),
			CompletionTokens: int(completionTokens),
			TotalTokens:      int(promptTokens + completionTokens),
		},
	}

	return response, nil
}

func (c *OllamaClient) parseCompletionResponse(resp *http.Response, model string) (*domain.CompletionResponse, error) {
	var ollamaResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, errors.Wrap(errors.CodeLLMService, "failed to decode completion response", err)
	}

	responseText, err := errors.SafeStringAssert(ollamaResp["response"], "response")
	if err != nil {
		return nil, err
	}

	promptTokens, _ := errors.SafeFloat64Assert(ollamaResp["prompt_eval_count"], "prompt_eval_count")
	completionTokens, _ := errors.SafeFloat64Assert(ollamaResp["eval_count"], "eval_count")

	response := &domain.CompletionResponse{
		ID:      c.generateID(),
		Object:  "text_completion",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []domain.Choice{
			{
				Index: 0,
				Message: domain.Message{
					Role:    "assistant",
					Content: responseText,
				},
			},
		},
		Usage: domain.Usage{
			PromptTokens:     int(promptTokens),
			CompletionTokens: int(completionTokens),
			TotalTokens:      int(promptTokens + completionTokens),
		},
	}

	return response, nil
}

func (c *OllamaClient) parseEmbeddingResponse(resp *http.Response, model string) (*domain.EmbeddingResponse, error) {
	var ollamaResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, errors.Wrap(errors.CodeLLMService, "failed to decode embedding response", err)
	}

	embeddingSlice, err := errors.SafeSliceAssert(ollamaResp["embedding"], "embedding")
	if err != nil {
		return nil, err
	}

	floatSlice := make([]float64, len(embeddingSlice))
	for i, v := range embeddingSlice {
		f, err := errors.SafeFloat64Assert(v, fmt.Sprintf("embedding[%d]", i))
		if err != nil {
			c.logger.Warn("skipping invalid embedding value", err)
			floatSlice[i] = 0
			continue
		}
		floatSlice[i] = f
	}

	embedding := domain.Embedding{
		Object:    "embedding",
		Index:     0,
		Embedding: floatSlice,
	}

	response := &domain.EmbeddingResponse{
		Object: "list",
		Data:   []domain.Embedding{embedding},
		Model:  model,
		Usage: domain.Usage{
			PromptTokens: 1,
			TotalTokens:  1,
		},
	}

	return response, nil
}

func (c *OllamaClient) generateID() string {
	return "chatcmpl-" + strconv.FormatInt(time.Now().UnixNano(), 10)
}

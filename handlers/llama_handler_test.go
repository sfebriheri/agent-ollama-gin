package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"agent-ollama-gin/models"
	"agent-ollama-gin/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLlamaService is a mock implementation of LlamaServiceInterface for testing
type MockLlamaService struct {
	mock.Mock
}

// Ensure MockLlamaService implements the interface
var _ services.LlamaServiceInterface = (*MockLlamaService)(nil)

func (m *MockLlamaService) Chat(request models.ChatRequest) (*models.ChatResponse, error) {
	args := m.Called(request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ChatResponse), args.Error(1)
}

func (m *MockLlamaService) Completion(request models.CompletionRequest) (*models.CompletionResponse, error) {
	args := m.Called(request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CompletionResponse), args.Error(1)
}

func (m *MockLlamaService) Embedding(request models.EmbeddingRequest) (*models.EmbeddingResponse, error) {
	args := m.Called(request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EmbeddingResponse), args.Error(1)
}

func (m *MockLlamaService) ListModels() ([]models.Model, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Model), args.Error(1)
}

func (m *MockLlamaService) SignIn(username, password string) (*models.AuthResponse, error) {
	args := m.Called(username, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AuthResponse), args.Error(1)
}

func (m *MockLlamaService) SignOut() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockLlamaService) PullModel(modelName string) error {
	args := m.Called(modelName)
	return args.Error(0)
}

func (m *MockLlamaService) StreamChat(request models.ChatRequest, responseChan chan<- string) {
	m.Called(request, responseChan)
}

func setupRouter(handler *LlamaHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	api := router.Group("/api/v1/llama")
	{
		api.POST("/chat", handler.Chat)
		api.POST("/completion", handler.Completion)
		api.POST("/embedding", handler.Embedding)
		api.GET("/models", handler.ListModels)
		api.POST("/chat/stream", handler.StreamChat)
		api.POST("/cloud/signin", handler.SignIn)
		api.POST("/cloud/signout", handler.SignOut)
		api.POST("/models/:model/pull", handler.PullModel)
		api.GET("/cloud/models", handler.ListCloudModels)
	}

	return router
}

func TestChat_Success(t *testing.T) {
	mockService := new(MockLlamaService)
	handler := NewLlamaHandler(mockService)
	router := setupRouter(handler)

	expectedResponse := &models.ChatResponse{
		ID:      "test-id",
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   "llama2",
		Choices: []models.Choice{
			{
				Index: 0,
				Message: models.Message{
					Role:    "assistant",
					Content: "Hello! How can I help you?",
				},
			},
		},
	}

	chatRequest := models.ChatRequest{
		Messages: []models.Message{
			{Role: "user", Content: "Hello"},
		},
		Model: "llama2",
	}

	mockService.On("Chat", chatRequest).Return(expectedResponse, nil)

	body, _ := json.Marshal(chatRequest)
	req, _ := http.NewRequest("POST", "/api/v1/llama/chat", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestChat_InvalidRequest(t *testing.T) {
	mockService := new(MockLlamaService)
	handler := NewLlamaHandler(mockService)
	router := setupRouter(handler)

	req, _ := http.NewRequest("POST", "/api/v1/llama/chat", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestChat_EmptyMessages(t *testing.T) {
	mockService := new(MockLlamaService)
	handler := NewLlamaHandler(mockService)
	router := setupRouter(handler)

	chatRequest := models.ChatRequest{
		Messages: []models.Message{},
		Model:    "llama2",
	}

	body, _ := json.Marshal(chatRequest)
	req, _ := http.NewRequest("POST", "/api/v1/llama/chat", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestChat_ServiceError(t *testing.T) {
	mockService := new(MockLlamaService)
	handler := NewLlamaHandler(mockService)
	router := setupRouter(handler)

	chatRequest := models.ChatRequest{
		Messages: []models.Message{
			{Role: "user", Content: "Hello"},
		},
	}

	mockService.On("Chat", chatRequest).Return(nil, errors.New("service error"))

	body, _ := json.Marshal(chatRequest)
	req, _ := http.NewRequest("POST", "/api/v1/llama/chat", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestCompletion_Success(t *testing.T) {
	mockService := new(MockLlamaService)
	handler := NewLlamaHandler(mockService)
	router := setupRouter(handler)

	expectedResponse := &models.CompletionResponse{
		ID:      "test-id",
		Object:  "text_completion",
		Created: time.Now().Unix(),
		Model:   "llama2",
		Choices: []models.Choice{
			{
				Index: 0,
				Message: models.Message{
					Role:    "assistant",
					Content: "The future of AI is bright",
				},
			},
		},
	}

	completionRequest := models.CompletionRequest{
		Prompt: "The future of AI is",
		Model:  "llama2",
	}

	mockService.On("Completion", completionRequest).Return(expectedResponse, nil)

	body, _ := json.Marshal(completionRequest)
	req, _ := http.NewRequest("POST", "/api/v1/llama/completion", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestCompletion_EmptyPrompt(t *testing.T) {
	mockService := new(MockLlamaService)
	handler := NewLlamaHandler(mockService)
	router := setupRouter(handler)

	completionRequest := models.CompletionRequest{
		Prompt: "",
		Model:  "llama2",
	}

	body, _ := json.Marshal(completionRequest)
	req, _ := http.NewRequest("POST", "/api/v1/llama/completion", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEmbedding_Success(t *testing.T) {
	mockService := new(MockLlamaService)
	handler := NewLlamaHandler(mockService)
	router := setupRouter(handler)

	expectedResponse := &models.EmbeddingResponse{
		Object: "list",
		Data: []models.Embedding{
			{
				Object:    "embedding",
				Embedding: []float64{0.1, 0.2, 0.3},
				Index:     0,
			},
		},
		Model: "llama2",
	}

	embeddingRequest := models.EmbeddingRequest{
		Input: "Test input",
		Model: "llama2",
	}

	mockService.On("Embedding", embeddingRequest).Return(expectedResponse, nil)

	body, _ := json.Marshal(embeddingRequest)
	req, _ := http.NewRequest("POST", "/api/v1/llama/embedding", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestListModels_Success(t *testing.T) {
	mockService := new(MockLlamaService)
	handler := NewLlamaHandler(mockService)
	router := setupRouter(handler)

	expectedModels := []models.Model{
		{
			ID:      "llama2",
			Object:  "model",
			Created: time.Now().Unix(),
			OwnedBy: "ollama",
		},
	}

	mockService.On("ListModels").Return(expectedModels, nil)

	req, _ := http.NewRequest("GET", "/api/v1/llama/models", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestSignIn_Success(t *testing.T) {
	mockService := new(MockLlamaService)
	handler := NewLlamaHandler(mockService)
	router := setupRouter(handler)

	expectedResponse := &models.AuthResponse{
		Success: true,
		Token:   "test-token",
		Message: "Successfully signed in",
	}

	authRequest := models.AuthRequest{
		Username: "test@example.com",
		Password: "password123",
	}

	mockService.On("SignIn", "test@example.com", "password123").Return(expectedResponse, nil)

	body, _ := json.Marshal(authRequest)
	req, _ := http.NewRequest("POST", "/api/v1/llama/cloud/signin", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestSignIn_MissingCredentials(t *testing.T) {
	mockService := new(MockLlamaService)
	handler := NewLlamaHandler(mockService)
	router := setupRouter(handler)

	authRequest := models.AuthRequest{
		Username: "",
		Password: "",
	}

	body, _ := json.Marshal(authRequest)
	req, _ := http.NewRequest("POST", "/api/v1/llama/cloud/signin", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSignOut_Success(t *testing.T) {
	mockService := new(MockLlamaService)
	handler := NewLlamaHandler(mockService)
	router := setupRouter(handler)

	mockService.On("SignOut").Return(nil)

	req, _ := http.NewRequest("POST", "/api/v1/llama/cloud/signout", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestPullModel_Success(t *testing.T) {
	mockService := new(MockLlamaService)
	handler := NewLlamaHandler(mockService)
	router := setupRouter(handler)

	mockService.On("PullModel", "llama2").Return(nil)

	req, _ := http.NewRequest("POST", "/api/v1/llama/models/llama2/pull", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestListCloudModels_Success(t *testing.T) {
	handler := NewLlamaHandler(nil) // No mock needed for this simple handler
	router := setupRouter(handler)

	req, _ := http.NewRequest("GET", "/api/v1/llama/cloud/models", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "models")

	// Check that we have models in the response
	models, ok := response["models"]
	assert.True(t, ok)
	assert.NotNil(t, models)
}

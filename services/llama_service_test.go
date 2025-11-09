package services

import (
	"testing"

	"agent-ollama-gin/models"

	"github.com/stretchr/testify/assert"
)

func TestNewLlamaService(t *testing.T) {
	service := NewLlamaService()

	assert.NotNil(t, service)
	assert.NotNil(t, service.config)
	assert.NotNil(t, service.httpClient)
}

func TestIsCloudModel(t *testing.T) {
	service := NewLlamaService()

	tests := []struct {
		name       string
		modelName  string
		isCloud    bool
	}{
		{
			name:      "Cloud model",
			modelName: "gpt-oss:120b-cloud",
			isCloud:   true,
		},
		{
			name:      "Local model",
			modelName: "llama2",
			isCloud:   false,
		},
		{
			name:      "Another cloud model",
			modelName: "qwen3-coder:480b-cloud",
			isCloud:   true,
		},
		{
			name:      "Local model with version",
			modelName: "llama3.2:1b",
			isCloud:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.IsCloudModel(tt.modelName)
			assert.Equal(t, tt.isCloud, result)
		})
	}
}

func TestSignIn(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
	}{
		{
			name:     "Valid credentials",
			username: "test@example.com",
			password: "password123",
		},
		{
			name:     "Empty credentials",
			username: "",
			password: "",
		},
		{
			name:     "Only username",
			username: "test@example.com",
			password: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewLlamaService()
			response, err := service.SignIn(tt.username, tt.password)

			assert.NoError(t, err)
			assert.NotNil(t, response)
			// Response will depend on cloud configuration
		})
	}
}

func TestSignOut(t *testing.T) {
	service := NewLlamaService()

	// Sign out should work regardless of current state
	err := service.SignOut()
	assert.NoError(t, err)
}

func TestCloudModelsAvailability(t *testing.T) {
	assert.NotEmpty(t, CloudModels)

	// Verify cloud model structure
	for _, model := range CloudModels {
		assert.NotEmpty(t, model.Name)
		assert.NotEmpty(t, model.ID)
		assert.NotEmpty(t, model.Size)
		assert.NotEmpty(t, model.Description)
		assert.True(t, model.Available)

		// Verify cloud models have -cloud suffix
		assert.Contains(t, model.ID, "-cloud")
	}
}

func TestGenerateID(t *testing.T) {
	id1 := generateID()
	id2 := generateID()

	// IDs should not be empty
	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, id2)

	// IDs should be unique
	assert.NotEqual(t, id1, id2)

	// IDs should have the expected prefix
	assert.Contains(t, id1, "chatcmpl-")
	assert.Contains(t, id2, "chatcmpl-")
}

func TestConvertToFloat64Slice(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		expected []float64
	}{
		{
			name:     "Valid floats",
			input:    []interface{}{1.0, 2.5, 3.7},
			expected: []float64{1.0, 2.5, 3.7},
		},
		{
			name:     "Empty slice",
			input:    []interface{}{},
			expected: []float64{},
		},
		{
			name:     "Mixed types",
			input:    []interface{}{1.0, "invalid", 3.0},
			expected: []float64{1.0, 0.0, 3.0}, // Invalid types become 0.0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertToFloat64Slice(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractUsage(t *testing.T) {
	service := NewLlamaService()

	tests := []struct {
		name     string
		response map[string]interface{}
		expected models.Usage
	}{
		{
			name: "Valid usage data",
			response: map[string]interface{}{
				"prompt_eval_count": 10.0,
				"eval_count":        20.0,
			},
			expected: models.Usage{
				PromptTokens:     10,
				CompletionTokens: 20,
				TotalTokens:      30,
			},
		},
		{
			name:     "Empty response",
			response: map[string]interface{}{},
			expected: models.Usage{
				PromptTokens:     0,
				CompletionTokens: 0,
				TotalTokens:      0,
			},
		},
		{
			name: "Partial data",
			response: map[string]interface{}{
				"prompt_eval_count": 15.0,
			},
			expected: models.Usage{
				PromptTokens:     15,
				CompletionTokens: 0,
				TotalTokens:      15,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.extractUsage(tt.response)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractContent(t *testing.T) {
	service := NewLlamaService()

	tests := []struct {
		name     string
		response map[string]interface{}
		expected string
	}{
		{
			name: "Valid content",
			response: map[string]interface{}{
				"message": map[string]interface{}{
					"content": "Hello, world!",
				},
			},
			expected: "Hello, world!",
		},
		{
			name:     "Missing message",
			response: map[string]interface{}{},
			expected: "",
		},
		{
			name: "Missing content",
			response: map[string]interface{}{
				"message": map[string]interface{}{},
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.extractContent(tt.response)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractResponse(t *testing.T) {
	service := NewLlamaService()

	tests := []struct {
		name     string
		response map[string]interface{}
		expected string
	}{
		{
			name: "Valid response",
			response: map[string]interface{}{
				"response": "This is a response",
			},
			expected: "This is a response",
		},
		{
			name:     "Missing response",
			response: map[string]interface{}{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.extractResponse(tt.response)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetModel(t *testing.T) {
	service := NewLlamaService()

	tests := []struct {
		name           string
		requestedModel string
		expectedModel  string
	}{
		{
			name:           "Explicit model",
			requestedModel: "llama2",
			expectedModel:  "llama2",
		},
		{
			name:           "Empty model uses default",
			requestedModel: "",
			expectedModel:  service.config.DefaultModel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.getModel(tt.requestedModel)
			assert.Equal(t, tt.expectedModel, result)
		})
	}
}

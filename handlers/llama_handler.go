package handlers

import (
	"net/http"

	"agent-ollama-gin/models"

	"github.com/gin-gonic/gin"

	"agent-ollama-gin/services"
)

type LlamaHandler struct {
	llamaService *services.LlamaService
}

func NewLlamaHandler(llamaService *services.LlamaService) *LlamaHandler {
	return &LlamaHandler{
		llamaService: llamaService,
	}
}

// Chat handles chat completion requests
func (h *LlamaHandler) Chat(c *gin.Context) {
	var request models.ChatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate request
	if len(request.Messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "At least one message is required",
		})
		return
	}

	response, err := h.llamaService.Chat(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to process chat request",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Completion handles text completion requests
func (h *LlamaHandler) Completion(c *gin.Context) {
	var request models.CompletionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate request
	if request.Prompt == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Prompt is required",
		})
		return
	}

	response, err := h.llamaService.Completion(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to process completion request",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Embedding handles text embedding requests
func (h *LlamaHandler) Embedding(c *gin.Context) {
	var request models.EmbeddingRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate request
	if request.Input == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input text is required",
		})
		return
	}

	response, err := h.llamaService.Embedding(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to process embedding request",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ListModels returns available Llama models
func (h *LlamaHandler) ListModels(c *gin.Context) {
	models, err := h.llamaService.ListModels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve models",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"models": models,
	})
}

// StreamChat handles streaming chat responses
func (h *LlamaHandler) StreamChat(c *gin.Context) {
	var request models.ChatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Set headers for streaming
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// Create a channel for streaming responses
	responseChan := make(chan string)

	go func() {
		h.llamaService.StreamChat(request, responseChan)
	}()

	// Stream responses
	for response := range responseChan {
		c.SSEvent("message", response)
		c.Writer.Flush()
	}
}

// SignIn handles Ollama cloud authentication
func (h *LlamaHandler) SignIn(c *gin.Context) {
	var request models.AuthRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate request
	if request.Username == "" || request.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username and password are required",
		})
		return
	}

	response, err := h.llamaService.SignIn(request.Username, request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to sign in",
			"details": err.Error(),
		})
		return
	}

	if !response.Success {
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

// SignOut handles Ollama cloud sign out
func (h *LlamaHandler) SignOut(c *gin.Context) {
	err := h.llamaService.SignOut()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to sign out",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully signed out",
	})
}

// PullModel handles model pulling (local or cloud)
func (h *LlamaHandler) PullModel(c *gin.Context) {
	modelName := c.Param("model")
	if modelName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Model name is required",
		})
		return
	}

	err := h.llamaService.PullModel(modelName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to pull model",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Model pulled successfully",
		"model":   modelName,
	})
}

// ListCloudModels returns available cloud models
func (h *LlamaHandler) ListCloudModels(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"models": services.CloudModels,
	})
}

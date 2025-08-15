package handlers

import (
	"net/http"

	"llama-api/minimal/services"
	"llama-api/models"

	"github.com/gin-gonic/gin"
)

type LlamaHandler struct {
	llamaService *services.OptimizedLlamaService
}

func NewLlamaHandler(llamaService *services.OptimizedLlamaService) *LlamaHandler {
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

package http

import (
	"net/http"

	"agent-ollama-gin/internal/domain"
	"agent-ollama-gin/pkg/errors"
	"agent-ollama-gin/pkg/logger"

	"github.com/gin-gonic/gin"
)

// LLMHandler handles HTTP requests for LLM operations
type LLMHandler struct {
	usecase domain.LLMUsecase
	logger  *logger.Logger
}

// NewLLMHandler creates a new LLM handler
func NewLLMHandler(usecase domain.LLMUsecase, logger *logger.Logger) *LLMHandler {
	return &LLMHandler{
		usecase: usecase,
		logger:  logger,
	}
}

// Chat handles chat completion requests
func (h *LLMHandler) Chat(c *gin.Context) {
	var request domain.LLMRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Warn("invalid chat request", err)
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

	// Get context with timeout from gin
	ctx := c.Request.Context()

	response, err := h.usecase.Chat(ctx, request)
	if err != nil {
		h.logger.Error("chat request failed", err)
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// Completion handles text completion requests
func (h *LLMHandler) Completion(c *gin.Context) {
	var request domain.CompletionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Warn("invalid completion request", err)
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

	ctx := c.Request.Context()

	response, err := h.usecase.Completion(ctx, request)
	if err != nil {
		h.logger.Error("completion request failed", err)
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// Embedding handles text embedding requests
func (h *LLMHandler) Embedding(c *gin.Context) {
	var request domain.EmbeddingRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Warn("invalid embedding request", err)
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

	ctx := c.Request.Context()

	response, err := h.usecase.Embedding(ctx, request)
	if err != nil {
		h.logger.Error("embedding request failed", err)
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// ListModels returns available Llama models
func (h *LLMHandler) ListModels(c *gin.Context) {
	ctx := c.Request.Context()

	models, err := h.usecase.ListModels(ctx)
	if err != nil {
		h.logger.Error("list models request failed", err)
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"models": models,
	})
}

// StreamChat handles streaming chat responses
func (h *LLMHandler) StreamChat(c *gin.Context) {
	var request domain.LLMRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Warn("invalid stream chat request", err)
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

	// Get context from gin
	ctx := c.Request.Context()

	go func() {
		defer close(responseChan)
		err := h.usecase.StreamChat(ctx, request, responseChan)
		if err != nil {
			h.logger.Error("stream chat failed", err)
		}
	}()

	// Stream responses
	for response := range responseChan {
		c.SSEvent("message", response)
		c.Writer.Flush()
	}
}

// Health returns health status
func (h *LLMHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "llama",
		"message": "LLM service is running",
		"version": "2.0.0",
	})
}

// handleError handles application errors
func (h *LLMHandler) handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		statusCode := h.errorCodeToHTTPStatus(appErr.Code)
		c.JSON(statusCode, gin.H{
			"error":   appErr.Message,
			"code":    appErr.Code,
			"details": appErr.Details,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "Internal server error",
		"code":  errors.CodeInternal,
	})
}

// errorCodeToHTTPStatus maps error codes to HTTP status codes
func (h *LLMHandler) errorCodeToHTTPStatus(code string) int {
	switch code {
	case errors.CodeInvalidInput:
		return http.StatusBadRequest
	case errors.CodeUnauthorized:
		return http.StatusUnauthorized
	case errors.CodeForbidden:
		return http.StatusForbidden
	case errors.CodeNotFound:
		return http.StatusNotFound
	case errors.CodeTimeout:
		return http.StatusGatewayTimeout
	case errors.CodeServiceUnavail:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

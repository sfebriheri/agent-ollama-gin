package http

import (
	"net/http"

	"agent-ollama-gin/internal/domain"
	"agent-ollama-gin/pkg/errors"
	"agent-ollama-gin/pkg/logger"

	"github.com/gin-gonic/gin"
)

// EncyclopediaHandler handles HTTP requests for encyclopedia operations
type EncyclopediaHandler struct {
	usecase    domain.EncyclopediaUsecase
	llmUsecase domain.LLMUsecase
	logger     *logger.Logger
}

// NewEncyclopediaHandler creates a new encyclopedia handler
func NewEncyclopediaHandler(
	usecase domain.EncyclopediaUsecase,
	llmUsecase domain.LLMUsecase,
	logger *logger.Logger,
) *EncyclopediaHandler {
	return &EncyclopediaHandler{
		usecase:    usecase,
		llmUsecase: llmUsecase,
		logger:     logger,
	}
}

// SearchEncyclopedia handles encyclopedia search requests
func (h *EncyclopediaHandler) SearchEncyclopedia(c *gin.Context) {
	var request domain.EncyclopediaSearchRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Warn("invalid search request", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate request
	if request.Query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Search query is required",
		})
		return
	}

	ctx := c.Request.Context()

	response, err := h.usecase.SearchEncyclopedia(ctx, request)
	if err != nil {
		h.logger.Error("search encyclopedia request failed", err)
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetArticle handles requests to retrieve specific encyclopedia articles
func (h *EncyclopediaHandler) GetArticle(c *gin.Context) {
	var request domain.EncyclopediaArticleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Warn("invalid article request", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate request - either title or URL must be provided
	if request.Title == "" && request.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Either title or URL is required",
		})
		return
	}

	ctx := c.Request.Context()

	response, err := h.usecase.GetArticle(ctx, request)
	if err != nil {
		h.logger.Error("get article request failed", err)
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GeneratePrompt handles requests to generate encyclopedia-style prompts
func (h *EncyclopediaHandler) GeneratePrompt(c *gin.Context) {
	var request domain.EncyclopediaPromptRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Warn("invalid prompt request", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate request
	if request.Topic == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Topic is required",
		})
		return
	}

	ctx := c.Request.Context()

	response, err := h.usecase.GeneratePrompt(ctx, request, h.llmUsecase)
	if err != nil {
		h.logger.Error("generate prompt request failed", err)
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetSources returns available encyclopedia sources
func (h *EncyclopediaHandler) GetSources(c *gin.Context) {
	sources := []gin.H{
		{
			"name":                "Wikipedia",
			"description":         "Free online encyclopedia with articles in multiple languages",
			"url":                 "https://www.wikipedia.org",
			"supported_languages": []string{"en", "es", "fr", "de", "it", "pt", "ru", "ja", "zh", "ar"},
		},
		{
			"name":                "Britannica",
			"description":         "Professional encyclopedia with expert-curated content",
			"url":                 "https://www.britannica.com",
			"supported_languages": []string{"en"},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"sources": sources,
		"total":   len(sources),
	})
}

// GetLanguages returns supported languages for encyclopedia content
func (h *EncyclopediaHandler) GetLanguages(c *gin.Context) {
	languages := []gin.H{
		{"code": "en", "name": "English", "native_name": "English"},
		{"code": "es", "name": "Spanish", "native_name": "Español"},
		{"code": "fr", "name": "French", "native_name": "Français"},
		{"code": "de", "name": "German", "native_name": "Deutsch"},
		{"code": "it", "name": "Italian", "native_name": "Italiano"},
		{"code": "pt", "name": "Portuguese", "native_name": "Português"},
		{"code": "ru", "name": "Russian", "native_name": "Русский"},
		{"code": "ja", "name": "Japanese", "native_name": "日本語"},
		{"code": "zh", "name": "Chinese", "native_name": "中文"},
		{"code": "ar", "name": "Arabic", "native_name": "العربية"},
	}

	c.JSON(http.StatusOK, gin.H{
		"languages": languages,
		"total":     len(languages),
		"default":   "en",
	})
}

// Health returns health status
func (h *EncyclopediaHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "encyclopedia",
		"message": "Encyclopedia service is running",
		"version": "2.0.0",
		"features": []string{
			"search",
			"article_retrieval",
			"prompt_generation",
			"multi_source_support",
			"multi_language_support",
			"parallel_search",
			"caching",
		},
	})
}

// handleError handles application errors
func (h *EncyclopediaHandler) handleError(c *gin.Context, err error) {
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
func (h *EncyclopediaHandler) errorCodeToHTTPStatus(code string) int {
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

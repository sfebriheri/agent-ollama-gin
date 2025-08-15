package handlers

import (
	"net/http"

	"llama-api/models"
	"llama-api/services"

	"github.com/gin-gonic/gin"
)

type EncyclopediaHandler struct {
	encyclopediaService *services.EncyclopediaService
	llamaService        *services.LlamaService
}

func NewEncyclopediaHandler(encyclopediaService *services.EncyclopediaService, llamaService *services.LlamaService) *EncyclopediaHandler {
	return &EncyclopediaHandler{
		encyclopediaService: encyclopediaService,
		llamaService:        llamaService,
	}
}

// SearchEncyclopedia handles encyclopedia search requests
func (h *EncyclopediaHandler) SearchEncyclopedia(c *gin.Context) {
	var request models.EncyclopediaSearchRequest
	if err := c.ShouldBindJSON(&request); err != nil {
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

	response, err := h.encyclopediaService.SearchEncyclopedia(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to search encyclopedia",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetArticle handles requests to retrieve specific encyclopedia articles
func (h *EncyclopediaHandler) GetArticle(c *gin.Context) {
	var request models.EncyclopediaArticleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
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

	response, err := h.encyclopediaService.GetArticle(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve article",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GeneratePrompt handles requests to generate encyclopedia-style prompts
func (h *EncyclopediaHandler) GeneratePrompt(c *gin.Context) {
	var request models.EncyclopediaPromptRequest
	if err := c.ShouldBindJSON(&request); err != nil {
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

	response, err := h.encyclopediaService.GenerateEncyclopediaPrompt(request, h.llamaService)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate prompt",
			"details": err.Error(),
		})
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

// Health check for encyclopedia service
func (h *EncyclopediaHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "encyclopedia",
		"message": "Encyclopedia service is running",
		"version": "1.0.0",
		"features": []string{
			"search",
			"article_retrieval",
			"prompt_generation",
			"multi_source_support",
			"multi_language_support",
		},
	})
}

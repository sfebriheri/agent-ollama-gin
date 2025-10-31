package http

import (
	"agent-ollama-gin/internal/container"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all HTTP routes
func RegisterRoutes(r *gin.Engine, c *container.Container) {
	// Initialize handlers
	llmHandler := NewLLMHandler(c.LLMUsecase, c.Logger)
	encyclopediaHandler := NewEncyclopediaHandler(c.EncyclopediaUsecase, c.LLMUsecase, c.Logger)

	// API routes
	api := r.Group("/api/v1")
	{
		// Health check
		api.GET("/health", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"status":  "ok",
				"message": "Encyclopedia Agent API is running",
				"version": "1.0.0",
			})
		})

		// LLM endpoints
		llama := api.Group("/llama")
		{
			llama.POST("/chat", llmHandler.Chat)
			llama.POST("/completion", llmHandler.Completion)
			llama.POST("/embedding", llmHandler.Embedding)
			llama.GET("/models", llmHandler.ListModels)
		}

		// Encyclopedia endpoints
		encyclopedia := api.Group("/encyclopedia")
		{
			encyclopedia.GET("/health", encyclopediaHandler.Health)
			encyclopedia.GET("/sources", encyclopediaHandler.GetSources)
			encyclopedia.GET("/languages", encyclopediaHandler.GetLanguages)
			encyclopedia.POST("/search", encyclopediaHandler.SearchEncyclopedia)
			encyclopedia.POST("/article", encyclopediaHandler.GetArticle)
			encyclopedia.POST("/prompt", encyclopediaHandler.GeneratePrompt)
		}
	}

	// Root route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":     "Welcome to Encyclopedia Agent API",
			"version":     "1.0.0",
			"description": "AI-powered encyclopedia agent with Llama integration",
			"endpoints": gin.H{
				"health": "/api/v1/health",
				"llama": gin.H{
					"chat":       "/api/v1/llama/chat",
					"completion": "/api/v1/llama/completion",
					"embedding":  "/api/v1/llama/embedding",
					"models":     "/api/v1/llama/models",
				},
				"encyclopedia": gin.H{
					"health":    "/api/v1/encyclopedia/health",
					"sources":   "/api/v1/encyclopedia/sources",
					"languages": "/api/v1/encyclopedia/languages",
					"search":    "/api/v1/encyclopedia/search",
					"article":   "/api/v1/encyclopedia/article",
					"prompt":    "/api/v1/encyclopedia/prompt",
				},
			},
			"features": []string{
				"AI-powered encyclopedia search",
				"Multi-source support (Wikipedia, Britannica)",
				"Multi-language support",
				"Intelligent prompt generation",
				"Llama LLM integration",
			},
			"docs": "Check README.md for full API documentation",
		})
	})
}

package main

import (
	"log"

	"llama-api/minimal/config"
	"llama-api/minimal/handlers"
	"llama-api/minimal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load minimal configuration
	cfg := config.LoadMinimal()

	// Initialize services
	llamaService := services.NewOptimizedLlamaService()

	// Initialize handlers
	llamaHandler := handlers.NewLlamaHandler(llamaService)

	// Create minimal Gin router
	r := gin.New() // Use gin.New() instead of gin.Default() for minimal logging

	// Configure minimal CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"GET", "POST"}
	corsConfig.AllowHeaders = []string{"Content-Type"}
	r.Use(cors.New(corsConfig))

	// Essential routes only
	api := r.Group("/api/v1")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"version": "1.0.0",
			})
		})

		// Core LLM endpoints only
		llama := api.Group("/llama")
		{
			llama.POST("/chat", llamaHandler.Chat)
			llama.POST("/completion", llamaHandler.Completion)
			llama.POST("/embedding", llamaHandler.Embedding)
			llama.GET("/models", llamaHandler.ListModels)
		}
	}

	// Start server
	log.Printf("Starting minimal Llama API on %s:%s", cfg.Host, cfg.Port)
	if err := r.Run(cfg.Host + ":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

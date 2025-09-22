package main

import (
	"log"
	"os"

	"agent-ollama-gin/handlers"
	"agent-ollama-gin/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize services
	llamaService := services.NewLlamaService()

	// Initialize handlers
	llamaHandler := handlers.NewLlamaHandler(llamaService)

	// Create Gin router
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// Root route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Llama API with Ollama Cloud Support",
			"version": "2.0.0",
			"endpoints": gin.H{
				"health":       "/api/v1/health",
				"chat":         "/api/v1/llama/chat",
				"completion":   "/api/v1/llama/completion",
				"embedding":    "/api/v1/llama/embedding",
				"models":       "/api/v1/llama/models",
				"cloud_models": "/api/v1/llama/cloud/models",
				"signin":       "/api/v1/llama/cloud/signin",
				"signout":      "/api/v1/llama/cloud/signout",
				"pull_model":   "/api/v1/llama/models/:model/pull",
				"stream_chat":  "/api/v1/llama/chat/stream",
			},
			"docs": "Check README.md for full API documentation",
			"features": []string{
				"Local Ollama models",
				"Ollama cloud models",
				"Authentication",
				"Streaming responses",
			},
		})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "Llama API is running",
				"version": "1.0.0",
			})
		})

		// Llama LLM endpoints
		llama := api.Group("/llama")
		{
			// Core endpoints
			llama.POST("/chat", llamaHandler.Chat)
			llama.POST("/completion", llamaHandler.Completion)
			llama.POST("/embedding", llamaHandler.Embedding)
			llama.GET("/models", llamaHandler.ListModels)

			// Streaming endpoints
			llama.POST("/chat/stream", llamaHandler.StreamChat)

			// Model management
			llama.POST("/models/:model/pull", llamaHandler.PullModel)

			// Cloud endpoints
			cloud := llama.Group("/cloud")
			{
				cloud.POST("/signin", llamaHandler.SignIn)
				cloud.POST("/signout", llamaHandler.SignOut)
				cloud.GET("/models", llamaHandler.ListCloudModels)
			}
		}

	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Llama API server with Ollama Cloud support on port %s", port)

	// Start the server
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

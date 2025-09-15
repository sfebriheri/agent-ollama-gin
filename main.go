package main

import (
	"context"
	"log"
	"os"

	"agent-ollama-gin/handlers"
	"agent-ollama-gin/services"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize Genkit context
	ctx := context.Background()
	g := genkit.Init(ctx, genkit.WithPlugins(&googlegenai.GoogleAI{}))

	// Define Genkit flows
	defineGenkitFlows(g, ctx)

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
			"message": "Welcome to Llama API",
			"version": "1.0.0",
			"endpoints": gin.H{
				"health":     "/api/v1/health",
				"chat":       "/api/v1/llama/chat",
				"completion": "/api/v1/llama/completion",
				"embedding":  "/api/v1/llama/embedding",
				"models":     "/api/v1/llama/models",
			},
			"docs": "Check README.md for full API documentation",
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
			llama.POST("/chat", llamaHandler.Chat)
			llama.POST("/completion", llamaHandler.Completion)
			llama.POST("/embedding", llamaHandler.Embedding)
			llama.GET("/models", llamaHandler.ListModels)
		}


	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Llama API server on port %s", port)
	
	// Start the server
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// defineGenkitFlows defines Genkit flows for AI operations
func defineGenkitFlows(g *genkit.Genkit, ctx context.Context) {
	// Basic text generation flow
	genkit.DefineFlow(g, "basicInferenceFlow",
		func(ctx context.Context, topic string) (string, error) {
			response, err := genkit.Generate(ctx, g,
				ai.WithModelName("googleai/gemini-2.5-flash"),
				ai.WithPrompt("Write a short, creative paragraph about %s.", topic),
				ai.WithConfig(&genai.GenerateContentConfig{
					Temperature: genai.Ptr[float32](0.8),
				}),
			)
			if err != nil {
				return "", err
			}
			return response.Text(), nil
		},
	)

	// Chat completion flow
	genkit.DefineFlow(g, "chatCompletionFlow",
		func(ctx context.Context, input struct {
			Messages []map[string]string `json:"messages"`
			Model    string              `json:"model,omitempty"`
		}) (string, error) {
			model := input.Model
			if model == "" {
				model = "googleai/gemini-2.5-flash"
			}

			// Convert messages to prompt
			prompt := ""
			for _, msg := range input.Messages {
				prompt += msg["role"] + ": " + msg["content"] + "\n"
			}

			response, err := genkit.Generate(ctx, g,
				ai.WithModelName(model),
				ai.WithPrompt(prompt),
				ai.WithConfig(&genai.GenerateContentConfig{
					Temperature: genai.Ptr[float32](0.7),
				}),
			)
			if err != nil {
				return "", err
			}
			return response.Text(), nil
		},
	)

	// Text completion flow
	genkit.DefineFlow(g, "textCompletionFlow",
		func(ctx context.Context, input struct {
			Prompt      string  `json:"prompt"`
			Model       string  `json:"model,omitempty"`
			Temperature float32 `json:"temperature,omitempty"`
		}) (string, error) {
			model := input.Model
			if model == "" {
				model = "googleai/gemini-2.5-flash"
			}

			temperature := input.Temperature
			if temperature == 0 {
				temperature = 0.7
			}

			response, err := genkit.Generate(ctx, g,
				ai.WithModelName(model),
				ai.WithPrompt(input.Prompt),
				ai.WithConfig(&genai.GenerateContentConfig{
					Temperature: genai.Ptr[float32](temperature),
				}),
			)
			if err != nil {
				return "", err
			}
			return response.Text(), nil
		},
	)
}

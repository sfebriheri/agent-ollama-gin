package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"agent-ollama-gin/internal/container"
	"agent-ollama-gin/internal/delivery/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize DI container
	c := container.New()

	// Create Gin router
	r := gin.Default()

	// Register middleware
	http.RegisterMiddleware(r, c)

	// Register routes
	http.RegisterRoutes(r, c)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting Encyclopedia Agent API server on port %s", port)
		if err := r.Run(":" + port); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")
}

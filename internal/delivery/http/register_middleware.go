package http

import (
	"agent-ollama-gin/internal/container"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// RegisterMiddleware registers all HTTP middleware
func RegisterMiddleware(r *gin.Engine, c *container.Container) {
	// CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// Logging middleware
	r.Use(LoggingMiddleware(c.Logger))

	// Recovery middleware
	r.Use(gin.Recovery())

	// Serve static files
	r.Static("/examples", "./examples")
}

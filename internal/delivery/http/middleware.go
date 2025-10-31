package http

import (
	"context"
	"net/http"
	"sync"
	"time"

	"agent-ollama-gin/pkg/logger"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware logs all HTTP requests
func LoggingMiddleware(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		duration := time.Since(startTime)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.RequestURI
		clientIP := c.ClientIP()

		if statusCode >= 400 {
			logger.Error(
				"HTTP Request",
				"method", method,
				"path", path,
				"status", statusCode,
				"ip", clientIP,
				"duration", duration.String(),
			)
		} else {
			logger.Info(
				"HTTP Request",
				"method", method,
				"path", path,
				"status", statusCode,
				"ip", clientIP,
				"duration", duration.String(),
			)
		}
	}
}

// ErrorHandlingMiddleware handles panics and errors
func ErrorHandlingMiddleware(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Internal server error",
					"details": "An unexpected error occurred",
				})
			}
		}()
		c.Next()
	}
}

// RateLimiter implements token bucket rate limiting
type RateLimiter struct {
	requestsPerSecond int
	burst             int
	tokens            float64
	lastRefillTime    time.Time
	mu                sync.RWMutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerSecond, burst int) *RateLimiter {
	return &RateLimiter{
		requestsPerSecond: requestsPerSecond,
		burst:             burst,
		tokens:            float64(burst),
		lastRefillTime:    time.Now(),
	}
}

// Allow checks if a request is allowed
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastRefillTime).Seconds()
	rl.tokens += elapsed * float64(rl.requestsPerSecond)

	if rl.tokens > float64(rl.burst) {
		rl.tokens = float64(rl.burst)
	}

	rl.lastRefillTime = now

	if rl.tokens >= 1 {
		rl.tokens--
		return true
	}

	return false
}

// RateLimitingMiddleware applies rate limiting to requests
func RateLimitingMiddleware(limiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequestValidationMiddleware validates request content type
func RequestValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			contentType := c.GetHeader("Content-Type")
			if contentType == "" || contentType != "application/json" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Content-Type must be application/json",
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// CORSMiddleware configures CORS headers
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// TimeoutMiddleware adds timeout to requests
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// RecoveryMiddleware recovers from panics and logs them
func RecoveryMiddleware(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic in request handler", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
			}
		}()
		c.Next()
	}
}

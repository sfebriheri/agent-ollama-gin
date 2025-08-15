package config

import (
	"os"
	"strconv"
)

// MinimalConfig contains only essential configuration
type MinimalConfig struct {
	Port         string
	Host         string
	LlamaBaseURL string
	LlamaTimeout int
	DefaultModel string
}

// LoadMinimal loads only essential configuration
func LoadMinimal() *MinimalConfig {
	return &MinimalConfig{
		Port:         getEnv("PORT", "8080"),
		Host:         getEnv("HOST", "0.0.0.0"),
		LlamaBaseURL: getEnv("LLAMA_BASE_URL", "http://localhost:11434"),
		LlamaTimeout: getEnvAsInt("LLAMA_TIMEOUT", 120),
		DefaultModel: getEnv("LLAMA_DEFAULT_MODEL", "phi3:mini"),
	}
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

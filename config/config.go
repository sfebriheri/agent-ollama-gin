package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Llama    LlamaConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port         string
	Host         string
	ReadTimeout  int
	WriteTimeout int
}

type LlamaConfig struct {
	BaseURL      string
	APIKey       string
	DefaultModel string
	Timeout      int
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			Host:         getEnv("HOST", "0.0.0.0"),
			ReadTimeout:  getEnvAsInt("READ_TIMEOUT", 30),
			WriteTimeout: getEnvAsInt("WRITE_TIMEOUT", 30),
		},
		Llama: LlamaConfig{
			BaseURL:      getEnv("LLAMA_BASE_URL", "http://localhost:11434"),
			APIKey:       getEnv("LLAMA_API_KEY", ""),
			DefaultModel: getEnv("LLAMA_DEFAULT_MODEL", "llama2"),
			Timeout:      getEnvAsInt("LLAMA_TIMEOUT", 60),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "llama_api"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
	}
}

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

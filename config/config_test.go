package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad_Defaults(t *testing.T) {
	// Clear environment variables
	os.Clearenv()

	config := Load()

	assert.NotNil(t, config)
	assert.Equal(t, "8080", config.Server.Port)
	assert.Equal(t, "0.0.0.0", config.Server.Host)
	assert.Equal(t, 30, config.Server.ReadTimeout)
	assert.Equal(t, 30, config.Server.WriteTimeout)

	assert.Equal(t, "http://localhost:11434", config.Llama.BaseURL)
	assert.Equal(t, "llama2", config.Llama.DefaultModel)
	assert.Equal(t, 60, config.Llama.Timeout)
	assert.False(t, config.Llama.CloudEnabled)
	assert.Equal(t, "https://api.ollama.com", config.Llama.CloudAPIURL)
}

func TestLoad_WithEnvironmentVariables(t *testing.T) {
	// Set environment variables
	os.Setenv("PORT", "9090")
	os.Setenv("HOST", "127.0.0.1")
	os.Setenv("READ_TIMEOUT", "60")
	os.Setenv("WRITE_TIMEOUT", "60")
	os.Setenv("LLAMA_BASE_URL", "http://localhost:11435")
	os.Setenv("LLAMA_DEFAULT_MODEL", "llama3")
	os.Setenv("LLAMA_TIMEOUT", "120")
	os.Setenv("LLAMA_CLOUD_ENABLED", "true")
	os.Setenv("LLAMA_CLOUD_API_URL", "https://custom.api.com")
	os.Setenv("LLAMA_CLOUD_API_KEY", "test-key")
	os.Setenv("LLAMA_SIGNED_IN", "true")

	defer func() {
		os.Clearenv()
	}()

	config := Load()

	assert.Equal(t, "9090", config.Server.Port)
	assert.Equal(t, "127.0.0.1", config.Server.Host)
	assert.Equal(t, 60, config.Server.ReadTimeout)
	assert.Equal(t, 60, config.Server.WriteTimeout)

	assert.Equal(t, "http://localhost:11435", config.Llama.BaseURL)
	assert.Equal(t, "llama3", config.Llama.DefaultModel)
	assert.Equal(t, 120, config.Llama.Timeout)
	assert.True(t, config.Llama.CloudEnabled)
	assert.Equal(t, "https://custom.api.com", config.Llama.CloudAPIURL)
	assert.Equal(t, "test-key", config.Llama.CloudAPIKey)
	assert.True(t, config.Llama.SignedIn)
}

func TestLoad_DatabaseConfig(t *testing.T) {
	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_SSL_MODE", "require")

	defer os.Clearenv()

	config := Load()

	assert.Equal(t, "db.example.com", config.Database.Host)
	assert.Equal(t, "5433", config.Database.Port)
	assert.Equal(t, "testuser", config.Database.User)
	assert.Equal(t, "testpass", config.Database.Password)
	assert.Equal(t, "testdb", config.Database.DBName)
	assert.Equal(t, "require", config.Database.SSLMode)
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "Returns environment value when set",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "custom",
			expected:     "custom",
		},
		{
			name:         "Returns default when not set",
			key:          "UNSET_KEY",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnv(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetEnvAsInt(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue int
		envValue     string
		expected     int
	}{
		{
			name:         "Returns parsed int when valid",
			key:          "TEST_INT",
			defaultValue: 10,
			envValue:     "25",
			expected:     25,
		},
		{
			name:         "Returns default when not set",
			key:          "UNSET_INT",
			defaultValue: 10,
			envValue:     "",
			expected:     10,
		},
		{
			name:         "Returns default when invalid int",
			key:          "INVALID_INT",
			defaultValue: 10,
			envValue:     "not-a-number",
			expected:     10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnvAsInt(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfig_Structs(t *testing.T) {
	// Test that all structs are properly defined
	config := &Config{
		Server: ServerConfig{
			Port:         "8080",
			Host:         "localhost",
			ReadTimeout:  30,
			WriteTimeout: 30,
		},
		Llama: LlamaConfig{
			BaseURL:      "http://localhost:11434",
			APIKey:       "test-key",
			DefaultModel: "llama2",
			Timeout:      60,
			CloudEnabled: false,
			CloudAPIURL:  "https://api.ollama.com",
			CloudAPIKey:  "cloud-key",
			SignedIn:     false,
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "password",
			DBName:   "testdb",
			SSLMode:  "disable",
		},
	}

	assert.NotNil(t, config)
	assert.Equal(t, "8080", config.Server.Port)
	assert.Equal(t, "llama2", config.Llama.DefaultModel)
	assert.Equal(t, "localhost", config.Database.Host)
}

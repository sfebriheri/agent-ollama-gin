variable "app_name" {
  description = "Name of the application"
  type        = string
  default     = "agent-ollama-gin"
}

variable "environment" {
  description = "Environment name (dev, staging, production)"
  type        = string
  default     = "dev"

  validation {
    condition     = contains(["dev", "staging", "production", "test"], var.environment)
    error_message = "Environment must be one of: dev, staging, production, test"
  }
}

variable "app_port" {
  description = "Port for the application to listen on"
  type        = number
  default     = 8080

  validation {
    condition     = var.app_port > 0 && var.app_port < 65536
    error_message = "Port must be between 1 and 65535"
  }
}

variable "ollama_base_url" {
  description = "Base URL for Ollama API"
  type        = string
  default     = "http://localhost:11434"
}

variable "ollama_default_model" {
  description = "Default Ollama model to use"
  type        = string
  default     = "llama2"
}

variable "ollama_timeout" {
  description = "Timeout for Ollama API requests (in seconds)"
  type        = number
  default     = 300

  validation {
    condition     = var.ollama_timeout > 0
    error_message = "Timeout must be greater than 0"
  }
}

variable "ollama_cloud_enabled" {
  description = "Enable Ollama cloud functionality"
  type        = bool
  default     = false
}

variable "ollama_cloud_api_url" {
  description = "Ollama cloud API URL"
  type        = string
  default     = "https://api.ollama.com"
}

variable "ollama_cloud_api_key" {
  description = "Ollama cloud API key (sensitive)"
  type        = string
  default     = ""
  sensitive   = true
}

variable "gin_mode" {
  description = "Gin framework mode (debug, release, test)"
  type        = string
  default     = "release"

  validation {
    condition     = contains(["debug", "release", "test"], var.gin_mode)
    error_message = "Gin mode must be one of: debug, release, test"
  }
}

variable "generate_config" {
  description = "Generate application configuration file"
  type        = bool
  default     = false
}

variable "setup_test_env" {
  description = "Setup test environment"
  type        = bool
  default     = false
}

variable "run_health_check" {
  description = "Run health check validation"
  type        = bool
  default     = false
}

variable "generate_ci_config" {
  description = "Generate CI/CD configuration file"
  type        = bool
  default     = false
}

variable "enable_monitoring" {
  description = "Enable monitoring and metrics"
  type        = bool
  default     = false
}

variable "log_level" {
  description = "Application log level"
  type        = string
  default     = "info"

  validation {
    condition     = contains(["debug", "info", "warn", "error"], var.log_level)
    error_message = "Log level must be one of: debug, info, warn, error"
  }
}

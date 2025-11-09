output "app_name" {
  description = "Application name"
  value       = var.app_name
}

output "environment" {
  description = "Environment name"
  value       = var.environment
}

output "app_port" {
  description = "Application port"
  value       = var.app_port
}

output "ollama_base_url" {
  description = "Ollama base URL"
  value       = var.ollama_base_url
}

output "ollama_default_model" {
  description = "Default Ollama model"
  value       = var.ollama_default_model
}

output "ollama_cloud_enabled" {
  description = "Ollama cloud enabled status"
  value       = var.ollama_cloud_enabled
}

output "gin_mode" {
  description = "Gin framework mode"
  value       = var.gin_mode
}

output "health_check_endpoint" {
  description = "Health check endpoint URL"
  value       = "http://localhost:${var.app_port}/api/v1/health"
}

output "api_base_url" {
  description = "API base URL"
  value       = "http://localhost:${var.app_port}/api/v1"
}

output "configuration_summary" {
  description = "Configuration summary"
  value = {
    application = {
      name        = var.app_name
      environment = var.environment
      port        = var.app_port
      mode        = var.gin_mode
      log_level   = var.log_level
    }
    ollama = {
      base_url      = var.ollama_base_url
      default_model = var.ollama_default_model
      timeout       = var.ollama_timeout
      cloud_enabled = var.ollama_cloud_enabled
      cloud_api_url = var.ollama_cloud_api_url
    }
    features = {
      monitoring      = var.enable_monitoring
      health_check    = var.run_health_check
      test_env_setup  = var.setup_test_env
      config_gen      = var.generate_config
      ci_config_gen   = var.generate_ci_config
    }
  }
}

terraform {
  required_version = ">= 1.5.0"

  required_providers {
    null = {
      source  = "hashicorp/null"
      version = "~> 3.2"
    }
    local = {
      source  = "hashicorp/local"
      version = "~> 2.4"
    }
  }

  # Uncomment for remote backend (e.g., S3, GCS, Azure)
  # backend "s3" {
  #   bucket = "your-terraform-state-bucket"
  #   key    = "agent-ollama-gin/terraform.tfstate"
  #   region = "us-east-1"
  # }
}

# Provider configuration
# Note: This is a minimal example. Add cloud provider configs as needed.

# Local variables
locals {
  app_name    = var.app_name
  environment = var.environment
  common_tags = {
    Application = local.app_name
    Environment = local.environment
    ManagedBy   = "Terraform"
    Repository  = "agent-ollama-gin"
  }
}

# Example: Application configuration file generation
resource "local_file" "app_config" {
  count = var.generate_config ? 1 : 0

  filename = "${path.module}/../.env.${var.environment}"
  content  = <<-EOF
    # Auto-generated configuration for ${var.environment}
    OLLAMA_BASE_URL=${var.ollama_base_url}
    OLLAMA_DEFAULT_MODEL=${var.ollama_default_model}
    OLLAMA_TIMEOUT=${var.ollama_timeout}
    OLLAMA_CLOUD_ENABLED=${var.ollama_cloud_enabled}
    OLLAMA_CLOUD_API_URL=${var.ollama_cloud_api_url}
    PORT=${var.app_port}
    GIN_MODE=${var.gin_mode}
  EOF

  file_permission = "0644"
}

# Example: Test environment setup
resource "null_resource" "test_setup" {
  count = var.setup_test_env ? 1 : 0

  triggers = {
    config_hash = md5(jsonencode({
      app_name              = var.app_name
      environment           = var.environment
      app_port              = var.app_port
      ollama_base_url       = var.ollama_base_url
      ollama_default_model  = var.ollama_default_model
      ollama_timeout        = var.ollama_timeout
      ollama_cloud_enabled  = var.ollama_cloud_enabled
      ollama_cloud_api_url  = var.ollama_cloud_api_url
      gin_mode              = var.gin_mode
    }))
  }

  provisioner "local-exec" {
    command = <<-EOT
      echo "Setting up test environment for ${var.environment}"
      echo "Application: ${var.app_name}"
      echo "Port: ${var.app_port}"
    EOT
  }
}

# Example: Health check validation
resource "null_resource" "health_check" {
  count = var.run_health_check ? 1 : 0

  depends_on = [null_resource.test_setup]

  provisioner "local-exec" {
    command = <<-EOT
      echo "Health check configuration ready"
      echo "Endpoint: http://localhost:${var.app_port}/api/v1/health"
    EOT
  }
}

# Example: CI/CD pipeline configuration
resource "local_file" "ci_config" {
  count = var.generate_ci_config ? 1 : 0

  filename = "${path.module}/ci-config.json"
  content = jsonencode({
    app_name    = var.app_name
    environment = var.environment
    port        = var.app_port
    healthcheck = {
      enabled  = var.run_health_check
      endpoint = "/api/v1/health"
      timeout  = 30
    }
    ollama = {
      base_url      = var.ollama_base_url
      default_model = var.ollama_default_model
      timeout       = var.ollama_timeout
      cloud_enabled = var.ollama_cloud_enabled
    }
  })

  file_permission = "0644"
}

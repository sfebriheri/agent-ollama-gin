variable "github_owner" {
  description = "GitHub repository owner (username or organization)"
  type        = string
}

variable "github_token" {
  description = "GitHub personal access token for authentication"
  type        = string
  sensitive   = true
}

variable "repository_name" {
  description = "GitHub repository name"
  type        = string
}

variable "go_version" {
  description = "Go version to use in CI/CD"
  type        = string
  default     = "1.25.1"
}

variable "docker_registry" {
  description = "Docker registry URL"
  type        = string
  default     = "docker.io"
}

variable "docker_image_name" {
  description = "Docker image name"
  type        = string
  default     = "agent-ollama-gin"
}

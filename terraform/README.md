# Terraform Configuration for Agent Ollama Gin

This directory contains Terraform configuration for managing infrastructure and configuration for the Agent Ollama Gin application.

## Overview

The Terraform configuration provides:
- Application configuration generation
- Environment setup automation
- CI/CD pipeline configuration
- Health check validation
- Infrastructure as Code (IaC) for deployment

## Prerequisites

- [Terraform](https://www.terraform.io/downloads.html) >= 1.5.0
- Go 1.23+ (for the application)

## Quick Start

1. **Initialize Terraform**
   ```bash
   cd terraform
   terraform init
   ```

2. **Copy and customize variables**
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   # Edit terraform.tfvars with your settings
   ```

3. **Validate configuration**
   ```bash
   terraform validate
   ```

4. **Plan changes**
   ```bash
   terraform plan
   ```

5. **Apply configuration**
   ```bash
   terraform apply
   ```

## Configuration Files

- `main.tf` - Main Terraform configuration
- `variables.tf` - Variable definitions
- `outputs.tf` - Output values
- `terraform.tfvars.example` - Example variable values

## Variables

### Application Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `app_name` | Application name | `agent-ollama-gin` |
| `environment` | Environment (dev/staging/production/test) | `dev` |
| `app_port` | Application port | `8080` |
| `gin_mode` | Gin mode (debug/release/test) | `release` |
| `log_level` | Log level (debug/info/warn/error) | `info` |

### Ollama Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `ollama_base_url` | Ollama API base URL | `http://localhost:11434` |
| `ollama_default_model` | Default model | `llama2` |
| `ollama_timeout` | API timeout (seconds) | `300` |
| `ollama_cloud_enabled` | Enable cloud features | `false` |
| `ollama_cloud_api_url` | Cloud API URL | `https://api.ollama.com` |
| `ollama_cloud_api_key` | Cloud API key (sensitive) | `""` |

### Feature Flags

| Variable | Description | Default |
|----------|-------------|---------|
| `generate_config` | Generate app config file | `false` |
| `setup_test_env` | Setup test environment | `false` |
| `run_health_check` | Run health checks | `false` |
| `generate_ci_config` | Generate CI config | `false` |
| `enable_monitoring` | Enable monitoring | `false` |

## Outputs

After applying, Terraform will output:
- Application configuration summary
- Health check endpoint
- API base URL
- Environment settings

View outputs:
```bash
terraform output
```

## Environments

### Development
```hcl
environment = "dev"
gin_mode    = "debug"
log_level   = "debug"
```

### Test/CI
```hcl
environment        = "test"
gin_mode          = "test"
setup_test_env    = true
run_health_check  = true
```

### Production
```hcl
environment          = "production"
gin_mode            = "release"
log_level           = "info"
enable_monitoring   = true
ollama_cloud_enabled = true
```

## CI/CD Integration

The Terraform configuration integrates with GitHub Actions:
- Validates on every push
- Plans on pull requests
- Can auto-apply on main branch (optional)

See `.github/workflows/terraform.yml` for details.

## Best Practices

1. **Never commit sensitive data**
   - Use environment variables for secrets
   - Keep `terraform.tfvars` out of version control
   - Use sensitive variable marking

2. **Use remote state**
   - Configure S3, GCS, or Azure backend
   - Enable state locking
   - Use workspaces for environments

3. **Version control**
   - Pin provider versions
   - Use semantic versioning
   - Document changes

4. **Testing**
   - Validate before apply
   - Use `terraform plan` in CI
   - Test in dev/staging first

## Common Commands

```bash
# Initialize
terraform init

# Validate syntax
terraform validate

# Format code
terraform fmt

# Plan changes
terraform plan

# Apply changes
terraform apply

# Show current state
terraform show

# List outputs
terraform output

# Destroy resources
terraform destroy
```

## Troubleshooting

### Configuration Issues
```bash
# Check validation
terraform validate

# View plan details
terraform plan -out=plan.tfplan
terraform show plan.tfplan
```

### State Issues
```bash
# Refresh state
terraform refresh

# View state
terraform state list
terraform state show <resource>
```

## Extending

To add cloud infrastructure (AWS, GCP, Azure):

1. Add provider configuration in `main.tf`
2. Define cloud resources
3. Update variables and outputs
4. Configure backend for state management

Example AWS provider:
```hcl
provider "aws" {
  region = var.aws_region
}

resource "aws_instance" "app" {
  # Your EC2 configuration
}
```

## Support

For issues or questions:
- Check Terraform documentation: https://www.terraform.io/docs
- Review application README: ../README.md
- Open an issue on GitHub

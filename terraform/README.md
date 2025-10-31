# Terraform CI/CD Configuration

This directory contains Terraform configuration for managing GitHub Actions workflows and branch protection rules for the Agent Ollama Gin project.

**For comprehensive CI/CD documentation, see**: `../docs/CI_CD.md`

## Overview

The Terraform configuration automates the setup of:
- **Build Workflow**: Compiles Go application and CLI
- **Test Workflow**: Runs unit tests and generates coverage reports
- **Lint Workflow**: Performs code quality checks (fmt, vet, golangci-lint, staticcheck)
- **Docker Workflow**: Builds and pushes Docker images
- **Security Workflow**: Runs security scanners (Gosec, Trivy, Nancy, CodeQL)
- **Branch Protection**: Enforces required checks on main and feat/encyclopedia branches

## Prerequisites

1. **Terraform**: Version 1.0 or higher
   ```bash
   terraform --version
   ```

2. **GitHub Token**: Personal access token with repo and workflow permissions
   - Create token at: https://github.com/settings/tokens/new
   - Required scopes: `repo`, `workflow`, `admin:repo_hook`

3. **GitHub Repository**: Push access to configure workflows

## Setup

### 1. Initialize Terraform

```bash
cd terraform
terraform init
```

### 2. Create terraform.tfvars

Copy the example file and fill in your values:

```bash
cp terraform.tfvars.example terraform.tfvars
```

Edit `terraform.tfvars`:
```hcl
github_owner      = "your-github-username"
repository_name   = "agent-ollama-gin"
go_version        = "1.25.1"
docker_image_name = "agent-ollama-gin"
docker_registry   = "docker.io"
```

### 3. Set GitHub Token

Set the GitHub token as an environment variable:

```bash
export TF_VAR_github_token="ghp_your_token_here"
```

Or add to `terraform.tfvars`:
```hcl
github_token = "ghp_your_token_here"
```

### 4. Plan and Apply

Preview changes:
```bash
terraform plan
```

Apply configuration:
```bash
terraform apply
```

## Workflows

### Build Workflow (.github/workflows/build.yml)
- Runs on: Push to main/feat/encyclopedia, Pull requests
- Tests: Go 1.25.1, 1.24.0
- Outputs: Binary artifacts

### Test Workflow (.github/workflows/test.yml)
- Runs on: Push to main/feat/encyclopedia, Pull requests
- Coverage: Unit tests with race detector
- Reports: Coverage HTML and Codecov upload

### Lint Workflow (.github/workflows/lint.yml)
- Runs on: Push to main/feat/encyclopedia, Pull requests
- Checks:
  - gofmt formatting
  - go vet analysis
  - golangci-lint
  - staticcheck
  - Module sync (go mod tidy)

### Docker Workflow (.github/workflows/docker.yml)
- Runs on: Push (main/feat/encyclopedia) and tags
- Registries: Docker Hub, GitHub Container Registry
- Caching: GitHub Actions cache

### Security Workflow (.github/workflows/security.yml)
- Runs on: Push, Pull requests, Weekly schedule
- Scanners:
  - Gosec (Go security)
  - Trivy (Vulnerability)
  - Nancy (Dependency)
  - CodeQL (SAST)

## Branch Protection

The configuration enforces:
- **main**: All status checks required (build, test, lint, security)
- **feat/encyclopedia**: Required checks (build, test, lint)

Admins can bypass these rules.

## Secrets

The Docker workflow requires these GitHub repository secrets:

```
DOCKER_USERNAME = your-docker-username
DOCKER_PASSWORD = your-docker-password (or token)
```

Set these in: Settings → Secrets and variables → Actions → New repository secret

## File Structure

```
terraform/
├── main.tf                      # Primary configuration
├── variables.tf                 # Input variables
├── outputs.tf                   # Output values
├── terraform.tfvars.example     # Example variables
├── workflows/
│   ├── build.yml               # Build workflow
│   ├── test.yml                # Test workflow
│   ├── lint.yml                # Lint workflow
│   ├── docker.yml              # Docker workflow
│   └── security.yml            # Security workflow
└── README.md                    # This file
```

## Usage

### View Current State

```bash
terraform state list
terraform state show 'github_repository_file.build_workflow'
```

### Update Workflows

Edit workflow files in `terraform/workflows/` and apply:

```bash
terraform apply
```

### Destroy Configuration

Remove all GitHub Actions workflows:

```bash
terraform destroy
```

## Troubleshooting

### Token Issues
```bash
# Verify token is set
echo $TF_VAR_github_token

# Regenerate token if needed
# https://github.com/settings/tokens
```

### Workflow Not Creating
```bash
# Check Terraform plan
terraform plan

# Verify file paths exist
ls -la workflows/
```

### Branch Protection Errors
- Ensure branch exists in repository
- Verify token has required scopes
- Check required status check names match workflow job names

### Docker Push Fails
- Verify Docker credentials in GitHub secrets
- Check image name format
- Ensure repository exists in registry

## Best Practices

1. **Keep Token Secure**
   - Use environment variables, not version control
   - Regenerate token if accidentally exposed
   - Rotate regularly

2. **Review Changes**
   - Always run `terraform plan` before `apply`
   - Review workflow changes in GitHub UI
   - Test workflows on feature branches first

3. **Monitor Workflows**
   - Check GitHub Actions tab for failures
   - Review artifact storage and retention
   - Clean up old artifacts periodically

4. **Update Dependencies**
   - Keep GitHub provider updated
   - Update action versions regularly
   - Test new Go versions in matrix strategy

## References

- [Terraform GitHub Provider](https://registry.terraform.io/providers/integrations/github/latest)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Go Testing](https://golang.org/pkg/testing/)
- [Docker Build Action](https://github.com/docker/build-push-action)

## Support

For issues or questions:
1. Check Terraform state: `terraform state list`
2. Review GitHub Actions logs
3. Verify credentials and permissions
4. Check GitHub provider version compatibility

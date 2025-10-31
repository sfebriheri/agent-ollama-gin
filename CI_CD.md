# CI/CD Pipeline Documentation

This document explains the continuous integration and continuous deployment (CI/CD) pipeline for the Agent Ollama Gin project.

## Overview

The CI/CD pipeline is managed through Terraform and GitHub Actions. It automates:
- Building and testing Go applications
- Code quality checks and linting
- Security vulnerability scanning
- Docker image building and pushing
- Branch protection enforcement

## Architecture

```
┌─────────────────┐
│  Git Push/PR    │
└────────┬────────┘
         │
    ┌────▼────────────────────────────────────┐
    │  GitHub Actions Workflows                │
    ├────────────────────────────────────────┤
    │  ✓ Build (Go 1.25.1, 1.24.0)          │
    │  ✓ Test (Unit tests + Coverage)        │
    │  ✓ Lint (fmt, vet, golangci-lint)     │
    │  ✓ Security (Gosec, Trivy, CodeQL)    │
    │  ✓ Docker (Build & Push)               │
    └────┬───────────────────────────────────┘
         │
    ┌────▼────────────────────────┐
    │  Branch Protection Rules     │
    ├────────────────────────────┤
    │  main: All checks required │
    │  feat/*: Selected checks   │
    └────────────────────────────┘
```

## Workflows

### 1. Build Workflow

**File**: `.github/workflows/build.yml`

Compiles the Go application for multiple versions.

**Triggers**:
- Push to main or feat/encyclopedia
- Pull requests to main or feat/encyclopedia

**Jobs**:
- Compile for Go 1.25.1 and 1.24.0
- Build main binary: `bin/agent-ollama-gin`
- Build CLI binary: `bin/encyclopedia`
- Store artifacts for 7 days

**Status Check**: `build`

```bash
# Manual trigger (if needed)
gh workflow run build.yml --ref feat/encyclopedia
```

### 2. Test Workflow

**File**: `.github/workflows/test.yml`

Runs unit tests with coverage analysis.

**Triggers**:
- Push to main or feat/encyclopedia
- Pull requests to main or feat/encyclopedia

**Jobs**:
- Run tests: `go test -v -race -coverprofile=coverage.out ./...`
- Generate coverage HTML report
- Upload to Codecov
- Test on Go 1.25.1 and 1.24.0

**Status Check**: `test`

**Coverage**:
- View HTML report in artifacts
- Track trends on Codecov dashboard

```bash
# View test results
gh run view <run-id> --log
```

### 3. Lint Workflow

**File**: `.github/workflows/lint.yml`

Performs code quality analysis.

**Triggers**:
- Push to main or feat/encyclopedia
- Pull requests to main or feat/encyclopedia

**Checks**:
- `gofmt`: Code formatting
- `go vet`: Static analysis
- `go mod tidy`: Dependency consistency
- `golangci-lint`: Comprehensive linting
- `staticcheck`: Additional static checks

**Status Check**: `lint`

**Failure Causes**:
- Unformatted code → Run `go fmt ./...`
- Unused imports → Run `go mod tidy`
- Lint violations → Fix reported issues

### 4. Docker Workflow

**File**: `.github/workflows/docker.yml`

Builds and pushes Docker images.

**Triggers**:
- Push to main or feat/encyclopedia
- Tags matching `v*` (releases)
- Pull requests (build only, no push)

**Registries**:
- Docker Hub: `${{ secrets.DOCKER_USERNAME }}/agent-ollama-gin`
- GitHub Container Registry: `ghcr.io/your-org/agent-ollama-gin`

**Tags**:
- `main` → `latest`
- `feat/encyclopedia` → `feat-encyclopedia`
- `v1.2.3` → `1.2.3`, `1.2`, `1`
- Commit SHA → `sha-<short-hash>`

**Requirements**:
- Set `DOCKER_USERNAME` secret
- Set `DOCKER_PASSWORD` secret

### 5. Security Workflow

**File**: `.github/workflows/security.yml`

Runs security vulnerability and code analysis scans.

**Triggers**:
- Push to main or feat/encyclopedia
- Pull requests to main or feat/encyclopedia
- Weekly schedule (Sundays at 00:00 UTC)

**Scanners**:

1. **Gosec**: Go security scanner
   - Detects security issues in Go code
   - Report: `gosec-report.json`

2. **Trivy**: Container vulnerability scanner
   - Scans filesystem for vulnerabilities
   - Format: SARIF (uploaded to GitHub Security tab)

3. **Nancy**: Go dependency vulnerabilities
   - Checks known vulnerable dependencies
   - Part of SBOM analysis

4. **CodeQL**: Semantic code analysis
   - SAST (Static Application Security Testing)
   - Custom queries for Go
   - Results in GitHub Security tab

**Status Check**: `security` (for scheduled runs)

**Viewing Results**:
- GitHub UI: Security → Code scanning
- Artifacts: Download reports
- GitHub Issues: Auto-created for vulnerabilities

## Branch Protection Rules

### Main Branch

**Protection**: `.github/workflows/build.yml`, `.github/workflows/test.yml`, `.github/workflows/lint.yml`, `.github/workflows/security.yml`

**Requirements**:
- ✓ All status checks must pass
- ✓ Branches must be up to date
- ✓ Dismiss stale PR reviews
- ✓ Require pull request

**Bypass**: Admins only

### Feat/Encyclopedia Branch

**Protection**: `.github/workflows/build.yml`, `.github/workflows/test.yml`, `.github/workflows/lint.yml`

**Requirements**:
- ✓ Build, test, lint checks required
- ✓ Branches must be up to date

## Setting Up CI/CD

### Prerequisites

1. **Repository**: Push access to agent-ollama-gin
2. **GitHub Token**: With repo and workflow scopes
3. **Terraform**: Version 1.0+
4. **Docker Hub** (optional): For pushing images

### Installation Steps

1. **Clone and navigate to terraform directory**:
```bash
cd terraform
terraform init
```

2. **Create terraform.tfvars**:
```bash
cp terraform.tfvars.example terraform.tfvars
# Edit with your values
```

3. **Set GitHub token**:
```bash
export TF_VAR_github_token="ghp_your_token_here"
```

4. **Apply Terraform**:
```bash
terraform plan
terraform apply
```

5. **Configure Docker secrets** (if using Docker workflow):
   - Go to GitHub repository Settings
   - Secrets and variables → Actions
   - Add `DOCKER_USERNAME` and `DOCKER_PASSWORD`

6. **Verify workflows**:
   - Go to repository Actions tab
   - Confirm workflows are present and enabled

## Workflow Status

### Check Status

```bash
# View all workflow runs
gh run list --repo sfebriheri/agent-ollama-gin

# View specific workflow
gh run list --workflow=build.yml

# View detailed run
gh run view <run-id> --log
```

### Debugging Failures

1. **Build Fails**:
   - Check Go version compatibility
   - Verify dependencies: `go mod tidy`
   - Check import paths are correct

2. **Test Fails**:
   - Run locally: `go test -v -race ./...`
   - Check race condition detection
   - Verify test data and fixtures

3. **Lint Fails**:
   - Format code: `go fmt ./...`
   - Run vet: `go vet ./...`
   - Fix lint issues: `golangci-lint run ./...`

4. **Security Fails**:
   - Review Gosec warnings
   - Update vulnerable dependencies
   - Check CodeQL results

## Local Development

### Run Tests Locally

```bash
# Unit tests
go test -v -race ./...

# With coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Code Quality Checks

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run ./...

# Static analysis
go vet ./...
staticcheck ./...
```

### Build Locally

```bash
# Build main binary
go build -o bin/agent-ollama-gin ./

# Build CLI
go build -o bin/encyclopedia ./cmd/encyclopedia

# Docker build
docker build -t agent-ollama-gin:local .
```

## Performance & Optimization

### Caching

All workflows use GitHub Actions cache for:
- Go modules: `~/go/pkg/mod`
- Docker layers: GitHub Actions cache

**Cache Key**:
- Go modules: Hash of `go.sum`
- Invalidates when dependencies change

### Parallel Testing

Tests run in parallel where possible:
- Multiple Go versions tested simultaneously
- Go test parallelization enabled (race detector)

### Build Times

Typical workflow execution times:
- Build: 2-3 minutes
- Test: 3-5 minutes
- Lint: 2-3 minutes
- Docker: 5-10 minutes
- Security: 5-8 minutes

## Monitoring & Alerts

### GitHub Notifications

- Workflow failures email author
- Status checks in PR
- Check suite details in conversation

### Metrics to Monitor

- Build success rate
- Test coverage percentage
- Lint warning count
- Security vulnerability count
- Workflow execution time

### Integration with Tools

- **Codecov**: Coverage reports and trends
- **GitHub Security**: Vulnerability tracking
- **GitHub Projects**: Automated task management

## Troubleshooting

### Workflows Not Running

1. Check if workflows are enabled:
   - Settings → Actions → General
   - Ensure "Allow all actions and reusable workflows" is selected

2. Check branch is correct:
   - Workflows only run on pushed branches
   - Test on feature branch first

3. Verify Terraform applied successfully:
   - Check `.github/workflows/` directory exists
   - Verify files are committed

### Status Checks Failing

1. View detailed logs:
```bash
gh run view <run-id>
```

2. Check required status checks in branch protection:
   - Must match workflow job names exactly
   - Case-sensitive

3. Rerun failed workflow:
```bash
gh run rerun <run-id>
```

### Docker Push Fails

1. Verify secrets are set:
```bash
# In repository settings, check:
# - DOCKER_USERNAME
# - DOCKER_PASSWORD
```

2. Check credentials are valid:
```bash
docker login -u $DOCKER_USERNAME
```

3. Ensure image name is correct:
   - Must match Docker Hub username
   - Lowercase characters only

## Security Best Practices

1. **Token Management**:
   - Use environment variables for tokens
   - Regenerate if exposed
   - Limit token scope

2. **Secrets Management**:
   - Store in GitHub Secrets, not code
   - Rotate regularly
   - Use specific service accounts

3. **Dependency Security**:
   - Review dependency updates
   - Run security scans regularly
   - Keep tools updated

## CI/CD Pipeline Customization

### Add New Workflow

1. Create workflow file in `terraform/workflows/`
2. Update `main.tf` with new resource
3. Apply Terraform: `terraform apply`

### Modify Existing Workflow

1. Edit workflow file in `terraform/workflows/`
2. Test changes locally
3. Apply: `terraform apply`

### Disable Workflow

```bash
terraform state rm github_repository_file.workflow_name
# Or manually delete from .github/workflows/
```

## References

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Terraform GitHub Provider](https://registry.terraform.io/providers/integrations/github/latest)
- [Go Testing](https://golang.org/pkg/testing/)
- [golangci-lint](https://golangci-lint.run/)
- [Docker GitHub Action](https://github.com/docker/build-push-action)

## Support

For issues or questions:
1. Check workflow logs in GitHub Actions
2. Review this documentation
3. Check Terraform state: `terraform state list`
4. Review GitHub provider configuration

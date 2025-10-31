# Documentation

This directory contains comprehensive documentation for the Agent Ollama Gin project.

## Contents

### CI/CD Pipeline Documentation
- **File**: `CI_CD.md`
- **Purpose**: Complete guide to the GitHub Actions CI/CD pipeline
- **Topics**:
  - Workflow overview and architecture
  - 5 GitHub Actions workflows (build, test, lint, docker, security)
  - Branch protection rules
  - Setup and configuration
  - Debugging and troubleshooting
  - Local development
  - Performance optimization
  - Security best practices

### API Documentation
- **Location**: Root level `README.md`
- **Includes**: API endpoints, examples, configuration

### Architecture Documentation
- **Location**: Root level `README.md`
- **Includes**: Clean architecture layers, design patterns

## Quick Links

### For CI/CD Setup
1. Read: `docs/CI_CD.md`
2. Navigate: `terraform/` directory
3. Follow: `terraform/README.md`

### For Development
1. Check: Root `README.md` for API endpoints
2. Run: `make dev-setup` for local environment
3. Test: `go test -v ./...`

### For Deployment
1. Review: `docs/CI_CD.md` - Docker section
2. Configure: `terraform/terraform.tfvars`
3. Apply: `terraform apply`

## Documentation Files Structure

```
docs/
├── README.md                    # This file
└── CI_CD.md                     # CI/CD pipeline documentation
```

```
terraform/
├── README.md                    # Terraform setup guide
├── main.tf                      # Primary configuration
├── variables.tf                 # Input variables
├── outputs.tf                   # Output values
├── terraform.tfvars.example     # Configuration template
├── .gitignore                   # Git ignore rules
└── workflows/
    ├── build.yml                # Build workflow
    ├── test.yml                 # Test workflow
    ├── lint.yml                 # Lint workflow
    ├── docker.yml               # Docker workflow
    └── security.yml             # Security workflow
```

## Getting Started

### New Developers
1. Read `README.md` in root directory
2. Run `make dev-setup` for environment setup
3. Check `docs/CI_CD.md` for CI/CD information

### DevOps/SRE
1. Read `terraform/README.md`
2. Follow Terraform setup in `docs/CI_CD.md`
3. Configure GitHub secrets for Docker

### Contributors
1. Understand clean architecture in root `README.md`
2. Follow `docs/CI_CD.md` for workflow expectations
3. Ensure local tests pass before submitting PR

## Key Documentation Links

| Document | Purpose | Audience |
|----------|---------|----------|
| `CI_CD.md` | GitHub Actions workflows | DevOps, Maintainers |
| `terraform/README.md` | Infrastructure as Code setup | DevOps, SRE |
| Root `README.md` | Project overview & API | All developers |
| `Makefile` | Build & dev commands | All developers |

## Workflow Overview

```
┌──────────────┐
│ Contributing │
└──────┬───────┘
       │
       ├─► `CI/CD.md` - Understand CI/CD
       ├─► `Makefile` - Local build/test
       └─► `README.md` - API & architecture
       
       │
    ┌──▼──────────────────┐
    │ GitHub CI/CD Runs   │
    │ (5 workflows)       │
    └──┬───────────────────┘
       │
       ├─► Build: Go 1.25.1, 1.24.0
       ├─► Test: Unit tests + coverage
       ├─► Lint: Code quality checks
       ├─► Docker: Image build/push
       └─► Security: Vulnerability scanning
```

## Important Notes

### CI/CD Configuration
- Workflows are managed via Terraform
- Do not manually edit `.github/workflows/` files
- All changes go through `terraform/workflows/` → `terraform apply`

### Documentation Maintenance
- Keep `CI_CD.md` updated with workflow changes
- Update `terraform/README.md` for Terraform changes
- Sync `docs/` with main `README.md` where applicable

### Access & Permissions
- Terraform requires GitHub token with `repo` and `workflow` scopes
- Docker secrets require repository admin access
- Branch protection rules enforced on `main` and `feat/encyclopedia`

## Support & Troubleshooting

### For CI/CD Issues
1. Check `docs/CI_CD.md` - Troubleshooting section
2. Review workflow logs in GitHub Actions tab
3. Run local tests: `go test -v -race ./...`

### For Terraform Issues
1. Check `terraform/README.md` - Troubleshooting section
2. Verify terraform state: `terraform state list`
3. Review GitHub provider logs: `TF_LOG=DEBUG terraform apply`

### For Build/Test Issues
1. Follow `Makefile` targets locally
2. Check Go version: `go version`
3. Run tests: `go test -v -race ./...`
4. Format code: `go fmt ./...`

## Contributing to Documentation

When updating documentation:
1. Maintain clear, concise language
2. Include examples where applicable
3. Update table of contents if needed
4. Keep file paths relative to project root
5. Verify all links work correctly

## See Also

- Main README: `../README.md`
- Terraform Guide: `../terraform/README.md`
- Makefile: `../Makefile`
- GitHub Actions: `.github/workflows/`

---

Last Updated: 2025-10-31
For latest updates, check the git history: `git log --oneline docs/`

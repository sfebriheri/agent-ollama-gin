terraform {
  required_version = ">= 1.0"
  required_providers {
    github = {
      source  = "integrations/github"
      version = "~> 6.0"
    }
  }
}

provider "github" {
  owner = var.github_owner
  token = var.github_token
}

# Enable Actions for the repository
resource "github_repository_environment" "github_actions" {
  environment = "github-actions"
  repository  = var.repository_name
}

# Create workflow files
resource "github_repository_file" "build_workflow" {
  repository          = var.repository_name
  branch              = "feat/encyclopedia"
  file                = ".github/workflows/build.yml"
  content             = file("${path.module}/../build.yml")
  commit_message      = "Add CI/CD: Build workflow"
  commit_author       = "Terraform"
  commit_email        = "terraform@github.com"
  overwrite_on_create = true
}

resource "github_repository_file" "test_workflow" {
  repository          = var.repository_name
  branch              = "feat/encyclopedia"
  file                = ".github/workflows/test.yml"
  content             = file("${path.module}/../test.yml")
  commit_message      = "Add CI/CD: Test workflow"
  commit_author       = "Terraform"
  commit_email        = "terraform@github.com"
  overwrite_on_create = true
}

resource "github_repository_file" "lint_workflow" {
  repository          = var.repository_name
  branch              = "feat/encyclopedia"
  file                = ".github/workflows/lint.yml"
  content             = file("${path.module}/../lint.yml")
  commit_message      = "Add CI/CD: Lint workflow"
  commit_author       = "Terraform"
  commit_email        = "terraform@github.com"
  overwrite_on_create = true
}

resource "github_repository_file" "docker_workflow" {
  repository          = var.repository_name
  branch              = "feat/encyclopedia"
  file                = ".github/workflows/docker.yml"
  content             = file("${path.module}/../docker.yml")
  commit_message      = "Add CI/CD: Docker build workflow"
  commit_author       = "Terraform"
  commit_email        = "terraform@github.com"
  overwrite_on_create = true
}

resource "github_repository_file" "security_workflow" {
  repository          = var.repository_name
  branch              = "feat/encyclopedia"
  file                = ".github/workflows/security.yml"
  content             = file("${path.module}/../security.yml")
  commit_message      = "Add CI/CD: Security checks workflow"
  commit_author       = "Terraform"
  commit_email        = "terraform@github.com"
  overwrite_on_create = true
}

# Configure branch protection rule
resource "github_branch_protection" "main" {
  repository_id            = var.repository_name
  pattern                  = "main"
  enforce_admins           = false
  require_status_checks    = true
  require_branches_up_to_date = true

  required_status_checks {
    strict   = true
    contexts = [
      "build",
      "test",
      "lint",
      "security"
    ]
  }

  required_pull_request_reviews {
    dismiss_stale_reviews = true
    require_code_owner_reviews = false
  }
}

resource "github_branch_protection" "feat_encyclopedia" {
  repository_id            = var.repository_name
  pattern                  = "feat/encyclopedia"
  enforce_admins           = false
  require_status_checks    = true
  require_branches_up_to_date = true

  required_status_checks {
    strict   = true
    contexts = [
      "build",
      "test",
      "lint"
    ]
  }
}

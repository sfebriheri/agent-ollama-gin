# Outputs disabled for testing
# Uncomment these after testing is complete

# output "workflows_created" {
#   description = "List of created GitHub Actions workflows"
#   value = [
#     github_repository_file.build_workflow.file,
#     github_repository_file.test_workflow.file,
#     github_repository_file.lint_workflow.file,
#     github_repository_file.docker_workflow.file,
#     github_repository_file.security_workflow.file,
#   ]
# }

# output "branch_protection_main" {
#   description = "Branch protection rule for main branch"
#   value       = github_branch_protection.main.pattern
# }

# output "branch_protection_feat_encyclopedia" {
#   description = "Branch protection rule for feat/encyclopedia branch"
#   value       = github_branch_protection.feat_encyclopedia.pattern
# }

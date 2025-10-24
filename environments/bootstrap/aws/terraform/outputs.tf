output "bucket_name" {
  value = aws_s3_bucket.tf_state_bucket.bucket 
}

output "dynamodb_table_name" {
  value = aws_dynamodb_table.tf_lock_db.name
}

locals {
  repo_root = "${path.module}/../../../.."
}

resource "local_file" "backend_info_json" {
    filename = "${local.repo_root}/environments/aws/backend.json"
    content = jsonencode({
        bucket = aws_s3_bucket.tf_state_bucket.bucket
        dynamodb_table = aws_dynamodb_table.tf_lock_db.name
        region = var.region
    })
}

output "gh_oidc_provider_arn" {
  value = aws_iam_openid_connect_provider.gh_oidc_provider.arn
}

output "gh_actions_role_arn" {
  value = aws_iam_role.gh_actions_role.arn
}

resource "local_file" "gh_actions_ecs_change" {
    filename = "${local.repo_root}/.github/workflows/ecs-image-change.yml"
    content = templatefile("${path.module}/templates/ecs-image-change.yml.tmpl", {
        envs = var.envs
        region = var.region
        included_branches = join(" || ", [for env in var.envs : "github.ref == 'refs/heads/${env.branch}'"])
        role_arn = aws_iam_role.gh_actions_role.arn
    })
}

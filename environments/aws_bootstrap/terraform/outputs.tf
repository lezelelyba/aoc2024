output "bucket_name" {
  value = aws_s3_bucket.tf_state_bucket.bucket 
}

output "dynamodb_table_name" {
  value = aws_dynamodb_table.tf_lock_db.name
}

resource "local_file" "backend_info" {
  for_each = var.envs  

    filename = "${path.module}/../../aws/${each.value.prefix}/terraform/backend.tf"
    content = templatefile("${path.module}/templates/backend.tf.tmpl", {
        bucket = aws_s3_bucket.tf_state_bucket.bucket
        lock_table = aws_dynamodb_table.tf_lock_db.name
        region = var.region
        environment = each.value.prefix
    })
}
resource "local_file" "bootstrap_variables" {
  for_each = var.envs  

    filename = "${path.module}/../../aws/${each.value.prefix}/terraform/boostrap-variables.tf"
    content = templatefile("${path.module}/templates/bootstrap-variables.tf.tmpl", {
        region = var.region
    })
}

output "gh_oidc_provider_arn" {
  value = aws_iam_openid_connect_provider.gh_oidc_provider.arn
}

output "gh_actions_role_arn" {
  value = aws_iam_role.gh_actions_role.arn
}

resource "local_file" "gh_actions_deploy" {

  for_each = var.envs
    filename = "${path.module}/../../../.github/workflows/${each.value.prefix}_deploy.yml"
    content = templatefile("${path.module}/templates/deploy.yml.tmpl", {
        region = var.region
        role_arn = aws_iam_role.gh_actions_role.arn
        branch = each.value.branch
    })
}
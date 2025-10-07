resource "aws_s3_bucket" "tf_state_bucket" {
    provider = aws.prov1
    bucket = "${var.tf_state_bucket}-${random_string.bucket_number.result}"

    force_destroy = true

    tags = {
        managed-by = var.mngd
    }
}

resource "aws_dynamodb_table" "tf_lock_db" {
    provider = aws.prov1
    name = "${var.tf_lock_db}"
    billing_mode = "PAY_PER_REQUEST"

    hash_key = "LockID"
    
    attribute {
        name = "LockID"
        type = "S"
    }

    tags = {
        managed-by = var.mngd
    }
}

resource "aws_iam_openid_connect_provider" "gh_oidc_provider" {
    url = "https://token.actions.githubusercontent.com"

    client_id_list = ["sts.amazonaws.com"]

    thumbprint_list = ["6938fd4d98bab03faadb97b34396831e3780aea1"]

    tags = {
        managed-by = var.mngd
    }
}

resource "aws_iam_role" "gh_actions_role" {
    provider = aws.prov1
    name = "gh-actions-tf-build-role"

    assume_role_policy = jsonencode({
        Version =  "2012-10-17",
        Statement = [
            {
                Effect = "Allow",
                Principal = {
                    Federated = aws_iam_openid_connect_provider.gh_oidc_provider.arn
                },
                Action = "sts:AssumeRoleWithWebIdentity",
                Condition = {
                    StringEquals = {
                        "token.actions.githubusercontent.com:sub" = "repo:${var.repo_name}:ref:refs/heads/*",
                        "token.actions.githubusercontent.com:aud" = "sts.amazonaws.com"
                    }
                }
            }
        ]
    })

    tags = {
        managed-by = var.mngd
    }
}

resource "aws_iam_policy" "gh_actions_policy" {
    name = "gh-actions-tf-build-policy" 

    policy = jsonencode({
      Version = "2012-10-17",
      Statement = [
        {
          Effect   = "Allow",
          Action   = [
            "iam:PassRole",
            "s3:*",
            "dynamodb:*",
            "ecs:*",
            "elasticloadbalancing:*"
          ],
          Resource = "*"
        }
      ]
    })
}

resource "aws_iam_role_policy_attachment" "gh_actions_policy_attach" {
    role = aws_iam_role.gh_actions_role.name
    policy_arn = aws_iam_policy.gh_actions_policy.arn
}

resource "aws_iam_user" "gh_actions_user" {
    provider = aws.prov1
    name = "gh-actions-user"

    tags = {
        managed-by = var.mngd
    }
}


resource "aws_s3_bucket_ownership_controls" "tf_state_ctrl" {
    bucket = aws_s3_bucket.tf_state_bucket.id
    rule {
        object_ownership = "BucketOwnerPreferred"
    }
}

resource "aws_s3_bucket_acl" "tf_state_acl" {
    depends_on = [aws_s3_bucket_ownership_controls.tf_state_ctrl]

    bucket = aws_s3_bucket.tf_state_bucket.id
    acl = "private"
}

resource "aws_s3_bucket_versioning" "tf_state_versioning" {
    bucket = aws_s3_bucket.tf_state_bucket.id
    versioning_configuration {
      status = "Enabled" 
    }
}
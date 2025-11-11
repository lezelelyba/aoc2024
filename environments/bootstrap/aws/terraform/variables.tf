variable "region" {
    description = "AWS region"
    default = "eu-central-1" 
}
variable "repo_name" {
    description = "name of github repo"
    default = "lezelelyba/aoc2024"
}
variable "mngd" {
    description = "tag added to resources created by this script"
    default = "TF-Bootstrap"
}

resource "random_string" "bucket_number" {
    length = 6
    upper = false
    special = false
}

variable "tf_state_bucket" {
    description = "S3 Bucket for terraform states - prefix"
    default = "tf-state-bucket"
}

variable "tf_lock_db" {
    description = "Dynamo DB Table name for terraform state locks"
    default = "tf-locks-db"
}
variable "envs" {
    description = "map of environments in the github repo"
    default = {
        dev = {
            prefix = "dev"
            branch = "dev"
        },
        stage = {
            prefix = "stage"
            branch = "stage"
        },
        prod = {
            prefix = "prod"
            branch = "master"
        }
    }
}
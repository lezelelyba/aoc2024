variable "mngd" {
    default = "TF-Bootstrap"
}

variable "resource-group-name" {
    default = "aoc-bootstrap"
}

variable "region" {
    default = "westeurope"
}

resource "random_string" "storage_account_number" {
    length = 6
    upper = false
    special = false
}

// for OIDC setup
variable "repo_name" {
    default = "lezelelyba/aoc2024"
}
variable "envs" {
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
locals {
    github_subs = [
        for e in var.envs :
        "repo:${var.repo_name}:ref:refs/heads/${e.branch}"
    ]
    
    github_sub_wildcard = "repo:${var.repo_name}:ref:refs/heads/*"
}
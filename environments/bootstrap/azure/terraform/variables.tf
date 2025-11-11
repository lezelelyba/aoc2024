variable "mngd" {
    description = "tag added to resources created by this script"
    default = "TF-Bootstrap"
}

variable "resource-group-name" {
    description = "azure resource group under which the resources will be created"
    default = "aoc-bootstrap"
}

variable "region" {
    description = "azure region"
    default = "westeurope"
}

resource "random_string" "storage_account_number" {
    length = 6
    upper = false
    special = false
}

// OIDC setup
variable "gh_actions_application_name" {
    description = "OIDC - application name for GitHub actions - needs to be same as in bootstrap"
    default = "gh-actions"
}
variable "repo_name" {
    description = "name of github repo"
    default = "lezelelyba/aoc2024"
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
locals {
    // github references used in principal creation
    github_subs = [
        for e in var.envs :
        "repo:${var.repo_name}:ref:refs/heads/${e.branch}"
    ]
    
    github_sub_wildcard = "repo:${var.repo_name}:ref:refs/heads/*"
}
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
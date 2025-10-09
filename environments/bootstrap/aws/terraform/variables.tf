variable "mngd" {
    default = "TF-Bootstrap"
}

resource "random_string" "bucket_number" {
    length = 6
    upper = false
    special = false
}

variable "tf_state_bucket" {
    default = "tf-state-bucket"
}

variable "tf_lock_db" {
    default = "tf-locks-db"
}

variable "state_filepath" {
    default = "/var/tmp/advent2024-bootstrap.tfstate"
}

variable "repo_name" {
    default = "lezelelyba/aoc2024"
}

variable "region" {
    default = "eu-central-1" 
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
variable "sshpubkeypath" {
    default = "~/.ssh/id_rsa.pub"
}

variable "sshprivkeypath" {
    default = "~/.ssh/id_rsa"
}

variable "ssm_path" {
    default = "/cd/prod/config"
}

variable "region" {
    default = "eu-central-1"
}

variable "docker_image" {
    default = "jsafar/advent2024.web:latest"
}

variable "env" {
    default = "prod"
}
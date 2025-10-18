variable "sshpubkeypath" {
    default = "~/.ssh/id_rsa.pub"
}

variable "sshprivkeypath" {
    default = "~/.ssh/id_rsa"
}

variable "ssm_path" {
    default = "/cd/dev/config"
}

variable "region" {
    default = "eu-central-1"
}

variable "docker_image" {
    default = "jsafar/advent2024.web:dev"
}

variable "env" {
    default = "dev"
}

variable "dns_zone" {
    default = "lezeleprojects.org."
}
variable "alb_dns_name" {
    default = "dev.aoc2024.lezeleprojects.org"
}
variable "ecs_app_env_map" {
    type = map(string)
    default = {}
}
variable "ecs_app_env_map_secret" {
    type = map(string)
    default = {}
    sensitive = true
}
variable "ecs_app_env_map_secret_keys" {
    type = list(string)
    default = []
}
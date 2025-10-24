variable "region" {
    default = "eu-central-1"
}
variable "env" {
    default = "prod"
}
variable "ssm_path" {
    default = "/cd/prod/config"
}
variable "dns_zone" {
    default = "lezeleprojects.org."
}
variable "alb_dns_name" {
    default = "aoc2024.lezeleprojects.org"
}

variable "docker_image" {
    default = "jsafar/advent2024.web:latest"
}
variable "sshpubkeypath" {
    default = "~/.ssh/id_rsa.pub"
}

variable "sshprivkeypath" {
    default = "~/.ssh/id_rsa"
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
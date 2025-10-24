variable "region" {
    default = "eu-central-1"
}
variable "env" {
    default = "dev"
}
variable "ssm_path" {
    default = "/cd/dev/config"
}
variable "dns_zone" {
    default = "lezeleprojects.org."
}
variable "alb_dns_name" {
    default = "dev.aoc2024.lezeleprojects.org"
}

variable "docker_image" {
    default = "jsafar/advent2024.web:dev"
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
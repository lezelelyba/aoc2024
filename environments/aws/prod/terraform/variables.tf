variable "region" {
    description = "AWS region"
    default = "eu-central-1"
}
variable "env" {
    description = "deployed environment"
    default = "prod"
}
variable "docker_image" {
    description = "Image for the ECS deployment"
    default = "jsafar/advent2024.web:latest"
}
variable "alb_dns_name" {
    description = "FQDN where the app will be deployed"
}
variable "dns_zone" {
    description = "DNS zone for the FQDN"
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
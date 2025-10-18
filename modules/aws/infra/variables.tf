variable "vpc_cidr" {
    default = "10.192.0.0/20"
}

variable "public_cidr" {
    default = "10.192.0.0/24"
}

variable "private_cidr" {
    default = "10.192.8.0/24"
}

variable "app_tcp_port" {
    default = "8080"
}
variable "alb_http_port" {
    default = "80"
}
variable "alb_http" {
    type = bool
    default = true 
}
variable "alb_https" {
    type = bool
    default = false
}

variable "alb_https_port" {
    default = "443"
}

variable "docker_image" {
    default = "jsafar/advent2024.web:latest"
}

variable "health_check_path" {
    default = "/healthcheck"
}

variable "sshpubkeypath" {
    default = "~/.ssh/id_rsa.pub"
}
variable "sshprivkeypath" {
    default = "~/.ssh/id_rsa"
}

variable "region" {
    default = "eu-central-1" 
}

variable "env" {
    default = "prod"
}

variable "ssm_path" {
    default = "/cd/prod/config"
}

variable "bastion" {
    default = true
}

variable "private_host" {
    default = true
}

variable "dns_zone" {
    default = "example.com."
}
variable "alb_dns_name" {
    default = "app.example.com"
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
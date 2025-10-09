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

variable "lb_tcp_port" {
    default = "80"
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
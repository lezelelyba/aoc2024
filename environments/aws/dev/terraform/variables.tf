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
    default = "jsafar/advent2024.web:dev"
}

variable "health_check_path" {
    default = "/healthcheck"
}

variable "sshpubkeypath" {
    default = "~/id_rsa.pub"
}
variable "sshprivkeypath" {
    default = "~/id_rsa"
}
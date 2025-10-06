variable "vpc_cidr" {
    default = "10.192.0.0/20"
}

variable "public_cidr" {
    default = "10.192.1.0/24"
}

variable "private_cidr" {
    default = "10.192.129.0/24"
}

variable "app_tcp_port" {
    default = "8080"
}
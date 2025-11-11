variable "region" {
    description = "AWS region"
}

variable "env" {
    description = "environment name"
}
variable "vpc_cidr" {
    description = "range for the whole VPC"
    default = "10.192.0.0/20"
}

variable "public_cidr" {
    description = "range for public subnets"
    default = "10.192.0.0/24"
}

variable "private_cidr" {
    description = "range for private subnets"
    default = "10.192.8.0/24"
}

variable "app_tcp_port" {
    description = "TCP port for the application"
    default = "8080"
}
variable "alb_http_port" {
    description = "HTTP port for the ALB frontend"
    default = "80"
}
variable "alb_https_port" {
    description = "HTTPS port for the ALB frontend"
    default = "443"
}
variable "alb_http" {
    description = "enable HTTP on the ALB"
    type = bool
    default = true 
}
variable "alb_https" {
    description = "enable HTTPS on the ALB"
    type = bool
    default = false
}
locals {
    redirect = var.alb_http && var.alb_https
}
variable "alb_dns_name" {
    description = "alb FQDN, for certification generation and CNAME registration"
}
variable "dns_zone" {
    description = "dnz zone for CNAME registration"
}

variable "docker_image" {
   description = "docker image for the ecs containers" 
}

variable "health_check_path" {
    description = "healthcheck url"
    default = "/healthcheck"
}
variable "ecs_app_env_map" {
    description = "ENV variables for the ECS container"
    type = map(string)
    default = {}
}

variable "ecs_app_env_map_secret" {
    description = "ENV secrets for the ECS container - will be stored as secrets"
    type = map(string)
    default = {}
    sensitive = true
}

// key map for the secrets - nonsensitive so it can be referenced in the code
locals {
    ecs_app_env_map_secret_keys = nonsensitive(keys(var.ecs_app_env_map_secret))
}

// key for information about ECS for gh action to pull
locals {
    ssm_path = "/cd/${var.env}/config"
}


variable "bastion" {
    description = "spin up bastion host"
    default = false
}

locals {
    keys_required = var.bastion || var.private_host
}
variable "sshpubkeypath" {
    description = "ssh public key for bastion host"
    default = ""

    validation {
        condition = !local.keys_required || ( local.keys_required && var.sshpubkeypath != "" )
        error_message = "bastion to be configured but missing public key" 
    }
}
variable "sshprivkeypath" {
    description = "ssh private key for bastion host"
    default = ""
    validation {
        condition = !local.keys_required || ( local.keys_required && var.sshprivkeypath != "" )
        error_message = "bastion to be configured but missing private key" 
    }
}
variable "private_host" {
    description = "spin up host in private subnet"
    default = false
}
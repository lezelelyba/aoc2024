variable "region" {
    default = "westeurope"
}

variable "env" {
    default = "dev"
}

variable "app_tcp_port" {
    default = 8080
}

variable "dns_zone" {
}

variable "app_dns_name" {
}

variable "docker_image" {
    default = "jsafar/advent2024.web:dev"
}

locals {
    app_port = contains(keys(var.aci_app_env_map), "PORT") ? var.aci_app_env_map["PORT"] : 8080
    cert = contains(keys(var.aci_app_env_map), "TLS_CERT_FILE") ? file(var.aci_app_env_map["TLS_CERT_FILE"]) : ""
    key = contains(keys(var.aci_app_env_map), "TLS_KEY_FILE") ? file(var.aci_app_env_map["TLS_KEY_FILE"]) : ""
    enable_https = contains(keys(var.aci_app_env_map), "ENABLE_HTTPS") ? var.aci_app_env_map["ENABLE_HTTPS"] : "false"
    use_acme = local.enable_https == "true" && ( local.cert == "" || local.key == "" )
}

resource "random_string" "dns_name_suffix" {
    length = 6
    upper = false
    special = false
}

variable "aci_app_env_map" {
    type = map(string)
    default = {}
}

variable "aci_app_env_map_secret" {
    type = map(string)
    default = {}
    sensitive = true
}

variable "tenant_id" {
    default = {}
}

variable "dns_provider" {
}

variable "email" {
}

locals {
    acr = jsondecode(file("../../acr.json"))
}

locals {
    aci_app_env_map_secret_keys = nonsensitive(keys(var.aci_app_env_map_secret))
}
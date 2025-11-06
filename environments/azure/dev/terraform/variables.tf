variable "region" {
    default = "westeurope"
}

variable "env" {
    default = "dev"
}

variable "app_tcp_port" {
    default = 8080
}

variable "docker_image" {
    default = "jsafar/advent2024.web:dev"
}

locals {
    app_port = contains(keys(var.aci_app_env_map), "PORT") ? var.aci_app_env_map["PORT"] : 8080
    cert = contains(keys(var.aci_app_env_map), "TLS_CERT_FILE") ? file(var.aci_app_env_map["TLS_CERT_FILE"]) : ""
    key = contains(keys(var.aci_app_env_map), "TLS_KEY_FILE") ? file(var.aci_app_env_map["TLS_KEY_FILE"]) : ""
    enable_https = contains(keys(var.aci_app_env_map), "ENABLE_HTTPS") ? var.aci_app_env_map["ENABLE_HTTPS"] : "false"
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

locals {
    acr = jsondecode(file("../../acr.json"))
}

locals {
    aci_app_env_map_secret_keys = nonsensitive(keys(var.aci_app_env_map_secret))
}
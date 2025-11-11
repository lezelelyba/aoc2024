variable "dns_zone" {
    description = "dns zone ending with ."
}

variable "domain" {
    description = "domain to be registered"
}

variable "alias" {
    description = "alias, ex. azure aci name, to be registered with the domain"
}
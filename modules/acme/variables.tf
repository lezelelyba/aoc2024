variable "subject" {
    description = "fqdn of the certificate"
}

variable "dns_provider" {
    description = "provider owning the dns domain"
}

variable "email" {
    description = "email to be registered with the let's encrypt server"
}

variable "algo" {
    default = "RSA"
    description = "algorithm to be used for the key pair"
}
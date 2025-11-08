terraform {
  required_providers {
    acme = {
      source  = "vancluever/acme"
      version = ">= 2.7.0"
    }
  }
}

provider "acme" {
  server_url = "https://acme-v02.api.letsencrypt.org/directory"
}
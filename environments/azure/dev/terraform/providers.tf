terraform {
  backend "azurerm" {
    key = "dev/terraform.tfstate"
    use_azuread_auth = true
  }

  required_providers {
    aws = {
        source = "hashicorp/aws"
        version = "~> 5.92"
    }
    acme = {
      source  = "vancluever/acme"
      version = ">= 2.7.0"
    }
    azurerm = {
        source = "hashicorp/azurerm"
        version = "~> 4.0"
    }

    azuread = {
        source = "hashicorp/azuread"
        version = "~> 3.0"
    }
  }
}

provider "acme" {
  server_url = "https://acme-v02.api.letsencrypt.org/directory"
}
provider "azurerm" {
    features {}
}
provider "azuread" {
}
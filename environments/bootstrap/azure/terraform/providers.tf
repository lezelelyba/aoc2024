terraform {
    required_providers {
        azurerm = {
            source = "hashicorp/azurerm"
            version = "~> 4.0"
        }

        azuread = {
            source = "hashicorp/azuread"
            version = "~> 3.0"
        }
    }
    // local backend for bootstrap
    // TODO: migration to azure?
    backend "local" {
        path = "/var/tmp/advent2024-azure-bootstrap.tfstate"
    }
}

provider "azurerm" {
    features {}
}

provider "azuread" {
}
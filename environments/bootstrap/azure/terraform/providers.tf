terraform {
    required_providers {
        azurerm = {
            source = "hashicorp/azurerm"
            version = "~> 4.0"
        }
    }
    backend "local" {
        path = "/var/tmp/advent2024-azure-bootstrap.tfstate"
    }
}

provider "azurerm" {
    features {}
}
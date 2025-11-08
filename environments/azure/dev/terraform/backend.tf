terraform {
  backend "azurerm" {
    key = "dev/terraform.tfstate"
    use_azuread_auth = true
  }
}
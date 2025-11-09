resource "azurerm_resource_group" "bootstrap" {
    name = var.resource-group-name
    location = var.region
    
    tags = {
        managed-by = var.mngd
    }
}

resource "azurerm_storage_account" "bootstrap" {
    name = "bootstrap${random_string.storage_account_number.result}" 
    resource_group_name = azurerm_resource_group.bootstrap.name
    location = var.region
    account_tier = "Standard"
    account_replication_type = "LRS"

    https_traffic_only_enabled = "true"
    public_network_access_enabled = "true"
    default_to_oauth_authentication = "true"

    identity { 
        type = "SystemAssigned"
    }

    tags = {
        managed-by = var.mngd
    }
}

resource "azurerm_storage_container" "tfstates" {
    name = "tfstates"
    storage_account_id = azurerm_storage_account.bootstrap.id
    container_access_type = "private"
}

resource "azurerm_container_registry" "acr" {
    name = "aocacr${random_string.storage_account_number.result}"
    resource_group_name = azurerm_resource_group.bootstrap.name
    location = var.region
    sku = "Basic"
    admin_enabled = "true"
}

resource "azuread_application" "gh_actions" {
    display_name = "gh-actions"    
}

resource "azuread_service_principal" "gh_actions" {
    client_id = azuread_application.gh_actions.client_id
}

resource "azuread_application_federated_identity_credential" "gh_oidc" {
    for_each = var.envs
    
    application_id = azuread_application.gh_actions.id
    display_name = "gh-actions-federation-${each.value.branch}"
    issuer = "https://token.actions.githubusercontent.com"
    subject = "repo:${var.repo_name}:ref:refs/heads/${each.value.branch}"
    audiences = ["api://AzureADTokenExchange"]
}

data "azurerm_subscription" "current" {}
resource "azurerm_role_assignment" "gh_actions_role" {
    scope = data.azurerm_subscription.current.id
    role_definition_name = "Contributor"
    principal_id = azuread_service_principal.gh_actions.object_id
}

resource "azurerm_app_configuration" "config_store" {
    name = "aoc-config-store"
    resource_group_name = azurerm_resource_group.bootstrap.name
    location = azurerm_resource_group.bootstrap.location
    sku = "free"
    public_network_access = "Enabled"
}

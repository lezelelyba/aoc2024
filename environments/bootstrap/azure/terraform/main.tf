// resource group under which the resources will be created 
resource "azurerm_resource_group" "bootstrap" {
    name = var.resource-group-name
    location = var.region
    
    tags = {
        managed-by = var.mngd
    }
}

// storage for the terraform states for each environment
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

// container for tfstates
resource "azurerm_storage_container" "tfstates" {
    name = "tfstates"
    storage_account_id = azurerm_storage_account.bootstrap.id
    container_access_type = "private"
}

// container registry for the aci as docker is rate limited
resource "azurerm_container_registry" "acr" {
    name = "aocacr${random_string.storage_account_number.result}"
    resource_group_name = azurerm_resource_group.bootstrap.name
    location = var.region
    sku = "Basic"
    admin_enabled = "true"
}

// OIDC
// same app name needs to be used in each environment
// TODO: required with contributor?
resource "azuread_application" "gh_actions" {
    display_name = var.gh_actions_application_name
}

resource "azuread_service_principal" "gh_actions" {
    client_id = azuread_application.gh_actions.client_id
}

// needs to be create for each branch as azure doesn't support wildcards
resource "azuread_application_federated_identity_credential" "gh_oidc" {
    for_each = var.envs
    
    application_id = azuread_application.gh_actions.id
    display_name = "gh-actions-federation-${each.value.branch}"
    issuer = "https://token.actions.githubusercontent.com"
    subject = "repo:${var.repo_name}:ref:refs/heads/${each.value.branch}"
    audiences = ["api://AzureADTokenExchange"]
}


// config store for gh actions to pull ACI id for redeploy for each environment
resource "azurerm_app_configuration" "config_store" {
    name = "aoc-config-store"
    resource_group_name = azurerm_resource_group.bootstrap.name
    location = azurerm_resource_group.bootstrap.location
    sku = "free"
    public_network_access = "Enabled"
}

// principal for gh action to allow read from config store

// didn't work
// resource "azurerm_role_assignment" "gh_actions_role" {
//     scope = azurerm_app_configuration.config_store.id
//     role_definition_name = "App Configuration Data Reader"
//     principal_id = azuread_service_principal.gh_actions.object_id
// }

// didn't work
// resource "azurerm_role_assignment" "gh_actions_role" {
//     scope = azurerm_app_configuration.config_store.id
//     role_definition_name = "Reader"
//     principal_id = azuread_service_principal.gh_actions.object_id
// }

// didn't work
// resource "azurerm_role_assignment" "gh_actions_role" {
//     scope = data.azurerm_subscription.current.id
//     role_definition_name = "App Configuration Data Reader"
//     principal_id = azuread_service_principal.gh_actions.object_id
// }

// didn't work
// resource "azurerm_role_assignment" "gh_actions_role" {
//     scope = data.azurerm_subscription.current.id
//     role_definition_name = "Reader"
//     principal_id = azuread_service_principal.gh_actions.object_id
// }

// didn't work
// resource "azurerm_role_assignment" "gh_actions_role" {
//     scope = azurerm_resource_group.bootstrap.id
//     role_definition_name = "Reader"
//     principal_id = azuread_service_principal.gh_actions.object_id
// }


// works
resource "azurerm_role_assignment" "gh_actions_role" {
    scope = azurerm_app_configuration.config_store.id
    role_definition_name = "Contributor"
    principal_id = azuread_service_principal.gh_actions.object_id
}
// path to environment modules
locals {
  repo_root = "${path.module}/../../../.."
}

// creates file with backend information for environments to use
resource "local_file" "backend_info_json" {
    filename = "${local.repo_root}/environments/azure/backend.json"
    content = jsonencode({
        resource_group_name   = azurerm_resource_group.bootstrap.name
        storage_account_name  = azurerm_storage_account.bootstrap.name
        container_name        = azurerm_storage_container.tfstates.name
    })
}

// creates file with container registry information for environments to use
resource "local_file" "acr_json" {
    filename = "${local.repo_root}/environments/azure/acr.json"
    content = jsonencode({
        login_server = azurerm_container_registry.acr.login_server
        admin_username = azurerm_container_registry.acr.admin_username
        admin_password = azurerm_container_registry.acr.admin_password
    })
}

// for image push to ACR
output "gh_secret_acr_username" {
    description = "ACR_USERNAME - setup as secret in GitHub"
    value = azurerm_container_registry.acr.admin_username
}
output "gh_secret_acr_password" {
    description = "ACR_PASSWORD - setup as secret in GitHub"
    value = nonsensitive(azurerm_container_registry.acr.admin_password)
    sensitive = false
}

// for ACI redeploy
output "gh_secret_azure_client_id" {
    description = "AZURE_CLIENT_ID - setup as secret in GitHub"
    value = azuread_application.gh_actions.client_id
}

data "azurerm_client_config" "current" {}
output "gh_secret_azure_tenant_id" {
    description = "AZURE_TENANT_ID - setup as secret in GitHub"
    value = data.azurerm_client_config.current.tenant_id
}
output "gh_secret_azure_subscription_id" {
    description = "AZURE_SUBSCRIPTION_ID - setup as secret in GitHub"
    value = data.azurerm_client_config.current.subscription_id
}

// creates file with config store information for environments to use
resource "local_file" "config_store_json" {
    filename = "${local.repo_root}/environments/azure/config_store.json"
    content = jsonencode({
        id = azurerm_app_configuration.config_store.id
    })
}

// creates gh action for the ACI redeploy
resource "local_file" "gh_actions_aci_change" {
    filename = "${local.repo_root}/.github/workflows/aci-image-change.yml"
    content = templatefile("${path.module}/templates/aci-image-change.yml.tmpl", {
        included_branches = join(" || ", [for env in var.envs : "github.ref == 'refs/heads/${env.branch}'"])
        envs = var.envs
        config_store_name = azurerm_app_configuration.config_store.name
    })
}
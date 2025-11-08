output "resource_group_name" {
    value = azurerm_resource_group.bootstrap.name
}

output "resource_storage_account" {
    value = azurerm_storage_account.bootstrap.name
}

output "resource_storage_container" {
    value = azurerm_storage_container.tfstates.name
}

locals {
  repo_root = "${path.module}/../../../.."
}

resource "local_file" "backend_info_json" {
    filename = "${local.repo_root}/environments/azure/backend.json"
    content = jsonencode({
        resource_group_name   = azurerm_resource_group.bootstrap.name
        storage_account_name  = azurerm_storage_account.bootstrap.name
        container_name        = azurerm_storage_container.tfstates.name
    })
}

resource "local_file" "acr_json" {
    filename = "${local.repo_root}/environments/azure/acr.json"
    content = jsonencode({
        login_server = azurerm_container_registry.acr.login_server
        admin_username = azurerm_container_registry.acr.admin_username
        admin_password = azurerm_container_registry.acr.admin_password
    })
}

output "gh_secret_azure_client_id" {
    value = azuread_application.gh_actions.client_id
}

data "azurerm_client_config" "current" {}
output "gh_secret_azure_tenant_id" {
    value = data.azurerm_client_config.current.tenant_id
}
output "gh_secret_azure_subscription_id" {
    value = data.azurerm_client_config.current.subscription_id
}
output "config_store_id" {
    value = azurerm_app_configuration.config_store.id
}
resource "local_file" "config_store_json" {
    filename = "${local.repo_root}/environments/azure/config_store.json"
    content = jsonencode({
        id = azurerm_app_configuration.config_store.id
    })
}
resource "local_file" "gh_actions_aci_change" {
    filename = "${local.repo_root}/.github/workflows/aci-image-change.yml"
    content = templatefile("${path.module}/templates/aci-image-change.yml.tmpl", {
        included_branches = join(" || ", [for env in var.envs : "github.ref == 'refs/heads/${env.branch}'"])
        envs = var.envs
        config_store_name = azurerm_app_configuration.config_store.name
    })
}
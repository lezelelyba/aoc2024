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

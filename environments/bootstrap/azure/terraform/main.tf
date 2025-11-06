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

    network_rules {
        default_action = "Allow"
    } 

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
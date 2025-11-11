output "public_url" {
  description = "Public URL of the ACI"
  value       = azurerm_container_group.app.fqdn
}
output "public_url_cname" {
  description = "Public URL of the ACI (CNAME)"
  value       = local.enable_https == "true" ? "https://${var.app_dns_name}" : "http://${var.app_dns_name}"
}

// save information about ACI for gh action to pull
resource "azurerm_app_configuration_key" "container_service_id" {
    configuration_store_id = local.config_store.id
    key = "/cd/${var.env}/aci_id"
    value = azurerm_container_group.app.id
}
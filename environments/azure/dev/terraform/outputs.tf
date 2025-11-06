output "public_url" {
  description = "Public URL of the ACI"
  value       = azurerm_container_group.app.fqdn
}
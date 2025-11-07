output "public_url" {
  description = "Public URL of the ACI"
  value       = azurerm_container_group.app.fqdn
}
output "public_url_cname" {
  description = "Public URL of the ACI (CNAME)"
  value       = local.enable_https == "true" ? "https://${var.app_dns_name}" : "http://${var.app_dns_name}"
}
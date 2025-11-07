module acme {
    source = "../../../../modules/acme"
    dns_provider = var.dns_provider
    email = var.email
    domain = var.app_dns_name
}


resource "azurerm_resource_group" "group" {
    name = "aoc-${var.env}-rg"
    location = var.region

    tags = {
        environment = var.env
    }
}

resource "azurerm_key_vault" "kv" {
    name = "aoc-${var.env}-kv"
    location = var.region
    resource_group_name = azurerm_resource_group.group.name
    sku_name = "standard"
    tenant_id = var.tenant_id

    rbac_authorization_enabled = true
}

resource "azurerm_key_vault_secret" "container_env_secret" {
    for_each = {for k in local.aci_app_env_map_secret_keys: replace(k, "_", "-") => var.aci_app_env_map_secret[k] }

    name = each.key
    value = each.value
    key_vault_id = azurerm_key_vault.kv.id
}

// import cert
// can't use cert as azure cert in key vault is always pfx
resource "azurerm_key_vault_secret" "cert" {
    count = !local.use_acme ? 1 : 0
    name = "imported-cert"
    key_vault_id = azurerm_key_vault.kv.id
    value = base64encode(local.cert)
}

resource "azurerm_key_vault_secret" "key" {
    count = !local.use_acme ? 1 : 0
    name = "imported-key"
    key_vault_id = azurerm_key_vault.kv.id
    value = base64encode(local.key)
}

// create cert via acme
resource "azurerm_key_vault_secret" "certacme" {
    count = local.use_acme ? 1 : 0
    name = "imported-cert-acme"
    key_vault_id = azurerm_key_vault.kv.id
    value = base64encode(module.acme.certificate)
}

resource "azurerm_key_vault_secret" "keyacme" {
    count = local.use_acme ? 1 : 0
    name = "imported-key-acme"
    key_vault_id = azurerm_key_vault.kv.id
    value = base64encode(module.acme.private_key)
}


resource "azurerm_container_group" "app" {
    name = "aoc-${var.env}-solver"
    location = var.region
    resource_group_name = azurerm_resource_group.group.name
    os_type = "Linux"

    image_registry_credential {
        server = local.acr.login_server
        username = local.acr.admin_username
        password = local.acr.admin_password
    }

    identity {
        type = "SystemAssigned"        
    }

    container {
       name = "aoc-solver"
       image = "${local.acr.login_server}/${var.docker_image}"
       cpu = "0.5"
       memory = "1.0"

        ports {
            port = tonumber(local.app_port)
            protocol = "TCP"
        }

        // env variables
        environment_variables = merge (
            var.aci_app_env_map,
            // if https is enabled, pass the corresponding env variables
            local.enable_https == "true" ? {
                "ENABLE_HTTPS" = "true",
                "TLS_CERT_FILE" = "/files/cert.pem",
                "TLS_KEY_FILE" = "/files/key.pem"
            } : {}
        )

        // secure env variables
        secure_environment_variables = { for k in local.aci_app_env_map_secret_keys: k => azurerm_key_vault_secret.container_env_secret[replace(k, "_", "-")].value }

        // create volume only if cert and key has to be passed to container
        dynamic "volume" {
            for_each = local.enable_https == "true" ? ["1"] : []
            content {
                name = "files"

                mount_path = "/files"
                read_only = true

                secret = local.use_acme ? {
                    "key.pem" = azurerm_key_vault_secret.keyacme[0].value
                    "cert.pem" = azurerm_key_vault_secret.certacme[0].value
                } : {
                    "key.pem" = azurerm_key_vault_secret.key[0].value
                    "cert.pem" = azurerm_key_vault_secret.cert[0].value
                }
            }
        }
    }

    ip_address_type = "Public"

    dns_name_label = "aoc${var.env}${random_string.dns_name_suffix.result}"

    tags = {
        environment = var.env
    }
}
module route53registration {
    source = "../../../../modules/aws/route53"
    alias = azurerm_container_group.app.fqdn
    domain = var.app_dns_name
    dns_zone = var.dns_zone
}

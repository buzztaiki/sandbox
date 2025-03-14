locals {
  tags = {
    name       = var.name
    managed_by = "Terraform"
  }
}

data "azuread_client_config" "current" {}
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "this" {
  name     = var.name
  location = "japaneast"
  tags     = local.tags
}

resource "random_password" "this" {
  length      = 16
  min_lower   = 1
  min_upper   = 1
  min_numeric = 1
  min_special = 1
}

resource "azurerm_mysql_flexible_server" "this" {
  name                = var.name
  location            = azurerm_resource_group.this.location
  resource_group_name = azurerm_resource_group.this.name
  sku_name            = var.mysql_sku_name

  administrator_login    = var.mysql_admin_user
  administrator_password = random_password.this.result
  backup_retention_days  = 1
  tags                   = local.tags

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.this.id]
  }

  lifecycle {
    ignore_changes = [
      administrator_password,
      zone,
    ]
  }
}

resource "azapi_resource_action" "mysql_set_public_access" {
  type        = "Microsoft.DBforMySQL/flexibleServers@2023-06-30"
  resource_id = azurerm_mysql_flexible_server.this.id
  method      = "PATCH"

  body = {
    properties = {
      network = {
        publicNetworkAccess = "Enabled"
      }
    }
  }
}

resource "azurerm_mysql_flexible_server_firewall_rule" "this" {
  name                = "allow-all-public-access"
  resource_group_name = azurerm_resource_group.this.name
  server_name         = azurerm_mysql_flexible_server.this.name
  start_ip_address    = "0.0.0.0"
  end_ip_address      = "255.255.255.255"
}

resource "azurerm_mysql_flexible_server_configuration" "this" {
  for_each = {
    "log_output"            = "FILE",
    "audit_log_enabled"     = "ON",
    "audit_log_events"      = var.mysql_audit_log_events,
    "performance_schema"    = "ON",
    "slow_query_log"        = "ON",
    "error_server_log_file" = "ON",
  }

  resource_group_name = azurerm_resource_group.this.name
  server_name         = azurerm_mysql_flexible_server.this.name
  name                = each.key
  value               = each.value
}


resource "azurerm_user_assigned_identity" "this" {
  name                = var.name
  resource_group_name = azurerm_resource_group.this.name
  location            = azurerm_resource_group.this.location
  tags                = local.tags
}

data "azuread_application_published_app_ids" "well_known" {}

data "azuread_service_principal" "msgraph" {
  client_id = data.azuread_application_published_app_ids.well_known.result.MicrosoftGraph
}

resource "azuread_app_role_assignment" "this" {
  for_each = toset([
    "User.Read.All",
    "GroupMember.Read.All",
    "Application.Read.All",
  ])
  app_role_id         = data.azuread_service_principal.msgraph.app_role_ids[each.value]
  principal_object_id = azurerm_user_assigned_identity.this.principal_id
  resource_object_id  = data.azuread_service_principal.msgraph.object_id
}

resource "azuread_group" "mysql_admin" {
  display_name     = "${var.name}-admin"
  security_enabled = true
  members = [
    data.azuread_client_config.current.object_id,
  ]
}

resource "azuread_group" "mysql_maintainer" {
  display_name     = "${var.name}-maintainer"
  security_enabled = true
  members = [
    data.azuread_client_config.current.object_id,
  ]
}

resource "azuread_group" "mysql_reader" {
  display_name     = "${var.name}-reader"
  security_enabled = true
  members = [
    data.azuread_client_config.current.object_id,
  ]
}

resource "azurerm_mysql_flexible_server_active_directory_administrator" "this" {
  server_id   = azurerm_mysql_flexible_server.this.id
  identity_id = azurerm_user_assigned_identity.this.id
  login       = azuread_group.mysql_admin.display_name
  object_id   = azuread_group.mysql_admin.object_id
  tenant_id   = data.azurerm_client_config.current.tenant_id
}

resource "azurerm_log_analytics_workspace" "this" {
  name                = var.name
  resource_group_name = azurerm_resource_group.this.name
  location            = azurerm_resource_group.this.location
  tags                = local.tags
}

resource "azurerm_monitor_diagnostic_setting" "mysql-all" {
  name                       = "all"
  target_resource_id         = azurerm_mysql_flexible_server.this.id
  log_analytics_workspace_id = azurerm_log_analytics_workspace.this.id

  enabled_log {
    category_group = "audit"
  }

  metric {
    category = "AllMetrics"
    enabled  = false
  }
}

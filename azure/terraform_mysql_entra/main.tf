locals {
  name = "mysql-entra"
  tags = {
    name       = local.name
    managed_by = "Terraform"
  }
}

data "azuread_client_config" "current" {}
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "this" {
  name     = local.name
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
  name                = local.name
  location            = azurerm_resource_group.this.location
  resource_group_name = azurerm_resource_group.this.name
  sku_name            = "B_Standard_B1s"

  administrator_login    = "mysql_fake_admin"
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

resource "azurerm_mysql_flexible_server_firewall_rule" "this" {
  name                = local.name
  resource_group_name = azurerm_resource_group.this.name
  server_name         = azurerm_mysql_flexible_server.this.name
  start_ip_address    = "0.0.0.0"
  end_ip_address      = "255.255.255.255"
}

resource "azurerm_mysql_flexible_server_configuration" "this" {
  for_each = {
    "log_output"            = "FILE",
    "audit_log_enabled"     = "ON",
    "audit_log_events"      = "CONNECTION,DDL,DCL",
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
  name                = local.name
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
  display_name     = "${local.name}-admin"
  security_enabled = true
  members = [
    data.azuread_client_config.current.object_id,
  ]
}

resource "azuread_group" "mysql_maintainer" {
  display_name     = "${local.name}-maintainer"
  security_enabled = true
  members = [
    data.azuread_client_config.current.object_id,
  ]
}

resource "azuread_group" "mysql_reader" {
  display_name     = "${local.name}-reader"
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
  name                = local.name
  resource_group_name = azurerm_resource_group.this.name
  location            = azurerm_resource_group.this.location
  tags                = local.tags
}

data "azurerm_monitor_diagnostic_categories" "mysql" {
  resource_id = azurerm_mysql_flexible_server.this.id
}

resource "azurerm_monitor_diagnostic_setting" "mysql-all" {
  name                       = "all"
  target_resource_id         = azurerm_mysql_flexible_server.this.id
  log_analytics_workspace_id = azurerm_log_analytics_workspace.this.id

  dynamic "enabled_log" {
    for_each = data.azurerm_monitor_diagnostic_categories.mysql.log_category_types
    content {
      category = enabled_log.value
    }
  }

  metric {
    category = "AllMetrics"
    enabled  = false
  }
}

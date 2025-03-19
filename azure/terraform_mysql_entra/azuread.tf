data "azuread_application_published_app_ids" "well_known" {}

data "azuread_service_principal" "msgraph" {
  client_id = data.azuread_application_published_app_ids.well_known.result.MicrosoftGraph
}

resource "azuread_app_role_assignment" "mysql" {
  # var.mysql_identity_grant_directory_readers が false の場合だけ有効
  for_each = toset([for x in [
    "User.Read.All",
    "GroupMember.Read.All",
    "Application.Read.All",
  ] : x if !var.mysql_identity_grant_directory_readers])
  app_role_id         = data.azuread_service_principal.msgraph.app_role_ids[each.value]
  principal_object_id = azurerm_user_assigned_identity.mysql.principal_id
  resource_object_id  = data.azuread_service_principal.msgraph.object_id
}

resource "azuread_directory_role" "directory_readers" {
  display_name = "Directory Readers"
}

resource "azuread_directory_role_assignment" "mysql" {
  # var.mysql_identity_grant_directory_readers が true の場合だけ有効
  count               = var.mysql_identity_grant_directory_readers ? 1 : 0
  role_id             = azuread_directory_role.directory_readers.template_id
  principal_object_id = azurerm_user_assigned_identity.mysql.principal_id
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

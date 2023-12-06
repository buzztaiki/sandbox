data "azurerm_client_config" "current" {}

data "azuread_domains" "default" {
  only_default = true
}

locals {
  name           = "terraform-azuread-example"
  default_domain = data.azuread_domains.default.domains[0].domain_name
}

resource "random_password" "this" {
  length = 16
}

resource "azuread_application" "this" {
  display_name = local.name
  owners       = [data.azurerm_client_config.current.object_id]
}

resource "azuread_service_principal" "this" {
  client_id = azuread_application.this.client_id
  owners    = [data.azurerm_client_config.current.object_id]
}

resource "azuread_application_federated_identity_credential" "this" {
  application_id = azuread_application.this.id
  display_name   = azuread_application.this.display_name
  audiences      = ["api://AzureADTokenExchange"]
  issuer         = "https://token.actions.githubusercontent.com"
  subject        = "repo:buzztaiki/sandbox:pull_request"
}

resource "azuread_user" "this" {
  user_principal_name   = "${local.name}@${local.default_domain}"
  display_name          = local.name
  password              = random_password.this.result
  force_password_change = true

  lifecycle {
    ignore_changes = [
      password,
      force_password_change,
    ]
  }
}

resource "azuread_group" "this" {
  display_name     = local.name
  security_enabled = true
}

resource "azuread_group_member" "this" {
  group_object_id  = azuread_group.this.object_id
  member_object_id = azuread_user.this.object_id
}

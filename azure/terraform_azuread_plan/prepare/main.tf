data "azurerm_client_config" "current" {}

locals {
  subscription_scope = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
}

resource "time_rotating" "month" {
  rotation_days = 30
}

resource "azuread_application" "plan" {
  display_name = "terraform-plan"
  owners       = [data.azurerm_client_config.current.object_id]
}

resource "azuread_service_principal" "plan" {
  client_id = azuread_application.plan.client_id
  owners    = [data.azurerm_client_config.current.object_id]
}

resource "azuread_service_principal_password" "plan" {
  service_principal_id = azuread_service_principal.plan.id
  rotate_when_changed = {
    rotation = time_rotating.month.id
  }
}

resource "azurerm_role_assignment" "plan_is_subscription_contributor" {
  scope                = local.subscription_scope
  role_definition_name = "Contributor"
  principal_id         = azuread_service_principal.plan.id
}

resource "local_sensitive_file" "plan_tfvars" {
  content         = <<-EOT
    tenant_id       = "${data.azurerm_client_config.current.tenant_id}"
    subscription_id = "${data.azurerm_client_config.current.subscription_id}"
    client_id       = "${azuread_service_principal.plan.client_id}"
    client_secret   = "${azuread_service_principal_password.plan.value}"
  EOT
  filename        = "${path.module}/../main/plan.tfvars"
  file_permission = "0600"
}

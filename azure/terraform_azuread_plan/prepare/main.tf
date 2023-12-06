data "azurerm_client_config" "current" {}

locals {
  subscription_scope = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
}

data "azuread_application_published_app_ids" "well_known" {}

resource "azuread_service_principal" "msgraph" {
  # "MicrosoftGraph" = "00000003-0000-0000-c000-000000000000"
  client_id    = data.azuread_application_published_app_ids.well_known.result.MicrosoftGraph
  use_existing = true
}

resource "time_rotating" "month" {
  rotation_days = 30
}

resource "azuread_application" "plan" {
  display_name = "terraform-plan"
  owners       = [data.azurerm_client_config.current.object_id]

  # アプリが利用できるロール
  required_resource_access {
    resource_app_id = data.azuread_application_published_app_ids.well_known.result.MicrosoftGraph
    resource_access {
      # https://learn.microsoft.com/en-us/graph/permissions-reference#directoryreadall
      id   = azuread_service_principal.msgraph.app_role_ids["Directory.Read.All"]
      type = "Role"
    }
  }
}

resource "azuread_service_principal" "plan" {
  client_id = azuread_application.plan.client_id
  owners    = [data.azurerm_client_config.current.object_id]
}

# アプリが利用できるロールのうち、SPが要求するロール
resource "azuread_app_role_assignment" "plan_has_directory_read_all" {
  resource_object_id  = azuread_service_principal.msgraph.object_id
  app_role_id         = azuread_service_principal.msgraph.app_role_ids["Directory.Read.All"]
  principal_object_id = azuread_service_principal.plan.object_id
}

resource "azuread_service_principal_password" "plan" {
  service_principal_id = azuread_service_principal.plan.id
  rotate_when_changed = {
    rotation = time_rotating.month.id
  }
}

resource "azurerm_role_assignment" "plan_has_subscription_contributor" {
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

# output "well_known_apps" {
#   value = keys(data.azuread_application_published_app_ids.well_known.result)
# }
# output "msgraph_app_role_ids" {
#   value = keys(azuread_service_principal.msgraph.app_role_ids)
# }
# output "msgraph_oauth2_permission_scope_ids" {
#   value = keys(azuread_service_principal.msgraph.oauth2_permission_scope_ids)
# }

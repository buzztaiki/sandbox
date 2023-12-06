data "azurerm_client_config" "current" {}

data "azuread_application_published_app_ids" "well_known" {}

data "azuread_service_principal" "msgraph" {
  # "MicrosoftGraph" = "00000003-0000-0000-c000-000000000000"
  client_id = data.azuread_application_published_app_ids.well_known.result.MicrosoftGraph
}

resource "time_rotating" "month" {
  rotation_days = 30
}

resource "azuread_application" "this" {
  display_name = var.name
  owners       = var.owners

  # アプリが利用できるロール
  required_resource_access {
    resource_app_id = data.azuread_application_published_app_ids.well_known.result.MicrosoftGraph
    dynamic "resource_access" {
      for_each = var.msgraph_roles
      content {
        id   = data.azuread_service_principal.msgraph.app_role_ids[resource_access.value]
        type = "Role"
      }
    }
  }
}

resource "azuread_service_principal" "this" {
  client_id = azuread_application.this.client_id
  owners    = var.owners
}

# アプリが利用できるロールのうちSPが要求するロール
resource "azuread_app_role_assignment" "this" {
  for_each            = var.msgraph_roles
  resource_object_id  = data.azuread_service_principal.msgraph.object_id
  app_role_id         = data.azuread_service_principal.msgraph.app_role_ids[each.value]
  principal_object_id = azuread_service_principal.this.object_id
}

resource "azuread_application_password" "this" {
  display_name   = var.name
  application_id = azuread_application.this.id
  rotate_when_changed = {
    rotation = time_rotating.month.id
  }
}

# こっちを使うと Azure Portal から見れない。多分 multi tenant の時はこれを使う事になる。
# resource "azuread_service_principal_password" "this" {
#   service_principal_id = azuread_service_principal.this.object_id
#   rotate_when_changed = {
#     rotation = time_rotating.month.id
#   }
# }


resource "azurerm_role_assignment" "this" {
  # for_each = { for x in var.msgraph_roles : (x.scope) => x }
  for_each             = { for x in var.azurerm_roles : "${x.scope}:${x.role}" => x }
  scope                = each.value.scope
  role_definition_name = each.value.role
  principal_id         = azuread_service_principal.this.id
}

resource "local_sensitive_file" "tfvars" {
  content         = <<-EOT
    tenant_id       = "${data.azurerm_client_config.current.tenant_id}"
    subscription_id = "${data.azurerm_client_config.current.subscription_id}"
    client_id       = "${azuread_application.this.client_id}"
    client_secret   = "${azuread_application_password.this.value}"
  EOT
  filename        = var.tfvars
  file_permission = "0600"
}

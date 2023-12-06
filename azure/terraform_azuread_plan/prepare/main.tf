data "azurerm_client_config" "current" {}

locals {
  subscription_resid = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
}

module "plan_app" {
  source = "./azuread_app"
  name   = "terraform-azuread-plan"
  tfvars = "${path.module}/../main/plan.tfvars"
  owners = [data.azurerm_client_config.current.object_id]
  msgraph_roles = [
    # https://learn.microsoft.com/en-us/graph/permissions-reference#directoryreadall
    "Directory.Read.All",
  ]
  azurerm_roles = [
    { scope = local.subscription_resid, role = "Contributor" },
  ]
}

module "apply_app" {
  source = "./azuread_app"
  name   = "terraform-azuread-apply"
  tfvars = "${path.module}/../main/apply.tfvars"
  owners = [data.azurerm_client_config.current.object_id]
  msgraph_roles = [
    # https://learn.microsoft.com/en-us/graph/permissions-reference#directoryreadall
    "Directory.Read.All",
    # https://learn.microsoft.com/en-us/graph/permissions-reference#applicationreadwriteall
    "Application.ReadWrite.All",
    # https://learn.microsoft.com/en-us/graph/permissions-reference#userreadwriteall
    "User.ReadWrite.All",
    # https://learn.microsoft.com/en-us/graph/permissions-reference#groupreadwriteall
    "Group.ReadWrite.All",
  ]
  azurerm_roles = [
    { scope = local.subscription_resid, role = "Contributor" }
  ]
}

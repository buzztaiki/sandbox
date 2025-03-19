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

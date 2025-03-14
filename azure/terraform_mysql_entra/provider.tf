terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 4.23"
    }
    azuread = {
      source  = "hashicorp/azuread"
      version = "~> 3.1"
    }
    azapi = {
      source  = "Azure/azapi"
      version = "~> 2.3"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.7"
    }
  }
}

provider "azurerm" {
  features {}
  subscription_id = var.azurerm_subscription_id
}

provider "azuread" {
  tenant_id = var.azuread_tenant_id
}

provider "azapi" {
}

provider "random" {
}

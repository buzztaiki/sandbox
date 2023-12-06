terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.83.0"
    }
    azuread = {
      source  = "hashicorp/azuread"
      version = "~> 2.46.0"
    }
    time = {
      source  = "hashicorp/time"
      version = "~> 0.9.2"
    }
    local = {
      source  = "hashicorp/local"
      version = "~> 2.4.0"
    }
  }
}

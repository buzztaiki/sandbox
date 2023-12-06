variable "name" {
  type = string
}

variable "owners" {
  type = set(string)
}

variable "msgraph_roles" {
  type = set(string)
}

variable "azurerm_roles" {
  type = list(object({
    scope = string
    role  = string
  }))
}

variable "tfvars" {
  type = string
}

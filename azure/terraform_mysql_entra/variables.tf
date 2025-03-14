variable "azurerm_subscription_id" {
  type = string
}
variable "azuread_tenant_id" {
  type = string
}

variable "name" {
  type    = string
  default = "mysql-entra"
}

variable "mysql_admin_user" {
  type    = string
  default = "mysql_fake_admin"
}

variable "mysql_sku_name" {
  type    = string
  default = "B_Standard_B1ms"
}

variable "mysql_audit_log_events" {
  type    = string
  default = "CONNECTION,DDL,DCL,DML"
}

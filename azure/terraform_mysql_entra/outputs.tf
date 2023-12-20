output "mysql_server_fqdn" {
  value = azurerm_mysql_flexible_server.this.fqdn
}

output "mysql_admin" {
  value = azuread_group.mysql_admin.display_name
}

output "mysql_reader" {
  value = azuread_group.mysql_reader.display_name
}

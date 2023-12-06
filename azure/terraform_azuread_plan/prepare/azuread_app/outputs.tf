output "client_id" {
  value = azuread_application.this.client_id
}

output "client_secret" {
  value     = azuread_application_password.this.value
  sensitive = true
}

output "service_principal_id" {
  value = azuread_service_principal.this.id
}


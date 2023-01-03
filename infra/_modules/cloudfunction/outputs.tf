output "function_name" {
  value = google_cloudfunctions2_function.this.name
}

output "service_account_email" {
  value = google_cloudfunctions2_function.this.service_config[0].service_account_email
}

output "endpoint" {
  value = google_cloudfunctions2_function.this.service_config[0].uri
}

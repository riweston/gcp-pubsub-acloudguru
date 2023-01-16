output "pubsub_topic" {
  value = try(google_pubsub_topic.this[0].name, "")
}

output "service_account_email" {
  value = google_service_account.this.email
}

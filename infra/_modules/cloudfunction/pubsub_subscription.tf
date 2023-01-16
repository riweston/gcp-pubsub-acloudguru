resource "google_pubsub_subscription" "this" {
  count = var.pubsub_subscription_topic == "" ? 0 : 1

  name  = google_cloudfunctions2_function.this.name
  topic = var.pubsub_subscription_topic

  push_config {
    push_endpoint = google_cloudfunctions2_function.this.service_config[0].uri
    oidc_token {
      service_account_email = google_service_account.this.email
    }
  }
}

resource "google_pubsub_subscription" "this" {
  name  = "lookup-slack-id-subscription"
  topic = var.topic

  push_config {
    push_endpoint = var.push_endpoint
    oidc_token {
      service_account_email = var.function_service_account_email
    }
  }
}

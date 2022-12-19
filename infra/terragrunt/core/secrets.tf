data "google_secret_manager_secret" "slack_signing_secret" {
  secret_id = "SLACK_SIGNING_SECRET"
}

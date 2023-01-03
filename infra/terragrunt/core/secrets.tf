data "google_secret_manager_secret" "slack_signing_secret" {
  secret_id = "SLACK_SIGNING_SECRET"
}

data "google_secret_manager_secret" "slack_bot_token" {
  secret_id = "SLACK_BOT_TOKEN"
}

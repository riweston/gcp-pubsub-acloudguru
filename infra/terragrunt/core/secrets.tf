data "google_secret_manager_secret" "slack_signing_secret" {
  secret_id = "SLACK_SIGNING_SECRET"
}

data "google_secret_manager_secret" "slack_bot_token" {
  secret_id = "SLACK_BOT_TOKEN"
}

data "google_secret_manager_secret" "acloudguru_api_key" {
  secret_id = "ACLOUDGURU_API_KEY"
}

data "google_secret_manager_secret" "acloudguru_consumer_id" {
  secret_id = "ACLOUDGURU_CONSUMER_ID"
}

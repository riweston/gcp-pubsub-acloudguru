resource "google_secret_manager_secret" "slack_signing_secret" {
  secret_id = "SLACK_SIGNING_SECRET"
  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "slack_signing_secret" {
  secret      = google_secret_manager_secret.slack_signing_secret.id
  secret_data = var.slack_signing_secret
}

resource "google_secret_manager_secret" "slack_bot_token" {
  secret_id = "SLACK_BOT_TOKEN"
  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "slack_bot_token" {
  secret      = google_secret_manager_secret.slack_bot_token.id
  secret_data = var.slack_bot_token
}

resource "google_secret_manager_secret" "acloudguru_api_key" {
  secret_id = "ACLOUDGURU_API_KEY"
  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "acloudguru_api_key" {
  secret      = google_secret_manager_secret.acloudguru_api_key.id
  secret_data = var.acloudguru_api_key
}

resource "google_secret_manager_secret" "acloudguru_consumer_id" {
  secret_id = "ACLOUDGURU_CONSUMER_ID"
  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "acloudguru_consumer_id" {
  secret      = google_secret_manager_secret.acloudguru_consumer_id.id
  secret_data = var.acloudguru_consumer_id
}

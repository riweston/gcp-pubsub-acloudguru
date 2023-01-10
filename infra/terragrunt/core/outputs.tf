output "bucket_name_cloud_function" {
  value = google_storage_bucket.cloud_functions.name
}

output "location" {
  value = var.location
}

output "project_id" {
  value = data.google_project.this.number
}

output "secret_id_slack_signing_secret" {
  value = data.google_secret_manager_secret.slack_signing_secret.secret_id
}

output "secret_id_slack_bot_token" {
  value = data.google_secret_manager_secret.slack_bot_token.secret_id
}

output "secret_id_acloudguru_api_key" {
  value = data.google_secret_manager_secret.acloudguru_api_key.secret_id
}

output "secret_id_acloudguru_consumer_id" {
  value = data.google_secret_manager_secret.acloudguru_consumer_id.secret_id
}

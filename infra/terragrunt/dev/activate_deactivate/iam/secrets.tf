resource "google_secret_manager_secret_iam_member" "this" {
  for_each = toset([
    var.secret_id_acloudguru_api_key,
    var.secret_id_acloudguru_consumer_id,
    var.secret_id_slack_bot_token,
  ])

  secret_id = each.value
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${var.function_service_account_email}"
}

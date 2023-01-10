resource "google_secret_manager_secret_iam_binding" "this" {
  for_each = toset([
    var.secret_id_acloudguru_api_key,
    var.secret_id_acloudguru_consumer_id,
  ])

  secret_id = each.value
  role      = "roles/secretmanager.secretAccessor"
  members = [
    "serviceAccount:${var.function_service_account_email}",
  ]
}
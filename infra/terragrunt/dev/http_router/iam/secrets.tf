resource "google_secret_manager_secret_iam_member" "this" {
  secret_id = var.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${var.function_service_account_email}"
}

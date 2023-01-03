resource "google_secret_manager_secret_iam_binding" "this" {
  secret_id = var.secret_id
  role      = "roles/secretmanager.secretAccessor"
  members   = [
    "serviceAccount:${var.function_service_account_email}",
  ]
}

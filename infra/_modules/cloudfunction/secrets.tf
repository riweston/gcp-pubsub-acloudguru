data "google_secret_manager_secret" "secret_environment_variables" {
  for_each = toset(var.secret_environment_variables)

  secret_id = each.value
}

resource "google_secret_manager_secret_iam_member" "secret_environment_variables" {
  for_each = data.google_secret_manager_secret.secret_environment_variables

  secret_id = each.value.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.this.email}"
}

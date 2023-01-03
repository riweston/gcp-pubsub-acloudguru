resource "google_cloud_run_service_iam_member" "this" {
  location = var.function_location
  service  = var.function_name
  role     = "roles/run.invoker"
  member   = "allUsers"
}

resource "google_cloudfunctions2_function_iam_member" "this" {
  location       = var.function_location
  cloud_function = var.function_name
  role           = "roles/cloudfunctions.invoker"
  member         = "allUsers"
}

resource "google_cloud_run_service_iam_member" "this" {
  location = var.function_location
  service  = var.function_name
  role     = "roles/run.invoker"
  member   = "serviceAccount:${var.function_service_account_email}"
}

resource "google_cloudfunctions2_function_iam_member" "this" {
  location       = var.function_location
  cloud_function = var.function_name
  role           = "roles/cloudfunctions.invoker"
  member         = "serviceAccount:${var.function_service_account_email}"
}

resource "google_pubsub_topic_iam_member" "this" {
  topic  = var.topic_name
  role   = "roles/pubsub.editor"
  member = "serviceAccount:${var.function_service_account_email}"
}

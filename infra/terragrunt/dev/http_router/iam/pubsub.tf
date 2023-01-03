resource "google_pubsub_topic_iam_member" "this" {
  topic  = var.topic_name
  role   = "roles/pubsub.editor"
  member = "serviceAccount:${var.function_service_account_email}"
}

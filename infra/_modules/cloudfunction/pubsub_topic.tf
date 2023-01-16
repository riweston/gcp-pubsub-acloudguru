resource "google_pubsub_topic" "this" {
  count = var.pubsub_topic == "" ? 0 : 1

  name = var.pubsub_topic
}

resource "google_pubsub_topic_iam_member" "this" {
  count = try(var.environment_variables["TOPIC_NAME"], null) == null ? 0 : 1

  topic  = var.environment_variables["TOPIC_NAME"]
  role   = "roles/pubsub.editor"
  member = "serviceAccount:${google_service_account.this.email}"
}

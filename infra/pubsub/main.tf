resource "google_pubsub_topic" "this" {
  for_each = var.topics

  name    = each.value
  project = var.project_id
}

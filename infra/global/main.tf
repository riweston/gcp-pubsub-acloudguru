resource "google_storage_bucket" "this" {
  name     = var.project_number
  project  = var.project_id
  location = var.location
}

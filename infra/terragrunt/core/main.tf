resource "google_storage_bucket" "state" {
  name     = data.google_project.this.number
  location = var.location
  versioning {
    enabled = true
  }
}

resource "google_storage_bucket" "cloud_functions" {
  name     = "${data.google_project.this.number}-cloud-functions"
  location = var.location
}

data "google_project" "this" {}

resource "google_project_service" "this" {
  for_each = toset(jsondecode(var.apis))

  service = each.value
}

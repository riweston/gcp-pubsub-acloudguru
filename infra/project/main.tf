locals {
  location = "europe-west2"
  apis = [
    "cloudfunctions.googleapis.com",
    "run.googleapis.com",
    "artifactregistry.googleapis.com",
    "cloudbuild.googleapis.com",
    "logging.googleapis.com",
    "secretmanager.googleapis.com",
  ]
}

resource "google_storage_bucket" "state" {
  name     = data.google_project.this.number
  location = local.location
  versioning {
    enabled = true
  }
}

data "google_project" "this" {}

resource "google_storage_bucket" "cloud_functions" {
  name     = "${data.google_project.this.number}-cloud-functions"
  location = local.location
}

resource "google_project_service" "this" {
  for_each = toset(local.apis)

  service = each.value
}

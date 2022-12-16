locals {
  apis = [
    "cloudfunctions.googleapis.com",
    "run.googleapis.com",
    "artifactregistry.googleapis.com",
    "cloudbuild.googleapis.com",
    "logging.googleapis.com",
    "secretmanager.googleapis.com",
  ]
}

resource "google_project_service" "this" {
  for_each = toset(local.apis)

  service                    = each.value
  disable_dependent_services = true
}

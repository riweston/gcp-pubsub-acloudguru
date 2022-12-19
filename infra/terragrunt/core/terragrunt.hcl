include {
  path = find_in_parent_folders()
}

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

inputs = {
  location = local.location
  apis     = local.apis
}

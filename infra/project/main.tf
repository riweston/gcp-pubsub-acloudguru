locals {
  location = "europe-west2"
  apis     = [
    "cloudfunctions.googleapis.com",
    "run.googleapis.com",
    "artifactregistry.googleapis.com",
    "cloudbuild.googleapis.com",
    "logging.googleapis.com",
    "secretmanager.googleapis.com",
    "iamcredentials.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "iam.googleapis.com",
  ]
  repo = "cloudreach/gcp-pubsub-acloudguru"
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

resource "google_iam_workload_identity_pool" "github_pool" {
  project                   = data.google_project.this.name
  workload_identity_pool_id = "github-pool"
  display_name              = "GitHub pool"
  description               = "Identity pool for GitHub deployments"
}

resource "google_iam_workload_identity_pool_provider" "github" {
  project                            = data.google_project.this.name
  workload_identity_pool_id          = google_iam_workload_identity_pool.github_pool.workload_identity_pool_id
  workload_identity_pool_provider_id = "github-provider"
  attribute_mapping                  = {
    "google.subject"       = "assertion.sub"
    "attribute.actor"      = "assertion.actor"
    "attribute.aud"        = "assertion.aud"
    "attribute.repository" = "assertion.repository"
  }
  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }
}

resource "google_service_account" "github_actions" {
  project      = data.google_project.this.name
  account_id   = "github-actions"
  display_name = "Service Account used for GitHub Actions"
}

resource "google_service_account_iam_member" "workload_identity_user" {
  service_account_id = google_service_account.github_actions.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.github_pool.name}/attribute.repository/${local.repo}"
}

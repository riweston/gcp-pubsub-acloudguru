provider "google" {
  project = "cr-lab-rweston-2811223416"
}

terraform {
  backend "gcs" {
    bucket = "68717893997"
    prefix = "terraform/state"
  }
}

data "google_project" "this" {}

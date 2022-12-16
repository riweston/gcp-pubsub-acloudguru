provider "google" {
  project = "cr-lab-rweston-2811223416"
}

terraform {
  backend "gcs" {
    bucket = "68717893997"
    prefix = "terraform/state"
  }
}

resource "google_storage_bucket" "state" {
  name     = data.google_project.this.number
  location = local.location
}


resource "google_storage_bucket" "cloud_functions" {
  name          = "${data.google_project.this.number}-cloud-functions"
  location      = local.location
  force_destroy = true
}

data "google_project" "this" {}

provider "google" {
  project = "cr-lab-rweston-0301234201"
}

terraform {
  backend "gcs" {
    bucket = "529307337096"
    prefix = "terraform/state"
  }
  required_version = ">= 1.3.6"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 4.48.0"
    }
  }
}

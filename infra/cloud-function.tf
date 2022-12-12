locals {
  function_name_http_trigger = "http-trigger"
  location                   = "europe-west2"
}

resource "google_storage_bucket" "this" {
  name                        = data.google_project.this.number
  location                    = local.location
  uniform_bucket_level_access = true
  force_destroy               = true
}

data "archive_file" "http_trigger" {
  type        = "zip"
  source_dir  = "${path.module}/../src/${local.function_name_http_trigger}"
  output_path = "/tmp/${local.function_name_http_trigger}.zip"
}

resource "google_storage_bucket_object" "http_trigger" {
  name   = "${local.function_name_http_trigger}.zip"
  bucket = google_storage_bucket.this.name
  source = data.archive_file.http_trigger.output_path
}

resource "google_cloudfunctions2_function" "http_trigger" {
  name        = "http-trigger"
  location    = local.location
  description = "http trigger function"

  build_config {
    runtime     = "go119"
    entry_point = "HelloWorld" # Set the entry point
    source {
      storage_source {
        bucket = google_storage_bucket.this.name
        object = google_storage_bucket_object.http_trigger.name
      }
    }
  }
  service_config {
    max_instance_count = 1
    available_memory   = "256M"
    timeout_seconds    = 60
  }

  lifecycle {
    replace_triggered_by = [
      google_storage_bucket_object.http_trigger.md5hash
    ]
  }
}

resource "google_cloud_run_service_iam_binding" "http_trigger" {
  location = google_cloudfunctions2_function.http_trigger.location
  service  = google_cloudfunctions2_function.http_trigger.name
  role     = "roles/run.invoker"
  members  = [
    "allUsers"
  ]
}

output "function_uri" {
  value = google_cloudfunctions2_function.http_trigger.service_config[0].uri
}

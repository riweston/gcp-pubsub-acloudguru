locals {
  function_name_http_trigger = "http_router"
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
  source_dir  = "${path.module}/../src/build/${local.function_name_http_trigger}"
  output_path = "/tmp/${local.function_name_http_trigger}.zip"
}

resource "google_storage_bucket_object" "http_trigger" {
  name   = "${local.function_name_http_trigger}.zip"
  bucket = google_storage_bucket.this.name
  source = data.archive_file.http_trigger.output_path
}

resource "google_cloudfunctions2_function" "http_trigger" {
  name        = replace(local.function_name_http_trigger, "_", "-")
  location    = local.location
  description = "http trigger function"

  build_config {
    runtime     = "go119"
    entry_point = "Router" # Set the entry point
    source {
      storage_source {
        bucket = google_storage_bucket.this.name
        object = google_storage_bucket_object.http_trigger.name
      }
    }
  }
  service_config {
    max_instance_count    = 1
    available_memory      = "256M"
    timeout_seconds       = 60
    environment_variables = {
      PROJECT_ID           = data.google_project.this.number
      TOPIC_NAME           = google_pubsub_topic.this.name
      SLACK_SIGNING_SECRET = "a7c64cd4d17448bc3964c1473e03bb6a"
    }
    service_account_email = google_service_account.http_trigger.email
  }

  lifecycle {
    replace_triggered_by = [
      google_storage_bucket_object.http_trigger.md5hash
    ]
  }
}

resource "google_service_account" "http_trigger" {
  account_id = replace(local.function_name_http_trigger, "_", "-")
}

resource "google_cloud_run_service_iam_binding" "http_trigger" {
  location = google_cloudfunctions2_function.http_trigger.location
  service  = google_cloudfunctions2_function.http_trigger.name
  role     = "roles/run.invoker"
  members  = [
    "allUsers"
  ]
}

resource "google_cloudfunctions2_function_iam_member" "member" {
  project        = google_cloudfunctions2_function.http_trigger.project
  location       = google_cloudfunctions2_function.http_trigger.location
  cloud_function = google_cloudfunctions2_function.http_trigger.name

  role   = "roles/cloudfunctions.invoker"
  member = "allUsers"
}

output "function_uri" {
  value = google_cloudfunctions2_function.http_trigger.service_config[0].uri
}

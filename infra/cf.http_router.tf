locals {
  function_name_http_router = "http_router"
  location                  = "europe-west2"
}

data "archive_file" "http_router" {
  type        = "zip"
  source_dir  = "${path.module}/../src/build/${local.function_name_http_router}"
  output_path = "/tmp/${local.function_name_http_router}.zip"
}

resource "google_storage_bucket_object" "http_router" {
  name   = "${local.function_name_http_router}.zip"
  bucket = google_storage_bucket.cloud_functions.name
  source = data.archive_file.http_router.output_path
}

data "google_secret_manager_secret" "http_router" {
  secret_id = "SLACK_SIGNING_SECRET"
}

resource "google_secret_manager_secret_iam_binding" "http_router" {
  secret_id = data.google_secret_manager_secret.http_router.secret_id
  role      = "roles/secretmanager.secretAccessor"
  members   = [
    "serviceAccount:${google_service_account.http_router.email}",
  ]
}

resource "google_cloudfunctions2_function" "http_router" {
  name        = replace(local.function_name_http_router, "_", "-")
  location    = local.location
  description = "http trigger function"

  build_config {
    runtime     = "go119"
    entry_point = "Router"
    source {
      storage_source {
        bucket = google_storage_bucket.cloud_functions.name
        object = google_storage_bucket_object.http_router.name
      }
    }
  }
  service_config {
    max_instance_count    = 1
    available_memory      = "256M"
    timeout_seconds       = 60
    environment_variables = {
      PROJECT_ID = data.google_project.this.number
      TOPIC_NAME = google_pubsub_topic.request.name
    }
    secret_environment_variables {
      key        = "SLACK_SIGNING_SECRET"
      project_id = data.google_project.this.number
      secret     = data.google_secret_manager_secret.http_router.secret_id
      version    = "latest"
    }

    service_account_email = google_service_account.http_router.email
  }

  lifecycle {
    replace_triggered_by = [
      google_storage_bucket_object.http_router.md5hash
    ]
  }
}

resource "google_service_account" "http_router" {
  account_id = replace(local.function_name_http_router, "_", "-")
}

resource "google_cloud_run_service_iam_member" "http_router" {
  location = google_cloudfunctions2_function.http_router.location
  service  = google_cloudfunctions2_function.http_router.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}

resource "google_cloudfunctions2_function_iam_member" "http_router" {
  location       = google_cloudfunctions2_function.http_router.location
  cloud_function = google_cloudfunctions2_function.http_router.name
  role           = "roles/cloudfunctions.invoker"
  member         = "allUsers"
}

resource "google_pubsub_topic" "request" {
  name = "acg-request"
}

resource "google_pubsub_topic_iam_member" "request" {
  topic  = google_pubsub_topic.request.name
  role   = "roles/pubsub.editor"
  member = "serviceAccount:${google_service_account.http_router.email}"
}

output "function_uri" {
  value = google_cloudfunctions2_function.http_router.service_config[0].uri
}

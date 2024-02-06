data "archive_file" "this" {
  type        = "zip"
  source_dir  = var.source_dir
  output_path = "/tmp/${var.function_name}.zip"
}

resource "google_storage_bucket_object" "this" {
  name   = "${var.function_name}.${data.archive_file.this.output_md5}.zip"
  bucket = var.bucket_name
  source = data.archive_file.this.output_path
}

resource "google_cloudfunctions2_function" "this" {
  name        = replace(var.function_name, "_", "-")
  location    = var.function_location
  description = var.function_description

  build_config {
    runtime     = var.function_runtime
    entry_point = var.entry_point
    source {
      storage_source {
        bucket = var.bucket_name
        object = google_storage_bucket_object.this.name
      }
    }
  }

  service_config {
    max_instance_count    = 1
    available_memory      = "256M"
    timeout_seconds       = 60
    service_account_email = google_service_account.this.email
    environment_variables = var.environment_variables

    dynamic "secret_environment_variables" {
      for_each = data.google_secret_manager_secret.secret_environment_variables
      content {
        key        = secret_environment_variables.value.secret_id
        project_id = var.project_id
        secret     = secret_environment_variables.value.secret_id
        version    = "latest"
      }
    }
  }
}

resource "google_service_account" "this" {
  account_id = replace(var.function_name, "_", "-")
}

// IAM for Cloud Run Service and Cloud Function
// Only grant public HTTP access to the service if there is no pubsub subscription

resource "google_cloud_run_service_iam_member" "this" {
  depends_on = [
    google_cloudfunctions2_function.this
  ]

  lifecycle {
    replace_triggered_by = [
      google_cloudfunctions2_function.this
    ]
  }

  location = google_cloudfunctions2_function.this.location
  service  = google_cloudfunctions2_function.this.name
  role     = "roles/run.invoker"
  member   = var.pubsub_subscription_topic != "" ? "serviceAccount:${google_service_account.this.email}" : "allUsers"
}

resource "google_cloudfunctions2_function_iam_member" "this" {
  depends_on = [
    google_cloudfunctions2_function.this
  ]

  lifecycle {
    replace_triggered_by = [
      google_cloudfunctions2_function.this
    ]
  }

  location       = google_cloudfunctions2_function.this.location
  cloud_function = google_cloudfunctions2_function.this.name
  role           = "roles/cloudfunctions.invoker"
  member         = var.pubsub_subscription_topic != "" ? "serviceAccount:${google_service_account.this.email}" : "allUsers"
}

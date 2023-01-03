data "archive_file" "this" {
  type        = "zip"
  source_dir  = var.source_dir
  output_path = "/tmp/${var.function_name}.zip"
}

resource "google_storage_bucket_object" "this" {
  name   = "${var.function_name}.zip"
  bucket = var.bucket_name
  source = data.archive_file.this.output_path
}

resource "google_cloudfunctions2_function" "this" {
  name        = replace(var.function_name, "_", "-")
  location    = var.location
  description = "http trigger function"

  build_config {
    runtime     = "go119"
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
      for_each = var.secret_environment_variables
      content {
        key        = secret_environment_variables.value.key
        project_id = var.project_id
        secret     = secret_environment_variables.value.secret_value
        version    = "latest"
      }
    }
  }

  lifecycle {
    replace_triggered_by = [
      google_storage_bucket_object.this.md5hash
    ]
  }
}

resource "google_service_account" "this" {
  account_id = replace(var.function_name, "_", "-")
}

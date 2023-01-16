locals {
  location   = "europe-west2"
  source_dir = "../../src/build"

  function_name_http_router         = "http_router"
  function_name_lookup_slack_id     = "lookup_slack_id"
  function_name_house_keeping       = "house_keeping"
  function_name_activate_deactivate = "activate_deactivate"

  topic_name_requests   = "requests"
  topic_name_work_queue = "work-queue"
}

data "google_project" "this" {}

data "google_storage_bucket" "this" {
  name = "${data.google_project.this.number}-cloud-functions"
}

module "http_router" {
  source = "./../_modules/cloudfunction"

  source_dir    = "${local.source_dir}/${local.function_name_http_router}"
  function_name = local.function_name_http_router
  bucket_name   = data.google_storage_bucket.this.name
  location      = local.location
  environment_variables = {
    TOPIC_NAME = local.topic_name_requests
  }
  secret_environment_variables = [
    "SLACK_SIGNING_SECRET",
  ]
  pubsub_topic = local.topic_name_requests
}

module "lookup_slack_id" {
  source = "./../_modules/cloudfunction"

  source_dir    = "${local.source_dir}/${local.function_name_lookup_slack_id}"
  function_name = local.function_name_lookup_slack_id
  bucket_name   = data.google_storage_bucket.this.name
  location      = local.location
  environment_variables = {
    TOPIC_NAME = local.topic_name_work_queue
  }
  secret_environment_variables = [
    "SLACK_BOT_TOKEN",
    "ACLOUDGURU_API_KEY",
    "ACLOUDGURU_CONSUMER_ID",
  ]
  pubsub_topic              = local.topic_name_work_queue
  pubsub_subscription_topic = module.http_router.pubsub_topic
}

module "house_keeping" {
  source = "./../_modules/cloudfunction"

  source_dir    = "${local.source_dir}/${local.function_name_house_keeping}"
  function_name = local.function_name_house_keeping
  bucket_name   = data.google_storage_bucket.this.name
  location      = local.location
  environment_variables = {
    TOPIC_NAME  = module.lookup_slack_id.pubsub_topic
    DAYS_CAP    = 45
    LICENSE_CAP = 250
  }
  secret_environment_variables = [
    "ACLOUDGURU_API_KEY",
    "ACLOUDGURU_CONSUMER_ID",
  ]
  pubsub_subscription_topic = module.http_router.pubsub_topic
}

module "activate_deactivate" {
  source = "./../_modules/cloudfunction"

  source_dir    = "${local.source_dir}/${local.function_name_activate_deactivate}"
  function_name = local.function_name_activate_deactivate
  bucket_name   = data.google_storage_bucket.this.name
  location      = local.location
  secret_environment_variables = [
    "SLACK_BOT_TOKEN",
    "ACLOUDGURU_API_KEY",
    "ACLOUDGURU_CONSUMER_ID",
  ]
  pubsub_subscription_topic = module.lookup_slack_id.pubsub_topic
}

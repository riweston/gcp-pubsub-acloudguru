include {
  path = find_in_parent_folders()
}

locals {
  function_name = basename(dirname(get_terragrunt_dir()))
}

terraform {
  source = "${path_relative_from_include()}../../../..//_modules/cloudfunction"
}

dependency "core" {
  config_path = find_in_parent_folders("core")
}

dependency "pubsub" {
  config_path = "${find_in_parent_folders("lookup_slack_id")}//pubsub_topic"
}

inputs = {
  source_dir            = "${find_in_parent_folders("src")}/build/${local.function_name}"
  entry_point           = "HouseKeeping"
  function_name         = local.function_name
  bucket_name           = dependency.core.outputs.bucket_name_cloud_function
  location              = dependency.core.outputs.location
  project_id            = dependency.core.outputs.project_id
  environment_variables = {
    PROJECT_ID  = dependency.core.outputs.project_id
    TOPIC_NAME  = dependency.pubsub.outputs.topic_name
    DAYS_CAP    = 45
    LICENSE_CAP = 250
  }
  secret_environment_variables = [
    {
      key          = "ACLOUDGURU_API_KEY"
      secret_value = dependency.core.outputs.secret_id_acloudguru_api_key
    },
    {
      key          = "ACLOUDGURU_CONSUMER_ID"
      secret_value = dependency.core.outputs.secret_id_acloudguru_consumer_id
    },
  ]
}

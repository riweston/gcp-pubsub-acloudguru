include {
  path = find_in_parent_folders()
}

dependency "core" {
  config_path = find_in_parent_folders("core")
}

dependency "cloudfunction" {
  config_path = find_in_parent_folders("cloudfunction")
}

dependency "pubsub" {
  config_path = find_in_parent_folders("pubsub_topic")
}

inputs = {
  function_name                  = dependency.cloudfunction.outputs.function_name
  function_location              = dependency.core.outputs.location
  function_service_account_email = dependency.cloudfunction.outputs.service_account_email
  secret_id                      = dependency.core.outputs.secret_id_slack_bot_token
  topic_name                     = dependency.pubsub.outputs.topic_name
}

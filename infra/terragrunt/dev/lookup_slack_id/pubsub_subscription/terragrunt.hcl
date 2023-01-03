include {
  path = find_in_parent_folders()
}

dependency "cloudfunction" {
  config_path = find_in_parent_folders("cloudfunction")
}

dependency "pubsub" {
  config_path = "${find_in_parent_folders("http_router")}//pubsub"
}

inputs = {
  topic         = dependency.pubsub.outputs.topic_name
  push_endpoint = dependency.cloudfunction.outputs.endpoint
  function_service_account_email = dependency.cloudfunction.outputs.service_account_email
}

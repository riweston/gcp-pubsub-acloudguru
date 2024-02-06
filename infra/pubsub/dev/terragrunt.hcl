locals {
  environment = basename(get_terragrunt_dir())
  project_id  = read_terragrunt_config(find_in_parent_folders("vars.hcl")).locals.project_id
  topics      = read_terragrunt_config(find_in_parent_folders()).locals.topics
}

terraform {
  source = "../"
}

include {
  path = find_in_parent_folders()
}

inputs = {
  project_id = local.project_id
  topics     = [for topic in local.topics : "${local.environment}-${topic}"]
}

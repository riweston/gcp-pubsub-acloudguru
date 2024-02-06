locals {
  # Automatically load environment-level variables
  project_id     = read_terragrunt_config(find_in_parent_folders("vars.hcl")).locals.project_id
  project_number = read_terragrunt_config(find_in_parent_folders("vars.hcl")).locals.project_number
}

terraform {
  source = "./"
}

include {
  path = find_in_parent_folders()
}

inputs = {
  project_id     = local.project_id
  project_number = local.project_number
  location       = "eu"
}

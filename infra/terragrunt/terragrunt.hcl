remote_state {
  backend  = "gcs"
  generate = {
    path      = "core_generated.tf"
    if_exists = "overwrite_terragrunt"
  }
  config = {
    bucket = "529307337096"
    prefix = "${basename(get_parent_terragrunt_dir())}/${path_relative_to_include()}.tfstate"
  }
}

generate "providers" {
  path      = "providers_generated.tf"
  if_exists = "overwrite"
  contents  = <<EOF
provider "google" {
 project = "cr-lab-rweston-0301234201"
}
EOF
}

variable "source_dir" {}

variable "entry_point" {}

variable "function_name" {
  default = ""
}

variable "bucket_name" {
  default = ""
}

variable "location" {
  default = ""
}

variable "project_id" {
  default = ""
}

variable "environment_variables" {
  default = null
  type    = map(string)
}

variable "secret_environment_variables" {
  default = []
  type = list(object({
    key          = string,
    secret_value = string,
  }))
}

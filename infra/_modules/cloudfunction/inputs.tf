variable "source_dir" {
  default     = ""
  description = "The path to the source code directory"
}

variable "entry_point" {
  default     = "EntryPoint"
  description = "The name of a function (as defined in source code) that will be executed. Defaults to 'EntryPoint'"
}

variable "function_name" {
  default     = ""
  description = "The name of the function"
}

variable "bucket_name" {
  default     = ""
  description = "The name of the bucket where the function's deployment package will be stored"
}

variable "location" {
  default     = ""
  description = "The location in which the function should be created"
}

variable "project_id" {
  default     = ""
  description = "The ID of the project in which the function will be created"
}

variable "environment_variables" {
  default     = null
  type        = map(string)
  description = "A map of environment variables to pass to the function."
}

variable "secret_environment_variables" {
  default     = []
  type        = list(string)
  description = "A list of secret environment variables to pass to the function that are sourced from Secret Manager."
}

variable "pubsub_topic" {
  default     = ""
  description = "The name of the PubSub topic to which messages will be published after a function is triggered."
}

variable "pubsub_subscription_topic" {
  default     = ""
  description = "The name of the PubSub topic to which the function will subscribe."
}

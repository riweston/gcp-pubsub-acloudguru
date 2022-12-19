include {
  path = find_in_parent_folders()
}

locals {
  topic_name = "requests"
}

terraform {
  source = "${path_relative_from_include()}../../../..//_modules/pubsub_topic"
}

inputs = {
  name = "${local.topic_name}"
}

/*resource "google_pubsub_schema" "this" {
  name       = "acg-request"
  type       = "AVRO"
  definition = jsonencode(
    {
      type   = "record"
      name   = "request"
      fields = [
        {
          name = "user"
          type = "string"
        },
        {
          name = "activate"
          type = "string"
        }
      ]
    }
  )
}*/

resource "google_pubsub_topic" "this" {
  name = "acg-request"
  /*  schema_settings {
      schema   = google_pubsub_schema.this.id
      encoding = "JSON"
    }*/
}

resource "google_pubsub_subscription" "this" {
  name                 = "acg-request"
  topic                = google_pubsub_topic.this.name
  ack_deadline_seconds = 600
}

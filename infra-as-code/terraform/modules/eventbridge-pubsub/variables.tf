variable "bus_name" {
  description = "The name of event bus that acts to support the pub/sub exchange of events."
  type        = string
}

variable "source_id" {
  description = "The source id string to filter events on"
  type        = string
}

variable "lambda_subscribers" {
  description = "A list of lambda ARNs that need to be subscribed to events on the event bus."
  type        = map(string)
}

variable "topic_name" {
  description = "The name of the topic."
  type        = string
}

variable "lambda_subscribers" {
  description = "A list of lambda ARNs to subscribe to the topic"
  type        = map(string)
}

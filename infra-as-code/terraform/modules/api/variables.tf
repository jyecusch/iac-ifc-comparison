variable "name" {
  description = "The name of the API Gateway"
  type        = string
}

variable "spec" {
  description = "Open API spec"
  type        = string
}

variable "target_lambda_functions" {
  description = "The names of the target lambda functions"
  type        = map(string)
}

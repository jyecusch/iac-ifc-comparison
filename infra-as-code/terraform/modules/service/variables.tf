variable "service_name" {
  type        = string
  description = "The name of the service"
}

variable "image" {
  type        = string
  description = "The docker image to deploy"
}

# environment variables
variable "environment" {
  type        = map(string)
  description = "Environment variables to set on the lambda function"
}

variable "memory" {
  description = "The amount of memory to allocate to the lambda function"
  type        = number
  default     = 128
}

variable "timeout" {
  description = "The timeout for the lambda function"
  type        = number
  default     = 10
}

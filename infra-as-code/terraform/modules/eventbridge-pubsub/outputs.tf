output "bus_arn" {
  description = "The ARN of the deployed event bus"
  value       = aws_cloudwatch_event_bus.pubsub_event_bus.arn
}

output "bus_name" {
  description = "The name of the deployed event bus"
  value       = aws_cloudwatch_event_bus.pubsub_event_bus.name
}

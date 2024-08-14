output "topic_arn" {
  description = "The ARN of the deployed topic"
  value       = aws_sns_topic.topic.arn
}

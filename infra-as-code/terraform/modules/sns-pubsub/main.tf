# AWS SNS Topic
resource "aws_sns_topic" "topic" {
  name = "${var.topic_name}"
}

# Loop over the subsribers and deploy subscriptions and permissions
resource "aws_sns_topic_subscription" "subscription" {
  for_each = var.lambda_subscribers

  topic_arn = aws_sns_topic.topic.arn
  protocol  = "lambda"
  endpoint  = each.value
}

resource "aws_lambda_permission" "sns" {
  for_each = var.lambda_subscribers

  action        = "lambda:InvokeFunction"
  function_name = each.value
  principal     = "sns.amazonaws.com"
  source_arn    = aws_sns_topic.topic.arn
}

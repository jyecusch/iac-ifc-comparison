resource "aws_cloudwatch_event_bus" "pubsub_event_bus" {
  name = "pubsub-event-bus"
}

resource "aws_cloudwatch_event_rule" "pubsub_event_rule" {
  name           = "pubsub-event-rule"
  event_bus_name = aws_cloudwatch_event_bus.pubsub_event_bus.name
  event_pattern = jsonencode({
    "source" : [var.source_id]
  })
}

resource "aws_cloudwatch_event_target" "custom_event_target" {
  for_each = var.lambda_subscribers

  event_bus_name = aws_cloudwatch_event_bus.pubsub_event_bus.name
  rule      = aws_cloudwatch_event_rule.pubsub_event_rule.name
  target_id = "subscriberLambda${each.key}"
  arn       = each.value
}

resource "aws_lambda_permission" "allow_eventbridge_invoke" {
  for_each = var.lambda_subscribers

  statement_id  = "AllowExecutionFromEventBridge${each.key}"
  action        = "lambda:InvokeFunction"
  function_name = each.value
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.pubsub_event_rule.arn
}

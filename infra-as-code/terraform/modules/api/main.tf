resource "aws_apigatewayv2_api" "api_gateway" {
  name          = var.name
  protocol_type = "HTTP"
  body          = var.spec
}

resource "aws_apigatewayv2_stage" "stage" {
  api_id      = aws_apigatewayv2_api.api_gateway.id
  name        = "$default"
  auto_deploy = true
}

# deploy lambda permissions for execution
resource "aws_lambda_permission" "apigw_lambda" {
  for_each      = var.target_lambda_functions
  action        = "lambda:InvokeFunction"
  function_name = each.value
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.api_gateway.execution_arn}/*/*/*"
}

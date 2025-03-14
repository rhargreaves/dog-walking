output "api_url" {
  value = aws_apigatewayv2_api.api.api_endpoint
}

output "api_domain" {
  value = aws_apigatewayv2_domain_name.api.domain_name
}

output "lambda_function_name" {
  value = aws_lambda_function.api.function_name
}
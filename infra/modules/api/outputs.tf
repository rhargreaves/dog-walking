output "lambda_function_name" {
  value = aws_lambda_function.api.function_name
}

output "api_id" {
  value       = aws_apigatewayv2_api.api.id
  description = "The ID of the API Gateway (for CloudWatch dashboards)"
}

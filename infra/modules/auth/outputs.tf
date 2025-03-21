output "cognito_user_pool_id" {
  description = "The ID of the Cognito User Pool"
  value       = aws_cognito_user_pool.pool.id
}

output "cognito_user_pool_arn" {
  description = "The ARN of the Cognito User Pool"
  value       = aws_cognito_user_pool.pool.arn
}

output "cognito_client_ids" {
  description = "The IDs of the Cognito User Pool Clients that can be used to authenticate with the API"
  value       = [
    aws_cognito_user_pool_client.api.id,
    aws_cognito_user_pool_client.ui.id
    ]
}
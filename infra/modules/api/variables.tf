variable "environment" {
  description = "The deployment environment (e.g., uat, prod)"
  type        = string
}

variable "application_name" {
  description = "The name of the application"
  type        = string
}

variable "api_base_host" {
  description = "The base DNS name for the API"
  type        = string
}

variable "hosted_zone_id" {
  description = "The Route53 hosted zone ID for the domain"
  type        = string
}

variable "dynamodb_access_policy_arn" {
  description = "ARN of the policy for DynamoDB access"
  type        = string
}

variable "s3_access_policy_arn" {
  description = "ARN of the policy for S3 dog images access"
  type        = string
}

variable "dogs_table_name" {
  description = "Name of the DynamoDB table for dogs"
  type        = string
}

variable "dog_images_bucket" {
  description = "Name of the S3 bucket for dog images"
  type        = string
}

variable "bootstrap_path" {
  description = "Path to the bootstrap file"
  type        = string
}

variable "cognito_user_pool_id" {
  description = "The ID of the Cognito User Pool"
  type        = string
}

variable "cognito_client_ids" {
  description = "The IDs of the Cognito User Pool Clients that can be used to authenticate with the API"
  type        = list(string)
}

variable "cors_allowed_origin" {
  description = "The allowed origin for CORS"
  type        = string
}

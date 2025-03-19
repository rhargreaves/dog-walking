variable "application_name" {
  type        = string
  description = "Name of the application"
}

variable "environment" {
  type        = string
  description = "Environment name (e.g., dev, prod)"
}

variable "api_id" {
  type        = string
  description = "ID of the API Gateway"
}

variable "lambda_function_name" {
  type        = string
  description = "Name of the Lambda function"
}
variable "service_name" {
  type        = string
  description = "Name of the service"
}

variable "api_id" {
  type        = string
  description = "ID of the API Gateway"
}

variable "lambda_function_name" {
  type        = string
  description = "Name of the Lambda function"
}

variable "environment" {
  type        = string
  description = "Environment name (e.g., dev, prod)"
}

variable "aws_region" {
  type        = string
  description = "AWS region"
}
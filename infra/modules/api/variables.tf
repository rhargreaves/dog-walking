variable "environment" {
  description = "The deployment environment (e.g., uat, prod)"
  type        = string
}

variable "application_name" {
  description = "The name of the application"
  type        = string
}

variable "vpc_id" {
  description = "The ID of the VPC"
  type        = string
}

variable "private_subnet_ids" {
  description = "List of private subnet IDs"
  type        = list(string)
}

variable "domain_name" {
  description = "The base domain name for the application"
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

variable "dogs_table_name" {
  description = "Name of the DynamoDB table for dogs"
  type        = string
}

variable "bootstrap_path" {
  description = "Path to the bootstrap file"
  type        = string
}

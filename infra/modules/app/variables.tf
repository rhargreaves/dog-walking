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

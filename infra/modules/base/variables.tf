variable "environment" {
  description = "The deployment environment (e.g., uat, prod)"
  type        = string
}

variable "application_name" {
  description = "The name of the application"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block for the VPC"
  type        = string
}

variable "availability_zones" {
  description = "List of availability zones to use"
  type        = list(string)
}
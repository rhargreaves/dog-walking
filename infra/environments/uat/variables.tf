variable "aws_region" {
  type    = string
}

variable "environment" {
  type    = string
}

variable "application_name" {
  type    = string
}

variable "vpc_cidr" {
  type    = string
}

variable "availability_zones" {
  type    = list(string)
}

variable "domain_name" {
  type        = string
  description = "The base domain name for the application"
}

variable "hosted_zone_id" {
  type        = string
  description = "The Route53 hosted zone ID for the domain"
}

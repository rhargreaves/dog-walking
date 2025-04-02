variable "aws_region" {
  type = string
}

variable "environment" {
  type = string
}

variable "application_name" {
  type = string
}

variable "api_base_host" {
  type        = string
  description = "The base DNS name for the API"
}

variable "hosted_zone_id" {
  type        = string
  description = "The Route53 hosted zone ID for the domain"
}

variable "api_bootstrap_path" {
  description = "Path to the Lambda bootstrap binary"
  type        = string
}

variable "photo_moderation_bootstrap_path" {
  description = "Path to the Lambda bootstrap binary"
  type        = string
}

variable "sysadmin_username" {
  description = "Username for user with full admin access"
  type        = string
  sensitive   = true
}

variable "sysadmin_password" {
  description = "Password for user with full admin access"
  type        = string
  sensitive   = true
}

variable "cors_allowed_origin" {
  description = "The allowed origin for CORS"
  type        = string
}

variable "images_cdn_host" {
  description = "The fully qualified domain name for the images CDN"
  type        = string
}

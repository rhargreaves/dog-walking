variable "environment" {
  description = "The deployment environment (e.g., uat, prod)"
  type        = string
}

variable "hosted_zone_id" {
  description = "The ID of the Route53 hosted zone for DNS records"
  type        = string
}

variable "images_cdn_host" {
  description = "The fully qualified domain name for the images CDN"
  type        = string
}
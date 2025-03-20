variable "environment" {
  description = "Environment name (e.g., dev, prod)"
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
variable "environment" {
  description = "Environment name (e.g., dev, prod)"
  type        = string
}

variable "sysadmin_password" {
  description = "Password for the sysadmin user"
  type        = string
  sensitive   = true
}
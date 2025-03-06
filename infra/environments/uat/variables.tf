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

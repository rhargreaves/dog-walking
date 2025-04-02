
variable "environment" {
  type = string
}

variable "dynamodb_access_policy_arn" {
  type = string
}

variable "s3_access_policy_arn" {
  type = string
}

variable "bootstrap_path" {
  type = string
}

variable "pending_dog_images_bucket_name" {
  type = string
}

variable "dog_images_bucket" {
  type = string
}

variable "dogs_table_name" {
  type = string
}

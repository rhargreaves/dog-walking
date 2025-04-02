provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      env = var.environment
      app = var.application_name
    }
  }
}

module "data" {
  source = "../../modules/data"

  environment = var.environment
}

module "images" {
  source = "../../modules/images"

  environment     = var.environment
  hosted_zone_id  = var.hosted_zone_id
  images_cdn_host = var.images_cdn_host
}

module "api" {
  source = "../../modules/api"

  environment                    = var.environment
  application_name               = var.application_name
  api_base_host                  = var.api_base_host
  hosted_zone_id                 = var.hosted_zone_id
  dynamodb_access_policy_arn     = module.data.dynamodb_access_policy_arn
  s3_access_policy_arn           = module.images.s3_access_policy_arn
  dogs_table_name                = module.data.dogs_table_name
  dog_images_bucket              = module.images.bucket_name
  bootstrap_path                 = var.api_bootstrap_path
  cognito_client_ids             = module.auth.cognito_client_ids
  cognito_user_pool_id           = module.auth.cognito_user_pool_id
  cors_allowed_origin            = var.cors_allowed_origin
  images_cdn_host                = var.images_cdn_host
  pending_dog_images_bucket_name = module.images.pending_dog_images_bucket_name
}

module "auth" {
  source = "../../modules/auth"

  environment       = var.environment
  sysadmin_password = var.sysadmin_password
  sysadmin_username = var.sysadmin_username
}

module "monitoring" {
  source = "../../modules/monitoring"

  application_name     = var.application_name
  environment          = var.environment
  api_id               = module.api.api_id
  lambda_function_name = module.api.lambda_function_name
}

module "photo_moderation" {
  source = "../../modules/photo-moderation"

  environment                    = var.environment
  dynamodb_access_policy_arn     = module.data.dynamodb_access_policy_arn
  s3_access_policy_arn           = module.images.s3_access_policy_arn
  dogs_table_name                = module.data.dogs_table_name
  dog_images_bucket              = module.images.bucket_name
  pending_dog_images_bucket_name = module.images.pending_dog_images_bucket_name
  pending_dog_images_bucket_arn  = module.images.pending_dog_images_bucket_arn
  bootstrap_path                 = var.photo_moderation_bootstrap_path
}

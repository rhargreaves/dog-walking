provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      env = var.environment
      app = var.application_name
    }
  }
}

module "base" {
  source = "../../modules/base"

  environment        = var.environment
  application_name   = var.application_name
  vpc_cidr           = var.vpc_cidr
  availability_zones = var.availability_zones
}

module "data" {
  source = "../../modules/data"

  environment = var.environment
}

module "images" {
  source = "../../modules/images"

  environment = var.environment
}

module "api" {
  source = "../../modules/api"

  environment                = var.environment
  application_name           = var.application_name
  vpc_id                     = module.base.vpc_id
  private_subnet_ids         = module.base.private_subnet_ids
  domain_name                = var.domain_name
  hosted_zone_id             = var.hosted_zone_id
  dynamodb_access_policy_arn = module.data.dynamodb_access_policy_arn
  s3_access_policy_arn       = module.images.s3_access_policy_arn
  dogs_table_name            = module.data.dogs_table_name
  dog_images_bucket          = module.images.bucket_name
  bootstrap_path             = var.bootstrap_path
  cognito_client_id          = module.auth.cognito_client_id
  cognito_user_pool_id       = module.auth.cognito_user_pool_id
}

module "auth" {
  source = "../../modules/auth"

  environment       = var.environment
  sysadmin_username = var.sysadmin_username
  sysadmin_password = var.sysadmin_password
}

module "monitoring" {
  source = "../../modules/monitoring"

  application_name     = var.application_name
  environment          = var.environment
  api_id               = module.api.api_id
  lambda_function_name = module.api.lambda_function_name
}

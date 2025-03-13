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

  environment      = var.environment
  application_name = var.application_name
  vpc_cidr         = var.vpc_cidr
  availability_zones = var.availability_zones
}

module "data" {
  source = "../../modules/data"

  environment = var.environment
}

module "app" {
  source = "../../modules/app"

  environment        = var.environment
  application_name   = var.application_name
  vpc_id             = module.base.vpc_id
  private_subnet_ids = module.base.private_subnet_ids
  domain_name        = var.domain_name
  hosted_zone_id     = var.hosted_zone_id
  dynamodb_access_policy_arn = module.data.dynamodb_access_policy_arn
  dogs_table_name = module.data.dogs_table_name
}

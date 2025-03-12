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

module "app" {
  source = "../../modules/app"

  environment        = var.environment
  application_name   = var.application_name
  vpc_id             = module.base.vpc_id
  private_subnet_ids = module.base.private_subnet_ids
  domain_name        = var.domain_name
  hosted_zone_id     = var.hosted_zone_id
  api_zip_path       = "${path.module}/../../api.zip"
}

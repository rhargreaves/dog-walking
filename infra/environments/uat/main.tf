provider "aws" {
  region = var.aws_region
}

module "base" {
  source = "../../modules/base"

  environment      = var.environment
  application_name = var.application_name
  vpc_cidr         = var.vpc_cidr
  availability_zones = var.availability_zones
}

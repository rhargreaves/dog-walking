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

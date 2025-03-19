terraform {
  backend "s3" {
    bucket       = "rh-dog-walking-terraform-state"
    key          = "uat/dog-walking.tfstate"
    region       = "eu-west-1"
    encrypt      = true
    use_lockfile = true
  }
}

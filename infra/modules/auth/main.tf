resource "aws_cognito_user_pool" "pool" {
  name = "${var.environment}-dog-walking"

  password_policy {
    minimum_length    = 8
    require_lowercase = true
    require_numbers   = true
    require_symbols   = true
    require_uppercase = true
  }

  username_attributes      = ["email"]
  auto_verified_attributes = ["email"]

  schema {
    attribute_data_type = "String"
    name                = "email"
    required            = true
    mutable             = true

    string_attribute_constraints {
      min_length = 1
      max_length = 256
    }
  }

  verification_message_template {
    default_email_option = "CONFIRM_WITH_CODE"
    email_subject        = "Your verification code"
    email_message        = "Your verification code is {####}"
  }

  account_recovery_setting {
    recovery_mechanism {
      name     = "verified_email"
      priority = 1
    }
  }

  tags = {
    Name = "${var.environment}-dog-walking-user-pool"
  }
}

resource "aws_cognito_user_pool_client" "client" {
  name         = "${var.environment}-dog-walking-client"
  user_pool_id = aws_cognito_user_pool.pool.id

  explicit_auth_flows = [
    "ALLOW_USER_PASSWORD_AUTH",
    "ALLOW_REFRESH_TOKEN_AUTH"
  ]

  token_validity_units {
    access_token  = "hours"
    id_token      = "hours"
    refresh_token = "days"
  }

  access_token_validity  = 1
  id_token_validity      = 1
  refresh_token_validity = 30
}

resource "aws_cognito_user_group" "admin_group" {
  name         = "admins"
  user_pool_id = aws_cognito_user_pool.pool.id
  description  = "Administrator group"
  precedence   = 1
}

resource "aws_cognito_user" "sysadmin" {
  user_pool_id = aws_cognito_user_pool.pool.id
  username     = "sysadmin@dog-walking.roberthargreaves.com"
  password     = var.sysadmin_password
  attributes = {
    email          = "sysadmin@dog-walking.roberthargreaves.com"
    email_verified = true
  }
}

resource "aws_cognito_user_in_group" "sysadmin_admin" {
  user_pool_id = aws_cognito_user_pool.pool.id
  username     = aws_cognito_user.sysadmin.username
  group_name   = aws_cognito_user_group.admin_group.name
}
resource "aws_dynamodb_table" "dogs" {
  name         = "${var.environment}-dogs"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "id"

  attribute {
    name = "id"
    type = "S"
  }

  point_in_time_recovery {
    enabled = true
  }

  tags = {
    Name = "${var.environment}-dogs-table"
  }
}

data "aws_iam_policy_document" "dynamodb_access" {
  statement {
    actions = [
      "dynamodb:GetItem",
      "dynamodb:PutItem",
      "dynamodb:UpdateItem",
      "dynamodb:DeleteItem",
      "dynamodb:Scan",
      "dynamodb:Query",
      "dynamodb:BatchGetItem",
      "dynamodb:BatchWriteItem"
    ]
    resources = [
      aws_dynamodb_table.dogs.arn,
      "${aws_dynamodb_table.dogs.arn}/index/*"
    ]
  }
}

resource "aws_iam_policy" "dynamodb_access" {
  name        = "${var.environment}-dog-walking-dynamodb-access"
  description = "Policy for accessing DynamoDB dogs table"
  policy      = data.aws_iam_policy_document.dynamodb_access.json
}
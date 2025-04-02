
resource "aws_iam_role" "lambda_role" {
  name = "${var.environment}-dog-walking-photo-moderation-lambda-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })

  tags = {
    Name = "${var.environment}-dog-walking-photo-moderation-lambda-role"
  }
}

resource "aws_iam_policy" "rekognition_access" {
  name        = "${var.environment}-dog-walking-photo-moderation-rekognition-access"
  description = "Policy for accessing AWS Rekognition services"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "rekognition:DetectLabels"
        ]
        Effect   = "Allow"
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_rekognition" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = aws_iam_policy.rekognition_access.arn
}

resource "aws_iam_role_policy_attachment" "lambda_basic" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "lambda_dynamodb" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = var.dynamodb_access_policy_arn
}

resource "aws_iam_role_policy_attachment" "lambda_s3" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = var.s3_access_policy_arn
}

data "archive_file" "bootstrap" {
  type        = "zip"
  source_file = var.bootstrap_path
  output_path = "${var.bootstrap_path}.zip"
}

resource "aws_lambda_function" "photo_moderation" {
  function_name = "${var.environment}-dog-walking-photo-moderation"
  role          = aws_iam_role.lambda_role.arn
  handler       = "main"
  runtime       = "provided.al2023"
  architectures = ["arm64"]

  filename         = data.archive_file.bootstrap.output_path
  source_code_hash = data.archive_file.bootstrap.output_base64sha256

  logging_config {
    log_format = "Text"
    log_group  = aws_cloudwatch_log_group.lambda_logs.name
  }

  environment {
    variables = {
      ENVIRONMENT       = var.environment
      DOGS_TABLE_NAME   = var.dogs_table_name
      DOG_IMAGES_BUCKET = var.dog_images_bucket
    }
  }

  tags = {
    Name = "${var.environment}-dog-walking-photo-moderation"
  }
}

resource "aws_cloudwatch_log_group" "lambda_logs" {
  name              = "/aws/lambda/${var.environment}-dog-walking-photo-moderation"
  retention_in_days = 7

  tags = {
    Name = "${var.environment}-dog-walking-photo-moderation-logs"
  }
}

resource "aws_s3_bucket_notification" "bucket_notification" {
  bucket = var.pending_dog_images_bucket_name

  lambda_function {
    lambda_function_arn = aws_lambda_function.photo_moderation.arn
    events              = ["s3:ObjectCreated:*"]
  }

  depends_on = [aws_lambda_permission.allow_bucket_invoke]
}

resource "aws_lambda_permission" "allow_bucket_invoke" {
  statement_id  = "AllowExecutionFromS3Bucket"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.photo_moderation.arn
  principal     = "s3.amazonaws.com"
  source_arn    = var.pending_dog_images_bucket_arn
}
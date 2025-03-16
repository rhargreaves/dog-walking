resource "aws_s3_bucket" "dog_images" {
  bucket = "${var.environment}-dog-images"

  tags = {
    Name = "${var.environment}-dog-images"
  }
}

resource "aws_s3_bucket_public_access_block" "dog_images" {
  bucket = aws_s3_bucket.dog_images.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

data "aws_iam_policy_document" "s3_access" {
  statement {
    actions = [
      "s3:GetObject",
      "s3:PutObject",
      "s3:DeleteObject",
      "s3:ListBucket"
    ]
    resources = [
      aws_s3_bucket.dog_images.arn,
      "${aws_s3_bucket.dog_images.arn}/*"
    ]
  }
}

resource "aws_iam_policy" "s3_access" {
  name        = "${var.environment}-dog-walking-s3-access"
  description = "Policy for accessing S3 dog images bucket"
  policy      = data.aws_iam_policy_document.s3_access.json
}
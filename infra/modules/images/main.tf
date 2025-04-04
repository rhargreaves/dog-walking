provider "aws" { // for CloudFront certificates (must be in us-east-1)
  alias  = "us_east_1"
  region = "us-east-1"
}

resource "aws_s3_bucket" "pending_dog_images" {
  bucket = "${var.environment}-pending-dog-images"

  force_destroy = true

  tags = {
    Name = "${var.environment}-pending-dog-images"
  }
}

resource "aws_s3_bucket_public_access_block" "pending_dog_images" {
  bucket = aws_s3_bucket.pending_dog_images.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket" "dog_images" {
  bucket = "${var.environment}-dog-images"

  force_destroy = true

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
      "${aws_s3_bucket.dog_images.arn}/*",
      aws_s3_bucket.pending_dog_images.arn,
      "${aws_s3_bucket.pending_dog_images.arn}/*"
    ]
  }
}

resource "aws_iam_policy" "s3_access" {
  name        = "${var.environment}-dog-walking-s3-access"
  description = "Policy for accessing S3 dog images buckets"
  policy      = data.aws_iam_policy_document.s3_access.json
}

resource "aws_cloudfront_origin_access_identity" "dog_images" {
  comment = "OAI for ${var.environment} dog images buckets"
}

resource "aws_s3_bucket_policy" "dog_images" {
  bucket = aws_s3_bucket.dog_images.id
  policy = data.aws_iam_policy_document.cloudfront_access.json
}

data "aws_iam_policy_document" "cloudfront_access" {
  statement {
    actions   = ["s3:GetObject"]
    resources = ["${aws_s3_bucket.dog_images.arn}/*"]

    principals {
      type        = "AWS"
      identifiers = [aws_cloudfront_origin_access_identity.dog_images.iam_arn]
    }
  }
}

resource "aws_cloudfront_distribution" "dog_images" {
  origin {
    domain_name = aws_s3_bucket.dog_images.bucket_regional_domain_name
    origin_id   = aws_s3_bucket.dog_images.id

    s3_origin_config {
      origin_access_identity = aws_cloudfront_origin_access_identity.dog_images.cloudfront_access_identity_path
    }
  }

  enabled             = true
  is_ipv6_enabled     = true
  default_root_object = "index.html"
  aliases             = [var.images_cdn_host]

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD", "OPTIONS"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = aws_s3_bucket.dog_images.id

    forwarded_values {
      query_string            = true
      query_string_cache_keys = ["h"]

      cookies {
        forward = "none"
      }
    }

    viewer_protocol_policy = "https-only"
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400
  }

  price_class = "PriceClass_100"

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    acm_certificate_arn      = aws_acm_certificate_validation.cert_validation.certificate_arn
    ssl_support_method       = "sni-only"
    minimum_protocol_version = "TLSv1.2_2021"
  }

  tags = {
    Name = "${var.environment}-dog-images-cdn"
  }
}

resource "aws_acm_certificate" "cert" {
  provider          = aws.us_east_1
  domain_name       = var.images_cdn_host
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }

  tags = {
    Name = "${var.environment}-dog-images-cert"
  }
}

resource "aws_route53_record" "cert_validation" {
  for_each = {
    for dvo in aws_acm_certificate.cert.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  zone_id = var.hosted_zone_id
  name    = each.value.name
  type    = each.value.type
  records = [each.value.record]
  ttl     = 60
}

resource "aws_acm_certificate_validation" "cert_validation" {
  provider                = aws.us_east_1
  certificate_arn         = aws_acm_certificate.cert.arn
  validation_record_fqdns = [for record in aws_route53_record.cert_validation : record.fqdn]
}

resource "aws_route53_record" "dog_images" {
  zone_id = var.hosted_zone_id
  name    = var.images_cdn_host
  type    = "A"

  alias {
    name                   = aws_cloudfront_distribution.dog_images.domain_name
    zone_id                = aws_cloudfront_distribution.dog_images.hosted_zone_id
    evaluate_target_health = false
  }
}
output "bucket_name" {
  value = aws_s3_bucket.dog_images.bucket
}

output "pending_dog_images_bucket_name" {
  value = aws_s3_bucket.pending_dog_images.bucket
}

output "bucket_arn" {
  value = aws_s3_bucket.dog_images.arn
}

output "s3_access_policy_arn" {
  value = aws_iam_policy.s3_access.arn
}

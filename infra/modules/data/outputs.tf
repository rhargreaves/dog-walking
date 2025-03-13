output "dogs_table_name" {
  value = aws_dynamodb_table.dogs.name
}

output "dogs_table_arn" {
  value = aws_dynamodb_table.dogs.arn
}

output "dynamodb_access_policy_arn" {
  value = aws_iam_policy.dynamodb_access.arn
}
output "dashboard_name" {
  value       = aws_cloudwatch_dashboard.api_dashboard.dashboard_name
  description = "Name of the CloudWatch dashboard"
}

output "alerts_topic_arn" {
  value       = aws_sns_topic.alerts.arn
  description = "ARN of the SNS topic for alerts"
}
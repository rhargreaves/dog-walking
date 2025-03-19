resource "aws_cloudwatch_dashboard" "api_dashboard" {
  dashboard_name = "${var.environment}-${var.application_name}-dashboard"

  dashboard_body = jsonencode({
    widgets = [
      {
        type   = "metric"
        width  = 12
        height = 6
        properties = {
          metrics = [[{
            expression = "SELECT SUM(\"Count\") FROM SCHEMA(\"AWS/ApiGateway\", ApiId,\"Method\",Resource,Stage) GROUP BY \"Method\", Resource"
            id = "q1"
          }]]
          period = 300
          stat   = "Sum"
          region = var.aws_region
          view   = "timeSeries"
          title  = "Requests"
        }
      },
      {
        type   = "metric"
        width  = 12
        height = 6
        properties = {
          metrics = [[{
            expression = "SELECT SUM(\"4XXError\") FROM SCHEMA(\"AWS/ApiGateway\", ApiId,\"Method\",Resource,Stage) GROUP BY \"Method\", Resource"
            id = "q2"
          }]]
          period = 300
          stat   = "Sum"
          region = var.aws_region
          view   = "timeSeries"
          title  = "4XX Errors"
        }
      },
      {
        type   = "metric"
        width  = 12
        height = 6
        properties = {
          metrics = [[{
            expression = "SELECT SUM(\"5XXError\") FROM SCHEMA(\"AWS/ApiGateway\", ApiId,\"Method\",Resource,Stage) GROUP BY \"Method\", Resource"
            id = "q3"
          }]]
          period = 300
          stat   = "Sum"
          region = var.aws_region
          view   = "timeSeries"
          title  = "5XX Errors"
        }
      },
      {
        type   = "metric"
        width  = 12
        height = 6
        properties = {
          metrics = [[{
            expression = "SELECT AVG(\"Latency\") FROM SCHEMA(\"AWS/ApiGateway\", ApiId,\"Method\",Resource,Stage) GROUP BY \"Method\", Resource"
            id = "q4"
          }]]
          period = 300
          stat   = "Average"
          region = var.aws_region
          view   = "timeSeries"
          title  = "Latency"
        }
      }
    ]
  })
}

resource "aws_cloudwatch_metric_alarm" "api_errors" {
  alarm_name          = "${var.environment}-${var.application_name}-api-errors"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "5XXError"
  namespace           = "AWS/ApiGateway"
  period              = "300"
  statistic           = "Sum"
  threshold           = "5"
  alarm_description   = "This metric monitors API 5XX errors"
  alarm_actions       = [aws_sns_topic.alerts.arn]

  dimensions = {
    ApiId = var.api_id
    Stage = var.environment
  }
}

resource "aws_cloudwatch_metric_alarm" "lambda_errors" {
  alarm_name          = "${var.environment}-${var.application_name}-lambda-errors"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "Errors"
  namespace           = "AWS/Lambda"
  period              = "300"
  statistic           = "Sum"
  threshold           = "3"
  alarm_description   = "This metric monitors Lambda function errors"
  alarm_actions       = [aws_sns_topic.alerts.arn]

  dimensions = {
    FunctionName = var.lambda_function_name
  }
}

resource "aws_sns_topic" "alerts" {
  name = "${var.environment}-${var.application_name}-alerts"
}
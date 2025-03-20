data "aws_region" "current" { }

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
            expression = "SELECT SUM(\"Count\") FROM SCHEMA(\"AWS/ApiGateway\", ApiId,\"Method\",Resource,Stage) WHERE ApiId = '${var.api_id}' GROUP BY \"Method\", Resource"
            id = "q1"
          }]]
          period = 300
          stat   = "Sum"
          region = data.aws_region.current.name
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
            expression = "SELECT SUM(\"4XXError\") FROM SCHEMA(\"AWS/ApiGateway\", ApiId,\"Method\",Resource,Stage) WHERE ApiId = '${var.api_id}' GROUP BY \"Method\", Resource"
            id = "q2"
          }]]
          period = 300
          stat   = "Sum"
          region = data.aws_region.current.name
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
            expression = "SELECT SUM(\"5XXError\") FROM SCHEMA(\"AWS/ApiGateway\", ApiId,\"Method\",Resource,Stage) WHERE ApiId = '${var.api_id}' GROUP BY \"Method\", Resource"
            id = "q3"
          }]]
          period = 300
          stat   = "Sum"
          region = data.aws_region.current.name
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
            expression = "SELECT AVG(\"Latency\") FROM SCHEMA(\"AWS/ApiGateway\", ApiId,\"Method\",Resource,Stage) WHERE ApiId = '${var.api_id}' GROUP BY \"Method\", Resource"
            id = "q4"
          }]]
          period = 300
          stat   = "Average"
          region = data.aws_region.current.name
          view   = "timeSeries"
          title  = "Latency"
        }
      }
    ]
  })
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

resource "aws_cloudwatch_metric_alarm" "endpoint_error_rate" {
  for_each = {
    "create-dog" = { path = "/dogs", method = "POST" },
    "list-dogs" = { path = "/dogs", method = "GET" },
    "get-dog" = { path = "/dogs/{id}", method = "GET" },
    "update-dog" = { path = "/dogs/{id}", method = "PUT" },
    "delete-dog" = { path = "/dogs/{id}", method = "DELETE" },
    "upload-dog-photo" = { path = "/dogs/{id}/photo", method = "PUT" },
    "detect-dog-breed" = { path = "/dogs/{id}/photo/detect-breed", method = "POST" }
  }

  alarm_name                = "${var.environment}-${var.application_name}-${each.key}-error-rate"
  comparison_operator       = "GreaterThanOrEqualToThreshold"
  evaluation_periods        = 12
  threshold                 = 1
  alarm_description         = "Error rate > 1% for ${each.key}"
  alarm_actions       = [aws_sns_topic.alerts.arn]
  insufficient_data_actions = []

  metric_query {
    id          = "rate"
    expression  = "(errors*100)/requests"
    label       = "Error Rate"
    return_data = "true"
  }

  metric_query {
    id = "requests"

    metric {
      metric_name = "Count"
      namespace   = "AWS/ApiGateway"
      period      = 300
      stat        = "Sum"
      unit        = "Count"

      dimensions = {
        ApiId = var.api_id
        Resource = each.value.path
        Method = each.value.method
      }
    }
  }

  metric_query {
    id = "errors"

    metric {
      metric_name = "5XXError"
      namespace   = "AWS/ApiGateway"
      period      = 300
      stat        = "Sum"
      unit        = "Count"

      dimensions = {
        ApiId = var.api_id
        Resource = each.value.path
        Method = each.value.method
      }
    }
  }
}

resource "aws_sns_topic" "alerts" {
  name = "${var.environment}-${var.application_name}-alerts"
}
resource "aws_scheduler_schedule_group" "aws_cost_notify" {
  name = var.schedule_group_name
}

resource "aws_scheduler_schedule" "aws_cost_notify" {
  name       = var.eventbridge_schedule_name
  group_name = aws_scheduler_schedule_group.aws_cost_notify.name

  state                        = "ENABLED"
  schedule_expression          = "cron(0 09 * * ? *)"
  schedule_expression_timezone = "Asia/Tokyo"

  flexible_time_window {
    mode = "OFF"
  }

  target {
    arn      = aws_lambda_function.aws_cost_notify.arn
    role_arn = aws_iam_role.eventbridge_aws_cost_notify.arn

    retry_policy {
      maximum_event_age_in_seconds = 60
      maximum_retry_attempts       = 0
    }
  }
}
resource "aws_cloudwatch_log_group" "lambda_aws_cost_notify" {
  name              = "/aws/lambda/${var.aws_lambda_function_name}"
  retention_in_days = 3
}
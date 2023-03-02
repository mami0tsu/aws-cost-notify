resource "aws_lambda_function" "aws_cost_notify" {
  depends_on = [
    aws_cloudwatch_log_group.lambda_aws_cost_notify,
    null_resource.image_push,
  ]

  function_name = var.aws_lambda_function_name
  package_type  = "Image"
  image_uri     = "${aws_ecr_repository.aws_cost_notify.repository_url}:latest"
  role          = aws_iam_role.lambda_aws_cost_notify.arn
  publish       = true

  memory_size = 128
  timeout     = 3

  environment {
    variables = var.aws_lambda_environment
  }

  lifecycle {
    ignore_changes = [
      image_uri,
      last_modified,
    ]
  }
}

resource "aws_lambda_alias" "aws_cost_notify_prod" {
  name             = "prod"
  function_name    = aws_lambda_function.aws_cost_notify.arn
  function_version = aws_lambda_function.aws_cost_notify.version

  lifecycle {
    ignore_changes = [
      function_version,
    ]
  }
}

resource "aws_lambda_permission" "aws_cost_notify" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.aws_cost_notify.function_name
  principal     = "scheduler.amazonaws.com"
}
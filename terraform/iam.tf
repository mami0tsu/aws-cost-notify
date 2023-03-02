resource "aws_iam_role" "lambda_aws_cost_notify" {
  name               = var.aws_iam_role_name
  assume_role_policy = data.aws_iam_policy_document.lambda_aws_cost_notify_assume.json
}

data "aws_iam_policy_document" "lambda_aws_cost_notify_assume" {
  statement {
    effect = "Allow"

    actions = [
      "sts:AssumeRole",
    ]

    principals {
      type = "Service"
      identifiers = [
        "lambda.amazonaws.com",
      ]
    }
  }
}

resource "aws_iam_policy" "lambda_aws_cost_notify" {
  name   = var.aws_iam_policy_name
  policy = data.aws_iam_policy_document.lambda_aws_cost_notify_custom.json
}

data "aws_iam_policy_document" "lambda_aws_cost_notify_custom" {
  statement {
    effect = "Allow"

    actions = [
      "ce:GetCostAndUsage",
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]

    resources = [
      "*",
    ]
  }
}

resource "aws_iam_role_policy_attachment" "lambda_aws_cost_notify" {
  role       = aws_iam_role.lambda_aws_cost_notify.name
  policy_arn = aws_iam_policy.lambda_aws_cost_notify.arn
}

# EventBridge
resource "aws_iam_role" "eventbridge_aws_cost_notify" {
  name               = var.eventbridge_iam_role_name
  assume_role_policy = data.aws_iam_policy_document.eventbridge_aws_cost_notify_assume.json
}

data "aws_iam_policy_document" "eventbridge_aws_cost_notify_assume" {
  statement {
    effect = "Allow"

    actions = [
      "sts:AssumeRole",
    ]

    principals {
      type = "Service"
      identifiers = [
        "scheduler.amazonaws.com",
      ]
    }
  }
}

resource "aws_iam_policy" "eventbridge_aws_cost_notify" {
  name   = var.eventbridge_iam_policy_name
  policy = data.aws_iam_policy_document.eventbridge_aws_cost_notify_custom.json
}

data "aws_iam_policy_document" "eventbridge_aws_cost_notify_custom" {
  statement {
    effect = "Allow"

    actions = [
      "lambda:InvokeFunction",
    ]

    resources = [
      "*",
    ]
  }
}

resource "aws_iam_role_policy_attachment" "eventbridge_aws_cost_notify" {
  role       = aws_iam_role.eventbridge_aws_cost_notify.name
  policy_arn = aws_iam_policy.eventbridge_aws_cost_notify.arn
}
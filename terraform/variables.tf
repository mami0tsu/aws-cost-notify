variable "aws_ecr_repository_name" {
  type     = string
  nullable = false
}

variable "aws_lambda_function_name" {
  type     = string
  nullable = false
}

variable "aws_iam_role_name" {
  type     = string
  nullable = false
}

variable "aws_iam_policy_name" {
  type     = string
  nullable = false
}

variable "aws_lambda_environment" {
  type      = map(string)
  nullable  = false
  sensitive = true
}

variable "eventbridge_iam_policy_name" {
  type     = string
  nullable = false
}

variable "eventbridge_iam_role_name" {
  type     = string
  nullable = false
}

variable "eventbridge_schedule_name" {
  type     = string
  nullable = false
}

variable "schedule_group_name" {
  type     = string
  nullable = false
}
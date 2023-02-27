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
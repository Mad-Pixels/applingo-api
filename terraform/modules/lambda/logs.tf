resource "aws_cloudwatch_log_group" "logs" {
  name              = "/aws/lambda/${var.project}-${replace(var.function_name, "_", "-")}"
  retention_in_days = var.log_retention

  tags = var.shared_tags
}
resource "aws_iam_role_policy" "cloudwatch_metrics" {
  name = "${var.project}-${replace(var.function_name, "_", "-")}-cloudwatch-metrics"
  role = aws_iam_role.this.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "cloudwatch:PutMetricData",
          "cloudwatch:GetMetricStatistics",
          "cloudwatch:ListMetrics"
        ]
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_role_policy" "cloudwatch_logs" {
  name = "${var.project}-${replace(var.function_name, "_", "-")}-cloudwatch-logs"
  role = aws_iam_role.this.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_role_policy" "additional" {
  count  = var.policy != "" ? 1 : 0
  name   = "${var.project}-${replace(var.function_name, "_", "-")}-custom-lambda-policy"
  role   = aws_iam_role.this.id
  policy = var.policy
}

resource "aws_iam_role_policy_attachment" "base" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  role       = aws_iam_role.this.name
}

resource "aws_iam_role_policy" "cloudwatch_metrics" {
  name = "${var.project}-${var.function_name}-cloudwatch-metrics"
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

resource "aws_iam_role_policy" "additional" {
  count = var.policy != "" ? 1 : 0

  name   = "${var.project}-${var.function_name}-lambda-policy"
  role   = aws_iam_role.this.id
  policy = var.policy
}
resource "aws_iam_role_policy" "target_policy" {
  name = "${var.project}-${replace(var.scheduler_name, "_", "-")}-target-policy"
  role = aws_iam_role.scheduler.id

  policy = var.policy != "" ? var.policy : jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = var.target_type == "lambda" ? [
          "lambda:InvokeFunction"
          ] : var.target_type == "sqs" ? [
          "sqs:SendMessage"
          ] : var.target_type == "ecs" ? [
          "ecs:RunTask"
          ] : [
          "events:PutEvents"
        ]
        Resource = var.target_arn
      }
    ]
  })
}

resource "aws_iam_role_policy" "dead_letter_policy" {
  count = var.dead_letter_arn != null ? 1 : 0

  name = "${var.project}-${replace(var.scheduler_name, "_", "-")}-dlq-policy"
  role = aws_iam_role.scheduler.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "sqs:SendMessage"
        ]
        Resource = var.dead_letter_arn
      }
    ]
  })
}
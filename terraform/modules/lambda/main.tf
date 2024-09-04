resource "aws_lambda_function" "this" {
  function_name = "${var.project}-${var.function_name}"
  role          = aws_iam_role.this.arn
  image_uri     = var.image
  timeout       = var.timeout
  memory_size   = var.memory_size

  package_type  = "Image"
  architectures = ["arm64"]

  environment {
    variables = merge(
      var.environments,
      {
        LOG_LEVEL = var.log_level
      }
    )
  }

  dynamic "vpc_config" {
    for_each = var.vpc_config != null ? [var.vpc_config] : []
    content {
      subnet_ids         = vpc_config.value.subnet_ids
      security_group_ids = vpc_config.value.security_group_ids
    }
  }

  tags = merge(
    var.shared_tags,
    {
      "TF"      = "true",
      "Project" = var.project,
      "Github"  = "github.com/Mad-Pixels/lingocards-api",
    }
  )

  depends_on = [aws_cloudwatch_log_group.logs]
}

resource "aws_iam_role" "this" {
  name = "${var.project}-${var.function_name}-lambda-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })

  tags = merge(
    var.shared_tags,
    {
      "TF"      = "true",
      "Project" = var.project,
      "Github"  = "github.com/Mad-Pixels/lingocards-api",
    }
  )
}

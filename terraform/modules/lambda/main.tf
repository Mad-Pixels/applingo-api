resource "aws_lambda_function" "this" {
  function_name = "${var.project}-${replace(var.function_name, "_", "-")}"
  role          = aws_iam_role.this.arn
  image_uri     = var.image
  timeout       = var.timeout
  memory_size   = var.memory_size
  package_type  = "Image"
  architectures = [var.arch]

  environment {
    variables = merge(
      var.environments,
      {
        LOG_LEVEL        = var.log_level
        DEPLOY_TIMESTAMP = timestamp()
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
      "Type" = "image",
      "Arch" = var.arch,
    }
  )

  depends_on = [aws_cloudwatch_log_group.logs]
}

resource "aws_lambda_permission" "allow_all" {
  statement_id  = "AllowExecutionFromAnywhere"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.this.function_name
  principal     = "*"
}

resource "aws_iam_role" "this" {
  name = "${var.project}-${replace(var.function_name, "_", "-")}-lambda-role"

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

  tags = var.shared_tags
}

resource "aws_lambda_function_event_invoke_config" "this" {
  function_name                = aws_lambda_function.this.function_name
  maximum_event_age_in_seconds = var.maximum_event_age_in_seconds
  maximum_retry_attempts       = var.maximum_retry_attempts
}
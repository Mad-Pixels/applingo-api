variable "function_name" {
  description = "AWS Lambda function name"
  type        = string
}

variable "source_dir" {
  description = "Directory containing the Lambda function code"
  type        = string
}

resource "aws_iam_role" "this" {
  name = "${var.function_name}-role"

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
}

data "archive_file" "lambda_zip" {
  type        = "zip"
  source_dir  = var.source_dir
  output_path = "${path.module}/${var.function_name}.zip"
}

resource "aws_lambda_function" "function" {
  filename         = data.archive_file.lambda_zip.output_path
  function_name    = var.function_name
  role             = aws_iam_role.this.arn
  handler          = "main"
  source_code_hash = data.archive_file.lambda_zip.output_base64sha256
  runtime          = "go1.x"

  environment {
    variables = {}
  }
}

output "function_arn" {
  value = aws_lambda_function.function.arn
}
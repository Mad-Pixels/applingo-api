resource "aws_iam_role" "lambda_role" {
  name = "${var.name}-role"

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

resource "aws_iam_role_policy_attachment" "lambda_policy" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  role       = aws_iam_role.lambda_role.name
}

resource "aws_iam_role_policy" "additional_policy" {
  count  = var.additional_policy != "" ? 1 : 0
  name   = "${var.name}-additional-policy"
  role   = aws_iam_role.lambda_role.id
  policy = var.additional_policy
}

resource "aws_lambda_function" "container_function" {
  function_name = var.name
  role          = aws_iam_role.lambda_role.arn
  package_type  = "Image"
  image_uri     = var.image
  memory_size   = var.mem_size
  timeout       = var.timeout

  environment {
    variables = {
      CONTAINER_IMAGE = var.image
    }
  }
}
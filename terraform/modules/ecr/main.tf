resource "aws_ecr_repository" "this" {
  name                 = var.repository_name
  image_tag_mutability = var.image_tag_mutability
  
  image_scanning_configuration {
    scan_on_push = var.scan_on_push
  }
}

resource "aws_iam_policy" "ecr_policy" {
  name        = "${var.repository_name}-ecr-policy"
  description = "IAM policy to access the ECR repository"

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action   = [
            "ecr:GetDownloadUrlForLayer", 
            "ecr:BatchGetImage", 
            "ecr:BatchCheckLayerAvailability", 
            "ecr:PutImage", 
            "ecr:InitiateLayerUpload", 
            "ecr:UploadLayerPart", 
            "ecr:CompleteLayerUpload", 
            "ecr:DescribeImages"
        ],
        Effect   = "Allow",
        Resource = aws_ecr_repository.this.arn
      },
      {
        Action   = ["ecr:GetAuthorizationToken"],
        Effect   = "Allow",
        Resource = "*"
      }
    ]
  })
}
resource "aws_cloudwatch_log_group" "logs" {
  name              = "/aws/lambda/${var.project}-${var.function_name}"
  retention_in_days = var.log_retention

  tags = merge(
    var.shared_tags,
    {
      "TF"      = "true",
      "Project" = var.project,
      "Github"  = "github.com/Mad-Pixels/lingocards-api",
    }
  )
}
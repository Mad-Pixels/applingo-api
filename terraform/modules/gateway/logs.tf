resource "aws_cloudwatch_log_group" "this" {
  name              = "/aws/gateway/${var.project}-${var.api_name}-access-logs"
  retention_in_days = var.log_retention_days

  tags = merge(
    var.shared_tags,
    {
      "TF"      = "true",
      "Project" = var.project,
      "Github"  = "github.com/Mad-Pixels/lingocards-api",
    }
  )
}
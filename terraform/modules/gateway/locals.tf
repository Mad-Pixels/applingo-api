data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

locals {
  manifest = templatefile("../../../openapi-interface/.tmpl/openapi.yaml", {
    project        = var.project
    name           = var.api_name
    use_localstack = var.use_localstack
    region         = data.aws_region.current.name
    account_id     = try(data.aws_caller_identity.current.account_id, "*")

    api_subcategories = var.invoke_lambdas_arns["api-subcategories"].arn
    api_dictionaries  = var.invoke_lambdas_arns["api-dictionaries"].arn
    api_reports       = var.invoke_lambdas_arns["api-reports"].arn
    api_profile       = var.invoke_lambdas_arns["api-profile"].arn
    api_levels        = var.invoke_lambdas_arns["api-levels"].arn
    api_schema        = var.invoke_lambdas_arns["api-schema"].arn
    api_urls          = var.invoke_lambdas_arns["api-urls"].arn
    authorizer        = var.invoke_lambdas_arns["authorizer"].arn
  })
}
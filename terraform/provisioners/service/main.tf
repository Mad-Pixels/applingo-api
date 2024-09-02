terraform {}

provider "aws" {
  region = var.aws_region

  dynamic "endpoints" {
    for_each = var.use_localstack ? [1] : []
    content {
      lambda = var.localstack_endpoint
      iam    = var.localstack_endpoint
    }
  }

  skip_credentials_validation = var.use_localstack
  skip_metadata_api_check     = var.use_localstack
  skip_requesting_account_id  = var.use_localstack

  access_key = var.use_localstack ? "test" : null
  secret_key = var.use_localstack ? "test" : null
}

module "lambda_functions" {
    source   = "../../modules/lambda"
    for_each = local.lambda_functions

    name              = each.key
    #image             = "${data.terraform_remote_state.ecr.outputs.repository_url}:${each.key}"
    image = "000000000000.dkr.ecr.us-east-1.localhost.localstack.cloud:4566/lingocards-api:dictionary"
    mem_size          = try(local.lambda_configs[each.key].memory_size, 128)
    timeout           = try(local.lambda_configs[each.key].timeout, 30)
    additional_policy = try(local.lambda_configs[each.key].policy, "")
}
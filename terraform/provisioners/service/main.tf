terraform {}

data "terraform_remote_state" "ecr" {
  backend = var.use_localstack ? "local" : "s3"

  config = var.use_localstack ? {
    path = "../infra/terraform.tfstate"
  } : {
    bucket = local.state_bucket
    key    = local.tfstate_file
    region = var.aws_region
  }
}

provider "aws" {
  region = var.aws_region

  dynamic "endpoints" {
    for_each = var.use_localstack ? [1] : []
    content {
      lambda     = var.localstack_endpoint
      iam        = var.localstack_endpoint
      logs       = var.localstack_endpoint
      cloudwatch = var.localstack_endpoint
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
  for_each = local.lambdas

  function_name = each.key
  project       = "lingocards"
  image         = "${data.terraform_remote_state.ecr.outputs.repository_url}:${each.key}"
  memory_size   = try(each.value.memory_size, 128)
  timeout       = try(each.value.timeout, 30)
  policy        = try(jsonencode(each.value.policy), "")
}
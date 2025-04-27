data "aws_caller_identity" "current" {}

data "terraform_remote_state" "infra" {
  backend = var.use_localstack ? "local" : "s3"

  config = var.use_localstack ? {
    path = "../infra/terraform.tfstate"
    } : {
    bucket = var.infra_backend_bucket
    region = var.infra_backend_region
    key    = var.infra_backend_key
  }
}

module "lambda_functions" {
  source = "../../modules/lambda"

  for_each      = local.lambdas
  function_name = each.key
  project       = local.project
  image         = "${data.terraform_remote_state.infra.outputs.ecr-repository-api_url}:${each.key}"
  log_level     = var.use_localstack ? "DEBUG" : "ERROR"
  arch          = var.arch

  environments = try(each.value.envs, {})
  timeout      = try(each.value.timeout, 3)
  memory_size  = try(each.value.memory_size, 128)
  policy       = try(jsonencode(each.value.policy), "")
}

module "gateway" {
  source = "../../modules/gateway"

  api_name       = "api"
  project        = local.project
  use_localstack = var.use_localstack

  invoke_lambdas_arns = {
    for name, lambda in module.lambda_functions : name => {
      arn  = lambda.function_arn
      name = lambda.function_name
    }
  }
  depends_on = [module.lambda_functions]
}

resource "aws_lambda_event_source_mapping" "dynamo-stream-processing" {
  event_source_arn       = local.template_vars.processing_table_stream_arn
  function_name          = module.lambda_functions["trigger-processing-check"].function_arn
  starting_position      = "LATEST"
  maximum_retry_attempts = 0

  depends_on = [module.lambda_functions]
}

resource "aws_lambda_event_source_mapping" "dynamo-stream-dictionary" {
  event_source_arn       = local.template_vars.dictionary_table_stream_arn
  function_name          = module.lambda_functions["trigger-dictionary-check"].function_arn
  starting_position      = "LATEST"
  maximum_retry_attempts = 0

  depends_on = [module.lambda_functions]
}
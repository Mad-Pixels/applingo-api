data "terraform_remote_state" "infra" {
  backend = var.use_localstack ? "local" : "s3"

  config = var.use_localstack ? {
    path = "../infra/terraform.tfstate"
    } : {
    bucket = local.state_bucket
    key    = local.tfstate_file
    region = "eu-central-1"
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

resource "aws_lambda_event_source_mapping" "dynamo-queue" {
  event_source_arn              = local.template_vars.dictionary_table_stream_arn
  function_name                 = module.lambda_functions["trigger-dynamo-to-sqs-put-csv"].function_arn
  starting_position             = "LATEST"
  maximum_retry_attempts        = 3
  maximum_record_age_in_seconds = 120

  depends_on = [module.lambda_functions]
}

resource "aws_lambda_event_source_mapping" "queue-put-csv" {
  event_source_arn = local.template_vars.put_csv_sqs_queue_arn
  function_name    = module.lambda_functions["trigger-sqs-to-job-put-csv"].function_arn

  depends_on = [module.lambda_functions]
}

module "gateway" {
  source = "../../modules/gateway"

  project        = local.project
  api_name       = "api"
  use_localstack = var.use_localstack

  invoke_lambdas_arns = {
    for name, lambda in module.lambda_functions : name => {
      arn  = lambda.function_arn
      name = lambda.function_name
    }
  }
  depends_on = [module.lambda_functions]
}

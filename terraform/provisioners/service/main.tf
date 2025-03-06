data "aws_caller_identity" "current" {}

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

module "scheduler_events" {
  source = "../../modules/scheduler"

  for_each       = local.schedulers
  project        = local.project
  scheduler_name = each.key

  schedule_expression          = try(each.value.Config.schedule_expression, "rate(1 day)")
  flexible_time_window_mode    = try(each.value.Config.flexible_time_window_mode, "OFF")
  maximum_window_in_minutes    = try(each.value.Config.maximum_window_in_minutes, 5)
  maximum_retry_attempts       = try(each.value.Config.maximum_retry_attempts, null)
  maximum_event_age_in_seconds = try(each.value.Config.maximum_event_age_in_seconds, 3600)

  target_arn  = format(local.lambda_arn_template, each.value.Config.target_lambda_name)
  target_type = try(each.value.Config.target_type, "lambda")
  policy      = try(each.value.Config.policy != null ? jsonencode(each.value.Config.policy) : "", "")

  input_json = jsonencode({ Records = each.value.Records })
  depends_on = [module.lambda_functions]
}

resource "aws_lambda_event_source_mapping" "dynamo-stream-processing" {
  event_source_arn       = local.template_vars.processing_table_stream_arn
  function_name          = module.lambda_functions["scheduler-dictionary-craft"].function_arn
  starting_position      = "LATEST"
  maximum_retry_attempts = 0

  depends_on = [module.lambda_functions]
}

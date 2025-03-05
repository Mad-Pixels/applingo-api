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

resource "aws_lambda_event_source_mapping" "dynamo-stream-processing" {
  event_source_arn              = local.template_vars.processing_table_stream_arn
  function_name                 = module.lambda_functions["trigger-dictionary-check"].function_arn
  starting_position             = "LATEST"
  maximum_retry_attempts        = 0

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

// TODO: SQS was removed from the project, use directly the table stream.
# resource "aws_lambda_event_source_mapping" "queue-put-csv" {
#   event_source_arn = local.template_vars.put_csv_sqs_queue_arn
#   function_name    = module.lambda_functions["trigger-sqs-to-job-put-csv"].function_arn

#   depends_on = [module.lambda_functions]
# }


module "schedulers" {
  source = "../../modules/scheduler"
  
  for_each = local.schedulers
  
  project         = local.project
  scheduler_name  = each.key
  schedule_expression = lookup(each.value.Config, "schedule_expression", "rate(1 day)")
  
  flexible_time_window_mode = lookup(each.value.Config, "flexible_time_window_mode", "OFF")
  maximum_window_in_minutes = lookup(each.value.Config, "maximum_window_in_minutes", 5)
  
  maximum_retry_attempts = lookup(each.value.Config, "maximum_retry_attempts", null)
  maximum_event_age_in_seconds = lookup(each.value.Config, "maximum_event_age_in_seconds", 3600)
  
  target_type = lookup(each.value.Config, "target_type", "lambda")
  target_arn  = lookup(each.value.Config, "target_arn", module.lambda_functions["forge-dictionary"].function_arn)
  input_json  = jsonencode({ Records = each.value.Records })
  
  // Опциональные параметры для различных типов целей
  dead_letter_arn = lookup(each.value.Config, "dead_letter_arn", null)
  sqs_message_group_id = lookup(each.value.Config, "sqs_message_group_id", null)
  
  ecs_task_definition_arn = lookup(each.value.Config, "ecs_task_definition_arn", null)
  ecs_launch_type = lookup(each.value.Config, "ecs_launch_type", "FARGATE")
  ecs_network_config = lookup(each.value.Config, "ecs_network_config", null)
  
  policy = lookup(each.value.Config, "policy", "")
  
}
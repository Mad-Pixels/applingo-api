data "terraform_remote_state" "infra" {
  backend = var.use_localstack ? "local" : "s3"

  config = var.use_localstack ? {
    path = "../infra/terraform.tfstate"
    } : {
    bucket = local.state_bucket
    key    = local.tfstate_file
    region = var.aws_region
  }
}

module "lambda_functions" {
  source   = "../../modules/lambda"
  for_each = local.lambdas

  function_name = each.key
  project       = local.project
  image         = "${data.terraform_remote_state.infra.outputs.ecr-repository-api_url}:${each.key}"
  log_level     = var.use_localstack ? "DEBUG" : "ERROR"
  arch          = var.arch

  timeout     = try(each.value.timeout, 3)
  memory_size = try(each.value.memory_size, 128)
  policy      = try(jsonencode(each.value.policy), "")

  environments = {
    SERVICE_DICTIONARY_BUCKET : data.terraform_remote_state.infra.outputs.s3-dictionary-bucket_name,
    SERVICE_PROCESSING_BUCKET : data.terraform_remote_state.infra.outputs.s3-processing-bucket_name,
    SERVICE_DICTIONARY_DYNAMO : data.terraform_remote_state.infra.outputs.dynamo-dictionary-table_name
  }
}

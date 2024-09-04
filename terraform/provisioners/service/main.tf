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
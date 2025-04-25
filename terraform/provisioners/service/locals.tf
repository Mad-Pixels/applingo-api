# Base constants.
locals {
  project = "applingo"

}

# Base variables.
locals {
  template_vars = {
    var_jwt_secret              = var.jwt_secret
    var_openai_key              = var.openai_key
    var_device_api_token        = var.device_api_token
    log_errors_bucket_name      = data.terraform_remote_state.infra.outputs.s3-errors-bucket_name
    forge_bucket_name           = data.terraform_remote_state.infra.outputs.s3-forge-bucket_name
    forge_bucket_arn            = data.terraform_remote_state.infra.outputs.s3-forge-bucket_arn
    dictionary_bucket_name      = data.terraform_remote_state.infra.outputs.s3-dictionary-bucket_name
    processing_bucket_name      = data.terraform_remote_state.infra.outputs.s3-processing-bucket_name
    log_errors_bucket_arn       = data.terraform_remote_state.infra.outputs.s3-errors-bucket_arn
    dictionary_bucket_arn       = data.terraform_remote_state.infra.outputs.s3-dictionary-bucket_arn
    processing_bucket_arn       = data.terraform_remote_state.infra.outputs.s3-processing-bucket_arn
    dictionary_table_arn        = data.terraform_remote_state.infra.outputs.dynamo-dictionary-table_arn
    dictionary_table_stream_arn = data.terraform_remote_state.infra.outputs.dynamo-dictionary-stream_arn
    processing_table_arn        = data.terraform_remote_state.infra.outputs.dynamo-processing-table_arn
    processing_table_stream_arn = data.terraform_remote_state.infra.outputs.dynamo-processing-stream_arn
    profile_table_arn           = data.terraform_remote_state.infra.outputs.dynamo-profile-table_arn
  }
}

# Custom variables.
locals {
  lambda_arn_template = "arn:aws:lambda:${var.aws_region}:${data.aws_caller_identity.current.account_id}:function:${local.project}-%s"
}

# Lambdas configs.
locals {
  _root_dir = "${path.module}/../../../cmd"
  _entries  = fileset(local._root_dir, "**")

  _lambda_functions = distinct([
    for v in local._entries : split("/", v)[0]
    if length(split("/", v)) > 1
  ])

  _lambda_configs = {
    for func in local._lambda_functions :
    func => fileexists("${local._root_dir}/${func}/.infra/config.json") ?
    jsondecode(
      templatefile("${local._root_dir}/${func}/.infra/config.json", local.template_vars)
    ) : null
  }

  lambdas          = { for func in local._lambda_functions : func => local._lambda_configs[func] }
  lambda_functions = keys(local.lambdas)

  api_lambdas      = [for name in local.lambda_functions : name if startswith(name, "api")]
  trigger_lambdas  = [for name in local.lambda_functions : name if startswith(name, "trigger")]
  schedule_lambdas = [for name in local.lambda_functions : name if startswith(name, "scheduler")]
}

# Schedulers configs.
locals {
  _all_schedulers = flatten([
    for fn in local.schedule_lambdas : [
      for file in fileset("${local._root_dir}/${fn}/.infra", "scheduler*.json") : {
        function_name  = fn
        scheduler_name = replace(file, ".json", "")
        short_name     = replace(replace(file, "scheduler_", ""), ".json", "")
        config = jsondecode(
          templatefile("${local._root_dir}/${fn}/.infra/${file}",
            merge(local.template_vars, { target_lambda_name = "${fn}" })
          )
        )
      }
    ]
  ])

  schedulers = {
    for scheduler in local._all_schedulers :
    "${scheduler.function_name}-${scheduler.short_name}" => scheduler
  }
}

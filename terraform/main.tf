module "lambda_functions" {
  for_each = local.lambda_functions
  source   = "./modules/lambda"

  function_name = each.key
  source_dir    = "${path.root}/../cmd/${each.key}"
}
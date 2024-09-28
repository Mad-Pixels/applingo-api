output "lambda_functions" {
  value = {
    for name, lambda in module.lambda_functions : name => {
      arn  = lambda.function_arn
      name = lambda.function_name
    }
  }
}

output "execution_gateway_arn" {
  value = module.gateway.api_gateway_execution_arn
}
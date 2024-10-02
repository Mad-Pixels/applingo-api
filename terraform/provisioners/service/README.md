<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | ~> 5.65.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_terraform"></a> [terraform](#provider\_terraform) | n/a |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_gateway"></a> [gateway](#module\_gateway) | ../../modules/gateway | n/a |
| <a name="module_lambda_functions"></a> [lambda\_functions](#module\_lambda\_functions) | ../../modules/lambda | n/a |

## Resources

| Name | Type |
|------|------|
| [terraform_remote_state.infra](https://registry.terraform.io/providers/hashicorp/terraform/latest/docs/data-sources/remote_state) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_arch"></a> [arch](#input\_arch) | Set architecture which will be use in lambda services | `string` | n/a | yes |
| <a name="input_aws_region"></a> [aws\_region](#input\_aws\_region) | AWS region | `string` | n/a | yes |
| <a name="input_device_api_token"></a> [device\_api\_token](#input\_device\_api\_token) | Token which use for lambda request validate from device | `string` | n/a | yes |
| <a name="input_localstack_endpoint"></a> [localstack\_endpoint](#input\_localstack\_endpoint) | LocalStack endpoint | `string` | `"https://localhost.localstack.cloud:4566"` | no |
| <a name="input_use_localstack"></a> [use\_localstack](#input\_use\_localstack) | Whether to use LocalStack | `bool` | `false` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_execution_gateway_arn"></a> [execution\_gateway\_arn](#output\_execution\_gateway\_arn) | n/a |
| <a name="output_lambda_functions"></a> [lambda\_functions](#output\_lambda\_functions) | n/a |
<!-- END_TF_DOCS -->
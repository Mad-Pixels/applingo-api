<!-- BEGIN_TF_DOCS -->
## Requirements

No requirements.

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | n/a |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_api_gateway_account.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_account) | resource |
| [aws_api_gateway_base_path_mapping.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_base_path_mapping) | resource |
| [aws_api_gateway_deployment.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_deployment) | resource |
| [aws_api_gateway_domain_name.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_domain_name) | resource |
| [aws_api_gateway_method_settings.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_method_settings) | resource |
| [aws_api_gateway_rest_api.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_rest_api) | resource |
| [aws_api_gateway_stage.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_stage) | resource |
| [aws_cloudwatch_log_group.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_log_group) | resource |
| [aws_iam_role.cloudwatch_role](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role) | resource |
| [aws_iam_role_policy.cloudwatch_policy](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy) | resource |
| [aws_wafv2_web_acl_association.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/wafv2_web_acl_association) | resource |
| [aws_caller_identity.current](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/caller_identity) | data source |
| [aws_region.current](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/region) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_acm_certificate_arn"></a> [acm\_certificate\_arn](#input\_acm\_certificate\_arn) | AWS ACM certificate for api gateway custom domain | `string` | `""` | no |
| <a name="input_api_name"></a> [api\_name](#input\_api\_name) | Name of the API Gateway | `string` | n/a | yes |
| <a name="input_custom_domain"></a> [custom\_domain](#input\_custom\_domain) | GatewayAPI custom domain name | `string` | `""` | no |
| <a name="input_invoke_lambdas_arns"></a> [invoke\_lambdas\_arns](#input\_invoke\_lambdas\_arns) | Map of Lambda function ARNs and names | <pre>map(object({<br>    arn  = string<br>    name = string<br>  }))</pre> | n/a | yes |
| <a name="input_log_retention_days"></a> [log\_retention\_days](#input\_log\_retention\_days) | Optional set log retention | `number` | `3` | no |
| <a name="input_project"></a> [project](#input\_project) | Project name | `string` | n/a | yes |
| <a name="input_shared_tags"></a> [shared\_tags](#input\_shared\_tags) | Tags to add to all resources | `map` | `{}` | no |
| <a name="input_stage_name"></a> [stage\_name](#input\_stage\_name) | Option set API Gateway stage name | `string` | `"live"` | no |
| <a name="input_use_localstack"></a> [use\_localstack](#input\_use\_localstack) | Whether to use LocalStack | `bool` | `false` | no |
| <a name="input_wafv2_web_acl_arn"></a> [wafv2\_web\_acl\_arn](#input\_wafv2\_web\_acl\_arn) | WAF web acl arn | `string` | `""` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_api_gateway_execution_arn"></a> [api\_gateway\_execution\_arn](#output\_api\_gateway\_execution\_arn) | n/a |
<!-- END_TF_DOCS -->
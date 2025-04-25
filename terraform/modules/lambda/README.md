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
| [aws_cloudwatch_log_group.logs](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_log_group) | resource |
| [aws_iam_role.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role) | resource |
| [aws_iam_role_policy.additional](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy) | resource |
| [aws_iam_role_policy.cloudwatch_logs](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy) | resource |
| [aws_iam_role_policy.cloudwatch_metrics](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy) | resource |
| [aws_lambda_function.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_function) | resource |
| [aws_lambda_function_event_invoke_config.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_function_event_invoke_config) | resource |
| [aws_lambda_permission.allow_all](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_permission) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_arch"></a> [arch](#input\_arch) | Optional set architecture for Lambda function | `string` | `"arm64"` | no |
| <a name="input_environments"></a> [environments](#input\_environments) | Additional ENVs for the Lambda function | `map` | `{}` | no |
| <a name="input_function_name"></a> [function\_name](#input\_function\_name) | Name of the Lambda function | `string` | n/a | yes |
| <a name="input_image"></a> [image](#input\_image) | URI of the container image in ECR | `string` | n/a | yes |
| <a name="input_log_level"></a> [log\_level](#input\_log\_level) | The log level for the Lambda function | `string` | `"ERROR"` | no |
| <a name="input_log_retention"></a> [log\_retention](#input\_log\_retention) | Cloudwatch retention in days | `number` | `3` | no |
| <a name="input_maximum_event_age_in_seconds"></a> [maximum\_event\_age\_in\_seconds](#input\_maximum\_event\_age\_in\_seconds) | Maximum event age in seconds | `number` | `600` | no |
| <a name="input_maximum_retry_attempts"></a> [maximum\_retry\_attempts](#input\_maximum\_retry\_attempts) | Maximum retry attempts | `number` | `0` | no |
| <a name="input_memory_size"></a> [memory\_size](#input\_memory\_size) | Amount of memory in MB for the Lambda function | `number` | `128` | no |
| <a name="input_policy"></a> [policy](#input\_policy) | Additional IAM policy for the Lambda function | `string` | `""` | no |
| <a name="input_project"></a> [project](#input\_project) | Project name | `string` | n/a | yes |
| <a name="input_shared_tags"></a> [shared\_tags](#input\_shared\_tags) | Tags to add to all resources | `map` | `{}` | no |
| <a name="input_timeout"></a> [timeout](#input\_timeout) | Timeout for the Lambda function in seconds | `number` | `5` | no |
| <a name="input_vpc_config"></a> [vpc\_config](#input\_vpc\_config) | Optional VPC configuration for the Lambda function | <pre>object({<br>    subnet_ids         = list(string)<br>    security_group_ids = list(string)<br>  })</pre> | `null` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_function_arn"></a> [function\_arn](#output\_function\_arn) | ARN of the Lambda function |
| <a name="output_function_name"></a> [function\_name](#output\_function\_name) | Name of the Lambda function |
<!-- END_TF_DOCS -->
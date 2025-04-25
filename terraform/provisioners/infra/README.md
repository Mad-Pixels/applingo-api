<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | ~> 5.65.0 |

## Providers

No providers.

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_dynamo-dictionary-table"></a> [dynamo-dictionary-table](#module\_dynamo-dictionary-table) | ../../modules/dynamo | n/a |
| <a name="module_dynamo-processing-table"></a> [dynamo-processing-table](#module\_dynamo-processing-table) | ../../modules/dynamo | n/a |
| <a name="module_ecr-repository-api"></a> [ecr-repository-api](#module\_ecr-repository-api) | ../../modules/ecr | n/a |
| <a name="module_s3-dictionary-bucket"></a> [s3-dictionary-bucket](#module\_s3-dictionary-bucket) | ../../modules/s3 | n/a |
| <a name="module_s3-errors-bucket"></a> [s3-errors-bucket](#module\_s3-errors-bucket) | ../../modules/s3 | n/a |
| <a name="module_s3-forge-bucket"></a> [s3-forge-bucket](#module\_s3-forge-bucket) | ../../modules/s3 | n/a |
| <a name="module_s3-processing-bucket"></a> [s3-processing-bucket](#module\_s3-processing-bucket) | ../../modules/s3 | n/a |

## Resources

No resources.

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_arch"></a> [arch](#input\_arch) | Set architecture which will be use in lambda services | `string` | n/a | yes |
| <a name="input_aws_region"></a> [aws\_region](#input\_aws\_region) | AWS region | `string` | n/a | yes |
| <a name="input_localstack_endpoint"></a> [localstack\_endpoint](#input\_localstack\_endpoint) | LocalStack endpoint | `string` | `"https://localhost.localstack.cloud:4566"` | no |
| <a name="input_use_localstack"></a> [use\_localstack](#input\_use\_localstack) | Whether to use LocalStack | `bool` | `false` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_dynamo-dictionary-stream_arn"></a> [dynamo-dictionary-stream\_arn](#output\_dynamo-dictionary-stream\_arn) | n/a |
| <a name="output_dynamo-dictionary-table_arn"></a> [dynamo-dictionary-table\_arn](#output\_dynamo-dictionary-table\_arn) | n/a |
| <a name="output_dynamo-dictionary-table_name"></a> [dynamo-dictionary-table\_name](#output\_dynamo-dictionary-table\_name) | n/a |
| <a name="output_dynamo-processing-stream_arn"></a> [dynamo-processing-stream\_arn](#output\_dynamo-processing-stream\_arn) | n/a |
| <a name="output_dynamo-processing-table_arn"></a> [dynamo-processing-table\_arn](#output\_dynamo-processing-table\_arn) | n/a |
| <a name="output_dynamo-processing-table_name"></a> [dynamo-processing-table\_name](#output\_dynamo-processing-table\_name) | n/a |
| <a name="output_ecr-repository-api_url"></a> [ecr-repository-api\_url](#output\_ecr-repository-api\_url) | n/a |
| <a name="output_s3-dictionary-bucket_arn"></a> [s3-dictionary-bucket\_arn](#output\_s3-dictionary-bucket\_arn) | n/a |
| <a name="output_s3-dictionary-bucket_name"></a> [s3-dictionary-bucket\_name](#output\_s3-dictionary-bucket\_name) | n/a |
| <a name="output_s3-errors-bucket_arn"></a> [s3-errors-bucket\_arn](#output\_s3-errors-bucket\_arn) | n/a |
| <a name="output_s3-errors-bucket_name"></a> [s3-errors-bucket\_name](#output\_s3-errors-bucket\_name) | n/a |
| <a name="output_s3-forge-bucket_arn"></a> [s3-forge-bucket\_arn](#output\_s3-forge-bucket\_arn) | n/a |
| <a name="output_s3-forge-bucket_name"></a> [s3-forge-bucket\_name](#output\_s3-forge-bucket\_name) | n/a |
| <a name="output_s3-processing-bucket_arn"></a> [s3-processing-bucket\_arn](#output\_s3-processing-bucket\_arn) | n/a |
| <a name="output_s3-processing-bucket_name"></a> [s3-processing-bucket\_name](#output\_s3-processing-bucket\_name) | n/a |
<!-- END_TF_DOCS -->
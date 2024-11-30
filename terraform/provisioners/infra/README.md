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
| <a name="module_dictionary_put_csv_queue"></a> [dictionary\_put\_csv\_queue](#module\_dictionary\_put\_csv\_queue) | ../../modules/sqs | n/a |
| <a name="module_dynamo-dictionary-table"></a> [dynamo-dictionary-table](#module\_dynamo-dictionary-table) | ../../modules/dynamo | n/a |
| <a name="module_dynamo-level-table"></a> [dynamo-level-table](#module\_dynamo-level-table) | ../../modules/dynamo | n/a |
| <a name="module_dynamo-subcategory-table"></a> [dynamo-subcategory-table](#module\_dynamo-subcategory-table) | ../../modules/dynamo | n/a |
| <a name="module_ecr-repository-api"></a> [ecr-repository-api](#module\_ecr-repository-api) | ../../modules/ecr | n/a |
| <a name="module_s3-dictionary-bucket"></a> [s3-dictionary-bucket](#module\_s3-dictionary-bucket) | ../../modules/s3 | n/a |
| <a name="module_s3-errors-bucket"></a> [s3-errors-bucket](#module\_s3-errors-bucket) | ../../modules/s3 | n/a |
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
| <a name="output_dynamo-level-table_arn"></a> [dynamo-level-table\_arn](#output\_dynamo-level-table\_arn) | n/a |
| <a name="output_dynamo-level-table_name"></a> [dynamo-level-table\_name](#output\_dynamo-level-table\_name) | n/a |
| <a name="output_dynamo-subcategory-table_arn"></a> [dynamo-subcategory-table\_arn](#output\_dynamo-subcategory-table\_arn) | n/a |
| <a name="output_dynamo-subcategory-table_name"></a> [dynamo-subcategory-table\_name](#output\_dynamo-subcategory-table\_name) | n/a |
| <a name="output_ecr-repository-api_url"></a> [ecr-repository-api\_url](#output\_ecr-repository-api\_url) | n/a |
| <a name="output_s3-dictionary-bucket_arn"></a> [s3-dictionary-bucket\_arn](#output\_s3-dictionary-bucket\_arn) | n/a |
| <a name="output_s3-dictionary-bucket_name"></a> [s3-dictionary-bucket\_name](#output\_s3-dictionary-bucket\_name) | n/a |
| <a name="output_s3-errors-bucket_arn"></a> [s3-errors-bucket\_arn](#output\_s3-errors-bucket\_arn) | n/a |
| <a name="output_s3-errors-bucket_name"></a> [s3-errors-bucket\_name](#output\_s3-errors-bucket\_name) | n/a |
| <a name="output_s3-processing-bucket_arn"></a> [s3-processing-bucket\_arn](#output\_s3-processing-bucket\_arn) | n/a |
| <a name="output_s3-processing-bucket_name"></a> [s3-processing-bucket\_name](#output\_s3-processing-bucket\_name) | n/a |
| <a name="output_sqs-put-csv-dead-letter-queue_arn"></a> [sqs-put-csv-dead-letter-queue\_arn](#output\_sqs-put-csv-dead-letter-queue\_arn) | n/a |
| <a name="output_sqs-put-csv-dead-letter-queue_url"></a> [sqs-put-csv-dead-letter-queue\_url](#output\_sqs-put-csv-dead-letter-queue\_url) | n/a |
| <a name="output_sqs-put-csv-queue_arn"></a> [sqs-put-csv-queue\_arn](#output\_sqs-put-csv-queue\_arn) | n/a |
| <a name="output_sqs-put-csv-queue_url"></a> [sqs-put-csv-queue\_url](#output\_sqs-put-csv-queue\_url) | n/a |
<!-- END_TF_DOCS -->
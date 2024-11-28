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
| [aws_dynamodb_table.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/dynamodb_table) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_attributes"></a> [attributes](#input\_attributes) | List of nested attribute definitions. Only required for hash\_key and range\_key attributes | <pre>list(object({<br>    name = string<br>    type = string<br>  }))</pre> | n/a | yes |
| <a name="input_billing_mode"></a> [billing\_mode](#input\_billing\_mode) | Controls how you are charged for read and write throughput and how you manage capacity | `string` | `"PAY_PER_REQUEST"` | no |
| <a name="input_hash_key"></a> [hash\_key](#input\_hash\_key) | The attribute to use as the hash (partition) key | `string` | n/a | yes |
| <a name="input_project"></a> [project](#input\_project) | Project name | `string` | n/a | yes |
| <a name="input_range_key"></a> [range\_key](#input\_range\_key) | The attribute to use as the range (sort) key | `string` | `null` | no |
| <a name="input_secondary_index_list"></a> [secondary\_index\_list](#input\_secondary\_index\_list) | List of global secondary indexes | <pre>list(object({<br>    name               = string<br>    hash_key           = string<br>    range_key          = optional(string)<br>    projection_type    = string<br>    non_key_attributes = optional(list(string))<br>    read_capacity      = optional(number)<br>    write_capacity     = optional(number)<br>  }))</pre> | `null` | no |
| <a name="input_shared_tags"></a> [shared\_tags](#input\_shared\_tags) | Tags to add to all resources | `map` | `{}` | no |
| <a name="input_stream_enabled"></a> [stream\_enabled](#input\_stream\_enabled) | On/Off dynamo stream | `bool` | `false` | no |
| <a name="input_stream_type"></a> [stream\_type](#input\_stream\_type) | Type of streaming | `string` | `"NEW_IMAGE"` | no |
| <a name="input_table_name"></a> [table\_name](#input\_table\_name) | The name of the DynamoDB table | `string` | n/a | yes |
| <a name="input_ttl_attribute_name"></a> [ttl\_attribute\_name](#input\_ttl\_attribute\_name) | The name of the TTL attribute | `string` | `"ttl"` | no |
| <a name="input_ttl_enabled"></a> [ttl\_enabled](#input\_ttl\_enabled) | Whether to enable TTL for the DynamoDB table | `bool` | `false` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_stream_arn"></a> [stream\_arn](#output\_stream\_arn) | The ARM of the DynamoDB streaming |
| <a name="output_table_arn"></a> [table\_arn](#output\_table\_arn) | The ARN of the DynamoDB table |
| <a name="output_table_name"></a> [table\_name](#output\_table\_name) | The name of the DynamoDB table |
<!-- END_TF_DOCS -->
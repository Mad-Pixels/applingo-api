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
| [aws_s3_bucket.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket) | resource |
| [aws_s3_bucket_server_side_encryption_configuration.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_server_side_encryption_configuration) | resource |
| [aws_s3_bucket_versioning.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_versioning) | resource |
| [aws_s3_bucket_website_configuration.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_website_configuration) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_bucket_name"></a> [bucket\_name](#input\_bucket\_name) | Name of the S3 bucket. | `string` | n/a | yes |
| <a name="input_enable_versioning"></a> [enable\_versioning](#input\_enable\_versioning) | Enables versioning for the S3 bucket if set true. | `bool` | `false` | no |
| <a name="input_error_document"></a> [error\_document](#input\_error\_document) | Name of the error document for the S3 static web site. | `string` | `"error.html"` | no |
| <a name="input_force_destroy"></a> [force\_destroy](#input\_force\_destroy) | Allows Terraform to delete the bucket when removing the resource of set true. | `bool` | `true` | no |
| <a name="input_index_document"></a> [index\_document](#input\_index\_document) | Name of the index document for the S3 static web site. | `string` | `"index.html"` | no |
| <a name="input_is_website"></a> [is\_website](#input\_is\_website) | Specifies if the S3 bucket will host a static website when set to true. | `bool` | `false` | no |
| <a name="input_project"></a> [project](#input\_project) | Project name | `string` | n/a | yes |
| <a name="input_shared_tags"></a> [shared\_tags](#input\_shared\_tags) | Tags to add to all resources | `map` | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_s3_arn"></a> [s3\_arn](#output\_s3\_arn) | n/a |
| <a name="output_s3_domain"></a> [s3\_domain](#output\_s3\_domain) | n/a |
| <a name="output_s3_id"></a> [s3\_id](#output\_s3\_id) | n/a |
| <a name="output_s3_name"></a> [s3\_name](#output\_s3\_name) | n/a |
<!-- END_TF_DOCS -->
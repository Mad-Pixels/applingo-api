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
| <a name="module_vpc-infra"></a> [vpc-infra](#module\_vpc-infra) | ../../modules/vpc | n/a |

## Resources

No resources.

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_aws_region"></a> [aws\_region](#input\_aws\_region) | AWS region | `string` | n/a | yes |
| <a name="input_localstack_endpoint"></a> [localstack\_endpoint](#input\_localstack\_endpoint) | LocalStack endpoint | `string` | `"https://localhost.localstack.cloud:4566"` | no |
| <a name="input_use_localstack"></a> [use\_localstack](#input\_use\_localstack) | Whether to use LocalStack | `bool` | `false` | no |

## Outputs

No outputs.
<!-- END_TF_DOCS -->
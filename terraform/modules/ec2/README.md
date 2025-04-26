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
| [aws_instance.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/instance) | resource |
| [aws_ami.amazon_linux](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/ami) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_ami_id"></a> [ami\_id](#input\_ami\_id) | AMI ID for the instance | `string` | `""` | no |
| <a name="input_associate_public_ip_address"></a> [associate\_public\_ip\_address](#input\_associate\_public\_ip\_address) | Associate public ip address (ipv4) | `bool` | `false` | no |
| <a name="input_graviton_size"></a> [graviton\_size](#input\_graviton\_size) | Instance size (e.g. micro, small, medium) | `string` | `"nano"` | no |
| <a name="input_key_name"></a> [key\_name](#input\_key\_name) | Name of the EC2 Key Pair for SSH access | `string` | `""` | no |
| <a name="input_name"></a> [name](#input\_name) | Instance name | `string` | n/a | yes |
| <a name="input_security_group_ids"></a> [security\_group\_ids](#input\_security\_group\_ids) | List of security group IDs | `list(string)` | `[]` | no |
| <a name="input_shared_tags"></a> [shared\_tags](#input\_shared\_tags) | Tags to add to all resources | `map` | `{}` | no |
| <a name="input_subnet_id"></a> [subnet\_id](#input\_subnet\_id) | Subnet ID where the instance will be launched | `string` | n/a | yes |
| <a name="input_use_localstack"></a> [use\_localstack](#input\_use\_localstack) | Whether to use LocalStack | `bool` | `false` | no |
| <a name="input_user_data"></a> [user\_data](#input\_user\_data) | User data script to run at instance launch | `string` | `""` | no |
| <a name="input_volume_size"></a> [volume\_size](#input\_volume\_size) | Desired root volume size in GiB | `number` | `1` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_dns"></a> [dns](#output\_dns) | n/a |
| <a name="output_instance_id"></a> [instance\_id](#output\_instance\_id) | n/a |
| <a name="output_instance_public_dns"></a> [instance\_public\_dns](#output\_instance\_public\_dns) | IPv6 DNS name for the instance |
| <a name="output_private_ip"></a> [private\_ip](#output\_private\_ip) | n/a |
| <a name="output_public_ipv4"></a> [public\_ipv4](#output\_public\_ipv4) | n/a |
| <a name="output_public_ipv6"></a> [public\_ipv6](#output\_public\_ipv6) | n/a |
<!-- END_TF_DOCS -->
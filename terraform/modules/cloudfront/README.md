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
| [aws_cloudfront_distribution.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudfront_distribution) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_cache_policy"></a> [cache\_policy](#input\_cache\_policy) | Cache behavior settings | <pre>object({<br>    min_ttl     = number<br>    default_ttl = number<br>    max_ttl     = number<br>  })</pre> | <pre>{<br>  "default_ttl": 3600,<br>  "max_ttl": 86400,<br>  "min_ttl": 0<br>}</pre> | no |
| <a name="input_certificate_arn"></a> [certificate\_arn](#input\_certificate\_arn) | ARN of the ACM certificate | `string` | n/a | yes |
| <a name="input_domain_name"></a> [domain\_name](#input\_domain\_name) | Domain name for the distribution aliases | `string` | n/a | yes |
| <a name="input_forwarded_headers"></a> [forwarded\_headers](#input\_forwarded\_headers) | List of headers to forward to the origin | `list(string)` | <pre>[<br>  "Host",<br>  "Authorization"<br>]</pre> | no |
| <a name="input_name"></a> [name](#input\_name) | Name for the CloudFront distribution | `string` | n/a | yes |
| <a name="input_origin_domain_name"></a> [origin\_domain\_name](#input\_origin\_domain\_name) | Domain name of the origin server | `string` | n/a | yes |
| <a name="input_origin_protocol_policy"></a> [origin\_protocol\_policy](#input\_origin\_protocol\_policy) | Protocol policy for the origin | `string` | `"http-only"` | no |
| <a name="input_price_class"></a> [price\_class](#input\_price\_class) | CloudFront price class | `string` | `"PriceClass_100"` | no |
| <a name="input_project"></a> [project](#input\_project) | Project name | `string` | n/a | yes |
| <a name="input_shared_tags"></a> [shared\_tags](#input\_shared\_tags) | Tags to add to all resources | `map(string)` | `{}` | no |
| <a name="input_wait_for_deployment"></a> [wait\_for\_deployment](#input\_wait\_for\_deployment) | Whether to wait for the distribution to be deployed | `bool` | `false` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_aliases"></a> [aliases](#output\_aliases) | Custom domain names for the distribution |
| <a name="output_distribution_arn"></a> [distribution\_arn](#output\_distribution\_arn) | ARN of the CloudFront distribution |
| <a name="output_distribution_id"></a> [distribution\_id](#output\_distribution\_id) | ID of the CloudFront distribution |
| <a name="output_domain_name"></a> [domain\_name](#output\_domain\_name) | Domain name of the CloudFront distribution |
| <a name="output_hosted_zone_id"></a> [hosted\_zone\_id](#output\_hosted\_zone\_id) | Route 53 zone ID of the CloudFront distribution |
| <a name="output_status"></a> [status](#output\_status) | Current status of the distribution |
<!-- END_TF_DOCS -->
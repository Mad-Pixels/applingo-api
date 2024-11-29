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
| [aws_sqs_queue.dead_letter](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sqs_queue) | resource |
| [aws_sqs_queue.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sqs_queue) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_delay_seconds"></a> [delay\_seconds](#input\_delay\_seconds) | The time in seconds that the delivery of all messages in the queue will be delayed | `number` | `240` | no |
| <a name="input_dlq_message_retention_seconds"></a> [dlq\_message\_retention\_seconds](#input\_dlq\_message\_retention\_seconds) | The number of seconds Amazon SQS retains a message in the DLQ | `number` | `259200` | no |
| <a name="input_max_message_size"></a> [max\_message\_size](#input\_max\_message\_size) | The limit of how many bytes a message can contain before Amazon SQS rejects it | `number` | `262144` | no |
| <a name="input_max_receive_count"></a> [max\_receive\_count](#input\_max\_receive\_count) | Max attempts to receive message | `number` | `3` | no |
| <a name="input_message_retention_seconds"></a> [message\_retention\_seconds](#input\_message\_retention\_seconds) | The number of seconds Amazon SQS retains a message | `number` | `3600` | no |
| <a name="input_project"></a> [project](#input\_project) | Project name | `string` | n/a | yes |
| <a name="input_queue_name"></a> [queue\_name](#input\_queue\_name) | Name of the SQS queue. | `string` | n/a | yes |
| <a name="input_receive_wait_time_seconds"></a> [receive\_wait\_time\_seconds](#input\_receive\_wait\_time\_seconds) | The time for which a ReceiveMessage call will wait for a message to arrive | `number` | `0` | no |
| <a name="input_shared_tags"></a> [shared\_tags](#input\_shared\_tags) | Tags to add to all resources | `map` | `{}` | no |
| <a name="input_visibility_timeout_seconds"></a> [visibility\_timeout\_seconds](#input\_visibility\_timeout\_seconds) | The visibility timeout for the queue | `number` | `300` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_dead_letter_queue_arn"></a> [dead\_letter\_queue\_arn](#output\_dead\_letter\_queue\_arn) | The ARN of the created Amazon SQS dead-letter queue |
| <a name="output_dead_letter_queue_url"></a> [dead\_letter\_queue\_url](#output\_dead\_letter\_queue\_url) | The URL of the created Amazon SQS dead-letter queue |
| <a name="output_queue_arn"></a> [queue\_arn](#output\_queue\_arn) | The ARN of the created Amazon SQS queue |
| <a name="output_queue_url"></a> [queue\_url](#output\_queue\_url) | The URL of the created Amazon SQS queue |
<!-- END_TF_DOCS -->
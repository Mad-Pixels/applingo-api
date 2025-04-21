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
| [aws_iam_role.scheduler](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role) | resource |
| [aws_iam_role_policy.dead_letter_policy](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy) | resource |
| [aws_iam_role_policy.target_policy](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy) | resource |
| [aws_scheduler_schedule.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/scheduler_schedule) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_dead_letter_arn"></a> [dead\_letter\_arn](#input\_dead\_letter\_arn) | ARN of the dead letter queue for failed executions | `string` | `null` | no |
| <a name="input_ecs_launch_type"></a> [ecs\_launch\_type](#input\_ecs\_launch\_type) | Launch type for ECS tasks (EC2 or FARGATE) | `string` | `"FARGATE"` | no |
| <a name="input_ecs_network_config"></a> [ecs\_network\_config](#input\_ecs\_network\_config) | Network configuration for ECS tasks | <pre>object({<br>    subnets          = list(string)<br>    security_groups  = list(string)<br>    assign_public_ip = bool<br>  })</pre> | `null` | no |
| <a name="input_ecs_task_definition_arn"></a> [ecs\_task\_definition\_arn](#input\_ecs\_task\_definition\_arn) | ARN of the ECS task definition (for ECS targets) | `string` | `null` | no |
| <a name="input_flexible_time_window_mode"></a> [flexible\_time\_window\_mode](#input\_flexible\_time\_window\_mode) | The mode of the flexible time window (OFF or FLEXIBLE) | `string` | `"OFF"` | no |
| <a name="input_input_json"></a> [input\_json](#input\_input\_json) | JSON input passed to the target | `string` | `null` | no |
| <a name="input_maximum_event_age_in_seconds"></a> [maximum\_event\_age\_in\_seconds](#input\_maximum\_event\_age\_in\_seconds) | Maximum age of events to process in seconds | `number` | `3600` | no |
| <a name="input_maximum_retry_attempts"></a> [maximum\_retry\_attempts](#input\_maximum\_retry\_attempts) | Maximum number of retry attempts | `number` | `null` | no |
| <a name="input_maximum_window_in_minutes"></a> [maximum\_window\_in\_minutes](#input\_maximum\_window\_in\_minutes) | Maximum window in minutes for flexible time window | `number` | `5` | no |
| <a name="input_policy"></a> [policy](#input\_policy) | Custom IAM policy JSON for the scheduler role | `string` | `""` | no |
| <a name="input_project"></a> [project](#input\_project) | Project name used for resource naming | `string` | n/a | yes |
| <a name="input_schedule_expression"></a> [schedule\_expression](#input\_schedule\_expression) | Schedule expression (cron or rate expression) | `string` | `"rate(1 day)"` | no |
| <a name="input_scheduler_name"></a> [scheduler\_name](#input\_scheduler\_name) | Name of the scheduler | `string` | n/a | yes |
| <a name="input_shared_tags"></a> [shared\_tags](#input\_shared\_tags) | Tags to apply to all resources | `map(string)` | `{}` | no |
| <a name="input_sqs_message_group_id"></a> [sqs\_message\_group\_id](#input\_sqs\_message\_group\_id) | Message group ID for SQS FIFO queues | `string` | `null` | no |
| <a name="input_target_arn"></a> [target\_arn](#input\_target\_arn) | ARN of the target resource | `string` | n/a | yes |
| <a name="input_target_type"></a> [target\_type](#input\_target\_type) | Type of target resource (lambda, sqs, ecs, etc.) | `string` | `"lambda"` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_schedule_arn"></a> [schedule\_arn](#output\_schedule\_arn) | ARN of the created schedule |
| <a name="output_schedule_name"></a> [schedule\_name](#output\_schedule\_name) | Name of the created schedule |
<!-- END_TF_DOCS -->
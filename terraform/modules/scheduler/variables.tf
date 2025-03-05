variable "project" {
  description = "Project name used for resource naming"
  type        = string
}

variable "scheduler_name" {
  description = "Name of the scheduler"
  type        = string
}

variable "schedule_expression" {
  description = "Schedule expression (cron or rate expression)"
  type        = string
  default     = "rate(1 day)"
}

variable "flexible_time_window_mode" {
  description = "The mode of the flexible time window (OFF or FLEXIBLE)"
  type        = string
  default     = "OFF"
  validation {
    condition     = contains(["OFF", "FLEXIBLE"], var.flexible_time_window_mode)
    error_message = "Flexible time window mode must be OFF or FLEXIBLE."
  }
}

variable "maximum_window_in_minutes" {
  description = "Maximum window in minutes for flexible time window"
  type        = number
  default     = 5
}

variable "maximum_retry_attempts" {
  description = "Maximum number of retry attempts"
  type        = number
  default     = null
}

variable "maximum_event_age_in_seconds" {
  description = "Maximum age of events to process in seconds"
  type        = number
  default     = 3600
}

variable "target_type" {
  description = "Type of target resource (lambda, sqs, ecs, etc.)"
  type        = string
  default     = "lambda"
  validation {
    condition     = contains(["lambda", "sqs", "ecs", "events"], var.target_type)
    error_message = "Target type must be one of: lambda, sqs, ecs, events."
  }
}

variable "target_arn" {
  description = "ARN of the target resource"
  type        = string
}

variable "input_json" {
  description = "JSON input passed to the target"
  type        = string
  default     = null
}

variable "dead_letter_arn" {
  description = "ARN of the dead letter queue for failed executions"
  type        = string
  default     = null
}

variable "ecs_task_definition_arn" {
  description = "ARN of the ECS task definition (for ECS targets)"
  type        = string
  default     = null
}

variable "ecs_launch_type" {
  description = "Launch type for ECS tasks (EC2 or FARGATE)"
  type        = string
  default     = "FARGATE"
}

variable "ecs_network_config" {
  description = "Network configuration for ECS tasks"
  type = object({
    subnets          = list(string)
    security_groups  = list(string)
    assign_public_ip = bool
  })
  default = null
}

variable "sqs_message_group_id" {
  description = "Message group ID for SQS FIFO queues"
  type        = string
  default     = null
}

variable "policy" {
  description = "Custom IAM policy JSON for the scheduler role"
  type        = string
  default     = ""
}

variable "shared_tags" {
  description = "Tags to apply to all resources"
  type        = map(string)
  default     = {}
}

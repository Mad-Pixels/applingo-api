variable "use_localstack" {
  description = "Whether to use LocalStack"
  type        = bool
  default     = false
}

variable "localstack_endpoint" {
  description = "LocalStack endpoint"
  type        = string
  default     = "https://localhost.localstack.cloud:4566"
}

variable "aws_region" {
  description = "AWS region"
  type        = string
}

variable "arch" {
  description = "Set architecture which will be use in lambda services"
  type        = string
}
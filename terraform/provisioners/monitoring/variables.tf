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

variable "domain_name" {
  description = "Project monitoring service domain"
  type        = string
}

variable "domain_zone_id" {
  description = "Domain zone id"
  type        = string
}

variable "domain_acm_arn" {
  description = "ACM certificate ARN"
  type        = string
}
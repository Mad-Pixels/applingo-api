variable "use_localstack" {
  description = "Whether to use LocalStack"
  type        = bool
  default     = false
}

variable "localstack_endpoint" {
  description = "LocalStack endpoint"
  type        = string
  default     = "http://localhost:4566"
}

variable "aws_region" {
  description = "AWS region"
  default     = "us-east-1"
}
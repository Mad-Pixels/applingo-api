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

variable "device_api_token" {
  description = "Auth token which use for lambda request validate from device"
  type        = string
}

variable "openai_model" {
  description = "OpenAI GPT model name"
  type        = string
  default     = "gpt-3.5-turbo"
}

variable "openai_key" {
  description = "OpenAI request key"
  type        = string
}

variable "jwt_secret" {
  description = "Auth JWT secret which use for lambda request validate from external"
  type        = string
}
variable "project" {
  description = "Project name"
  type        = string
}

variable "function_name" {
  description = "Name of the Lambda function"
  type        = string
}

variable "image" {
  description = "URI of the container image in ECR"
  type        = string
}

variable "memory_size" {
  description = "Amount of memory in MB for the Lambda function"
  type        = number
  default     = 128
}

variable "timeout" {
  description = "Timeout for the Lambda function in seconds"
  type        = number
  default     = 5
}

variable "vpc_config" {
  description = "Optional VPC configuration for the Lambda function"
  type = object({
    subnet_ids         = list(string)
    security_group_ids = list(string)
  })
  default = null
}

variable "policy" {
  description = "Additional IAM policy for the Lambda function"
  type        = string
  default     = ""
}

variable "environments" {
  description = "Additional ENVs for the Lambda function"
  default     = {} 
}

variable "shared_tags" {
  description = "Tags to add to all resources"
  default     = {}
}

variable "log_level" {
  description = "The log level for the Lambda function"
  type        = string
  default     = "ERROR"
}

variable "log_retention" {
  description = "Cloudwatch retention in days"
  default     = 3 
}
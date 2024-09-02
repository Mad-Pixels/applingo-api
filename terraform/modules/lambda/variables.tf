variable "name" {
  description = "Name of the Lambda function"
  type        = string
}

variable "image" {
  description = "URI of the container image in ECR"
  type        = string
}

variable "mem_size" {
  description = "Amount of memory in MB for the Lambda function"
  type        = number
  default     = 128
}

variable "timeout" {
  description = "Timeout for the Lambda function in seconds"
  type        = number
  default     = 5
}

variable "additional_policy" {
  description = "Additional IAM policy for the Lambda function"
  type        = string
  default     = ""
}
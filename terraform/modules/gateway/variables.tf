variable "project" {
  description = "Project name"
  type        = string
}

variable "api_name" {
  description = "Name of the API Gateway"
  type        = string
}

variable "router_invoke_arn" {
  description = ""
  type        = string
}

variable "acm_certificate_arn" {
  description = "AWS ACM certificate for api gateway custom domain"
  type        = string
  default     = ""
}

variable "custom_domain" {
  description = "GatewayAPI custom domain name"
  type        = string
  default     = ""
}

variable "wafv2_web_acl_arn" {
  description = "WAF web acl arn"
  type        = string
  default     = ""
}

variable "stage_name" {
  description = "Option set API Gateway stage name"
  type        = string
  default     = "live"
}

variable "log_retention_days" {
  description = "Optional set log retention"
  type        = number
  default     = 3
}

variable "shared_tags" {
  description = "Tags to add to all resources"
  default     = {}
}
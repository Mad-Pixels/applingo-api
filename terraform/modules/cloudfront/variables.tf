variable "project" {
  description = "Project name"
  type        = string
}

variable "name" {
  description = "Name for the CloudFront distribution"
  type        = string
}

variable "domain_name" {
  description = "Domain name for the distribution aliases"
  type        = string
}

variable "origin_domain_name" {
  description = "Domain name of the origin server"
  type        = string
}

variable "certificate_arn" {
  description = "ARN of the ACM certificate"
  type        = string
}

variable "price_class" {
  description = "CloudFront price class"
  type        = string
  default     = "PriceClass_100"

  validation {
    condition     = contains(["PriceClass_100", "PriceClass_200", "PriceClass_All"], var.price_class)
    error_message = "Price class must be one of: PriceClass_100, PriceClass_200, PriceClass_All."
  }
}

variable "origin_protocol_policy" {
  description = "Protocol policy for the origin"
  type        = string
  default     = "http-only"

  validation {
    condition     = contains(["http-only", "match-viewer", "https-only"], var.origin_protocol_policy)
    error_message = "Origin protocol policy must be one of: http-only, match-viewer, https-only."
  }
}

variable "forwarded_headers" {
  description = "List of headers to forward to the origin"
  type        = list(string)
  default     = ["Host", "Authorization"]
}

variable "cache_policy" {
  description = "Cache behavior settings"
  type = object({
    min_ttl     = number
    default_ttl = number
    max_ttl     = number
  })
  default = {
    min_ttl     = 0
    default_ttl = 3600
    max_ttl     = 86400
  }
}

variable "shared_tags" {
  description = "Tags to add to all resources"
  type        = map(string)
  default     = {}
}

variable "wait_for_deployment" {
  description = "Whether to wait for the distribution to be deployed"
  type        = bool
  default     = false
}
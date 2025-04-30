variable "name" {
  description = "VPC name"
  type        = string
}

variable "region" {
  description = "AWS region"
  type        = string
}

variable "vpc_addr_block" {
  description = "Base address block for the VPC (e.g., 10.100.100.0), use cidr: /23 (512 addrs)"
  type        = string
}

variable "vpc_zones" {
  description = "Number of availability zones to use (1-3)"
  type        = number
  default     = 1
  validation {
    condition     = var.vpc_zones >= 1 && var.vpc_zones <= 3
    error_message = "vpc_zones must be between 1 and 3"
  }
}

variable "enable_dns_support" {
  description = "Enables DNS support"
  type        = bool
  default     = true
}

variable "enable_dns_hostnames" {
  description = "Enables DNS hostnames"
  type        = bool
  default     = true
}

variable "enable_internet_gateway" {
  description = "Enables access to/from internet"
  type        = bool
  default     = true
}

variable "use_public_subnets" {
  description = "Enables public subnets"
  type        = bool
  default     = true
}

variable "use_private_subnets" {
  description = "Enables private subnets"
  type        = bool
  default     = false
}

variable "enable_nat_gateway" {
  description = "Enables NAT gateway for private subnets"
  type        = bool
  default     = false
}

variable "shared_tags" {
  description = "Tags to add to all resources"
  type        = map(string)
  default     = {}
}

variable "enable_s3_endpoint" {
  description = "Add S3 endpoint"
  type        = bool
  default     = false
}

variable "enable_cloudwatch_endpoint" {
  description = "Enable CloudWatch VPC endpoint"
  type        = bool
  default     = false
}

variable "enable_sts_endpoint" {
  description = "Enable STS VPC endpoint"
  type        = bool
  default     = false
}
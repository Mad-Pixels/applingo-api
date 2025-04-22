variable "project" {
  description = "Project name"
  type        = string
}

variable "name" {
  description = "VPC name"
  type        = string
}

variable "aws_region" {
  description = "AWS region name"
  type        = string
}

variable "vpc_base_ip" {
  description = "Base IP address for the vps subnets (e.g., 10.216.0.0)"
  type        = string

  validation {
    condition     = can(regex("^10\\.\\d+\\.0\\.0$", var.vpc_base_ip))
    error_message = "vpc_base_ip should end with .0.0 (e.g., 10.216.0.0)"
  }
}

variable "vpc_zones" {
  description = "Availability zones count."
  type        = number

  validation {
    condition     = var.vpc_zones >= 1 && var.vpc_zones <= 3
    error_message = "vpc_zones should be in the range of 1 to 3."
  }
}

variable "use_public_subnets" {
  description = "Enables public subnets"
  type        = bool
  default     = false
}

variable "enable_dns_support" {
  description = "Enables DNS support"
  type        = bool
  default     = true
}

variable "enable_dns_hostnames" {
  description = "Enables DNS hostnames"
  type        = bool
  default     = false
}

variable "enable_internet_gateway" {
  description = "Enables access to/from internet"
  type        = bool
  default     = false
}

variable "shared_tags" {
  description = "Tags to add to all resources"
  type        = map(string)
  default     = {}
}

variable "ssh_allowed_cidr_blocks" {
  description = "CIDR blocks allowed for SSH access"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}

variable "create_endpoint_sg" {
  description = "Whether to create security group for VPC endpoints"
  type        = bool
  default     = false
}

variable "create_ssh_sg" {
  description = "Whether to create security group with port 22 open"
  type        = bool
  default     = false
}

variable "create_grafana_sg" {
  type    = bool
  default = false
}
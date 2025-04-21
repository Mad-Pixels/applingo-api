variable "project" {
  description = "Project name"
  type        = string
}

variable "ami_id" {
  description = "AMI ID for the instance"
  type        = string
  default     = ""
}

variable "graviton_size" {
  description = "Instance size (e.g. micro, small, medium)"
  type        = string
  default     = "micro"
}

variable "shared_tags" {
  description = "Tags to add to all resources"
  default     = {}
}

variable "subnet_id" {
  description = "Subnet ID where the instance will be launched"
  type        = string
}

variable "security_group_ids" {
  description = "List of security group IDs"
  type        = list(string)
  default     = []
}

variable "key_name" {
  description = "Name of the EC2 Key Pair for SSH access"
  type        = string
  default     = ""
}

variable "volume_size" {
  description = "Desired root volume size in GiB"
  type        = number
  default     = 1
}

variable "user_data" {
  description = "User data script to run at instance launch"
  type        = string
  default     = ""
}
variable "project" {
  description = "Project name"
  type        = string
}

variable "table_name" {
  description = "The name of the DynamoDB table"
  type        = string
}

variable "billing_mode" {
  description = "Controls how you are charged for read and write throughput and how you manage capacity"
  type        = string
  default     = "PAY_PER_REQUEST"
}

variable "hash_key" {
  description = "The attribute to use as the hash (partition) key"
  type        = string
}

variable "range_key" {
  description = "The attribute to use as the range (sort) key"
  type        = string
  default     = null
}

variable "attributes" {
  description = "List of nested attribute definitions. Only required for hash_key and range_key attributes"
  type = list(object({
    name = string
    type = string
  }))
}

variable "shared_tags" {
  description = "Tags to add to all resources"
  default     = {}
}
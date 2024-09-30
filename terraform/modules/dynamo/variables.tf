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

variable "secondary_index_list" {
  description = "List of global secondary indexes"
  type = list(object({
    name               = string
    hash_key           = string
    range_key          = optional(string)
    projection_type    = string
    non_key_attributes = optional(list(string))
    read_capacity      = optional(number)
    write_capacity     = optional(number)
  }))
  default = null
}

variable "ttl_enabled" {
  description = "Whether to enable TTL for the DynamoDB table"
  type        = bool
  default     = false
}

variable "ttl_attribute_name" {
  description = "The name of the TTL attribute"
  type        = string
  default     = "ttl"
}

variable "shared_tags" {
  description = "Tags to add to all resources"
  default     = {}
}

variable "stream_enabled" {
  description = "On/Off dynamo stream"
  default     = false
}

variable "stream_type" {
  description = "Type of streaming"
  default     = "NEW_IMAGE"
}
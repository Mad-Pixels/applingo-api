variable "project" {
  description = "Project name"
  type        = string
}

variable "bucket_name" {
  description = "Name of the S3 bucket."
  type        = string
}

variable "force_destroy" {
  description = "Allows Terraform to delete the bucket when removing the resource of set true."
  type        = bool
  default     = true
}

variable "is_website" {
  description = "Specifies if the S3 bucket will host a static website when set to true."
  type        = bool
  default     = false
}

variable "enable_versioning" {
  description = "Enables versioning for the S3 bucket if set true."
  type        = bool
  default     = false
}

variable "index_document" {
  description = "Name of the index document for the S3 static web site."
  type        = string
  default     = "index.html"
}

variable "error_document" {
  description = "Name of the error document for the S3 static web site."
  type        = string
  default     = "error.html"
}

variable "shared_tags" {
  description = "Tags to add to all resources"
  default     = {}
}

variable "rule" {
  description = "Optional lifecycle rule configuration"
  type = object({
    id     = string
    status = string
    filter = object({
      prefix = string
    })
    transition = optional(object({
      days          = number
      storage_class = string
    }))
    expiration = optional(object({
      days = number
    }))
  })
  default = null
}
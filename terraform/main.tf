terraform {
  backend "s3" {}
}

provider "aws" {
  region = var.region

  dynamic "endpoints" {
    for_each = var.use_localstack ? [1] : []
    content {
      iam    = var.localstack_endpoint
      lambda = var.localstack_endpoint
      s3     = var.localstack_endpoint
      sts    = var.localstack_endpoint
    }
  }

  s3_use_path_style           = var.use_localstack
  skip_credentials_validation = var.use_localstack
  skip_metadata_api_check     = var.use_localstack
  skip_requesting_account_id  = var.use_localstack

  access_key = var.use_localstack ? "test" : null
  secret_key = var.use_localstack ? "test" : null
}


variable "use_localstack" {
  description = "Whether to use LocalStack instead of real AWS"
  type        = bool
  default     = false
}

variable "localstack_endpoint" {
  description = "LocalStack endpoint"
  type        = string
  default     = "http://localhost:4566"
}

variable "region" {
  type    = string
  default = "us-east-1"
}

module "lambda_functions" {
  
  source   = "./modules/lambda"

  function_name = "dictionary"
  source_dir    = "${path.root}/../cmd/dictionary"
}
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.65.0"
    }
  }
}

provider "aws" {
  region = var.aws_region

  dynamic "endpoints" {
    for_each = var.use_localstack ? [1] : []
    content {
      s3       = var.localstack_endpoint
      ecr      = var.localstack_endpoint
      sqs      = var.localstack_endpoint
      iam      = var.localstack_endpoint
      dynamodb = var.localstack_endpoint
    }
  }

  skip_credentials_validation = var.use_localstack
  skip_metadata_api_check     = var.use_localstack
  skip_requesting_account_id  = var.use_localstack

  s3_use_path_style = var.use_localstack
  
  access_key = var.use_localstack ? "test" : null
  secret_key = var.use_localstack ? "test" : null
}
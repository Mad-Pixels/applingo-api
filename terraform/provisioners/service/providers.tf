terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.65.0"
    }
  }

  backend "s3" {
    bucket  = "tfstates-madpixels"
    key     = "lingocards-api/service.tfstate"
    region  = "eu-central-1"
    encrypt = true
  }
}

provider "aws" {
  region = var.aws_region

  dynamic "endpoints" {
    for_each = var.use_localstack ? [1] : []
    content {
      lambda     = var.localstack_endpoint
      iam        = var.localstack_endpoint
      sts        = var.localstack_endpoint
      logs       = var.localstack_endpoint
      cloudwatch = var.localstack_endpoint
      apigateway = var.localstack_endpoint
    }
  }

  skip_credentials_validation = var.use_localstack
  skip_metadata_api_check     = var.use_localstack
  skip_requesting_account_id  = var.use_localstack

  access_key = var.use_localstack ? "test" : null
  secret_key = var.use_localstack ? "test" : null
}
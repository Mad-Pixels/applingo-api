terraform {}

provider "aws" {
  region = var.aws_region

  dynamic "endpoints" {
    for_each = var.use_localstack ? [1] : []
    content {
      ecr = var.localstack_endpoint
    }
  }

  skip_credentials_validation = var.use_localstack
  skip_metadata_api_check     = var.use_localstack
  skip_requesting_account_id  = var.use_localstack

  access_key = var.use_localstack ? "test" : null
  secret_key = var.use_localstack ? "test" : null
}

module "ecr-lingocards-api" {
  source = "../../modules/ecr"

  repository_name = "lingocards-api"
}
module "ecr-repository-api" {
  source = "../../modules/ecr"

  project         = local.project
  shared_tags     = local.tags
  repository_name = "images"
}

module "s3-forge-bucket" {
  source = "../../modules/s3"

  project     = local.project
  shared_tags = local.tags
  bucket_name = "forge-${var.environment}"
}

module "s3-dictionary-bucket" {
  source = "../../modules/s3"

  project     = local.project
  shared_tags = local.tags
  bucket_name = "dictionary-${var.environment}"
}

module "s3-processing-bucket" {
  source = "../../modules/s3"

  project     = local.project
  shared_tags = local.tags
  bucket_name = "processing-${var.environment}"
}

module "s3-errors-bucket" {
  source = "../../modules/s3"

  project     = local.project
  shared_tags = local.tags
  bucket_name = "errors-${var.environment}"

  rule = {
    id     = "cleanup"
    status = "Enabled"
    filter = {
      prefix = "logs-"
    }
    expiration = {
      days = var.environment == "prd" ? 30 : 7
    }
  }
}

module "dynamo-dictionary-table" {
  source = "../../modules/dynamo"

  project              = local.project
  table_name           = local.dictionary_dynamo_schema.table_name
  hash_key             = local.dictionary_dynamo_schema.hash_key
  range_key            = local.dictionary_dynamo_schema.range_key
  attributes           = local.dictionary_dynamo_schema.attributes
  secondary_index_list = local.dictionary_dynamo_schema.secondary_indexes
  stream_enabled       = true

  stream_type = "NEW_AND_OLD_IMAGES"
  shared_tags = local.tags
}

module "dynamo-processing-table" {
  source = "../../modules/dynamo"

  project              = local.project
  table_name           = local.processing_dynamo_schema.table_name
  hash_key             = local.processing_dynamo_schema.hash_key
  range_key            = local.processing_dynamo_schema.range_key
  attributes           = local.processing_dynamo_schema.attributes
  secondary_index_list = local.processing_dynamo_schema.secondary_indexes
  stream_enabled       = true

  stream_type = "NEW_AND_OLD_IMAGES"
  shared_tags = local.tags
}

module "dynamo-profile-table" {
  source = "../../modules/dynamo"

  project              = local.project
  table_name           = local.profile_dynamo_schema.table_name
  hash_key             = local.profile_dynamo_schema.hash_key
  range_key            = local.profile_dynamo_schema.range_key
  attributes           = local.profile_dynamo_schema.attributes
  secondary_index_list = local.profile_dynamo_schema.secondary_indexes
  stream_enabled       = false

  shared_tags = local.tags
}
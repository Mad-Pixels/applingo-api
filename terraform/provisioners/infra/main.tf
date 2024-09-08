module "ecr-repository-api" {
  source = "../../modules/ecr"

  project         = local.project
  repository_name = "api"
}

module "s3-dictionary-bucket" {
  source = "../../modules/s3"

  project     = local.project
  bucket_name = "dictionary"
}

module "s3-processing-bucket" {
  source = "../../modules/s3"

  project     = local.project
  bucket_name = "processing"
}

module "dynamo-dictionary-table" {
  source = "../../modules/dynamo"

  project    = local.project
  table_name = "dictionary"
  hash_key   = "id"
  range_key  = "name"

  attributes = [
    { name = "id", type = "S" },
    { name = "name", type = "S" },
    { name = "author", type = "S" },
    { name = "category_main", type = "S" },
    { name = "is_private", type = "N" },
    { name = "is_publish", type = "N" }
  ]

  secondary_index_list = [
    {
      name            = "AuthorIndex"
      hash_key        = "author"
      range_key       = "name"
      projection_type = "ALL"
    },
    {
      name            = "CategoryMainIndex"
      hash_key        = "category_main"
      range_key       = "name"
      projection_type = "ALL"
    },
    {
      name               = "IsPrivateIndex"
      hash_key           = "is_private"
      range_key          = "name"
      projection_type    = "INCLUDE"
      non_key_attributes = ["author", "category_main"]
    },
    {
      name               = "IsPublishIndex"
      hash_key           = "is_publish"
      range_key          = "name"
      projection_type    = "INCLUDE"
      non_key_attributes = ["author", "category_main"]
    }
  ]
}
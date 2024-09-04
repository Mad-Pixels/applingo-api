module "ecr-lingocards-api" {
  source = "../../modules/ecr"

  project         = "lingocards"
  repository_name = "api"
}

module "s3-dictionary-bucket" {
  source = "../../modules/s3"

  project     = "lingocards"
  bucket_name = "dictionary"
}

module "dynamo-dictionary-table" {
  source = "../../modules/dynamo"

  project    = "lingocards"
  table_name = "dictionary"
  hash_key   = "id"
  range_key  = "timestamp"

  attributes = [
    {
      name = "id"
      type = "S"
    },
    {
      name = "timestamp"
      type = "N"
    }
  ] 
}
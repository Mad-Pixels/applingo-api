output "ecr-repository-api_url" {
  value = module.ecr-repository-api.repository_url
}

output "s3-dictionary-bucket_name" {
  value = module.s3-dictionary-bucket.s3_name
}

output "dynamo-dictionary-table_name" {
  value = module.dynamo-dictionary-table.table_name
}
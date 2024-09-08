output "lambda_lingocards_repository_url" {
  value = module.ecr-lingocards-api.repository_url
}

output "lambda_lingocards_repository_arn" {
  value = module.ecr-lingocards-api.repository_arn
}

output "bucket_lingocards_name" {
  value = module.s3-dictionary-bucket.s3_name
}

output "dynamo_lingocards_table" {
  value = module.dynamo-dictionary-table.table_name
}
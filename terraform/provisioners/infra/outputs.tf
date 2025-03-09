output "ecr-repository-api_url" {
  value = module.ecr-repository-api.repository_url
}

output "s3-forge-bucket_name" {
  value = module.s3-forge-bucket.s3_name
}

output "s3-forge-bucket_arn" {
  value = module.s3-forge-bucket.s3_arn
}

output "s3-dictionary-bucket_name" {
  value = module.s3-dictionary-bucket.s3_name
}

output "s3-dictionary-bucket_arn" {
  value = module.s3-dictionary-bucket.s3_arn
}

output "s3-processing-bucket_name" {
  value = module.s3-processing-bucket.s3_name
}

output "s3-processing-bucket_arn" {
  value = module.s3-processing-bucket.s3_arn
}

output "s3-errors-bucket_name" {
  value = module.s3-errors-bucket.s3_name
}

output "s3-errors-bucket_arn" {
  value = module.s3-errors-bucket.s3_arn
}

output "dynamo-dictionary-table_name" {
  value = module.dynamo-dictionary-table.table_name
}

output "dynamo-dictionary-table_arn" {
  value = module.dynamo-dictionary-table.table_arn
}

output "dynamo-dictionary-stream_arn" {
  value = module.dynamo-dictionary-table.stream_arn
}

output "dynamo-processing-table_name" {
  value = module.dynamo-processing-table.table_name
}

output "dynamo-processing-table_arn" {
  value = module.dynamo-processing-table.table_arn
}

output "dynamo-processing-stream_arn" {
  value = module.dynamo-processing-table.stream_arn
}

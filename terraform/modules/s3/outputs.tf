output "s3_id" {
  value = aws_s3_bucket.this.id
}

output "s3_arn" {
  value = aws_s3_bucket.this.arn
}

output "s3_name" {
  value = "${var.project}-${var.bucket_name}"
}

output "s3_domain" {
  value = aws_s3_bucket.this.bucket_regional_domain_name
}
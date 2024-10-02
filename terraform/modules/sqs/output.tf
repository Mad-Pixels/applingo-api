output "queue_url" {
  description = "The URL of the created Amazon SQS queue"
  value       = aws_sqs_queue.this.id
}

output "queue_arn" {
  description = "The ARN of the created Amazon SQS queue"
  value       = aws_sqs_queue.this.arn
}

output "dead_letter_queue_url" {
  description = "The URL of the created Amazon SQS dead-letter queue"
  value       = aws_sqs_queue.dead_letter.id
}

output "dead_letter_queue_arn" {
  description = "The ARN of the created Amazon SQS dead-letter queue"
  value       = aws_sqs_queue.dead_letter.arn
}
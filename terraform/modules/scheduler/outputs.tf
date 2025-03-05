output "schedule_arn" {
  description = "ARN of the created schedule"
  value       = aws_scheduler_schedule.this.arn
}

output "schedule_name" {
  description = "Name of the created schedule"
  value       = aws_scheduler_schedule.this.name
}

output "role_arn" {
  description = "ARN of the IAM role created for the scheduler"
  value       = aws_iam_role.this.arn
}

output "role_name" {
  description = "Name of the IAM role created for the scheduler"
  value       = aws_iam_role.this.name
}
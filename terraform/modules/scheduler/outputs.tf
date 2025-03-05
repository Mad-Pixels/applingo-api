output "schedule_arn" {
  description = "ARN of the created schedule"
  value       = aws_scheduler_schedule.this.arn
}

output "schedule_name" {
  description = "Name of the created schedule"
  value       = aws_scheduler_schedule.this.name
}

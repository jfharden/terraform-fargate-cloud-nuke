output "task_arn" {
  description = "ARN for the ecs task"
  value       = ""
}

output "ecs_repository_arn" {
  description = "ARN for the ecs repository"
  value       = aws_ecr_repository.cloud_nuke.arn
}

output "cloudwatch_schedule_arn" {
  description = "ARN for the cloudwatch schedule, if it was created"
  value       = ""
}

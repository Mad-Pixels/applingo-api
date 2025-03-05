resource "aws_scheduler_schedule" "this" {
  name                = "${var.project}-${replace(var.scheduler_name, "_", "-")}"
  schedule_expression = var.schedule_expression

  flexible_time_window {
    mode                      = var.flexible_time_window_mode
    maximum_window_in_minutes = var.flexible_time_window_mode == "FLEXIBLE" ? var.maximum_window_in_minutes : null
  }

  target {
    arn      = var.target_arn
    role_arn = aws_iam_role.scheduler.arn
    input    = var.input_json

    retry_policy {
      maximum_event_age_in_seconds = var.maximum_event_age_in_seconds
      maximum_retry_attempts       = var.maximum_retry_attempts
    }

    dynamic "dead_letter_config" {
      for_each = var.dead_letter_arn != null ? [1] : []
      content {
        arn = var.dead_letter_arn
      }
    }

    dynamic "ecs_parameters" {
      for_each = var.target_type == "ecs" ? [1] : []
      content {
        task_definition_arn = var.ecs_task_definition_arn
        launch_type         = var.ecs_launch_type

        dynamic "network_configuration" {
          for_each = var.ecs_network_config != null ? [1] : []
          content {
            subnets          = var.ecs_network_config.subnets
            security_groups  = var.ecs_network_config.security_groups
            assign_public_ip = var.ecs_network_config.assign_public_ip
          }
        }
      }
    }

    dynamic "sqs_parameters" {
      for_each = var.target_type == "sqs" && var.sqs_message_group_id != null ? [1] : []
      content {
        message_group_id = var.sqs_message_group_id
      }
    }
  }
}

resource "aws_iam_role" "scheduler" {
  name = "${var.project}-${replace(var.scheduler_name, "_", "-")}-scheduler-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "scheduler.amazonaws.com"
        }
      }
    ]
  })

  tags = merge(
    var.shared_tags,
    {
      "TF"      = "true",
      "Project" = var.project,
      "Github"  = "github.com/Mad-Pixels/applingo-api",
    }
  )
}
resource "aws_scheduler_schedule" "this" {
  name                = "${var.project}-${replace(var.scheduler_name, "_", "-")}"
  schedule_expression = var.schedule_expression

  flexible_time_window {
    mode                      = var.flexible_time_window_mode
    maximum_window_in_minutes = var.flexible_time_window_mode == "FLEXIBLE" ? var.maximum_window_in_minutes : null
  }

  dynamic "retry_policy" {
    for_each = var.maximum_retry_attempts != null ? [1] : []
    content {
      maximum_retry_attempts       = var.maximum_retry_attempts
      maximum_event_age_in_seconds = var.maximum_event_age_in_seconds
    }
  }

  target {
    arn      = var.target_arn
    role_arn = aws_iam_role.this.arn

    input = var.input_json != null ? var.input_json : null

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

resource "aws_iam_role" "this" {
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

resource "aws_iam_role_policy" "target_policy" {
  name = "${var.project}-${replace(var.scheduler_name, "_", "-")}-target-policy"
  role = aws_iam_role.this.id

  policy = var.policy != "" ? var.policy : jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = var.target_type == "lambda" ? [
          "lambda:InvokeFunction"
          ] : var.target_type == "sqs" ? [
          "sqs:SendMessage"
          ] : var.target_type == "ecs" ? [
          "ecs:RunTask"
          ] : [
          "events:PutEvents"
        ]
        Resource = var.target_arn
      }
    ]
  })
}

resource "aws_iam_role_policy" "dead_letter_policy" {
  count = var.dead_letter_arn != null ? 1 : 0

  name = "${var.project}-${replace(var.scheduler_name, "_", "-")}-dlq-policy"
  role = aws_iam_role.this.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "sqs:SendMessage"
        ]
        Resource = var.dead_letter_arn
      }
    ]
  })
}
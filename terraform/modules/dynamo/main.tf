resource "aws_dynamodb_table" "this" {
  name           = "${var.project}-${var.table_name}" 
  billing_mode   = var.billing_mode
  hash_key       = var.hash_key
  range_key      = var.range_key

  dynamic "attribute" {
    for_each = var.attributes
    content {
      name = attribute.value.name
      type = attribute.value.type
    }
  }

  tags = merge(
    var.shared_tags,
    {
      "TF"      = "true",
      "Project" = var.project,
      "Github"  = "github.com/Mad-Pixels/lingocards-api", 
    }
  )

  point_in_time_recovery {
    enabled = true
  }

  server_side_encryption {
    enabled = true
  }
}
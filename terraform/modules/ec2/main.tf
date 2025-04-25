resource "aws_instance" "this" {
  vpc_security_group_ids = var.security_group_ids
  ami                    = local.selected_ami_id
  subnet_id              = var.subnet_id
  ipv6_address_count     = 1

  associate_public_ip_address = var.associate_public_ip_address
  instance_type               = "t4g.${var.graviton_size}"

  user_data = var.user_data != "" ? var.user_data : null
  key_name  = var.key_name != "" ? var.key_name : null

  dynamic "metadata_options" {
    for_each = var.use_localstack ? [] : [1]

    content {
      http_endpoint               = "enabled"
      http_tokens                 = "required"
      instance_metadata_tags      = "enabled"
      http_put_response_hop_limit = 1
    }
  }

  root_block_device {
    delete_on_termination = true
    encrypted             = true
    volume_type           = "gp3"
    iops                  = 3000
    throughput            = 125

    volume_size = max(var.volume_size, local.selected_root_volume_size)
  }

  tags = merge(
    var.shared_tags,
    {
      "Name" = var.name,
    }
  )
}

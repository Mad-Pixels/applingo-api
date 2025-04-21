data "aws_ami" "amazon_linux" {
  most_recent = true
  owners      = ["137112412989"]

  filter {
    name   = "name"
    values = ["amzn2-ami-hvm-*-arm64-gp2"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "architecture"
    values = ["arm64"]
  }
}

resource "aws_instance" "this" {
  ami           = local.selected_ami_id
  instance_type = "t4g.${var.graviton_size}"

  subnet_id              = var.subnet_id
  vpc_security_group_ids = var.security_group_ids

  associate_public_ip_address = false
  ipv6_address_count          = 1

  user_data = var.user_data != "" ? var.user_data : null
  key_name  = var.key_name  != "" ? var.key_name  : null

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
    throughput            = 125
    iops                  = 3000

    volume_size = max(var.volume_size, local.selected_root_volume_size)
  }

  tags = merge(
    var.shared_tags,
    {
      "TF"      = "true",
      "Project" = var.project,
      "Github"  = "github.com/Mad-Pixels/applingo-api",
    }
  )
}

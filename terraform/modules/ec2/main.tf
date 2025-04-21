data "aws_ami" "amazon_linux" {
  most_recent = true

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

  owners = ["137112412989"]
}

data "aws_ami" "selected" {
  count       = var.ami_id != "" ? 1 : 0
  owners      = ["amazon"]
  image_id    = var.ami_id
  most_recent = false
}

resource "aws_instance" "this" {
  ami           = var.ami_id != "" ? var.ami_id : data.aws_ami.amazon_linux.id
  instance_type = "t4g.${var.graviton_size}"

  vpc_security_group_ids      = var.security_group_ids
  subnet_id                   = var.subnet_id
  associate_public_ip_address = false
  ipv6_address_count          = 1

  user_data = var.user_data != "" ? var.user_data : null
  key_name  = var.key_name  != "" ? var.key_name  : null

  metadata_options {
    http_endpoint               = "enabled"
    http_tokens                 = "required"
    instance_metadata_tags      = "enabled"
    http_put_response_hop_limit = 1
  }

  root_block_device {
    delete_on_termination = true
    encrypted             = true
    volume_type           = "gp3"
    throughput            = 125
    iops                  = 3000

    volume_size = max(
      var.volume_size, 
      var.ami_id != "" 
        ? data.aws_ami.selected[0].block_device_mappings[0].ebs.volume_size 
        : data.aws_ami.amazon_linux.block_device_mappings[0].ebs.volume_size
    )
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
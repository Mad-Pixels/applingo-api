resource "aws_security_group" "endpoint_security_group" {
  count       = var.create_endpoint_sg ? 1 : 0
  name        = "${var.project}-sg"
  description = "Allow traffic for endpoints"
  vpc_id      = aws_vpc.this.id

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = -1
    cidr_blocks = [aws_vpc.this.cidr_block]
    description = "Allow internal VPC traffic (ingress)"
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = [aws_vpc.this.cidr_block]
    description = "Allow internal VPC traffic (egress)"
  }

  tags = merge(
    var.shared_tags,
    {
      "TF"      = "true",
      "Type"    = "endpoint",
      "Project" = var.project,
      "Github"  = "github.com/Mad-Pixels/applingo-api",
    }
  )
}

resource "aws_security_group" "ssh_security_group" {
  count       = var.create_ssh_sg ? 1 : 0
  name        = "${var.project}-ssh-sg"
  description = "Allow SSH access"
  vpc_id      = aws_vpc.this.id

  ingress {
    description = "SSH from allowed CIDRs"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = var.ssh_allowed_cidr_blocks
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow all outbound traffic"
  }

  tags = merge(
    var.shared_tags,
    {
      "TF"      = "true",
      "Type"    = "ssh",
      "Project" = var.project,
      "Github"  = "github.com/Mad-Pixels/applingo-api",
    }
  )
}

resource "aws_security_group_rule" "grafana_security_group" {
  count             = var.create_ssh_sg && var.create_grafana_sg ? 1 : 0
  type              = "ingress"
  from_port         = 3000
  to_port           = 3000
  protocol          = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.ssh_security_group[0].id
  description       = "Allow access to Grafana"
}
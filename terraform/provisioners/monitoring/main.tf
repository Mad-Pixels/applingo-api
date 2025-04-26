resource "aws_security_group" "egress_sg" {
  name        = "${local.project}-${local.provisioner}-egress"
  description = "Allow outbound traffic"
  vpc_id      = module.vpc.vpc_id

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = {
    "Type" = "egress",
    "Name" = "${local.project}-${local.provisioner}-egress",
  }

  depends_on = [module.vpc]
}

resource "aws_security_group" "ssh_sg" {
  count       = var.environment != "prd" ? 1 : 0
  name        = "${local.project}-${local.provisioner}-ssh"
  description = "Allow SSH access"
  vpc_id      = module.vpc.vpc_id

  ingress {
    from_port        = 22
    to_port          = 22
    protocol         = "tcp"
    cidr_blocks      = []
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = {
    "Type" = "ssh",
    "Name" = "${local.project}-${local.provisioner}-ssh"
  }

  depends_on = [module.vpc]
}

resource "aws_security_group" "ingress_sg" {
  name        = "${local.project}-${local.provisioner}-monitoring"
  vpc_id      = module.vpc.vpc_id
  description = "Allow access to monitoring services"

  ingress {
    from_port        = 80
    to_port          = 80
    protocol         = "tcp"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = {
    "Type" = "ingress",
    "Name" = "${local.project}-${local.provisioner}-ingress"
  }

  depends_on = [module.vpc]
}

module "vpc" {
  source = "../../modules/vpc"

  name        = "${local.project}-${local.provisioner}"
  shared_tags = local.tags

  vpc_addr_block = "10.100.100.0"
  vpc_zones      = 1
}

module "instance" {
  source = "../../modules/ec2"

  name        = "${local.project}-monitoring"
  shared_tags = local.tags

  key_name      = local.provisioner
  subnet_id     = element(module.vpc.public_subnets, 0)
  graviton_size = var.environment == "prd" ? "micro" : "nano"

  associate_public_ip_address = false

  security_group_ids = concat(
    [aws_security_group.ingress_sg.id],
    [aws_security_group.egress_sg.id],
    var.environment != "prd" ? [aws_security_group.ssh_sg[0].id] : []
  )

  user_data = file("${path.module}/scripts/init-instance.sh")
}
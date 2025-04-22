resource "aws_security_group" "egress_sg" {
  name        = "${local.project}-monitoring-egress"
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
    "TF"      = "true",
    "Type"    = "egress",
    "Project" = local.project,
    "Github"  = "github.com/Mad-Pixels/applingo-api",
    "Name"    = "${local.project}-monitoring-egress"
  }

  depends_on = [module.vpc]
}

resource "aws_security_group" "ssh_sg" {
  name        = "${local.project}-ssh"
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
    "TF"      = "true",
    "Type"    = "ssh",
    "Project" = local.project,
    "Github"  = "github.com/Mad-Pixels/applingo-api",
    "Name"    = "${local.project}-monitoring-ssh"
  }

  depends_on = [module.vpc]
}

resource "aws_security_group" "monitoring_sg" {
  name        = "${local.project}-monitoring"
  description = "Allow access to monitoring services"
  vpc_id      = module.vpc.vpc_id

  ingress {
    from_port        = 3000
    to_port          = 3000
    protocol         = "tcp"
    cidr_blocks      = []
    ipv6_cidr_blocks = ["::/0"]
  }

  ingress {
    from_port        = 9090
    to_port          = 9090
    protocol         = "tcp"
    cidr_blocks      = []
    ipv6_cidr_blocks = ["::/0"]
  }

  ingress {
    from_port        = 9100
    to_port          = 9100
    protocol         = "tcp"
    cidr_blocks      = []
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = {
    "TF"      = "true",
    "Type"    = "monitoring",
    "Project" = local.project,
    "Github"  = "github.com/Mad-Pixels/applingo-api",
    "Name"    = "${local.project}-monitoring"
  }

  depends_on = [module.vpc]
}

resource "aws_route53_record" "monitoring_aaaa" {
  zone_id = var.domain_zone_id 
  name    = var.domain_name
  type    = "AAAA"

  alias {
    name                   = module.distribution.domain_name
    zone_id                = module.distribution.hosted_zone_id
    evaluate_target_health = false
  }

  depends_on = [module.distribution]
}


module "vpc" {
  source = "../../modules/vpc"

  project = local.project
  name    = "${local.project}-monitoring"

  vpc_addr_block = "10.100.100.0"
  vpc_zones      = 1

  use_public_subnets      = true
  enable_dns_support      = true
  enable_dns_hostnames    = true
  enable_internet_gateway = true
  enable_nat_gateway      = false
  use_private_subnets     = false
}

module "instance" {
  source = "../../modules/ec2"

  project = local.project
  name    = "${local.project}-monitoring"

  use_localstack = var.use_localstack
  key_name       = local.project
  graviton_size  = "micro"

  security_group_ids = [
    aws_security_group.monitoring_sg.id,
    aws_security_group.egress_sg.id,
    aws_security_group.ssh_sg.id,
  ]

  subnet_id = element(module.vpc.public_subnets, 0)
  user_data = file("${path.module}/scripts/init-instance.sh")
}

module "distribution" {
  source = "../../modules/cloudfront"

  project = local.project
  name    = "monitoring"

  domain_name        = var.domain_name
  origin_domain_name = module.instance.instance_public_dns
  certificate_arn    = var.domain_acm_arn 
  forwarded_headers  = ["Host", "Authorization"]
}
module "vpc-infra" {
  source = "../../modules/vpc"

  project     = local.project
  aws_region  = var.aws_region
  name        = "platform-${local.project}"
  vpc_base_ip = "10.100.0.0"
  vpc_zones   = 1

  use_public_subnets      = true
  enable_dns_support      = true
  enable_dns_hostnames    = true
  enable_internet_gateway = true

  create_ssh_sg      = true
  create_grafana_sg  = true
  create_endpoint_sg = true
}

module "ec2-monitoring" {
  source = "../../modules/ec2"

  subnet_id = module.vpc-infra.subnet_ids[0]
  name      = "monitoring-${local.project}"
  project   = local.project
  key_name  = local.project

  security_group_ids = compact([
    module.vpc-infra.allow_ssh_ipv6,
  ])

  graviton_size  = "micro"
  use_localstack = var.use_localstack


  //user_data = file("${path.module}/scripts/ec2-monitoring.sh")
}
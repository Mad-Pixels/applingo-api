module "vpc-infra" {
  source = "../../modules/vpc"

  project     = local.project
  aws_region  = var.aws_region
  vpc_base_ip = "10.100.0.0"
  vpc_zones   = 1

  use_public_subnets      = true
  enable_dns_support      = true
  enable_dns_hostnames    = true
  enable_internet_gateway = true
  
  create_ssh_sg      = true
  create_endpoint_sg = true
}

module "ec2-monitoring" {
  source = "../../modules/ec2"

  project         = local.project
  subnet_id       = module.vpc-infra.subnet_ids[0]

  security_group_ids = compact([
    module.vpc-infra.ssh_security_group_id,
    module.vpc-infra.endpoint_security_group_id
  ])

  graviton_size   = "micro"
  use_localstack  = var.use_localstack
}
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
data "aws_availability_zones" "azs" {
  state = "available"
}

locals {
  vpc_cidr_block = "${var.vpc_addr_block}/23"

  public_blocks = [
    cidrsubnet(local.vpc_cidr_block, 3, 0),
    cidrsubnet(local.vpc_cidr_block, 3, 1),
    cidrsubnet(local.vpc_cidr_block, 3, 2),
  ]

  private_blocks = [
    cidrsubnet(local.vpc_cidr_block, 3, 3),
    cidrsubnet(local.vpc_cidr_block, 3, 4),
    cidrsubnet(local.vpc_cidr_block, 3, 5),
  ]

  used_public_cidrs  = slice(local.public_blocks, 0, var.vpc_zones)
  used_private_cidrs = var.use_private_subnets ? slice(local.private_blocks, 0, var.vpc_zones) : []

  igw_id            = var.enable_internet_gateway ? aws_internet_gateway.this[0].id : null
  nat_gateway_count = var.enable_nat_gateway && var.use_private_subnets ? 1 : 0
}
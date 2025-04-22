data "aws_availability_zones" "azs" {
  state = "available"
}

resource "aws_vpc" "this" {
  cidr_block           = local.vpc_cidr_block
  enable_dns_support   = var.enable_dns_support
  enable_dns_hostnames = var.enable_dns_hostnames

  assign_generated_ipv6_cidr_block = true

  tags = merge(
    var.shared_tags,
    {
      "TF"           = "true",
      "Name"         = var.name,
      "Project"      = var.project,
      "DnsSupport"   = var.enable_dns_support,
      "DnsHostnames" = var.enable_dns_hostnames,
      "Github"       = "github.com/Mad-Pixels/applingo-api",
    }
  )
}

resource "aws_internet_gateway" "this" {
  count  = var.enable_internet_gateway ? 1 : 0
  vpc_id = aws_vpc.this.id

  tags = merge(
    var.shared_tags,
    {
      "TF"      = "true",
      "Name"    = var.name,
      "Project" = var.project,
      "Github"  = "github.com/Mad-Pixels/applingo-api",
    }
  )
}
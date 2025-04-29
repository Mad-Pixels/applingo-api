resource "aws_subnet" "public" {
  count = var.use_public_subnets ? var.vpc_zones : 0

  ipv6_cidr_block   = cidrsubnet(aws_vpc.this.ipv6_cidr_block, 8, count.index)
  availability_zone = data.aws_availability_zones.azs.names[count.index]
  cidr_block        = local.used_public_cidrs[count.index]
  vpc_id            = aws_vpc.this.id

  assign_ipv6_address_on_creation = true

  tags = merge(
    var.shared_tags,
    {
      "Type" = "public",
      "Name" = "${var.name}-public",
    }
  )

  depends_on = [aws_vpc.this]
}

resource "aws_route_table" "public" {
  count = var.use_public_subnets && var.enable_internet_gateway ? 1 : 0

  vpc_id = aws_vpc.this.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = local.igw_id
  }

  route {
    ipv6_cidr_block = "::/0"
    gateway_id      = local.igw_id
  }

  tags = merge(
    var.shared_tags,
    {
      "Type" = "public",
      "Name" = "${var.name}-public",
    }
  )
}

resource "aws_route_table_association" "public" {
  count = var.use_public_subnets && var.enable_internet_gateway ? length(aws_subnet.public) : 0

  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public[0].id
}
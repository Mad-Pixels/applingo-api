resource "aws_subnet" "private" {
  count = var.use_private_subnets ? var.vpc_zones : 0

  vpc_id     = aws_vpc.this.id
  cidr_block = local.used_private_cidrs[count.index]

  assign_ipv6_address_on_creation = true

  ipv6_cidr_block   = cidrsubnet(aws_vpc.this.ipv6_cidr_block, 8, count.index + var.vpc_zones)
  availability_zone = data.aws_availability_zones.azs.names[count.index]

  tags = merge(
    var.shared_tags,
    {
      "Type" = "private",
      "Name" = "${var.name}-private",
    }
  )

  depends_on = [aws_vpc.this]
}

resource "aws_eip" "nat" {
  count  = local.nat_gateway_count
  domain = "vpc"

  tags = merge(
    var.shared_tags,
    {
      "Type" = "nat",
      "Name" = var.name,
    }
  )

  depends_on = [aws_internet_gateway.this]
}

resource "aws_nat_gateway" "this" {
  count         = local.nat_gateway_count
  allocation_id = aws_eip.nat[0].id
  subnet_id     = aws_subnet.public[0].id

  tags = merge(
    var.shared_tags,
    {
      "Type" = "nat",
      "Name" = "${var.name}-nat",
    }
  )

  depends_on = [aws_internet_gateway.this]
}

resource "aws_route_table" "private" {
  count  = var.use_private_subnets ? 1 : 0
  vpc_id = aws_vpc.this.id

  dynamic "route" {
    for_each = var.enable_nat_gateway ? [1] : []
    content {
      cidr_block     = "0.0.0.0/0"
      nat_gateway_id = aws_nat_gateway.this[0].id
    }
  }

  tags = merge(
    var.shared_tags,
    {
      "Type" = "private",
      "Name" = "${var.name}-private",
    }
  )
}

resource "aws_route_table_association" "private" {
  count          = var.use_private_subnets ? length(aws_subnet.private) : 0
  subnet_id      = aws_subnet.private[count.index].id
  route_table_id = aws_route_table.private[0].id
}
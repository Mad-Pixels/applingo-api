data "aws_availability_zones" "region_az_list" {
  state = "available"
}

resource "aws_subnet" "public" {
  count  = var.use_public_subnets ? var.vpc_zones : 0
  vpc_id = aws_vpc.this.id

  assign_ipv6_address_on_creation = true

  ipv6_cidr_block   = cidrsubnet(aws_vpc.this.ipv6_cidr_block, 8, count.index)
  cidr_block        = cidrsubnet("${var.vpc_base_ip}/16", 4, count.index)
  availability_zone = data.aws_availability_zones.region_az_list.names[count.index]

  tags = merge(
    var.shared_tags,
    {
      "TF"      = "true",
      "Type"    = "public",
      "Project" = var.project,
    }
  )

  depends_on = [aws_vpc.this]
}

resource "aws_route_table" "public" {
  count  = var.use_public_subnets && var.enable_internet_gateway ? 1 : 0
  vpc_id = aws_vpc.this.id

  route {
    cidr_block      = "0.0.0.0/0"
    ipv6_cidr_block = "::/0"
    gateway_id      = aws_internet_gateway.this[0].id
  }

  tags = merge(
    var.shared_tags,
    {
      "TF"      = "true",
      "Type"    = "public",
      "Project" = var.project,
      "Github"  = "github.com/Mad-Pixels/applingo-api",
    }
  )
}

resource "aws_route_table_association" "public" {
  count          = var.use_public_subnets && var.enable_internet_gateway ? length(aws_subnet.public) : 0
  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

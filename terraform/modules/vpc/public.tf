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


resource "aws_network_acl" "public" {
  vpc_id = aws_vpc.this.id

  tags = merge(
    var.shared_tags,
    {
      "Type" = "public-acl",
      "Name" = "${var.name}-public-acl",
    }
  )
}

resource "aws_network_acl_rule" "public_inbound_allow_all" {
  network_acl_id = aws_network_acl.public.id
  rule_number    = 100
  egress         = false
  protocol       = "-1"
  rule_action    = "allow"
  cidr_block     = "0.0.0.0/0"
}

resource "aws_network_acl_rule" "public_inbound_allow_ipv6_all" {
  network_acl_id  = aws_network_acl.public.id
  rule_number     = 101
  egress          = false
  protocol        = "-1"
  rule_action     = "allow"
  ipv6_cidr_block = "::/0"
}

resource "aws_network_acl_rule" "public_outbound_allow_all" {
  network_acl_id = aws_network_acl.public.id
  rule_number    = 100
  egress         = true
  protocol       = "-1"
  rule_action    = "allow"
  cidr_block     = "0.0.0.0/0"
}

resource "aws_network_acl_rule" "public_outbound_allow_ipv6_all" {
  network_acl_id  = aws_network_acl.public.id
  rule_number     = 101
  egress          = true
  protocol        = "-1"
  rule_action     = "allow"
  ipv6_cidr_block = "::/0"
}

resource "aws_network_acl_rule" "public_inbound_allow_metadata" {
  network_acl_id = aws_network_acl.public.id
  rule_number    = 110
  egress         = false
  protocol       = "-1"
  rule_action    = "allow"
  cidr_block     = "169.254.0.0/16"
}

resource "aws_network_acl_rule" "public_outbound_allow_metadata" {
  network_acl_id = aws_network_acl.public.id
  rule_number    = 110
  egress         = true
  protocol       = "-1"
  rule_action    = "allow"
  cidr_block     = "169.254.0.0/16"
}

resource "aws_network_acl_association" "public" {
  count = var.use_public_subnets ? length(aws_subnet.public) : 0

  subnet_id      = aws_subnet.public[count.index].id
  network_acl_id = aws_network_acl.public.id
}
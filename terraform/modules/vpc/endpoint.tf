resource "aws_vpc_endpoint" "s3" {
  count             = var.enable_s3_endpoint ? 1 : 0
  vpc_id            = aws_vpc.this.id
  service_name      = "com.amazonaws.${var.region}.s3"
  vpc_endpoint_type = "Gateway"

  route_table_ids = concat(
    var.use_public_subnets && var.enable_internet_gateway ? [aws_route_table.public[0].id] : [],
    var.use_private_subnets && var.enable_nat_gateway ? [aws_route_table.private[0].id] : []
  )

  tags = merge(var.shared_tags, { Name = "${var.name}-s3-endpoint" })
}

resource "aws_vpc_endpoint" "cloudwatch" {
  count             = var.enable_cloudwatch_endpoint ? 1 : 0
  vpc_id            = aws_vpc.this.id
  service_name      = "com.amazonaws.${var.region}.monitoring"
  vpc_endpoint_type = "Interface"

  subnet_ids = var.use_private_subnets ? aws_subnet.private[*].id : aws_subnet.public[*].id

  security_group_ids = [
    aws_security_group.vpc_endpoints[0].id
  ]

  private_dns_enabled = true

  tags = merge(var.shared_tags, { Name = "${var.name}-cloudwatch-endpoint" })
}

resource "aws_vpc_endpoint" "sts" {
  count             = var.enable_sts_endpoint ? 1 : 0
  vpc_id            = aws_vpc.this.id
  service_name      = "com.amazonaws.${var.region}.sts"
  vpc_endpoint_type = "Interface"

  subnet_ids = var.use_private_subnets ? aws_subnet.private[*].id : aws_subnet.public[*].id

  security_group_ids = [
    aws_security_group.vpc_endpoints[0].id
  ]

  private_dns_enabled = true

  tags = merge(var.shared_tags, { Name = "${var.name}-sts-endpoint" })
}

resource "aws_vpc_endpoint" "iam" {
  count             = var.enable_iam_endpoint ? 1 : 0
  vpc_id            = aws_vpc.this.id
  service_name      = "com.amazonaws.${var.region}.iam"
  vpc_endpoint_type = "Interface"

  subnet_ids = var.use_private_subnets ? aws_subnet.private[*].id : aws_subnet.public[*].id

  security_group_ids = [
    aws_security_group.vpc_endpoints[0].id
  ]

  private_dns_enabled = true

  tags = merge(var.shared_tags, { Name = "${var.name}-iam-endpoint" })
}

resource "aws_security_group" "vpc_endpoints" {
  count       = (var.enable_cloudwatch_endpoint || var.enable_sts_endpoint || var.enable_iam_endpoint) ? 1 : 0
  name        = "${var.name}-vpc-endpoints-sg"
  description = "Security group for VPC endpoints"
  vpc_id      = aws_vpc.this.id

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = [aws_vpc.this.cidr_block]
  }

  tags = merge(var.shared_tags, { Name = "${var.name}-vpc-endpoints-sg" })
}
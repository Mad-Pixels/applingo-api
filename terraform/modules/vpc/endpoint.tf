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

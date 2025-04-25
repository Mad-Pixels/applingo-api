output "vpc_id" {
  value = aws_vpc.this.id
}

output "vpc_ipv4_cidr_block" {
  value = aws_vpc.this.cidr_block
}

output "vpc_ipv6_cidr_block" {
  value = aws_vpc.this.ipv6_cidr_block
}

output "subnet_azs" {
  value = aws_subnet.public[*].availability_zone
}

output "public_subnets" {
  value = aws_subnet.public[*].id
}

output "public_subnet_cidrs" {
  value = aws_subnet.public[*].cidr_block
}

output "public_subnet_ipv6_cidrs" {
  value = aws_subnet.public[*].ipv6_cidr_block
}

output "private_subnets" {
  value = var.use_private_subnets ? aws_subnet.private[*].id : []
}

output "private_subnet_cidrs" {
  value = var.use_private_subnets ? aws_subnet.private[*].cidr_block : []
}

output "private_subnet_ipv6_cidrs" {
  value = var.use_private_subnets ? aws_subnet.private[*].ipv6_cidr_block : []
}

output "nat_gateway_id" {
  value = local.nat_gateway_count > 0 ? aws_nat_gateway.this[0].id : null
}

output "nat_gateway_public_ip" {
  value = local.nat_gateway_count > 0 ? aws_eip.nat[0].public_ip : null
}
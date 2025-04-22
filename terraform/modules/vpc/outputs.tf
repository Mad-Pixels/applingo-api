output "vpc_id" {
  value = aws_vpc.this.id
}

output "subnet_ids" {
  value = aws_subnet.public[*].id
}

output "vpc_cidr_block" {
  value = aws_vpc.this.cidr_block
}

output "public_subnets" {
  value = aws_subnet.public[*].id
}

output "vpc_ipv6_cidr_block" {
  value = aws_vpc.this.ipv6_cidr_block
}

output "subnet_azs" {
  value = aws_subnet.public[*].availability_zone
}

output "allow_ssh_ipv6" {
  value       = var.create_ssh_sg ? aws_security_group.ssh_security_group[0].id : ""
  description = "Security group ID for SSH over IPv6"
}
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

output "endpoint_security_group_id" {
  description = "ID of the endpoint security group (if created)"
  value       = var.create_endpoint_sg ? aws_security_group.endpoint_security_group[0].id : null
}

output "ssh_security_group_id" {
  description = "ID of the SSH security group (if created)"
  value       = var.create_ssh_sg ? aws_security_group.ssh_security_group[0].id : null
}
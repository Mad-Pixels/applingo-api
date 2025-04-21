output "instance_id" {
  value = aws_instance.this.id
}

output "public_ipv6" {
  value = length(aws_instance.this.ipv6_addresses) > 0 ? aws_instance.this.ipv6_addresses[0] : null
}

output "private_ip" {
  value = aws_instance.this.private_ip
}

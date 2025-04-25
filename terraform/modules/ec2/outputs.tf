output "instance_id" {
  value = aws_instance.this.id
}

output "public_ipv6" {
  value = length(aws_instance.this.ipv6_addresses) > 0 ? aws_instance.this.ipv6_addresses[0] : null
}

output "public_ipv4" {
  value = aws_instance.this.associate_public_ip_address ? aws_instance.this.public_ip : ""
}

output "dns" {
  value = aws_instance.this.public_dns
}

output "private_ip" {
  value = aws_instance.this.private_ip
}

output "instance_public_dns" {
  value       = length(aws_instance.this.ipv6_addresses) > 0 ? "${replace(aws_instance.this.ipv6_addresses[0], ":", "-")}.sslip.io" : null
  description = "IPv6 DNS name for the instance"
}
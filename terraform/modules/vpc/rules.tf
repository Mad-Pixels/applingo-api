resource "aws_security_group_rule" "allow_ssh_ipv6" {
  count             = var.create_ssh_sg ? 1 : 0
  type              = "ingress"
  from_port         = 22
  to_port           = 22
  protocol          = "tcp"
  cidr_blocks       = []
  ipv6_cidr_blocks  = ["::/0"]
  security_group_id = aws_security_group.ssh_security_group[0].id
}
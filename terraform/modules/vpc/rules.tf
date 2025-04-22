resource "aws_security_group_rule" "allow_ssh_ipv6" {
  count             = var.allow_ssh ? 1 : 0
  type              = "ingress"
  from_port         = 22
  to_port           = 22
  protocol          = "tcp"
  cidr_blocks       = []
  ipv6_cidr_blocks  = ["::/0"]
  security_group_id = aws_security_group.ssh_security_group[0].id
}

resource "aws_security_group_rule" "allow_all_outbound_ipv6" {
  count             = var.allow_egress ? 1 : 0
  type              = "egress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = []
  ipv6_cidr_blocks  = ["::/0"]
  security_group_id = aws_security_group.egress_security_group[0].id
}

resource "aws_security_group_rule" "allow_all_outbound_ipv4" {
  count             = var.allow_egress ? 1 : 0
  type              = "egress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
  ipv6_cidr_blocks  = []
  security_group_id = aws_security_group.egress_security_group[0].id
}
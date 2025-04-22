resource "aws_security_group" "ssh_security_group" {
  count       = var.create_ssh_sg ? 1 : 0
  name        = "${var.project}-ssh"
  description = "Allow SSH"
  vpc_id      = aws_vpc.this.id

  tags = merge(
    var.shared_tags,
    {
      "TF"      = "true",
      "Type"    = "ssh",
      "Project" = var.project,
      "Github"  = "github.com/Mad-Pixels/applingo-api",
      "Name"    = "${var.name}-ssh"
    }
  )
}

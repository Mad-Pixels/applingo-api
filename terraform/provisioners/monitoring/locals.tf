data "aws_availability_zones" "azs" {
  state = "available"
}

locals {
  project     = "applingo"
  provisioner = "monitoring"

  tags = {
    "TF"          = "true",
    "Project"     = local.project,
    "Environment" = var.environment,
    "Provisioner" = local.provisioner,
    "Github"      = "github.com/Mad-Pixels/applingo-api",
  }
}
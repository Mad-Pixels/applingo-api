resource "aws_security_group" "egress_sg" {
  name        = "${local.project}-${local.provisioner}-egress"
  description = "Allow outbound traffic"
  vpc_id      = module.vpc.vpc_id

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = {
    "Type" = "egress",
    "Name" = "${local.project}-${local.provisioner}-egress",
  }

  depends_on = [module.vpc]
}

resource "aws_security_group" "ssh_sg" {
  count       = var.environment != "prd" ? 1 : 0
  name        = "${local.project}-${local.provisioner}-ssh"
  description = "Allow SSH access"
  vpc_id      = module.vpc.vpc_id

  ingress {
    from_port        = 22
    to_port          = 22
    protocol         = "tcp"
    cidr_blocks      = []
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = {
    "Type" = "ssh",
    "Name" = "${local.project}-${local.provisioner}-ssh"
  }

  depends_on = [module.vpc]
}

resource "aws_security_group" "ingress_sg" {
  name        = "${local.project}-${local.provisioner}-monitoring"
  vpc_id      = module.vpc.vpc_id
  description = "Allow access to monitoring services"

  ingress {
    from_port        = 80
    to_port          = 80
    protocol         = "tcp"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = {
    "Type" = "ingress",
    "Name" = "${local.project}-${local.provisioner}-ingress"
  }

  depends_on = [module.vpc]
}

resource "aws_iam_role" "monitoring_instance_role" {
  name = "${local.project}-${local.provisioner}-ec2-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ec2.amazonaws.com"
      }
    }]
  })

  tags = local.tags
}

resource "aws_iam_policy" "monitoring_instance_policy" {
  name = "${local.project}-${local.provisioner}-ec2-policy"

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "ec2:DescribeTags"
        ],
        Resource = "*"
      },
      {
        Effect = "Allow",
        Action = [
          "s3:ListBucket"
        ],
        Resource = module.s3-monitoring-bucket.s3_arn
      },
      {
        Effect = "Allow",
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject"
        ],
        Resource = "${module.s3-monitoring-bucket.s3_arn}/*"
      },
      {
        Effect = "Allow",
        Action = [
          "cloudwatch:GetMetricData",
          "cloudwatch:ListMetrics",
          "cloudwatch:GetMetricStatistics"
        ],
        Resource = "*"
      }
    ]
  })

  tags = local.tags

  depends_on = [module.s3-monitoring-bucket]
}

resource "aws_iam_role_policy_attachment" "monitoring_instance_attach" {
  role       = aws_iam_role.monitoring_instance_role.name
  policy_arn = aws_iam_policy.monitoring_instance_policy.arn
}

resource "aws_iam_instance_profile" "monitoring_instance_profile" {
  name = "${local.project}-${local.provisioner}-ec2-profile"
  role = aws_iam_role.monitoring_instance_role.name
}

module "vpc" {
  source = "../../modules/vpc"

  name           = "${local.project}-${local.provisioner}"
  region         = var.aws_region
  shared_tags    = local.tags
  vpc_addr_block = "10.100.100.0"
  vpc_zones      = 1

  enable_s3_endpoint = true
}

module "s3-monitoring-bucket" {
  source = "../../modules/s3"

  project     = local.project
  shared_tags = local.tags
  bucket_name = "${local.provisioner}-${var.environment}"
}

module "instance" {
  source = "../../modules/ec2"

  name        = "${local.project}-monitoring"
  shared_tags = local.tags

  key_name      = local.provisioner
  subnet_id     = element(module.vpc.public_subnets, 0)
  graviton_size = var.environment == "prd" ? "micro" : "nano"

  associate_public_ip_address = false

  security_group_ids = concat(
    [aws_security_group.ingress_sg.id],
    [aws_security_group.egress_sg.id],
    var.environment != "prd" ? [aws_security_group.ssh_sg[0].id] : []
  )

  instance_profile = aws_iam_instance_profile.monitoring_instance_profile.name
  user_data        = file("${path.module}/scripts/init-instance.sh")

  depends_on = [module.s3-monitoring-bucket]
}
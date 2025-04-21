locals {
  selected_ami_id = var.ami_id != "" ? var.ami_id : data.aws_ami.amazon_linux.id
  selected_root_volume_size = var.ami_id != "" ? var.volume_size : element(tolist(data.aws_ami.amazon_linux.block_device_mappings), 0).ebs.volume_size
}
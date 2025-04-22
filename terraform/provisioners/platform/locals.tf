locals {
  project      = "applingo"
  state_bucket = "tfstates-${local.project}"
  tfstate_file = "platform.tfstates"
}
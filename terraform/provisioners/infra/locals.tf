locals {
  project      = "lingocards"
  state_bucket = "tfstates-${local.project}"
  tfstate_file = "infra.tfstates"
}
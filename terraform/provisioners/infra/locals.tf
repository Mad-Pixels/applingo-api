locals {
  project      = "lingocards"
  state_bucket = "tfstates-${local.project}"
  tfstate_file = "infra.tfstates"

  dictionary_dynamo_schema = jsondecode(
    file("${path.module}/../../../data/.tmpl/dynamo_dictionary_table.json")
  )
}
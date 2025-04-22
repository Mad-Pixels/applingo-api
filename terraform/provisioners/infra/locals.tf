locals {
  project      = "applingo"
  state_bucket = "tfstates-${local.project}"
  tfstate_file = "infra.tfstates"

  dictionary_dynamo_schema = jsondecode(
    file("${path.module}/../../../dynamodb-interface/.tmpl/dynamo_dictionary_table.json")
  )

  processing_dynamo_schema = jsondecode(
    file("${path.module}/../../../dynamodb-interface/.tmpl/dynamo_processing_table.json")
  )

  profile_dynamo_schema = jsondecode(
    file("${path.module}/../../../dynamodb-interface/.tmpl/dynamo_profile_table.json")
  )
}
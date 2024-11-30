locals {
  project      = "applingo"
  state_bucket = "tfstates-${local.project}"
  tfstate_file = "infra.tfstates"

  dictionary_dynamo_schema = jsondecode(
    file("${path.module}/../../../dynamodb-interface/.tmpl/dynamo_dictionary_table.json")
  )

  subcategory_dynamo_schema = jsondecode(
    file("${path.module}/../../../dynamodb-interface/.tmpl/dynamo_subcategory_table.json")
  )

  level_dynamo_schema = jsondecode(
    file("${path.module}/../../../dynamodb-interface/.tmpl/dynamo_level_table.json")
  )
}
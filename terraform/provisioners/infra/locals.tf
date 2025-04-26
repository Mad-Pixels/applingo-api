locals {
  project     = "applingo"
  provisioner = "infra"

  tags = {
    "TF"          = "true",
    "Project"     = local.project,
    "Environment" = var.environment,
    "Provisioner" = local.provisioner,
    "Github"      = "github.com/Mad-Pixels/applingo-api",
  }

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
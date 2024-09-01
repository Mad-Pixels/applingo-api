locals {
  lambda_functions = toset([for f in fileset("${path.root}/../cmd", "*") : f if fileexists("${path.root}/../cmd/${f}/main.go")])
  region           = "eu-central-1"
}
locals {
    root_dir = "${path.module}/../../../cmd"
    
    lambda_functions = toset([for d in fileset(local.root_dir, "*") : d if fileexists("${local.root_dir}/${d}/.infra/config.tf")])
    lambda_configs   = {
        for func in local.lambda_functions :
        func => jsondecode(file("${local.root_dir}/${func}/.infra/config.json"))
    }
}
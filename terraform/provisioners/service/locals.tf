locals {
    root_dir        = "${path.module}/../../../cmd"
    all_directories = fileset(local.root_dir, "*")
  
    lambda_functions = toset([
        for d in local.all_directories : d
        if fileexists("${local.root_dir}/${d}/.infra/config.json")
    ])
  
    lambda_configs = {
        for func in local.lambda_functions :
        func => jsondecode(file("${local.root_dir}/${func}/.infra/config.json"))
    }
}
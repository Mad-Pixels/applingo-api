locals {
    _root_dir = "${path.module}/../../../cmd"
    _entries  = fileset(local._root_dir, "**")

    _lambda_functions = distinct([
        for v in local._entries : split("/", v)[0]
        if length(split("/", v)) > 1
    ])
    _lambda_configs = {
        for func in local._lambda_functions : 
        func => fileexists("${local._root_dir}/${func}/.infra/config.json") ? jsondecode(file("${local._root_dir}/${func}/.infra/config.json")) : null
    }

    lambdas        = { for func in local._lambda_functions : func => local._lambda_configs[func]}
    state_bucket   = "tfstates-lingocards"
    tfstate_file   = "api.tfstates"
}
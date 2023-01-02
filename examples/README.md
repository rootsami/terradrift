# Examples

Those examples are how terradrift works in different environments style. Those resources are meaningless, however the structure is the same.
Terraform backend is local state file and you can test it with `terradrift-cli` or `terradrift-server`.

Stacks with tfvars file is considered as a single stack, and the stack name is the file name without extension.
Any environment with tfvars and hcl could be ran with regular terraform commands for example:

```bash
$ terraform init -backend-config=path/to/backend.hcl

# OR

$ terraform plan -var-file=path/to/env.tfvars

```


Terradrit is supposed to detect tfvars file and hcl and generate the configs based on that.

# terradrift-cli

terradrift-cli is a command line tool that can be used to detect drifts in terraform IaC. It discovers the stacks from the given workdir and then runs the `terraform plan` command to detect the drifts based on the plan output. The output format by default is table, however it can be changed to json or yaml by passing the `--output` flag.


## Installation

### build from source

```bash

$ git clone https://github.com/rootsami/terradrift.git
$ cd terradrift/terradrift-cli
$ go build -o terradrift-cli

```

### Binary

Download the binary from the [releases]()

## Usage
it can be used in multiple scenarios:
- Without a config file, where you only need to provide the workdir and it will discover any directory that has tf tfstack.
- With a config file, where you can define the stacks and their configurations in the config file. ex. name, path, tfvars, etc.
- Only to generate a config file based on a provided workdir. where it will try to generate the name and path and wheather the stack has tfvars file or not as it will consider each tfvars file as another tfstack.

```bash

usage: terradrift-cli [<flags>]

A command-line tool to detect drifts in terraform IaC

Flags:
  --help                  Show context-sensitive help (also try --help-long and --help-man).
  --workdir="./"        workdir of a project that contains all terraform directories
  --config=CONFIG         Path for configuration file holding the stack information
  --extra-backend-vars=EXTRA-BACKEND-VARS ...  
                          Extra backend environment variables ex. GOOGLE_CREDENTIALS, AWS_ACCESS_KEY or AWS_SECRET_KEY
  --debug                 Enable debug mode
  --generate-config-only  Generate a config file based on a provided worksapce
  --output=table          Output format supported: json, yaml and table

```

## Examples

### Run terradrift-cli

```bash

$ terradrift-cli --workdir ./examples/ --config examples/config.yaml        
STACK-NAME      DRIFT   ADD     CHANGE  DESTROY PATH                    TF-VERSION 
api-production  false   0       0       0       gcp/api                 1.2.7     
api-staging     false   0       0       0       gcp/api                 1.2.7     
core-production true    0       0       1       aws/core-production     1.2.7     
core-staging    true    1       0       0       gcp/core-staging        1.0.6

```

### Run terradrift-cli to generate config file only
`--generate-config-only` flag can be used to generate a config file based on the provided workdir. The generated config file can be used to run terradrift-cli or terradrift-server also to hand pick the stacks to be scanned.

```bash

$ terradrift-cli --workdir ./ --generate-config-only

stacks:
- name: aws-core-production
  path: aws/core-production
- name: gcp-api-environments-production
  path: gcp/api
  tfvars: environments/production.tfvars
  backend: environments/production.hcl
- name: gcp-api-environments-staging
  path: gcp/api
  tfvars: environments/staging.tfvars
  backend: environments/staging.hcl
- name: gcp-core-staging
  path: gcp/core-staging

```


### Run terradrift-cli with json output

```bash

$ terradrift-cli --workdir ./ --extra-backend-vars GOOGLE_CREDENTIALS=$SERVICE_ACCOUNT_PATH --output json
[
  {
    "name": "examples-gcp-api-environments-staging",
    "path": "examples/gcp/api",
    "drift": false,
    "add": 0,
    "change": 0,
    "destroy": 0,
    "tfver": "1.2.7"
  },
  {
    "name": "examples-gcp-api-environments-production",
    "path": "examples/gcp/api",
    "drift": false,
    "add": 0,
    "change": 0,
    "destroy": 0,
    "tfver": "1.2.7"
  },
  {
    "name": "examples-aws-core-production",
    "path": "examples/aws/core-production",
    "drift": true,
    "add": 0,
    "change": 0,
    "destroy": 1,
    "tfver": "1.2.7"
  },
  {
    "name": "examples-gcp-core-staging",
    "path": "examples/gcp/core-staging",
    "drift": true,
    "add": 1,
    "change": 0,
    "destroy": 0,
    "tfver": "1.0.6"
  }
]

```

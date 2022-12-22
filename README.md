# terradrift
A tool that will navigate through all terraform directories (stacks) to run terraform plan to detect the current drift between the committed code and applied infrastructure.

As we’re adding/modifying infrastructure pieces, sometimes we encounter a few terraform stacks (directory) that are drifting from the current infrastructure. That’s mainly caused by changes done manually or resources deleted as they’re no longer needed without any track on our code base. 

Some of these changes have been done for testing purposes and then forgotten, or have been sitting there for a longer period that the responsible team has lost the track of it.

# How it works?
Terradrift has two modes, CLI and Server. Both modes will scan all terraform stacks (directories) in a given workdir and run terraform plan to detect the drifts. The difference between the two modes is that the CLI mode will run the scan once and exit, while the Server mode will run the scan continuously based on a defined schedule and expose a rest api to query the drifts. It also exposes the drift results as prometheus metrics on /metrics endpoint.

## Server mode (terradrift-server)

### Installation / Deployment
The design of the service should be a running server that continuously runs those drift scans.
### Configurations
In [config.yaml](config.yaml), you have to define which repository that terradrift will scan those stacks from, also the stack's name, version and environment-specific variables.
example:
```yaml
stacks:
  - name: core-production
    path: aws/core-production

  - name: api-staging
    path: gcp/api
    tfvars: environments/staging.tfvars
    backend: environments/staging.hcl
```

## How to develop/use terradrift-server locally?
After setting up the configs in the above section, or if you need to generate a config file, you can use [terradrift-cli](terradrift-cli/README.md) to generate it.
you can run the server locally by running the following command:
```bash
$ git clone https://github.com/rootsami/terradrift.git
$ cd terradrift/terradrift-server
$ go build -o terradrift-server
$ ./terradrift-server --help

usage: terradrift-server --repository=REPOSITORY --git-token=GIT-TOKEN [<flags>]

A tool to detect drifts in terraform IaC, As a server mode it will expose a rest api to query the drifts and also as prometheus metrics on /metrics endpoint

Flags:
  --help                   Show context-sensitive help (also try --help-long and --help-man).
  --hostname="localhost"   hostname that apil will be exposed.
  --port="8080"            port of the service api is listening on
  --scheme="http"          The scheme of exposed endpoint http/https
  --repository=REPOSITORY  The git repository which include all terraform stacks
  --git-token=GIT-TOKEN    Personal access token to access git repositories
  --git-timeout=120        Wait timeout for git repoistory to clone or pull updates
  --interval=60            The interval for scan scheduler
  --config="config.yaml"   Path for configuration file holding the stack information
  --extra-backend-vars=EXTRA-BACKEND-VARS ...  
                           Extra backend environment variables ex. GOOGLE_CREDENTIALS, AWS_ACCESS_KEY or AWS_SECRET_KEY
  --debug                  Enable debug mode

```

## Examples
```bash
$ ./terradrift-server --repository https://github.com/username/reponame \
--git-token $GITHUB_AUTH_TOKEN \
--config ./config.yaml \
--extra-backend-vars GOOGLE_CREDENTIALS=$SERVICE_ACCOUNT_PATH \
--interval 600 \


```

It will start a local HTTP server on `http://localhost:8080`, where you can initiate HTTP requests and passing stackname in the URL. The response will be a JSON object with the drifts information.
```bash
$ curl http://localhost:8080/api/plan?stack=api-staging
{"drift":true,"add":1,"change":0,"destroy":0}

$ curl http://localhost:8080/api/plan?stack=core-production
{"drift":false,"add":0,"change":0,"destroy":1}
```

Retrieving the drifts as prometheus metrics
```bash
$ curl http://localhost:8080/metrics
# HELP terradrift_plan_add_resources Number of resources to be added based on tf plan
# TYPE terradrift_plan_add_resources gauge
terradrift_plan_add_resources{stack="api-staging"} 1
terradrift_plan_add_resources{stack="core-production"} 0
# HELP terradrift_plan_change_resources Number of resources to be changed based on tf plan
# TYPE terradrift_plan_change_resources gauge
terradrift_plan_change_resources{stack="api-staging"} 0
terradrift_plan_change_resources{stack="core-production"} 0
# HELP terradrift_plan_destroy_resources Number of resources to be destroyed based on tf plan
# TYPE terradrift_plan_destroy_resources gauge
terradrift_plan_destroy_resources{stack="api-staging"} 0
terradrift_plan_destroy_resources{stack="core-production"} 1
```

## CLI mode
See all details in [terradrift-cli](terradrift-cli/README.md)


## Roadmap
- [ ] Add support for multiple repositories

## Contributing

## License

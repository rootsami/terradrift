# terradrift-server

## Installation / Deployment
The design of the service should be a running server that continuously runs those drift scans.

### Helm
You can deploy the service using the helm chart in [deploy/helm/terradrift](../deploy/helm/terradrift/README.md)

### Compose
You can deploy the service using the docker-compose file in [deploy/docker-compose](../deploy/compose/README.md)


## How to develop/use terradrift-server locally?
You can run the server locally following the example below after setting the required environment variables for Github token and the cloud provider. ex. GOOGLE_CREDENTIALS, AWS_CONFIG_FILE, AWS_SHARED_CREDENTIALS_FILE, etc. similar to how you would run terraform.

After setting up the configs in the above section, or if you need to generate a config file, you can use [terradrift-cli](../terradrift-cli/README.md) to generate it.
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

### Configurations
In [config.yaml](config.yaml), you have to define which stacks that terradrift will scan, also the stack's name, version and environment-specific variables.
| Field | Description | Required |
| --- | --- | --- |
| stacks | List of stacks that terradrift will scan | Yes |
| stacks.name | Name of the stack | Yes |
| stacks.path | Path to the stack's terraform files | Yes |
| stacks.tfvars | Path to the stack's terraform variables file relative to the stack's path | No |
| stacks.backend | Path to the stack's terraform backend file relative to the stack's path | No |


Example:
```yaml
stacks:
  - name: core-production
    path: aws/core-production

  - name: api-staging
    path: gcp/api
    tfvars: environments/staging.tfvars
    backend: environments/staging.hcl
```

## Examples
```bash
$ ./terradrift-server --repository https://github.com/username/reponame \
--git-token $GITHUB_TOKEN \
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
# terradrift
A tool that will navigate through all terraform directories (stacks) to run terraform plan to detect the current drift between the committed code and applied infrastructure.

As we’re adding/modifying infrastructure pieces, sometimes we encounter a few terraform stacks (directory) that are drifting from the current infrastructure. That’s mainly caused by changes done manually or resources deleted as they’re no longer needed without any track on our code base. 

Some of these changes have been done for testing purposes and then forgotten, or have been sitting there for a longer period that the responsible team has lost the track

## Installation / Deployment
The design of the service should be a running server that continuously runs those drift scans.
### Configurations
In [config.yaml](config.yaml), you have to define which repository that terradrift will scan those stacks from, also the stack's name, version and environment-specific variables.
example:
```yaml
stacks:
  - name: stack-one
    path: gcp/stack-one

  - name: stack-three
    path: gcp/stack-multi
    tfvars: environments/three.tfvars
    backend: environments/three.hcl
```

## How to develop/use terradrift locally?
After setting up the configs in the above section
```bash
# go build .
./terradrift-server --help

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
                           Extra backend environment variables ex. GOOGLE_CREDENTIALS OR AWS_ACCESS_KEY
  --debug                  Enable debug mode
./terradrift-server --repository https://github.com/username/reponame --git-token $GITHUB_AUTH_TOKEN --config ./config.yaml --extra-backend-vars GOOGLE_CREDENTIALS=$SERVICE_ACCOUNT_PATH
```
It will start a local HTTP server on `http://localhost:8080`, where you can initiate terradrift calls HTTP request and passing stackname in the URL. 
```bash
# curl http://localhost:8080/api/plan?stack=stack-one
{"drift":true,"add":1,"change":0,"destroy":0}

# curl http://localhost:8080/api/plan?stack=stack-three
{"drift":false,"add":0,"change":0,"destroy":0}
```

## Roadmap
- [ ] Add support for multiple repositories


## Contributing

## License

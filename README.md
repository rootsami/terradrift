![terradrift-logo](https://user-images.githubusercontent.com/5733568/210888175-0b6d9e2a-e5fe-4e17-bf14-b6096705223a.png)

A tool that will navigate through all terraform directories (stacks) to run terraform plan to detect the current drift between the committed code and applied infrastructure.

As we’re adding/modifying infrastructure pieces, sometimes we encounter a few terraform stacks (directory) that are drifting from the current infrastructure. That’s mainly caused by changes done manually or resources deleted as they’re no longer needed without any track on our code base. 

Some of these changes have been done for testing purposes and then forgotten, or have been sitting there for a longer period that the responsible team has lost the track of it.

# How it works?
Terradrift has two modes, CLI and Server. Both modes will scan all terraform stacks (directories) in a given workdir and run terraform plan to detect the drifts. The difference between the two modes is that the CLI mode will run the scan once and exit, while the Server mode will run the scan continuously based on a defined schedule and expose a rest api to query the drifts. It also exposes the drift results as prometheus metrics on /metrics endpoint.

## Server mode (terradrift-server)
You can run the server following the example below after setting the required environment variables for Github token and the cloud provider.
### Example
```bash
$ ./terradrift-server --repository https://github.com/username/reponame \
--git-token $GITHUB_TOKEN \
--config ./config.yaml \
--interval 600 

```

It will start a local HTTP server on `http://localhost:8080`, where you can initiate HTTP requests and passing stackname in the URL. The response will be a JSON object with the drifts information.
```bash
$ curl http://localhost:8080/api/plan?stack=api-staging
{"drift":true,"add":1,"change":0,"destroy":0}

$ curl http://localhost:8080/api/plan?stack=core-production
{"drift":true,"add":0,"change":0,"destroy":1}
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
See all details in [terradrift-server](terradrift-server/README.md)

## CLI mode
terradrift-cli discovers the stacks from the given workdir and then runs the `terraform plan` command to detect the drifts based on the plan output.
### Example
```bash

$ terradrift-cli --workdir ./examples/ --config examples/config.yaml        
STACK-NAME      DRIFT   ADD     CHANGE  DESTROY PATH                    TF-VERSION 
api-production  false   0       0       0       gcp/api                 1.2.7     
api-staging     false   0       0       0       gcp/api                 1.2.7     
core-production true    0       0       1       aws/core-production     1.2.7     
core-staging    true    1       0       0       gcp/core-staging        1.0.6

```
See all details in [terradrift-cli](terradrift-cli/README.md)


## Roadmap
- [ ] Add support for multiple repositories for server mode
- [ ] Add support for multiple workdirs for cli mode

## License
[MIT](LICENSE)

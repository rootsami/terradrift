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
repository: "https://github.com/rootsami/terrad-examples"
stacks:
  - name: stack-one
    version: 1.0.6
    path: gcp/stack-one

  - name: stack-three
    version: 1.2.5
    path: gcp/stack-multi
    tfvars: environments/three.tfvars
    backend: environments/three.hcl
```

Also, you will have export `GITHUB_AUTH_TOKEN` as an environment variable to be able to checkout private repositories.

## How to develop/use terradrift locally?
After setting up the configs in the above section
```bash
# go build
# ./terradrift 
```
It will start a local HTTP server on `http://localhost:8080`, where you can initiate terradrift calls HTTP request and passing stackname in the URL. 
```bash
# curl http://localhost:8080/api/plan?stack=stack-one
"CHANGES DETECTED... Plan: 1 to add, 0 to change, 0 to destroy."

# curl http://localhost:8080/api/plan?stack=stack-three
"No changes. Infrastructure matches the configuration."
```

## Roadmap
This tool started as a fun project by Cumulus, but then it got INTERESTING! 🤩 

- Scheduled runs for all defined stacks.
- Once drift is detected, then what? For how long?
- Optimizations: No download/install for each run. it has to be once.


## Contributing

## License

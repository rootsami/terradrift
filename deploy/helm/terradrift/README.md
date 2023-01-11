# terradrift

A Helm chart for terradrift server mode - A tool that will navigate through all terraform directories (stacks) to run terraform plan to detect the current drift between the committed code and applied infrastructure.

## Getting Started

### Authentication to the git repository
In order to use this chart, you need to have a repository that contains all your terraform stacks.
You will need to provide a personal access token in order to clone the repository. You can create a personal access token in your github account
by going to Settings -> Developer settings -> Personal access tokens -> Generate new token.
To use the generated PAT, you will need to add it to your helm values file as follows:

```yaml
server:
  env:
    - name: GIT_TOKEN
      value: <your PAT>
```
Alternatively, you can use the `--set` flag in the helm install command.

Also you can use an already created secret in your cluster by adding the following to your helm values file:

```yaml

server:
  env:
    - name: GIT_TOKEN
      valueFrom:
        secretKeyRef:
          name: <your secret name>
          key: <your secret key>
```

### Authentication to the cloud provider
In order to be able to run terraform plan, you will need to have a cloud provider account and credentials. This part is flexible and you can be as creative as your authentication method requires.

You can use workloads-identity in GCP, or IAM roles in AWS. You could also use service accounts and static credentials. For example:

```yaml
server:
  env:
    - name: AWS_ACCESS_KEY_ID
      value: <your AWS access key>
    - name: AWS_SECRET_ACCESS_KEY
      value: <your AWS secret key>
    - name: AWS_DEFAULT_REGION
      value: <your AWS region>

    - name: GOOGLE_APPLICATION_CREDENTIALS
      value: <your GCP service account key file path>
```
Alternatively, you can use the `--set` flag in the helm install command. Bear in mind that the above is just an example and you can modify it to fit your needs. 
You might also need to mount extra volumes to the container in order to be able to access the credentials file. Detailed instructions can be found in the [values.yaml](values.yaml#L50) file.


### Scraping metrics
If you are using Prometheus Operator, you can scrape the metrics by enabling the serviceMonitor.

```yaml
serviceMonitor:
  enabled: true
  metricPath: /metrics
```

Alternatively, you can add the necessary static scrape configs to your Prometheus instance. For example:

```yaml
scrape_configs:
  - job_name: 'terradrift'
    scrape_interval: 300s
    static_configs:
      - targets: ['terradrift:8080']
```



## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` | affinity |
| fullnameOverride | string | `""` | Full name override |
| image.pullPolicy | string | `"IfNotPresent"` | Image pull policy |
| image.repository | string | `"rootsami/terradrift"` | Image repository |
| image.tag | string | `""` | Image tag |
| imagePullSecrets | list | `[]` | Image pull secrets |
| ingress.enabled | bool | `false` | Enable ingress |
| nameOverride | string | `""` | Name override |
| nodeSelector | object | `{}` | Node selector |
| podAnnotations | object | `{}` | Pod annotations |
| podSecurityContext | object | `{}` | Pod security context |
| replicaCount | int | `1` | Replica count |
| resources | object | `{}` | Resources |
| securityContext | object | `{}` | Security context |
| server.config | object | `{}` | A config that contains stack list and properties |
| server.debug | bool | `false` | Enable debug mode |
| server.env | list | `[]` | Environment variables to pass to the container |
| server.extraArgs | list | `[]` | Extra arguments to pass to the container |
| server.extraVolumeMounts | list | `[]` | Extra volume mounts to the container |
| server.extraVolumes | list | `[]` | Extra volumes to the pod |
| server.interval | int | `3600` | Interval in seconds to run the drift check |
| server.repository | string | `""` | Repository to use which contains terraform stacks |
| service.port | int | `8080` | Service port |
| service.type | string | `"ClusterIP"` | Service type |
| serviceAccount.annotations | object | `{}` | Annotations to add to the service account |
| serviceAccount.create | bool | `true` | Specifies whether a service account should be created |
| serviceAccount.name | string | `""` | The name of the service account to use. If not set and create is true, a name is generated using the fullname template |
| serviceMonitor.enabled | bool | `false` | Enable ServiceMonitor for Prometheus Operator |
| serviceMonitor.metricPath | string | `"/metrics"` | Path to metrics endpoint |
| tolerations | list | `[]` | Tolerations |

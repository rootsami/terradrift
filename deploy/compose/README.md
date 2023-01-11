# Compose

This directory contains the compose files for the terradrift-server mode. The compose files are used to run terradrift-server in a container and mounting the local directories into it including cloud provider credentials and exporting environment variables that would be used by terradrift-server in order to connect to the cloud provider.

Use .env file to export environment variables such as cloud provider credentials, github token, etc.

This part is flexible, you could use iam roles, workload-identity or other ways to authenticate to the cloud provider.

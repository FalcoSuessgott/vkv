# Azure Devops

Azure Devops Example for reading Secrets from Vault using `vkv`:

```yaml
resources:
  containers:
  - container: ghcr.io/falcosuessgott/vkv:latest
    image: vkv
    env:
      VAULT_ADDR: https://vault.server.de

      VKV_MODE: server
      VKV_SERVER_PATH: secrets
      VKV_LOGIN_COMMAND: |
        vault login -token-only -method=userpass username=admin password="${VAULT_PASSWORD}"
    ports:
      - 8080:8080

pool:
  vmImage: 'ubuntu-latest'

services:
  vkv: vkv

steps:
  - script: |
      eval $(curl http://vkv:8080)
      echo $secret
    displayName: Read secrets as env vars using vkv
```

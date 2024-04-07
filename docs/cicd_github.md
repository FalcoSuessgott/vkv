# Github Action

Github Action Example for reading Secrets from Vault using `vkv`:

```yaml
name: Vault Secrets using vkv
on: push

jobs:
  job_name:
    runs-on: ubuntu-latest
    services:
      vkv:
        image: ghcr.io/falcosuessgott/vkv:latest
        env:
          VAULT_ADDR: https://vault.server.de
          VKV_MODE: server
          VKV_SERVER_PATH: secrets
          VKV_LOGIN_COMMAND: |
            vault login -token-only -method=userpass username=admin password="${VAULT_PASSWORD}"
        ports:
          - 8080:8080
    steps:
      - name: read secrets from vkv server
        run: eval $(curl http://vkv:8080/export)
      - name: output secrets now available as env vars
        run: echo $secret
```
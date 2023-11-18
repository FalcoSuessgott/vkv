# Gitlab CI

Gitlab-CI Example for reading Secrets from Vault using vkv 

```yaml
variables:
  # vaults env vars
  # all of vault env vars are supported (https://developer.hashicorp.com/vault/docs/commands#environment-variables)
  # required:
  VAULT_ADDR: https://prod.vault.d4.sva.dev
  VAULT_NAMESPACE: "${CI_PROJECT_ROOT_NAMESPACE}"

  # command vkv uses to authenticate to vault, all vars are available
  VKV_LOGIN_COMMAND: vault write -field=token auth/jwt/login jwt="${VAULT_JWT_TOKEN}"

  # vault kv path to read secrets from
  VKV_SERVER_PATH: "secrets"

# default sets global default settings that are inherited to all jobs
default:
  # spin up a vkv service container in server mode, configure using variables/env vars
  services:
    - name: ghcr.io/falcosuessgott/vkv:v0.5.0
      command: ["server"]
      alias: vkv
  # global before_scripts block
  before_script: 
    # install curl, or wget in your job container   
    - apk add --no-cache curl

    # curl/wget vkv on /export, which will expot all secrets from VKV_SERVER_PATH, eval the output into your shell
    - eval $(curl http://vkv:8080/export)
  # global jwt token (https://docs.gitlab.com/ee/ci/examples/authenticating-with-hashicorp-vault/#example)
  id_tokens:
    # set jwt aud field to gitlab ci server host
    VAULT_JWT_TOKEN:
      aud: "${CI_SERVER_HOST}"

# job
# this job inherits the service container and before script block,
# hence all secrets in VKV_SERVER_PATH are available in your shell
test:
  stage: test
  script:
    - make test 
```

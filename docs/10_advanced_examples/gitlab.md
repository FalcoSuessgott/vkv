# Gitlab CI

Gitlab-CI Example for reading Secrets from Vault using vkv

```yaml
test:
  stage: test
  image: 
    name: ghcr.io/falcosuessgott/vkv:v0.2.2
    entrypoint: [""]
  variables:
    VAULT_ADDR: https://prod.vault.d4.sva.dev
    VAULT_TOKEN: hvs.xxx 
  before_script:
    # use eval in order to make all secrets avaible as env vars
    - eval $(vkv export --path secret --format=export)
  script:
    - make test 
```
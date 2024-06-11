# direnv
You can use `vkv` and [`direnv`](https://direnv.net/) to autimatically source KV secrets in your shell.

## Prerequisites
* Install [direnv](https://direnv.net/) and hook into your shell

## Demo

Create in a project a `.envrc` file:

```bash
> export VAULT_ADDR="https://vault:8200"
> export VAULT_TOKEN="$(cat ~/.vault-token)"
> eval $(vkv export -p kv/secrets -f export)
```

Now if you go into that directory and run `direnv allow`, 
you have the secrets under `kv/secrets` exported as env various:

```bash
> env | grep OS_
OS_USER=admin
OS_PASSWORD=pasword
```

## Upper Case the env var keys
`--export-upper` will upper-case all env var keys:

```bash
> vkv export -p secret/sub/demo -f export --export-upper
export DEMO='hello world'
export PASSWORD='s3cre5<'
```

## Include paths as env var key prefix
`--export-include-path` will prefix the secrets path and use `_` as the delimiter:

```bash
> vkv export -p secret/sub/demo -f export --export-upper --export-include-path
export SECRET_SUB_DEMO_DEMO='hello world'
export SECRET_SUB_DEMO_PASSWORD='s3cre5<'
export SECRET_SUB_DEMO_USER='admin'
```

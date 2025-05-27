# Authentication

`vkv` supports all of Vaults [environment variables](https://www.vaultproject.io/docs/commands#environment-variables) as well as any configured [Token helpers](https://developer.hashicorp.com/vault/docs/commands/token-helper).

In order to authenticate you will have to set at least one of the `VAULT_ADDR` or `VKV_LOGIN_COMMAND` and `VAULT_TOKEN` env vars.

## MacOS/Linux
```
export VAULT_ADDR="http://127.0.0.1:8200"
export VAULT_TOKEN="hvs.XXX"
vkv export --path <KVv2-path>
```

## Windows
```
SET VAULT_ADDR=http://127.0.0.1:8200
SET VAULT_TOKEN=s.XXX
vkv.exe export --path <KVv2-path>
```

## Special Env Var `VKV_LOGIN_COMMAND`
For advanced use cases, you can set `VKV_LOGIN_COMMAND`, that way `vkv` will first execute the specified command and use the output of the command as the token.
This is way you don't have to hardcode and set `VAULT_TOKEN`, this is especially useful when using `vkv` in CI. (See Gitlab Integration):

Example:

```bash
export VKV_LOGIN_COMMAND="vault write -field=token auth/jwt/login jwt=${CI_JOB_JWT_V2}"
vkv export -p
```

## Token Precedence
The following token precedence is applied (from highest to lowest):

1. `VAULT_TOKEN`
2. `VKV_LOGIN_COMMAND`
3. [Vault Token Helper](https://developer.hashicorp.com/vault/docs/commands/token-helper), where the token will be written to `~/.vault-token`.

If `vkv` detects **more than one possible token source**, warnings are shown as the following, indicating which token source will be used:

```bash
$> vkv export -p secret
[WARN] More than one token source configured (either VAULT_TOKEN, VKV_LOGIN_COMMAND or ~/.vault-token).
[WARN] See https://falcosuessgott.github.io/vkv/authentication for vkv's token precedence logic. Disable these warnings with VKV_DISABLE_WARNING.
[INFO] Using VAULT_TOKEN.

secret/ [desc=key/value secret storage] [type=kv2]
└── secret [v=1]
    └── key=*****
```

As described, one can disable these warning by setting `VKV_DISABLE_WARNING` to any value.

## Vault Token Lease Renewal
Depending on the number of Namespaces, KV mounts and secrets of your Vault and your Token TTL settings the lease of the `VAULT_TOKEN` being used may expire during a `vkv snapshot` operation (reported in [#363](https://github.com/FalcoSuessgott/vkv/issues/363)).

!!! important
    **To avoid that `vkv` automatically attempts to refresh the lease of the token (ref. [https://developer.hashicorp.com/vault/docs/concepts/lease](https://developer.hashicorp.com/vault/docs/concepts/lease)) being used.**
    
    This should only affect users of large enterprise Vaults.


Per default `vkv` will attempt to compare every `10s` (change with `VKV_RENEWAL_INTERVAL`) the current token TTL with the original creation TTL and if the current TTL less than half the creation TTL, a lease token renewal for another `30s` (change with `VKV_RENEWAL_INCREMENT`) is performed. `vkv` will error silently to not affect any JSON/YAML output.

You can find the exact implementation [here](https://github.com/FalcoSuessgott/vkv/blob/master/pkg/vault/lease.go).

!!! tip
    **You can always disable the token lease renewal by exporting `VKV_LEASE_REFRESHER_ENABLED`**
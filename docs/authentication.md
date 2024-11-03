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

1. `VKV_TOKEN`
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

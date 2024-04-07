# Authentication

`vkv` supports all of Vaults [environment variables](https://www.vaultproject.io/docs/commands#environment-variables). In order to authenticate you will have to set at least `VAULT_ADDR` and `VAULT_TOKEN`.

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
This is way you dont have to hardcode and set `VAULT_TOKEN`, this is especially useful when using `vkv` in CI. (See Gitlab Integration):

Example:

```bash
export VKV_LOGIN_COMMAND="vault write -field=token auth/jwt/login jwt=${CI_JOB_JWT_V2}"
vkv export -p
```